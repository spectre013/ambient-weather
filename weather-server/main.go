package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/IvanMenshykov/MoonPhase"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options
var db *sql.DB
var logger = logrus.New()
var loc *time.Location
var client *twitter.Client

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

	loc, err = time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Println(err)
		return
	}

	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	client, err = getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Println(err)
	}
	// Dont tweet if we are dev
	if os.Getenv("LOGLEVEL") != "Debug" {
		s := gocron.NewScheduler(loc)
		//s.Every(1).Hour().StartImmediately().Do(sendUpdate)
		if err != nil {
			logger.Error(err)
		}
		s.StartAsync()
	}
	// Set up Alerts
	//go startAlerts()

	// Setup Web Sockets
	hub := newHub()
	go hub.run()
	go broadcast(hub)

	r := mux.NewRouter()
	r.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Websocket connection")
		serveWs(hub, w, r)
	})
	//API
	r.HandleFunc("/api/current", loggingMiddleware(current))
	r.HandleFunc("/api/minmax/{field}", loggingMiddleware(minmax))
	r.HandleFunc("/api/trend/{type}", loggingMiddleware(trend))
	r.HandleFunc("/api/wind", loggingMiddleware(wind))
	r.HandleFunc("/api/forecast", loggingMiddleware(forecast))
	r.HandleFunc("/api/luna", loggingMiddleware(astro))
	r.HandleFunc("/api/chart/{type}/{period}", loggingMiddleware(chart))
	r.HandleFunc("/api/alltime/{calc}/{type}", loggingMiddleware(alltime))
	r.HandleFunc("/api/metar", loggingMiddleware(metar))
	r.HandleFunc("/api/apiin", loggingMiddleware(apiin))
	//Index
	r.HandleFunc("/", loggingMiddleware(home))

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + os.Getenv("PORT"),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting server on " + os.Getenv("PORT"))
	log.Fatal(srv.ListenAndServe())
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		logger.Infof(
			"%s\t%s\t%s",
			time.Now().Format("2006-01-02 15:04:05"),
			r.Method,
			r.RequestURI,
		)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func broadcast(hub *Hub) {
	for range time.NewTicker(30 * time.Second).C {
		m, err := getCurrent()
		if err != nil {
			log.Println(err)
			continue
		}
		hub.broadcast <- m
	}
}

func apiin(w http.ResponseWriter, r *http.Request) {
	rec := Record{}
	logger.Info("Incomming Ambient Data Processing ....")
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	inserted := insertRecord(rec)
	if inserted {
		go calculateStats(rec)
	}

	w.Write([]byte("Success"))
}

func forecast(w http.ResponseWriter, r *http.Request) {
	header := map[string]string{}
	url := fmt.Sprintf("https://api.darksky.net/forecast/%s/%s,%s", os.Getenv("DARKSKY"), os.Getenv("LAT"), os.Getenv("LON"))
	res, err := makeRequest(url, header)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(res)
}

func metar(w http.ResponseWriter, r *http.Request) {
	header := map[string]string{}
	header["X-API-Key"] = os.Getenv("METAR")
	url := "https://api.checkwx.com/metar/kcos/decoded"
	res, err := makeRequest(url, header)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(res)
}

func astro(w http.ResponseWriter, r *http.Request) {
	t := time.Now().UTC()
	m := MoonPhase.New(t)
	header := map[string]string{}
	url := fmt.Sprintf("https://api.ipgeolocation.io/astronomy?apiKey=%s&lat=%s&long=%s", os.Getenv("IPGEO"), os.Getenv("LAT"), os.Getenv("LON"))
	res, err := makeRequest(url, header)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	solar := Astro{}
	err = json.Unmarshal(res, &solar)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	d := time.Now().Add(48 * time.Hour).Format("2006-01-02")
	tomorrowURL := url + "&date=" + d
	tomor, err := makeRequest(tomorrowURL, header)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	tom := Tomorrow{}
	err = json.Unmarshal(tomor, &tom)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	solar.Tomorrow = tom
	solar.Newmoon = time.Unix(int64(m.NewMoon()), 0).UTC()
	solar.Nextnewmoon = time.Unix(int64(m.NextNewMoon()), 0).UTC()
	solar.Fullmoon = time.Unix(int64(m.FullMoon()), 0)
	solar.Phase = m.PhaseName()
	solar.Illuminated = math.Round(m.Illumination() * float64(100))
	solar.Age = math.Round(m.Age())

	b, err := json.Marshal(solar)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
}

