package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/IvanMenshykov/MoonPhase"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
var db *gorm.DB
var logger = logrus.New()
var loc *time.Location

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
	logger.SetLevel(logLevel)

	dburi := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_DATABASE"))
	db, err = gorm.Open("mysql", dburi)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if logLevel == logrus.DebugLevel {
		db.LogMode(true)
	}

	loc, err = time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Println(err)
		return
	}

	hub := newHub()
	go hub.run()
	go broadcast(hub)

	r := mux.NewRouter()
	r.HandleFunc("/api/ws", func(w http.ResponseWriter, r *http.Request) {
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
	tomorrowURL := url+"&date="+d
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
	dateformat := "DATE_FORMAT(date,'%Y-%m-%d')"
	if p == "day" || p == "yesterday" {
		dateformat = "DATE_FORMAT(date,'%H:%d:%s')"
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
	where := fmt.Sprintf("%s AS mmdd, max(%s) AS max, min(%s) AS min", dateformat, t, t)
	db.Table("records").
		Select(where).
		Where("date BETWEEN ? AND ?", dates[0], dates[1]).
		Group("mmdd").
		Order("mmdd").
		Scan(&rt)

	res := map[string][]Chart{}
	res["data1"] = make([]Chart, 0)
	res["data2"] = make([]Chart, 0)

	for _, data := range rt {
		res["data1"] = append(res["data1"], Chart{Date: data.Mmdd, Value: data.Max})
		res["data2"] = append(res["data2"], Chart{Date: data.Mmdd, Value: data.Min})
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
	sel := fmt.Sprintf("%s(%s) as value,`date`", c, t)
	group := fmt.Sprintf("`date`,%s", t)
	orderby := fmt.Sprintf("%s%s ", t, order)
	db.Table("records").
		Select(sel).
		Group(group).
		Order(orderby).
		Limit(1).
		Scan(&rt)

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
	sel := fmt.Sprintf("MAX(%s)", t)

	start := time.Now()
	end := start.Add(-30 * time.Minute)
	var avg float64
	current := Record{}
	row := db.Table("records").Select(sel).Where("date BETWEEN ? AND ?", end, start).Row()
	row.Scan(&avg)
	if err := db.Select(t).Order("date desc").First(&current).Error; err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	trend := Trend{}
	if strings.Contains(t, "temp") {
		if current.Tempf > avg {
			//trend up
			trend.Trend = "up"
			trend.By = toFixed(current.Tempf - avg)
		} else {
			//trend down
			trend.Trend = "down"
			trend.By = toFixed(avg - current.Tempf)
		}
	} else {
		if current.Baromrelin > avg {
			//trend up
			trend.Trend = "Steady"
			if (current.Baromrelin - avg) > .5 {
				trend.Trend = "Rising"
			}
		} else {
			//trend down
			trend.Trend = "Steady"
			if (avg - current.Baromrelin) > .5 {
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
	avg := Result{}
	avgdir := Result{}
	db.Table("records").Order("windspeedmph desc").Select("MAX(windspeedmph) as value, `date`").Where("date BETWEEN ? AND ?", end, start).Group("`date`,windspeedmph").Limit(1).Scan(&max)
	db.Table("records").Order("windgustmph desc").Select("MAX(windgustmph) as value, `date`").Where("date BETWEEN ? AND ?", end, start).Group("`date`,windgustmph").Limit(1).Scan(&maxgust)
	db.Table("records").Select("AVG(windspeedmph) as value").Where("date BETWEEN ? AND ?", end, start).Scan(&avg)
	db.Table("records").Select("AVG(winddir) as value").Where("date BETWEEN ? AND ?", end, start).Scan(&avgdir)

	res := map[string]Result{}

	res["dir"] = avgdir
	res["wind"] = max
	res["gust"] = maxgust
	res["avg"] = avg

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
	s := make([]Stat, 0)
	db.Table("stat").Find(&s)

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
	req.Header.Add("User-Agent", `Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`)
	if len(header) > 0 {
		for key, value := range header {
			req.Header.Add(key, value)
		}
	}
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
	rec := Record{}
	db.Select("`id`,`date`,`baromabsin`,`baromrelin`,`dailyrainin`,`dewpoint`,`eventrainin`,`feelslike`,`hourlyrainin`,`humidity`,`humidity1`,`humidity2`,`humidityin`,`lastrain`,`maxdailygust`,`monthlyrainin`,`solarradiation`,`tempf`,`temp1f`,`temp2f`,`tempinf`,`totalrainin`,`uv`,`weeklyrainin`,`winddir`,`windgustmph`,`windspeedmph`,`yearlyrainin`").Order("date desc").First(&rec)
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
	row := db.Table("records").Select("MAX(dailyrainin)").Where("date BETWEEN ? AND ?", end, start).Row()
	row.Scan(&max)
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
