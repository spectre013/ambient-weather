package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var (
	logger = logrus.New()
	db     *sql.DB
	client mqtt.Client
	loc    *time.Location

	// debounce: skip notifications we've already processed in this process.
	// The DB's stat_ingest_log handles cross-restart idempotency; this is
	// just an in-memory short-circuit for retries inside a single run.
	lastNotificationMu sync.Mutex
	lastNotificationTS time.Time
)

func init() {
	logger.Out = os.Stdout
	logger.SetLevel(logrus.InfoLevel)
}

func main() {
	var err error
	if os.Getenv("GO_ENV") != "production" {
		logger.Infof("Loading environment variables")
		if err = godotenv.Load(); err != nil {
			logger.Fatal("Error loading .env file")
		}
	}

	if os.Getenv("LOGLEVEL") == "Debug" {
		logger.SetLevel(logrus.DebugLevel)
	}
	logger.Info("Log level: ", logger.GetLevel())

	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		panic(err)
	}
	logger.Info("Connected to Postgres as ")

	// --- timezone ---
	loc, err = time.LoadLocation("America/Denver")
	if err != nil {
		logger.Fatal(err)
	}

	// --- one-shot backfill if requested ---
	if os.Getenv("STATS_BACKFILL") == "true" {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
		logger.Info("STATS_BACKFILL=true: rebuilding stats from records ...")
		if err := rebuildStatsFromRecords(ctx); err != nil {
			logger.Error("backfill failed: ", err)
		} else {
			logger.Info("Backfill complete")
		}
		cancel()
	}

	// --- MQTT ---
	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_BROKER"))
	opts.SetClientID(os.Getenv("MQTT_CLIENTID"))
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(onMQTTConnect)
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		logger.Warn("MQTT connection lost: ", err)
	})

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	logger.Info("Connected to MQTT broker ", os.Getenv("MQTT_HOST"))

	// --- scheduler ---
	s := gocron.NewScheduler(loc)

	if logger.GetLevel() == logrus.DebugLevel {
		logger.Info("Debug mode: every-minute schedules")
		if _, err := s.Every(1).Minute().Do(updateAlerts); err != nil {
			logger.Error(err)
		}
		if _, err := s.Every(1).Minute().Do(getForecast); err != nil {
			logger.Error(err)
		}
		if _, err := s.Every(5).Minute().Do(runDailyRebuild); err != nil {
			logger.Error(err)
		}
	} else {
		logger.Info("Alert cron: ", os.Getenv("ALERT_CRON"))
		if _, err := s.Cron(os.Getenv("ALERT_CRON")).Do(updateAlerts); err != nil {
			logger.Error(err)
		}
		logger.Info("Forecast cron: ", os.Getenv("FORECAST_CRON"))
		if _, err := s.Cron(os.Getenv("FORECAST_CRON")).Do(getForecast); err != nil {
			logger.Error(err)
		}
		rebuildCron := os.Getenv("STATS_REBUILD_CRON")
		if rebuildCron == "" {
			rebuildCron = "5 0 * * *" // 00:05 local time
		}
		logger.Info("Stats rebuild cron: ", rebuildCron)
		if _, err := s.Cron(rebuildCron).Do(runDailyRebuild); err != nil {
			logger.Error(err)
		}

		// startup
		go updateAlerts()
		go getForecast()
	}

	s.StartAsync()
	logger.Info("Service running; awaiting MQTT notifications")
	select {}
}

// runDailyRebuild rebuilds all stats from the records table and prunes the
// ingest log. Cheap insurance against drift, missed messages, or arithmetic
// edge cases.
func runDailyRebuild() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	start := time.Now()
	if err := rebuildStatsFromRecords(ctx); err != nil {
		logger.Error("daily stats rebuild: ", err)
		return
	}
	logger.Info("Daily stats rebuild complete in ", time.Since(start))

	if err := pruneIngestLog(ctx); err != nil {
		logger.Error("prune ingest log: ", err)
	}
}

// onMQTTConnect (re)subscribes whenever the client connects or reconnects.
func onMQTTConnect(c mqtt.Client) {
	topic := os.Getenv("MQTT_STATS_SUBSCRIBE")
	if token := c.Subscribe(topic, 0, statsNotificationHandler); token.Wait() && token.Error() != nil {
		logger.Error("MQTT subscribe: ", token.Error())
		return
	}
	logger.Info("Subscribed to MQTT topic ", topic)
}

// ----------------------------------------------------------------------------
// MQTT notification
// ----------------------------------------------------------------------------

