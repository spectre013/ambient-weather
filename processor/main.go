package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var db *sql.DB
var client mqtt.Client

func init() {
	logger.Out = os.Stdout
	logger.SetLevel(logrus.InfoLevel)
}
func main() {
	var err error
	if os.Getenv("GO_ENV") != "production" {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	logLevel := logrus.InfoLevel
	if os.Getenv("LOGLEVEL") == "Debug" {
		logLevel = logrus.DebugLevel
	}
	logger.Info("Setting Debug Level to ", logLevel)
	logger.SetLevel(logLevel)

	dburi := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_DATABASE"))
	db, err = sql.Open("postgres", dburi)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	logger.Info("Connected to Postgres as ", os.Getenv("DB_USER"))

	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Println(err)
		return
	}

	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_HOST"))
	opts.SetClientID(os.Getenv("MQTT_CLIENTID"))
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	logger.Info("Connected to MQTT broker as ", os.Getenv("MQTT_HOST"))
	if token := client.Subscribe(os.Getenv("MQTT_SUBSCRIBE"), 0, messageHandler); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}
	logger.Info("Subscribed to MQTT Channel", os.Getenv("MQTT_SUBSCRIBE"))
	s := gocron.NewScheduler(loc)
	var aj *gocron.Job
	var fj *gocron.Job

	if os.Getenv("LOGLEVEL") == "Debug" {
		fmt.Println("Starting every minute alert update")

		aj, err = s.Every(1).Minute().Do(updateAlerts)
		if err != nil {
			logger.Error(err)
		}

		fj, err = s.Every(1).Minute().Do(getForecast)
		if err != nil {
			logger.Error(err)
		}

	} else {
		fmt.Println("Starting cron alert update", os.Getenv("ALERT_CRON"))
		aj, err = s.Cron(os.Getenv("ALERT_CRON")).Do(updateAlerts)
		if err != nil {
			logger.Error(err)
		}

		fmt.Println("Starting cron Forecast update", os.Getenv("FORECAST_CRON"))
		fj, err = s.Cron(os.Getenv("FORECAST_CRON")).Do(getForecast)
		if err != nil {
			logger.Error(err)
		}

		//kick of an alert and a forecast when we start the server.
		go updateAlerts()
		go getForecast()
	}
	if err != nil {
		logger.Error(err)
	}
	logger.Info(aj.NextRun())
	logger.Info(fj.NextRun())

	s.StartAsync()
	//calculateStats()
	e := echo.New()
	e.GET("/", index)
	e.Logger.Fatal(e.Start(":8081"))
}
func index(c echo.Context) error {
	return c.String(http.StatusOK, "Ambient Weather Receiver ")
}

func CalculateFeelsLike(temp float64, windspeed float64, humidity int) float64 {
	feelslike := temp
	if temp <= 50 && windspeed > 3 {
		feelslike = calculateWindChill(temp, windspeed)
	} else if temp >= 80 {
		feelslike = calculateHeatIndex(temp, float64(humidity))
	} else {
		feelslike = temp
	}
	return feelslike
}

func cToF(c float64) float64     { return toFixed((c*1.8)+32.0, 2) }
func msToMph(ms float64) float64 { return toFixed(ms*2.23694, 2) }

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	ambientin(msg.Payload())
}

func ambientin(b []byte) {
	rec := Record{}
	err := json.Unmarshal(b, &rec)
	if err != nil {
		log.Println(err)
	}

	lastrainquery := "select recorded from records r where dailyrainin > 0 order by recorded desc limit 1"

	crows := db.QueryRow(lastrainquery)
	var lrain time.Time
	err = crows.Scan(&lrain)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", lastrainquery)
		} else {
			logger.Error("Scan:", err)
		}
	}
	rec.Lastrain = lrain

	rec.Dewpoint = dewpoint(rec.Tempf, rec.Humidity)
	rec.Feelslike = CalculateFeelsLike(rec.Tempf, rec.Windspeedmph, rec.Humidity)

	logger.Info("Received Record ", formatDate(rec.Recorded))
	inserted := insertRecord(rec)
	if inserted {
		go calculateStats()
	}

}

func calculateWindChill(t, v float64) float64 {
	return 35.74 + (0.6215 * t) - (35.75 * math.Pow(v, 0.16)) + (0.4275 * t * math.Pow(v, 0.16))
}