func chart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := vars["type"]
	p := vars["period"]
	t = cleanString(t)
	dates := getTimeframe(p)
	dformat := "TO_CHAR(recorded,'YYY-MM-DD')"
	if p == "day" || p == "yesterday" {
		dformat = "TO_CHAR(recorded,'YYY-MM-DD HH24:MI:SS')"
	}

	type Chart struct {
		Date  string  `json:"label"`
		Value float64 `json:"y"`
	}
	type Result struct {
		Mmdd string  `json:"date"`
		Max  float64 `json:"max"`
		Min  float64 `json:"min"`
	}
	rt := make([]Result, 0)
	where := ""
	if t == "dailyrainin" {
		where = fmt.Sprintf("%s AS mmdd, max(%s) AS max", dformat, t)
	} else {
		where = fmt.Sprintf("%s AS mmdd, max(%s) AS max, min(%s) AS min", dformat, t, t)
	}

	query := fmt.Sprintf("select %s from records where recorded BETWEEN '%s' AND '%s' group by mmdd order by mmdd", where, formatDate(dates[0]), formatDate(dates[1]))
	logger.Debug(query)
	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		r := Result{}
		if t == "dailyrainin" {
			err = rows.Scan(&r.Mmdd, &r.Max)
		} else {
			err = rows.Scan(&r.Mmdd, &r.Max, &r.Min)
		}
		if err != nil {
			logger.Error("Scan: %v", err)
		}
		rt = append(rt, r)
	}
	err = rows.Err()
	if err != nil {
		logger.Error(err)
	}
	res := map[string][]Chart{}
	res["data1"] = make([]Chart, 0)
	if t != "dailyrainin" {
		res["data2"] = make([]Chart, 0)
	}

	for _, data := range rt {
		res["data1"] = append(res["data1"], Chart{Date: data.Mmdd, Value: data.Max})
		if t != "dailyrainin" {
			res["data2"] = append(res["data2"], Chart{Date: data.Mmdd, Value: data.Min})
		}
	}

	b, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)

}

func alltime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := vars["type"]
	c := vars["calc"]
	t = cleanString(t)
	c = cleanString(c)
	order := ""
	if c == "max" {
		order = " desc"
	}
	type Result struct {
		Date  string  `json:"date"`
		Value float64 `json:"value"`
	}
	rt := Result{}
	sel := fmt.Sprintf("%s as value,recorded", t)
	orderby := fmt.Sprintf("%s%s ", t, order)

	query := fmt.Sprintf("select %s from records order by %s limit 1", sel, orderby)
	logger.Info(query)
	rows := db.QueryRow(query)
	err := rows.Scan(&rt.Value, &rt.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", query)
		} else {
			logger.Error("Scan: %v", err)
		}
	}

	b, err := json.Marshal(rt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
}

