package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var logger = logrus.New()
var db *sql.DB
var LastModified time.Time

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

	s := gocron.NewScheduler(loc)
	var aj *gocron.Job

	if os.Getenv("LOGLEVEL") == "Debug" {
		fmt.Println("Starting every minute alert update")

		aj, err = s.Every(1).Minute().Do(updateAlerts)
		if err != nil {
			logger.Error(err)
		}

	} else {
		fmt.Println("Starting cron alert update", os.Getenv("ALERT_CRON"))
		aj, err = s.Cron(os.Getenv("ALERT_CRON")).Do(updateAlerts)
		if err != nil {
			logger.Error(err)
		}
	}
	if err != nil {
		logger.Error(err)
	}
	logger.Info(aj.NextRun())

	s.StartAsync()

	e := echo.New()
	e.GET("/", index)
	e.GET("/api/receiver", ambientin)
	e.Logger.Fatal(e.Start(":8080"))
}
func index(c echo.Context) error {
	return c.String(http.StatusOK, "Ambient Weather Receiver ")
}
func ambientin(c echo.Context) error {
	logger.Info(c.Request().URL.String())
	output := map[string]interface{}{}

	in := map[string]string{}
	in["baromabsin"] = "float"
	in["baromrelin"] = "float"
	in["batt_co2"] = "int"
	in["battlightning"] = "int"
	in["batt1"] = "int"
	in["batt2"] = "int"
	in["batt3"] = "int"
	in["batt4"] = "int"
	in["battin"] = "int"
	in["battout"] = "int"
	in["dailyrainin"] = "float"
	in["dateutc"] = "string"
	in["eventrainin"] = "float"
	in["hourlyrainin"] = "float"
	in["humidity"] = "int"
	in["humidity1"] = "int"
	in["humidity2"] = "int"
	in["humidity3"] = "int"
	in["humidity4"] = "int"
	in["humidityin"] = "int"
	in["lightningday"] = "int"
	in["lightningdistance"] = "int"
	in["lightningtime"] = "string"
	in["maxdailygust"] = "float"
	in["monthlyrainin"] = "float"
	in["solarradiation"] = "float"
	in["temp1f"] = "float"
	in["temp2f"] = "float"
	in["temp3f"] = "float"
	in["temp4f"] = "float"
	in["tempf"] = "float"
	in["tempinf"] = "float"
	in["uv"] = "int"
	in["weeklyrainin"] = "float"
	in["winddir"] = "int"
	in["windgustmph"] = "float"
	in["windspeedmph"] = "float"
	in["yearlyrainin"] = "float"
	in["aqipm25"] = "int"
	in["aqipm2524h"] = "int"

	values := c.QueryParams()
	if len(values) == 0 {
		return errors.New("no values received")
	}
	for k, v := range values {
		k = strings.Replace(k, "_", "", -1)
		val := v[0]
		switch in[k] {
		case "int":
			i, err := strconv.Atoi(val)
			if err != nil {
				log.Printf("%s - %s\n", err, val)
			}
			output[k] = i
			logger.Debug(k, " - ", i)
		case "float":
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Printf("%s - %s\n", err, val)
			}
			output[k] = toFixed(f, 2)
			logger.Debug(k, " - ", toFixed(f, 2))
		default:
			if k == "PASSKEY" {
				k = "mac"
			}

			output[k] = v[0]

			if k == "lightningtime" {
				i, err := strconv.ParseInt(v[0], 10, 64)
				if err != nil {
					panic(err)
				}
				output[k] = time.Unix(i, 0)
			}
			logger.Debug(k, " - ", v[0])

		}
	}
	output["date"] = time.Now()
	lastrainquery := "select recorded from records r where dailyrainin > 0 order by recorded desc limit 1"

	crows := db.QueryRow(lastrainquery)
	var lrain time.Time
	err := crows.Scan(&lrain)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", lastrainquery)
		} else {
			logger.Error("Scan:", err)
		}
	}
	output["lastrain"] = lrain

	b, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
		return err
	}

	rec := Record{}
	err = json.Unmarshal(b, &rec)
	if err != nil {
		log.Println(err)
		return err
	}

	rec.Dewpoint = dewpoint(rec.Tempf, rec.Humidity)
	if rec.Tempf >= 70 {
		rec.Feelslike = heatIndex(rec.Tempf, rec.Humidity)
	} else {
		if rec.Windgustmph > 3 {
			rec.Feelslike = windChill(rec.Tempf, rec.Windspeedmph)
		} else {
			rec.Feelslike = rec.Tempf
		}
	}
	logger.Info("Received Record ", formatDate(rec.Recorded))
	inserted := insertRecord(rec)
	if inserted {
		go calculateStats()
	}

	return c.NoContent(http.StatusOK)
}