func calculateHeatIndex(t, rh float64) float64 {
	hi := 0.5 * (t + 61.0 + ((t - 68.0) * 1.2) + (rh * 0.094))
	if hi > 80 {
		hi = -42.379 + 2.04901523*t + 10.14333127*rh - 0.22475541*t*rh - 0.00683783*t*t - 0.05481717*rh*rh + 0.00122874*t*t*rh + 0.00085282*t*rh*rh - 0.00000199*t*t*rh*rh
	}
	return hi
}

func dewpoint(temp float64, humidity int) float64 {
	temp = (temp - 32) * 5 / 9
	const a, b = 17.625, 243.04
	alpha := math.Log(float64(humidity)/100.0) + ((a * temp) / (b + temp))
	dewPointC := (b * alpha) / (a - alpha)
	return toFixed(cToF(dewPointC), 2)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func calculateStats() {
	type Query struct {
		Query  string
		Params []time.Time
	}

	types := []string{"MAX", "MIN", "AVG"}
	periods := []string{"day", "yesterday", "month", "year"}
	fields := []string{"tempf", "tempinf", "temp1f", "temp2f", "temp3f", "baromrelin", "uv", "humidity", "windspeedmph", "windgustmph", "dewpoint", "humidityin", "humidity1", "humidity2", "humidity3", "dailyrainin", "lightning", "aqipm25"}
	queries := make(map[string]Query)
	for _, p := range periods {
		for _, t := range types {
			for _, f := range fields {
				key := p + "_" + strings.ToLower(t) + "_" + f
				order := ""
				if t == "MAX" {
					order = " DESC"
				}
				if (f == "dailyrainin" && t == "MAX") || f != "dailyrainin" {
					if !strings.Contains(f, "lightning") {

						q := fmt.Sprintf("select '%s' as id, CAST(COALESCE(%s,0.0) AS decimal(10,2)) as value, recorded from records where recorded between ? and ? order by %s%s limit 1", key, f, f, order)
						if t == "AVG" {
							q = fmt.Sprintf("select '%s' as id,  CAST(COALESCE(%s(%s),0.0) AS decimal(10,2)) as value from records where recorded between ? and ? limit 1", key, t, f)
						}
						queries[key] = Query{
							Query:  q,
							Params: getTimeframe(p),
						}
					}
				}

				if strings.Contains(f, "lightning") && t == "MAX" {
					q := ""
					if p == "month" || p == "year" {
						q = fmt.Sprintf(`
						SELECT '%s' as id, SUM(A.value) as value
						FROM (SELECT TO_CHAR(recorded,'YYY-MM-DD') as ldate, CAST(COALESCE(MAX(lightningday),0.0) AS decimal(10,2)) as value FROM records where recorded between ? and ? GROUP BY ldate) A
						`, key)
					} else {
						q = fmt.Sprintf(`SELECT '%s' as id, CAST(COALESCE(lightningday,0.0) AS decimal(10,2)) as value, recorded FROM records where recorded between ? and ? order by value desc limit 1`, key)
					}
					queries[key] = Query{
						Query:  q,
						Params: getTimeframe(p),
					}
				}

			}
		}
	}
	start := time.Now()
	tx, err := db.Begin()
	if err != nil {
		logger.Error(err)
	}
	for k, v := range queries {
		v.Query = strings.Replace(v.Query, "?", "'%s'", -1)
		v.Query = fmt.Sprintf(v.Query, formatDate(v.Params[0]), formatDate(v.Params[1]))
		d := "recorded = src.recorded,"
		if strings.Contains(k, "avg") || (strings.Contains(v.Query, "MAX(lightningday)")) {
			d = ""
		}

		update := checkStat(k)
		if update {
			updateQuery := fmt.Sprintf(`
				UPDATE stats set
				%s
				value = src.value
				from (
					%s
					) AS src
				WHERE
					stats.id = '%s';
			`, d, v.Query, k)
			logger.Debug(updateQuery)
			_, err := tx.Exec(updateQuery)
			if err != nil {
				logger.Debug(updateQuery)
				logger.Error(err)
				break
			}
		} else {
			insert := "insert into stats (id,value,recorded)"
			if strings.Contains(k, "avg") || (strings.Contains(v.Query, "MAX(lightningday)")) {
				insert = "insert into stats (id,value)"
			}
			updateQuery := fmt.Sprintf(`
				%s
				%s
			`, insert, v.Query)
			logger.Debug(updateQuery)
			_, err := tx.Exec(updateQuery)
			if err != nil {
				logger.Debug(updateQuery)
				logger.Error(err)
				break
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		logger.Error(err)
	}

	elapsed := time.Since(start)
	log.Printf("Update took %s", elapsed)

}

func getTimeframe(timeframe string) []time.Time {
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		logger.Error(err)
	}
	var dates []time.Time
	now := time.Now()

	if timeframe == "yesterday" {
		dates = append(dates, time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, loc))
		dates = append(dates, time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, loc))
	} else if timeframe == "day" {
		dates = append(dates, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc))
		dates = append(dates, time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc))
	} else if timeframe == "month" {
		dates = append(dates, time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc))
		dates = append(dates, time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, loc))
	} else if timeframe == "year" {
		dates = append(dates, time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc))
		dates = append(dates, time.Date(now.Year(), 12, 31, 23, 59, 59, 0, loc))
	}
	return dates
}