func getTimeframe(timeframe string) []time.Time {
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

func trend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := vars["type"]
	t = cleanString(t)
	sel := fmt.Sprintf("AVG(%s)", t)

	start := time.Now()
	end := start.Add(-30 * time.Minute)
	var avg float64
	cr := Record{}

	avgQuery := fmt.Sprintf("select %s from records where recorded BETWEEN '%s' AND '%s'", sel, formatDate(end), formatDate(start))
	logger.Debug(avgQuery)
	rows := db.QueryRow(avgQuery)
	err := rows.Scan(&avg)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", avgQuery)
		} else {
			logger.Error("Scan: %v", err)
		}
	}
	currentQuery := "select id,mac,recorded,baromabsin,baromrelin,battout,batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,hourlyrainin,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,hourlyrain,lightningday,lightninghour,lightningtime,lightningdistance,battlightning from records order by recorded desc limit 1"
	logger.Debug(currentQuery)
	crows := db.QueryRow(currentQuery)
	err = crows.Scan(&cr.ID, &cr.Mac, &cr.Date, &cr.Baromabsin, &cr.Baromrelin, &cr.Battout, &cr.Batt1, &cr.Batt2, &cr.Batt3, &cr.Batt4, &cr.Batt5, &cr.Batt6, &cr.Batt7, &cr.Batt8, &cr.Batt9, &cr.Batt10, &cr.Co2, &cr.Dailyrainin, &cr.Dewpoint, &cr.Eventrainin, &cr.Feelslike, &cr.Hourlyrainin, &cr.Humidity, &cr.Humidity1, &cr.Humidity2, &cr.Humidity3, &cr.Humidity4, &cr.Humidity5, &cr.Humidity6, &cr.Humidity7, &cr.Humidity8, &cr.Humidity9, &cr.Humidity10, &cr.Humidityin, &cr.Lastrain, &cr.Maxdailygust, &cr.Relay1, &cr.Relay2, &cr.Relay3, &cr.Relay4, &cr.Relay5, &cr.Relay6, &cr.Relay7, &cr.Relay8, &cr.Relay9, &cr.Relay10, &cr.Monthlyrainin, &cr.Solarradiation, &cr.Tempf, &cr.Temp1f, &cr.Temp2f, &cr.Temp3f, &cr.Temp4f, &cr.Temp5f, &cr.Temp6f, &cr.Temp7f, &cr.Temp8f, &cr.Temp9f, &cr.Temp10f, &cr.Tempinf, &cr.Totalrainin, &cr.Uv, &cr.Weeklyrainin, &cr.Winddir, &cr.Windgustmph, &cr.Windgustdir, &cr.Windspeedmph, &cr.Yearlyrainin, &cr.Hourlyrain, &cr.Lightningday, &cr.Lightninghour, &cr.Lightningtime, &cr.Lightningdistance, &cr.Battlightning)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", currentQuery)
		} else {
			logger.Error("Scan: %v", err)
		}
	}

	trend := Trend{}
	if strings.Contains(t, "temp") {
		if cr.Tempf > avg {
			//trend up
			trend.Trend = "up"
			trend.By = toFixed(cr.Tempf - avg)
		} else {
			//trend down
			trend.Trend = "down"
			trend.By = toFixed(avg - cr.Tempf)
		}
	} else {
		if cr.Baromrelin > avg {
			//trend up
			trend.Trend = "Steady"
			if (cr.Baromrelin - avg) > .5 {
				trend.Trend = "Rising"
			}
		} else {
			//trend down
			trend.Trend = "Steady"
			if (avg - cr.Baromrelin) > .5 {
				trend.Trend = "Falling"
			}
		}
	}

	b, err := json.Marshal(trend)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)

}
func wind(w http.ResponseWriter, r *http.Request) {
	type Result struct {
		Date  time.Time `json:"date"`
		Value float64   `json:"value"`
	}
	start := time.Now()
	end := start.Add(-30 * time.Minute)

	max := Result{}
	maxgust := Result{}
	var avg float64
	var avgdir float64

	maxSpeed := fmt.Sprintf("select windspeedmph as value, recorded from records where recorded BETWEEN '%s' AND '%s' order by windspeedmph desc limit 1", formatDate(end), formatDate(start))
	maxGust := fmt.Sprintf("select windgustmph as value, recorded from records where recorded BETWEEN '%s' AND '%s' order by windgustmph desc limit 1", formatDate(end), formatDate(start))
	avgSpeed := fmt.Sprintf("select AVG(windspeedmph) as value from records where recorded BETWEEN '%s' AND '%s'", formatDate(end), formatDate(start))
	avgDir := fmt.Sprintf("select AVG(winddir) as value from records where recorded BETWEEN '%s' AND '%s'", formatDate(end), formatDate(start))

	logger.Debug(maxSpeed)
	mrows := db.QueryRow(maxSpeed)
	err := mrows.Scan(&max.Value, &max.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows ", maxSpeed)
		} else {
			logger.Error("Scan: %v", err)
		}
	}

	logger.Debug(maxGust)
	mgrows := db.QueryRow(maxGust)
	err = mgrows.Scan(&maxgust.Value, &maxgust.Date)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", maxGust)
		} else {
			logger.Error("Scan: %v", err)
		}
	}

	logger.Debug(avgSpeed)
	asrows := db.QueryRow(avgSpeed)
	err = asrows.Scan(&avg)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", avgSpeed)
		} else {
			logger.Error("Scan: %v", err)
		}
	}

	logger.Debug(avgDir)
	crows := db.QueryRow(avgDir)
	err = crows.Scan(&avgdir)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", avgdir)
		} else {
			logger.Error("Scan: %v", err)
		}
	}

	res := map[string]Result{}
	res["dir"] = Result{Value: avgdir}
	res["wind"] = max
	res["gust"] = maxgust
	res["avg"] = Result{Value: avg}

	b, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)
}
func minmax(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	f := vars["field"]
	s := getStats()
	minmax := make(map[string]map[string]StatValue)

	for _, v := range s {
		if strings.Contains(v.ID, f) {
			parts := strings.Split(v.ID, "_")

			if _, ok := minmax[parts[1]]; !ok {
				minmax[parts[1]] = make(map[string]StatValue)
			}
			minmax[parts[1]][strings.ToLower(parts[0])] = StatValue{Date: v.Date, Value: v.Value}
		}
	}
	b, err := json.Marshal(minmax)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(b)

}
func current(w http.ResponseWriter, r *http.Request) {
	res, err := getCurrent()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write(res)
}