// StatsNotification is the payload published by the ingest service after it
// writes a record. Example:
//
//	{"id":"WS-5000","timestamp":"2026-05-18T15:24:36.578238-06:00"}
type StatsNotification struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

var statsNotificationHandler mqtt.MessageHandler = func(_ mqtt.Client, msg mqtt.Message) {
	var n StatsNotification
	if err := json.Unmarshal(msg.Payload(), &n); err != nil {
		logger.Error("invalid StatsNotification payload: ", err, " raw=", string(msg.Payload()))
		return
	}

	// in-memory debounce for retries within a single process lifetime
	lastNotificationMu.Lock()
	if !n.Timestamp.IsZero() && n.Timestamp.Equal(lastNotificationTS) {
		lastNotificationMu.Unlock()
		logger.Debug("Skipping duplicate notification for ", n.Timestamp)
		return
	}
	lastNotificationTS = n.Timestamp
	lastNotificationMu.Unlock()

	logger.Info("Stats notification: id=", n.ID, " ts=", n.Timestamp)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := foldRecord(ctx, n.Timestamp); err != nil {
			logger.Error("foldRecord: ", err)
		}
	}()
}

// ----------------------------------------------------------------------------
// Forecast
// ----------------------------------------------------------------------------

func getForecast() {
	includes := "days%2Chours"
	url := fmt.Sprintf(
		"https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Colorado%%20Springs?unitGroup=us&iconSets=icon2&include=%s&key=%s&contentType=json",
		includes, os.Getenv("WEATHER_API"),
	)
	res, err := makeRequest(url, nil)
	if err != nil {
		logger.Error("getForecast: ", err)
		return
	}

	var f Forecast
	if err := json.Unmarshal(res, &f); err != nil {
		logger.Error("getForecast unmarshal: ", err)
		return
	}
	logger.Info("Processing forecast days: ", len(f.Days))

	for _, v := range f.Days {
		day, err := convertDayToDB(v)
		if err != nil {
			logger.Error("convertDayToDB: ", err)
			continue
		}
		day.Summary, err = getForecastSummary(v)
		if err != nil {
			logger.Error("getForecastSummary: ", err)
		}
		if err := insertForecast(day); err != nil {
			logger.Error("insertForecast: ", err)
		}
	}
}

func getForecastSummary(day Day) (string, error) {
	prompt := fmt.Sprintf(
		`write a forecast summary based on the following forecast data do not use any
markdown or prepend any response before or after the summary.
High: %.2f°F, Low: %.2f°F, Dewpoint: %.2f°F, Humidity %.2f%%, Visibility %.2f mi,
Wind Speed: %.2f MPH, Wind Direction %.2f, Wind Gusts: %.2f MPH,
preciptype: %s, expected Precipitation: %.2f in, Precipitation Prob: %.2f%%`,
		day.Tempmax, day.Tempmin, day.Dew, day.Humidity, day.Visibility,
		day.Windspeed, day.Winddir, day.Windgust,
		strings.Join(day.Preciptype, ","), day.Precip, day.Precipprob,
	)
	return summerize(prompt)
}