func checkStat(id string) bool {
	rows := db.QueryRow(fmt.Sprintf("select id from stats where id = '%s'", id))
	var r string
	err := rows.Scan(&r)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			logger.Error("Scan:", err)
		}
	}
	return true
}

func formatDate(date time.Time) string {
	format := "2006-01-02 15:04:05 -0700"
	return date.Format(format)
}

func insertRecord(r Record) bool {
	query := fmt.Sprintf(`insert into records (id,mac,recorded,baromabsin,baromrelin,battout,batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,hourlyrainin,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,hourlyrain,lightningday,lightninghour,lightningtime,lightningdistance,battlightning, aqipm25, aqipm2524h) values
			(DEFAULT,'%s','%s',%f,%f,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%f,%f,%f,%f,%f,%f,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,'%s',%f,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%d,%f,%d,%f,%f,%f,%d,%d,'%s',%f,%d,%d,%d)				
	`, r.Mac, formatDate(r.Recorded), r.Baromabsin, r.Baromrelin, r.Battout, r.Batt1, r.Batt2, r.Batt3, r.Batt4, r.Batt5, r.Batt6, r.Batt7, r.Batt8, r.Batt9, r.Batt10, r.Co2, r.Dailyrainin, r.Dewpoint, r.Eventrainin, r.Feelslike, r.Hourlyrainin, r.Humidity, r.Humidity1, r.Humidity2, r.Humidity3, r.Humidity4, r.Humidity5, r.Humidity6, r.Humidity7, r.Humidity8, r.Humidity9, r.Humidity10, r.Humidityin, formatDate(r.Lastrain), r.Maxdailygust, r.Relay1, r.Relay2, r.Relay3, r.Relay4, r.Relay5, r.Relay6, r.Relay7, r.Relay8, r.Relay9, r.Relay10, r.Monthlyrainin, r.Solarradiation, r.Tempf, r.Temp1f, r.Temp2f, r.Temp3f, r.Temp4f, r.Temp5f, r.Temp6f, r.Temp7f, r.Temp8f, r.Temp9f, r.Temp10f, r.Tempinf, r.Totalrainin, r.Uv, r.Weeklyrainin, r.Winddir, r.Windgustmph, r.Windgustdir, r.Windspeedmph, r.Yearlyrainin, r.Hourlyrain, r.Lightningday, r.Lightninghour, formatDate(r.Lightningtime), r.Lightningdistance, r.Battlightning, r.Aqipm25, r.Aqipm2524h)
	logger.Debug(query)
	_, err := db.Exec(query)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func getForecast() {
	includes := "days%2Chours"
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Colorado%%20Springs?unitGroup=us&iconSets=icon2&include=%s&key=%s&contentType=json", includes, os.Getenv("WEATHER_API"))
	header := map[string]string{}
	res, err := makeRequest(url, header)
	if err != nil {
		logger.Error("Error in Get Forecast", err)
	}
	logger.Info("Received Forecast")
	f := Forecast{}
	err = json.Unmarshal(res, &f)
	if err != nil {
		logger.Error("Error in Unmarshall Forecast", err)
	}
	logger.Info("Processing Forecast days: ", len(f.Days))
	for _, v := range f.Days {
		day, err := convertDayToDB(v)
		if err != nil {
			logger.Error("Error in convertDayToDB", err)
		}
		day.Summary, err = getForecastSummary(v)
		if err != nil {
			logger.Error("Error in getForecastSummary", err)
		}
		err = insertForecast(day)
		if err != nil {
			logger.Error("Error in insertForecast", err)
		}
	}
}

func getForecastSummary(day Day) (string, error) {
	event := fmt.Sprintf(`write a forcast summary based on the following forecast data do not use any 
                                 markdown or prepend any response before or after the summary.\n
						        High: %.2f°F, Low: %.2f°F, Dewpoint: %.2f°F, Humidity %.2f%%, Visibility %.2f mi, 
                                Wind Speed: %.2f MPH, Wind Direction %.2f, Wind Gusts: %.2f MPH, 
                                preciptype: %s, expected Precipitation: %.2f in, Precipitation Prob: %.2f%`,
		day.Tempmax, day.Tempmin, day.Dew, day.Humidity, day.Visibility, day.Windspeed,
		day.Winddir, day.Windgust, strings.Join(day.Preciptype, ","), day.Precip,
		day.Precipprob)
	summary, err := summerize(event)
	if err != nil {
		logger.Error(err)
		summary = ""
	}
	return summary, nil
}

func insertForecast(f forecastDB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
	INSERT INTO forecast (
		datetime, datetime_epoch, tempmax, tempmin, temp, feelslikemax, feelslikemin, 
		feelslike, dew, humidity, precip, precipprob, precipcover, preciptype, 
		snow, snowdepth, windgust, windspeed, winddir, pressure, cloudcover, 
		visibility, solarradiation, solarenergy, uvindex, severerisk, sunrise, 
		sunrise_epoch, sunset, sunset_epoch, moonphase, conditions, description, 
		icon, stations, source, hours,summary
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, 
		$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, 
		$33, $34, $35, $36, $37,$38
	)
	ON CONFLICT (datetime) DO UPDATE SET
		datetime_epoch = EXCLUDED.datetime_epoch,
		tempmax = EXCLUDED.tempmax,
		tempmin = EXCLUDED.tempmin,
		temp = EXCLUDED.temp,
		feelslikemax = EXCLUDED.feelslikemax,
		feelslikemin = EXCLUDED.feelslikemin,
		feelslike = EXCLUDED.feelslike,
		dew = EXCLUDED.dew,
		humidity = EXCLUDED.humidity,
		precip = EXCLUDED.precip,
		precipprob = EXCLUDED.precipprob,
		precipcover = EXCLUDED.precipcover,
		preciptype = EXCLUDED.preciptype,
		snow = EXCLUDED.snow,
		snowdepth = EXCLUDED.snowdepth,
		windgust = EXCLUDED.windgust,
		windspeed = EXCLUDED.windspeed,
		winddir = EXCLUDED.winddir,
		pressure = EXCLUDED.pressure,
		cloudcover = EXCLUDED.cloudcover,
		visibility = EXCLUDED.visibility,
		solarradiation = EXCLUDED.solarradiation,
		solarenergy = EXCLUDED.solarenergy,
		uvindex = EXCLUDED.uvindex,
		severerisk = EXCLUDED.severerisk,
		sunrise = EXCLUDED.sunrise,
		sunrise_epoch = EXCLUDED.sunrise_epoch,
		sunset = EXCLUDED.sunset,
		sunset_epoch = EXCLUDED.sunset_epoch,
		moonphase = EXCLUDED.moonphase,
		conditions = EXCLUDED.conditions,
		description = EXCLUDED.description,
		icon = EXCLUDED.icon,
		stations = EXCLUDED.stations,
		source = EXCLUDED.source,
		hours = EXCLUDED.hours,
		summary = EXCLUDED.summary;`

	_, err := db.ExecContext(ctx, query,
		f.Datetime, f.DatetimeEpoch, f.TempMax, f.TempMin, f.Temp, f.FeelsLikeMax, f.FeelsLikeMin,
		f.FeelsLike, f.Dew, f.Humidity, f.Precip, f.PrecipProb, f.PrecipCover, f.PrecipType,
		f.Snow, f.SnowDepth, f.WindGust, f.WindSpeed, f.WindDir, f.Pressure, f.CloudCover,
		f.Visibility, f.SolarRadiation, f.SolarEnergy, f.UVIndex, f.SevereRisk, f.Sunrise,
		f.SunriseEpoch, f.Sunset, f.SunsetEpoch, f.MoonPhase, f.Conditions, f.Description,
		f.Icon, f.Stations, f.Source, f.Hours, f.Summary,
	)

	if err != nil {
		return fmt.Errorf("upsert failed for datetime %v: %w", f.Datetime, err)
	}

	return nil
}

func convertDayToDB(d Day) (forecastDB, error) {
	// Serialize the Hours slice to a JSON string
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

func updateAlerts() {
	fmt.Println("Updating Alerts ... ")
	insertSql := `
		INSERT INTO alerts (
			id, wxtype, areadesc, sent, effective, onset, expires, ends, status,
			messagetype, category, severity, certainty, urgency, event, sender, 
			senderName, headline, description, instruction, response, summary
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21,$22)
		ON CONFLICT (id) 
		DO UPDATE SET 
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
			senderName  = EXCLUDED.senderName,
			headline    = EXCLUDED.headline,
			description = EXCLUDED.description,
			instruction = EXCLUDED.instruction,
			response    = EXCLUDED.response;
	`

	iAlerts := getAlerts()
	for _, v := range iAlerts {
		checkSql := fmt.Sprintf(`SELECT EXISTS(SELECT 1 FROM alerts WHERE id = '%s')`, v.ID)
		rows := db.QueryRow(checkSql)
		var exists bool
		err := rows.Scan(&exists)
		if err != nil {
			logger.Error(fmt.Sprintf("Scan: %v", err))
		}
		if !exists || v.MessageType == "updated" {
			event := fmt.Sprintf("write a short summery of the following text only include the summary in your response: Severity: %s Event: %s Headline:%s Description: %s Instructions: %s", v.Severity, v.Event, v.Headline, v.Description, v.Instruction)
			summary, err := summerize(event)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			logger.Info(v.ID, summary)
			_, err = db.Exec(insertSql, v.ID, v.Type, v.AreaDesc, v.Sent, v.Effective, v.Onset,
				v.Expires, v.Ends, v.Status, v.MessageType, v.Category, v.Severity, v.Certainty, v.Urgency, v.Event,
				v.Sender, v.SenderName, v.Headline, v.Description, v.Instruction, v.Response, summary)
			logger.Info(fmt.Sprintf("Inserted Alert %s", v.Headline))
			if err != nil {
				logger.Error(err)
			}
		}

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
	alerts := Alerts{}
	err = json.Unmarshal(res, &alerts)
	if err != nil {
		logger.Error(err)
	}

	zones := []string{"COZ084", "COZ085"}
	for _, v := range alerts.Features {
		if hasCommon(zones, v.Properties.Geocode.UGC) {
			result = append(result, v.Properties)
		}
	}

	return result
}

func hasCommon(slice1, slice2 []string) bool {
	for _, v := range slice1 {
		if slices.Contains(slice2, v) {
			return true
		}
	}
	return false
}

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
		logger.Error("Error asking AI", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result OllamaResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	return result.Response, nil
}

func alertRequest(t string, url string) (body []byte, err error) {
	body = []byte("")
	client := &http.Client{}
	req, err := http.NewRequest(t, url, nil)
	if err != nil {
		logger.Error(err)
	}

	req.Header.Add("User-Agent", `Zoms Weather, wxcos@zoms.net`)
	logger.Debug(url)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		err = errors.New("server responded with an error")
		return body, err
	}

	logger.Debug("Alert Updates")
	if t != "head" {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			logger.Error(err)
			return
		}
		logger.Debug(fmt.Sprintf("%v", resp.Header))
	}

	return
}

func makeRequest(url string, header map[string]string) ([]byte, error) {
	logger.Debug(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if _, ok := header["User-Agent"]; !ok {
		req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
	}
	if len(header) > 0 {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}
	logger.Debug(req.Header)
	logger.Info("Sending request: ", url)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	logger.Info("Received Response with status code: ", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return body, nil
}