func makeRequest(url string, header map[string]string) ([]byte, error) {
	logger.Debug(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if _, ok := header["User-Agent"]; !ok {
		req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
	}
	if len(header) > 0 {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}
	logger.Debug(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return body, nil
}

func getCurrent() ([]byte, error) {
	query := `select id,mac,recorded,baromabsin,baromrelin,battout,Batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,
				maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,battlightning,lightningday,lightninghour,lightningtime,lightningdistance 
				from records order by recorded desc limit 1`
	rec := getRecord(query)

	hourlyRain := getHourlyRain()
	rec.Hourlyrain = hourlyRain
	b, err := json.Marshal(rec)
	if err != nil {
		log.Println(err)
		return []byte(""), err
	}
	return b, nil
}

func getHourlyRain() float64 {
	start := time.Now()
	end := start.Add(-60 * time.Minute)
	var max float64
	query := fmt.Sprintf("select dailyrainin from records where recorded BETWEEN '%s' AND '%s' order by dailyrainin desc limit 1", formatDate(end), formatDate(start))
	logger.Debug(query)
	crows := db.QueryRow(query)
	err := crows.Scan(&max)
	if err != nil {
		if err == sql.ErrNoRows {
			//logger.Error("Zero Rows Found", query)
		} else {
			logger.Error("Scan: %v", err)
		}
	}
	return max
}

func home(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{}
	res["message"] = "Zoms.net weather API visit https://weather.zoms.net for more information"
	b, _ := json.Marshal(r)
	w.Write(b)
}

func toFixed(x float64) float64 {
	return math.Round(x*100) / 100
}

func cleanString(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9 -]")
	replaceStr := reg.ReplaceAllString(s, "")
	return replaceStr
}

func calculateStats(r Record) {
	type Query struct {
		Query  string
		Params []time.Time
	}

	types := []string{"MAX", "MIN", "AVG"}
	periods := []string{"day", "yesterday", "month", "year"}
	fields := []string{"tempf", "tempinf", "temp1f", "temp2f", "temp3f", "baromrelin", "uv", "humidity", "windspeedmph", "windgustmph", "dewpoint", "humidityin", "humidity1", "humidity2", "humidity3", "dailyrainin", "lightning"}
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
						GROUP BY A.ldate order by A.ldate desc limit 1
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

func checkStat(id string) bool {
	rows := db.QueryRow(fmt.Sprintf("select id from stats where id = '%s'", id))
	var r string
	err := rows.Scan(&r)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			logger.Error("Scan: %v", err)
		}
	}
	return true
}