func insertForecast(f forecastDB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	const query = `
	INSERT INTO public.vcforecast (
	   datetime, datetime_epoch, tempmax, tempmin, temp, feelslikemax, feelslikemin,
	   feelslike, dew, humidity, precip, precipprob, precipcover, preciptype,
	   snow, snowdepth, windgust, windspeed, winddir, pressure, cloudcover,
	   visibility, solarradiation, solarenergy, uvindex, severerisk, sunrise,
	   sunrise_epoch, sunset, sunset_epoch, moonphase, conditions, description,
	   icon, stations, source, hours, summary
	) VALUES (
	   $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
	   $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32,
	   $33, $34, $35, $36, $37, $38
	)
	ON CONFLICT (datetime) DO UPDATE SET
	   datetime_epoch = EXCLUDED.datetime_epoch,
	   tempmax        = EXCLUDED.tempmax,
	   tempmin        = EXCLUDED.tempmin,
	   temp           = EXCLUDED.temp,
	   feelslikemax   = EXCLUDED.feelslikemax,
	   feelslikemin   = EXCLUDED.feelslikemin,
	   feelslike      = EXCLUDED.feelslike,
	   dew            = EXCLUDED.dew,
	   humidity       = EXCLUDED.humidity,
	   precip         = EXCLUDED.precip,
	   precipprob     = EXCLUDED.precipprob,
	   precipcover    = EXCLUDED.precipcover,
	   preciptype     = EXCLUDED.preciptype,
	   snow           = EXCLUDED.snow,
	   snowdepth      = EXCLUDED.snowdepth,
	   windgust       = EXCLUDED.windgust,
	   windspeed      = EXCLUDED.windspeed,
	   winddir        = EXCLUDED.winddir,
	   pressure       = EXCLUDED.pressure,
	   cloudcover     = EXCLUDED.cloudcover,
	   visibility     = EXCLUDED.visibility,
	   solarradiation = EXCLUDED.solarradiation,
	   solarenergy    = EXCLUDED.solarenergy,
	   uvindex        = EXCLUDED.uvindex,
	   severerisk     = EXCLUDED.severerisk,
	   sunrise        = EXCLUDED.sunrise,
	   sunrise_epoch  = EXCLUDED.sunrise_epoch,
	   sunset         = EXCLUDED.sunset,
	   sunset_epoch   = EXCLUDED.sunset_epoch,
	   moonphase      = EXCLUDED.moonphase,
	   conditions     = EXCLUDED.conditions,
	   description    = EXCLUDED.description,
	   icon           = EXCLUDED.icon,
	   stations       = EXCLUDED.stations,
	   source         = EXCLUDED.source,
	   hours          = EXCLUDED.hours,
	   summary        = EXCLUDED.summary`

	_, err := db.ExecContext(ctx, query,
		f.Datetime, f.DatetimeEpoch, f.TempMax, f.TempMin, f.Temp,
		f.FeelsLikeMax, f.FeelsLikeMin, f.FeelsLike, f.Dew, f.Humidity,
		f.Precip, f.PrecipProb, f.PrecipCover, f.PrecipType,
		f.Snow, f.SnowDepth, f.WindGust, f.WindSpeed, f.WindDir,
		f.Pressure, f.CloudCover, f.Visibility, f.SolarRadiation,
		f.SolarEnergy, f.UVIndex, f.SevereRisk, f.Sunrise,
		f.SunriseEpoch, f.Sunset, f.SunsetEpoch, f.MoonPhase,
		f.Conditions, f.Description, f.Icon, f.Stations, f.Source,
		f.Hours, f.Summary,
	)
	if err != nil {
		return fmt.Errorf("upsert failed for datetime %v: %w", f.Datetime, err)
	}
	return nil
}

func convertDayToDB(d Day) (forecastDB, error) {
	hoursJSON, err := json.Marshal(d.Hours)
	if err != nil {
		return forecastDB{}, err
	}
	datetime, err := time.Parse("2006-01-02", d.Datetime)
	if err != nil {
		datetime = time.Now()
	}
	return forecastDB{
		Datetime:       datetime,
		DatetimeEpoch:  d.DatetimeEpoch,
		TempMax:        d.Tempmax,
		TempMin:        d.Tempmin,
		Temp:           d.Temp,
		FeelsLikeMax:   d.Feelslikemax,
		FeelsLikeMin:   d.Feelslikemin,
		FeelsLike:      d.Feelslike,
		Dew:            d.Dew,
		Humidity:       d.Humidity,
		Precip:         d.Precip,
		PrecipProb:     d.Precipprob,
		PrecipCover:    d.Precipcover,
		PrecipType:     strings.Join(d.Preciptype, ","),
		Snow:           d.Snow,
		SnowDepth:      d.Snowdepth,
		WindGust:       d.Windgust,
		WindSpeed:      d.Windspeed,
		WindDir:        d.Winddir,
		Pressure:       d.Pressure,
		CloudCover:     d.Cloudcover,
		Visibility:     d.Visibility,
		SolarRadiation: d.Solarradiation,
		SolarEnergy:    d.Solarenergy,
		UVIndex:        d.Uvindex,
		SevereRisk:     d.Severerisk,
		Sunrise:        d.Sunrise,
		SunriseEpoch:   d.SunriseEpoch,
		Sunset:         d.Sunset,
		SunsetEpoch:    d.SunsetEpoch,
		MoonPhase:      d.Moonphase,
		Conditions:     d.Conditions,
		Description:    d.Description,
		Icon:           d.Icon,
		Stations:       strings.Join(d.Stations, ","),
		Source:         d.Source,
		Hours:          string(hoursJSON),
	}, nil
}

// ----------------------------------------------------------------------------
// Alerts
// ----------------------------------------------------------------------------