func heatIndex(T float64, humidity int) float64 {
	RH := float64(humidity)
	feelsLike := -42.379 + 2.04901523*T + 10.14333127*RH - .22475541*T*RH - .00683783*T*T - .05481717*RH*RH + .00122874*T*T*RH + .00085282*T*RH*RH - .00000199*T*T*RH*RH
	if RH < 13 && (T >= 80 && T <= 112) {
		feelsLike = feelsLike - ((13-RH)/4)*math.Sqrt((17-math.Abs(T-95.))/17)
		if RH > 85 && (T >= 80 && T <= 87) {
			feelsLike = feelsLike + ((RH-85)/10)*((87-RH)/5)
		}
	}
	return toFixed(feelsLike, 2)
}

func windChill(temperature float64, windSpeed float64) float64 {
	if windSpeed < 3 || temperature > 50 {
		return temperature
	}

	windChill := 35.74 + 0.6215*temperature - 35.75*math.Pow(windSpeed, 0.16) + 0.4275*temperature*math.Pow(windSpeed, 0.16)
	return toFixed(windChill, 2)
}

func dewpoint(temp float64, humidity int) float64 {
	tc := (temp - 32) * 5 / 9
	L := math.Log(float64(humidity) / 100)
	M := 17.27 * tc
	N := 237.3 + tc
	B := (L + (M / N)) / 17.27
	dp := (237.3 * B) / (1 - B)
	return toFixed((dp*9/5)+32, 2)
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
			//logger.Debug(updateQuery)
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
			//logger.Debug(updateQuery)
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

func updateAlerts() {
	fmt.Println("Updating Alerts ... ")
	insertSql := `
		insert into alerts (id, alertid, wxtype, areadesc, sent, effective, onset, expires, ends, status,
		messagetype, category, severity, certainty, urgency, event, sender, senderName, headline, description,
		instruction, response)
		values (DEFAULT,$1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20,$21)
	`

	iAlerts := getAlerts()
	for _, v := range iAlerts {
		checkSql := fmt.Sprintf(`select id from alerts where alertid = '%s'`, v.IDURI)
		rows := db.QueryRow(checkSql)
		var id int
		err := rows.Scan(&id)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Error("No Alerts Found for ID :", v.IDURI)
				id = -1
			} else {
				logger.Error(fmt.Sprintf("Scan: %v", err))
			}
		}
		if id == -1 {
			_, err := db.Exec(insertSql, v.IDURI, v.Type, v.AreaDesc, v.Sent, v.Effective, v.Onset,
				v.Expires, v.Ends, v.Status,
				v.MessageType, v.Category, v.Severity, v.Certainty, v.Urgency, v.Event, v.Sender, v.SenderName,
				v.Headline, v.Description, v.Instruction, v.Response)
			logger.Info(fmt.Sprintf("Inserted Alert %s", v.Headline))
			if err != nil {
				logger.Error(err)
			}
		}
	}

}

func getAlerts() []Property {
	uri := "https://api.weather.gov/alerts/active?area=CO"
	result := make([]Property, 0)
	if !LastModified.IsZero() {
		_, err := alertRequest("HEAD", uri)
		if err != nil {
			logger.Error(err)
			return result
		}
	}

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

	for _, v := range alerts.Features {
		if strings.Contains(v.Properties.AreaDesc, "El Paso") {
			result = append(result, v.Properties)
		}
	}

	return result
}