func formatDate(date time.Time) string {
	format := "2006-01-02 15:04:05"
	return date.Format(format)
}
func insertRecord(r Record) bool {
	query := fmt.Sprintf(`insert into records (id,mac,recorded,baromabsin,baromrelin,battout,batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,hourlyrainin,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,hourlyrain,lightningday,lightninghour,lightningtime,lightningdistance,battlightning) values
			(DEFAULT,'%s','%s',%f,%f,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%f,%f,%f,%f,%f,%f,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,'%s',%f,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%d,%f,%d,%f,%f,%f,%d,%d,'%s',%f,%d)				
	`, r.Mac, formatDate(r.Date), r.Baromabsin, r.Baromrelin, r.Battout, r.Batt1, r.Batt2, r.Batt3, r.Batt4, r.Batt5, r.Batt6, r.Batt7, r.Batt8, r.Batt9, r.Batt10, r.Co2, r.Dailyrainin, r.Dewpoint, r.Eventrainin, r.Feelslike, r.Hourlyrainin, r.Humidity, r.Humidity1, r.Humidity2, r.Humidity3, r.Humidity4, r.Humidity5, r.Humidity6, r.Humidity7, r.Humidity8, r.Humidity9, r.Humidity10, r.Humidityin, formatDate(r.Lastrain), r.Maxdailygust, r.Relay1, r.Relay2, r.Relay3, r.Relay4, r.Relay5, r.Relay6, r.Relay7, r.Relay8, r.Relay9, r.Relay10, r.Monthlyrainin, r.Solarradiation, r.Tempf, r.Temp1f, r.Temp2f, r.Temp3f, r.Temp4f, r.Temp5f, r.Temp6f, r.Temp7f, r.Temp8f, r.Temp9f, r.Temp10f, r.Tempinf, r.Totalrainin, r.Uv, r.Weeklyrainin, r.Winddir, r.Windgustmph, r.Windgustdir, r.Windspeedmph, r.Yearlyrainin, r.Hourlyrain, r.Lightningday, r.Lightninghour, formatDate(r.Lightningtime), r.Lightningdistance, r.Battlightning)
	logger.Debug(query)
	_, err := db.Exec(query)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func getRecord(sqlStatement string) Record {

	rows := db.QueryRow(sqlStatement)
	r := Record{}
	err := rows.Scan(&r.ID, &r.Mac, &r.Date, &r.Baromabsin, &r.Baromrelin, &r.Battout, &r.Batt1, &r.Batt2, &r.Batt3, &r.Batt4, &r.Batt5, &r.Batt6, &r.Batt7, &r.Batt8, &r.Batt9, &r.Batt10, &r.Co2, &r.Dailyrainin, &r.Dewpoint, &r.Eventrainin, &r.Feelslike, &r.Hourlyrainin, &r.Hourlyrain, &r.Humidity, &r.Humidity1, &r.Humidity2, &r.Humidity3, &r.Humidity4, &r.Humidity5, &r.Humidity6, &r.Humidity7, &r.Humidity8, &r.Humidity9, &r.Humidity10, &r.Humidityin, &r.Lastrain, &r.Maxdailygust, &r.Relay1, &r.Relay2, &r.Relay3, &r.Relay4, &r.Relay5, &r.Relay6, &r.Relay7, &r.Relay8, &r.Relay9, &r.Relay10, &r.Monthlyrainin, &r.Solarradiation, &r.Tempf, &r.Temp1f, &r.Temp2f, &r.Temp3f, &r.Temp4f, &r.Temp5f, &r.Temp6f, &r.Temp7f, &r.Temp8f, &r.Temp9f, &r.Temp10f, &r.Tempinf, &r.Totalrainin, &r.Uv, &r.Weeklyrainin, &r.Winddir, &r.Windgustmph, &r.Windgustdir, &r.Windspeedmph, &r.Yearlyrainin, &r.Battlightning, &r.Lightningday, &r.Lightninghour, &r.Lightningtime, &r.Lightningdistance)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found ", sqlStatement)
		} else {
			logger.Error("Scan: %v", err)
		}
	}
	fmt.Println(r.Date)

	return r
}
func getStats() []Stat {
	stats := make([]Stat, 0)
	rows, err := db.Query("Select * from stats")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		r := Stat{}
		err = rows.Scan(&r.ID, &r.Date, &r.Value)
		if err != nil {
			logger.Error("Scan: %v", err)
		}
		stats = append(stats, r)
	}
	err = rows.Err()
	if err != nil {
		logger.Error(err)
	}
	return stats
}
func getRecords(sqlStatement string) []Record {
	rec := make([]Record, 0)
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		r := Record{}
		err = rows.Scan(&r.ID, &r.Mac, &r.Date, &r.Baromabsin, &r.Baromrelin, &r.Battout, &r.Batt1, &r.Batt2, &r.Batt3, &r.Batt4, &r.Batt5, &r.Batt6, &r.Batt7, &r.Batt8, &r.Batt9, &r.Batt10, &r.Co2, &r.Dailyrainin, &r.Dewpoint, &r.Eventrainin, &r.Feelslike, &r.Hourlyrainin, &r.Hourlyrain, &r.Humidity, &r.Humidity1, &r.Humidity2, &r.Humidity3, &r.Humidity4, &r.Humidity5, &r.Humidity6, &r.Humidity7, &r.Humidity8, &r.Humidity9, &r.Humidity10, &r.Humidityin, &r.Lastrain, &r.Maxdailygust, &r.Relay1, &r.Relay2, &r.Relay3, &r.Relay4, &r.Relay5, &r.Relay6, &r.Relay7, &r.Relay8, &r.Relay9, &r.Relay10, &r.Monthlyrainin, &r.Solarradiation, &r.Tempf, &r.Temp1f, &r.Temp2f, &r.Temp3f, &r.Temp4f, &r.Temp5f, &r.Temp6f, &r.Temp7f, &r.Temp8f, &r.Temp9f, &r.Temp10f, &r.Tempinf, &r.Totalrainin, &r.Uv, &r.Weeklyrainin, &r.Winddir, &r.Windgustmph, &r.Windgustdir, &r.Windspeedmph, &r.Yearlyrainin, &r.Battlightning, &r.Lightningday, &r.Lightninghour, &r.Lightningtime, &r.Lightningdistance)
		if err != nil {
			logger.Error("Scan: %v", err)
		}
		rec = append(rec, r)
	}
	err = rows.Err()
	if err != nil {
		logger.Error(err)
	}
	return rec
}