func updateAlerts() {
	logger.Info("Updating alerts ...")

	const checkSQL = `SELECT EXISTS(SELECT 1 FROM alerts WHERE alertid = $1)`
	const insertSQL = `
        INSERT INTO alerts (
            alertid, wxtype, areadesc, sent, effective, onset, expires, ends, status,
            messagetype, category, severity, certainty, urgency, event, sender,
            sendername, headline, description, instruction, response, summary
        )
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
        ON CONFLICT (alertid) DO UPDATE SET
            wxtype      = EXCLUDED.wxtype,
            areadesc    = EXCLUDED.areadesc,
            sent        = EXCLUDED.sent,
            effective   = EXCLUDED.effective,
            onset       = EXCLUDED.onset,
            expires     = EXCLUDED.expires,
            ends        = EXCLUDED.ends,
            status      = EXCLUDED.status,
            messagetype = EXCLUDED.messagetype,
            category    = EXCLUDED.category,
            severity    = EXCLUDED.severity,
            certainty   = EXCLUDED.certainty,
            urgency     = EXCLUDED.urgency,
            event       = EXCLUDED.event,
            sender      = EXCLUDED.sender,
            sendername  = EXCLUDED.sendername,
            headline    = EXCLUDED.headline,
            description = EXCLUDED.description,
            instruction = EXCLUDED.instruction,
            response    = EXCLUDED.response,
            summary     = EXCLUDED.summary`

	for _, v := range getAlerts() {
		var exists bool
		if err := db.QueryRow(checkSQL, v.ID).Scan(&exists); err != nil {
			logger.Errorf("alert exists check %s: %v", v.ID, err)
			continue
		}
		if exists && !strings.EqualFold(v.MessageType, "Update") {
			continue
		}

		prompt := fmt.Sprintf(
			"write a short summary of the following text only include the summary in your response: Severity: %s Event: %s Headline: %s Description: %s Instructions: %s",
			v.Severity, v.Event, v.Headline, v.Description, v.Instruction,
		)
		summary, err := summerize(prompt)
		if err != nil {
			logger.Errorf("summarize %s: %v", v.ID, err)
			continue
		}

		if _, err := db.Exec(insertSQL,
			v.ID, v.Event, v.AreaDesc, v.Sent, v.Effective, v.Onset,
			v.Expires, v.Ends, v.Status, v.MessageType, v.Category, v.Severity,
			v.Certainty, v.Urgency, v.Event, v.Sender, v.SenderName,
			v.Headline, v.Description, v.Instruction, v.Response, summary,
		); err != nil {
			logger.Errorf("alert insert %s: %v", v.ID, err)
			continue
		}
		logger.Infof("Inserted/updated alert %s (%s)", v.ID, v.Headline)
	}
}

func getAlerts() []Property {
	uri := "https://api.weather.gov/alerts/active/area/CO"
	result := make([]Property, 0)

	res, err := alertRequest("GET", uri)
	if err != nil {
		logger.Error(err)
		return result
	}
	var alerts Alerts
	if err := json.Unmarshal(res, &alerts); err != nil {
		logger.Error(err)
		return result
	}

	zones := []string{"COZ084", "COZ085"}
	for _, v := range alerts.Features {
		if hasCommon(zones, v.Properties.Geocode.UGC) {
			result = append(result, v.Properties)
		}
	}
	return result
}

func hasCommon(a, b []string) bool {
	for _, v := range a {
		if slices.Contains(b, v) {
			return true
		}
	}
	return false
}

// ----------------------------------------------------------------------------
// HTTP clients (Ollama, NWS, VisualCrossing)
// ----------------------------------------------------------------------------

func summerize(prompt string) (string, error) {
	url := "http://10.10.1.120:11434/api/generate"
	payload := OllamaRequest{
		Model:  "llama3.1:8b",
		Prompt: prompt,
		Stream: false,
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("ollama: ", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result OllamaResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	return result.Response, nil
}

func alertRequest(method, url string) ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", `Zoms Weather, wxcos@zoms.net`)

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("server responded with status %d", resp.StatusCode)
	}
	if method == "head" {
		return nil, nil
	}
	return io.ReadAll(resp.Body)
}

func makeRequest(url string, header map[string]string) ([]byte, error) {
	c := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if _, ok := header["User-Agent"]; !ok {
		req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
	}
	for k, v := range header {
		req.Header.Add(k, v)
	}

	logger.Info("GET ", url)
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	logger.Info("status: ", resp.StatusCode)
	return io.ReadAll(resp.Body)
}