func alertRequest(t string, url string) (body []byte, err error) {
	body = []byte("")
	client := &http.Client{}
	req, err := http.NewRequest(t, url, nil)
	if err != nil {
		logger.Error(err)
	}

	req.Header.Add("User-Agent", `Zoms Weather, wxcos@zoms.net`)

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
		l, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 GMT", resp.Header["Last-Modified"][0])
		LastModified = l
	}

	return
}

type Alerts struct {
	Type     string `json:"type"`
	Features []struct {
		ID         string   `json:"id"`
		Type       string   `json:"type"`
		Geometry   Geometry `json:"geometry"`
		Properties Property `json:"properties"`
	} `json:"features"`
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
}

type Geometry struct {
	Type        string        `json:"type"`
	Coordinates []interface{} `json:"coordinates,omitempty"`
	Geometries  []*Geometry   `json:"geometries,omitempty"`
}
type Property struct {
	IDURI    string `json:"@id"`
	Type     string `json:"@type"`
	ID       string `json:"id"`
	AreaDesc string `json:"areaDesc"`
	Geocode  struct {
		UGC  []string `json:"UGC"`
		SAME []string `json:"SAME"`
	} `json:"geocode"`
	AffectedZones []string      `json:"affectedZones"`
	References    []interface{} `json:"references"`
	Sent          time.Time     `json:"sent"`
	Effective     time.Time     `json:"effective"`
	Onset         time.Time     `json:"onset"`
	Expires       time.Time     `json:"expires"`
	Ends          time.Time     `json:"ends"`
	Status        string        `json:"status"`
	MessageType   string        `json:"messageType"`
	Category      string        `json:"category"`
	Severity      string        `json:"severity"`
	Certainty     string        `json:"certainty"`
	Urgency       string        `json:"urgency"`
	Event         string        `json:"event"`
	Sender        string        `json:"sender"`
	SenderName    string        `json:"senderName"`
	Headline      string        `json:"headline"`
	Description   string        `json:"description"`
	Instruction   string        `json:"instruction"`
	Response      string        `json:"response"`
	Parameters    struct {
		NWSheadline  []string `json:"NWSheadline"`
		PIL          []string `json:"PIL"`
		BLOCKCHANNEL []string `json:"BLOCKCHANNEL"`
	} `json:"parameters"`
}

// Record data for main database  table
type Record struct {
	ID                int       `json:"id" db:"id"`
	Mac               string    `json:"mac" db:"mac"`
	Recorded          time.Time `json:"date" db:"recorded"`
	Baromabsin        float64   `json:"baromabsin" db:"baromabsin"`
	Baromrelin        float64   `json:"baromrelin" db:"baromrelin"`
	Battout           int       `json:"battout" db:"battout"`
	Batt1             int       `json:"batt1" db:"batt1"`
	Batt2             int       `json:"batt2" db:"batt2"`
	Batt3             int       `json:"batt3" db:"batt3"`
	Batt4             int       `json:"batt4" db:"batt4"`
	Batt5             int       `json:"batt5" db:"batt5"`
	Batt6             int       `json:"batt6" db:"batt6"`
	Batt7             int       `json:"batt7" db:"batt7"`
	Batt8             int       `json:"batt8" db:"batt8"`
	Batt9             int       `json:"batt9" db:"batt9"`
	Batt10            int       `json:"batt10" db:"batt10"`
	Battlightning     int       `json:"battlightning" db:"battlightning"`
	Co2               float64   `json:"co2" db:"co2"`
	Dailyrainin       float64   `json:"dailyrainin" db:"dailyrainin"`
	Dewpoint          float64   `json:"dewpoint" db:"dewpoint"`
	Eventrainin       float64   `json:"eventrain" db:"eventrainin"`
	Feelslike         float64   `json:"feelslike" db:"feelslike"`
	Hourlyrainin      float64   `json:"hourlyrainin" db:"hourlyrainin"`
	Hourlyrain        float64   `json:"hourlyrain" db:"hourlyrain"`
	Humidity          int       `json:"humidity" db:"humidity"`
	Humidity1         int       `json:"humidity1" db:"humidity1"`
	Humidity2         int       `json:"humidity2" db:"humidity2"`
	Humidity3         int       `json:"humidity3" db:"humidity3"`
	Humidity4         int       `json:"humidity4" db:"humidity4"`
	Humidity5         int       `json:"humidity5" db:"humidity5"`
	Humidity6         int       `json:"humidity6" db:"humidity6"`
	Humidity7         int       `json:"humidity7" db:"humidity7"`
	Humidity8         int       `json:"humidity8" db:"humidity8"`
	Humidity9         int       `json:"humidity9" db:"humidity9"`
	Humidity10        int       `json:"humidity10" db:"humidity10"`
	Humidityin        int       `json:"humidityin" db:"humidityin"`
	Lastrain          time.Time `json:"lastrain" db:"lastrain"`
	Maxdailygust      float64   `json:"maxdailygust" db:"maxdailygust"`
	Relay1            int       `json:"relay1" db:"relay1"`
	Relay2            int       `json:"relay2" db:"relay2"`
	Relay3            int       `json:"relay3" db:"relay3"`
	Relay4            int       `json:"relay4" db:"relay4"`
	Relay5            int       `json:"relay5" db:"relay5"`
	Relay6            int       `json:"relay6" db:"relay6"`
	Relay7            int       `json:"relay7" db:"relay7"`
	Relay8            int       `json:"relay8" db:"relay8"`
	Relay9            int       `json:"relay9" db:"relay9"`
	Relay10           int       `json:"relay10" db:"relay10"`
	Monthlyrainin     float64   `json:"monthlyrainin" db:"monthlyrainin"`
	Solarradiation    float64   `json:"solarradiation" db:"solarradiation"`
	Tempf             float64   `json:"tempf" db:"tempf"`
	Temp1f            float64   `json:"temp1f" db:"temp1f"`
	Temp2f            float64   `json:"temp2f" db:"temp2f"`
	Temp3f            float64   `json:"temp3f" db:"temp3f"`
	Temp4f            float64   `json:"temp4f" db:"temp4f"`
	Temp5f            float64   `json:"temp5f" db:"temp5f"`
	Temp6f            float64   `json:"temp6f" db:"temp6f"`
	Temp7f            float64   `json:"temp7f" db:"temp7f"`
	Temp8f            float64   `json:"temp8f" db:"temp8f"`
	Temp9f            float64   `json:"temp9f" db:"temp9f"`
	Temp10f           float64   `json:"temp10f" db:"temp10f"`
	Tempinf           float64   `json:"tempinf" db:"tempinf"`
	Totalrainin       float64   `json:"totalrainin" db:"totalrainin"`
	Uv                float64   `json:"uv" db:"uv"`
	Weeklyrainin      float64   `json:"weeklyrainin" db:"weeklyrainin"`
	Winddir           int       `json:"winddir" db:"winddir"`
	Windgustmph       float64   `json:"windgustmph" db:"windgustmph"`
	Windgustdir       int       `json:"windgustdir" db:"windgustdir"`
	Windspeedmph      float64   `json:"windspeedmph" db:"windspeedmph"`
	Yearlyrainin      float64   `json:"yearlyrainin" db:"yearlyrainin"`
	Lightningday      int       `json:"lightningday" db:"lightiningday"`
	Lightninghour     int       `json:"lightninghour" db:"lightininghour"`
	Lightningdistance float64   `json:"lightningdistance" db:"lightningdistance"`
	Lightningtime     time.Time `json:"lightningtime" db:"lightningtime"`
	Aqipm25           int       `json:"aqipm25" db:"aqi"`
	Aqipm2524h        int       `json:"aqipm2524h" db:"aqi24"`
}
