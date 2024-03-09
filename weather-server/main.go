package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/IvanMenshykov/MoonPhase"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"

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
var logger = logrus.New()
var loc *time.Location
var units = "imperial"

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
	db, err := sql.Open("postgres", dburi)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	weather := buildWeather(db)

	loc, err = time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Setup Web Sockets
	hub := newHub()
	hub.weather = weather
	go hub.run()
	go broadcast(hub)

	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Websocket connection")
		serveWs(hub, w, r)
	})
	//API
	r.HandleFunc("/current", loggingMiddleware(weather.current))
	r.HandleFunc("/temperature", loggingMiddleware(weather.temperature))
	r.HandleFunc("/forecast", loggingMiddleware(weather.forecast))
	r.HandleFunc("/temp/{time}", loggingMiddleware(weather.temp))
	r.HandleFunc("/alert/{id}", loggingMiddleware(weather.Alert))
	r.HandleFunc("/alertview/{id}", loggingMiddleware(weather.Alertview))
	r.HandleFunc("/wind", loggingMiddleware(weather.Wind))
	//Index
	r.HandleFunc("/", loggingMiddleware(weather.index))
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

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

func buildWeather(db *sql.DB) Weather {
	forecast, err := getForecast()
	if err != nil {
		logger.Error(err)
	}
	astro := astro()

	return Weather{
		DB:       db,
		Forecast: forecast,
		Astro:    astro,
		Updated:  time.Now(),
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return body, nil
}

func (w Weather) Minmax() map[string]map[string]map[string]StatValue {
	s := getStats(w.DB)

	minmax := make(map[string]map[string]map[string]StatValue)

	for _, v := range s {
		parts := strings.Split(v.ID, "_")
		if _, ok := minmax[parts[2]]; !ok {
			minmax[parts[2]] = make(map[string]map[string]StatValue)
		}
		if _, ok := minmax[parts[2]][parts[1]]; !ok {
			minmax[parts[2]][parts[1]] = make(map[string]StatValue)
		}
		minmax[parts[2]][parts[1]][strings.ToLower(parts[0])] = StatValue{Recorded: v.Recorded, Value: v.Value}
	}

	return minmax
}

func getForecast() (*ForecastImage, error) {
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Colorado%%20Springs?unitGroup=us&iconSets=icon2&include=days&key=%s&contentType=json", os.Getenv("WEATHER_API"))
	header := map[string]string{}
	res, err := makeRequest(url, header)
	if err != nil {
		logger.Error(err)
		return &ForecastImage{}, err
	}
	f := ForecastImage{}
	err = json.Unmarshal(res, &f)
	if err != nil {
		logger.Error(err)
		return &f, err
	}
	days := f.Days[1:8]
	f.Days = days
	return &f, err
}

func astro() *Astro {
	t := time.Now().UTC()
	m := MoonPhase.New(t)
	header := map[string]string{}
	url := fmt.Sprintf("https://api.ipgeolocation.io/astronomy?apiKey=%s&lat=%s&long=%s", os.Getenv("IPGEO"), os.Getenv("LAT"), os.Getenv("LON"))
	res, err := makeRequest(url, header)
	if err != nil {
		logger.Println(err)
	}
	solar := Astro{}
	err = json.Unmarshal(res, &solar)
	if err != nil {
		logger.Println(err)
	}
	d := time.Now().Local().Add(24 * time.Hour).Format("2006-01-02")
	tomorrowURL := url + "&date=" + d
	tomor, err := makeRequest(tomorrowURL, header)
	if err != nil {
		logger.Println(err)
	}
	tom := Tomorrow{}
	err = json.Unmarshal(tomor, &tom)
	if err != nil {
		logger.Println(err)
	}
	tom.Date = d

	solar.Tomorrow = tom
	solar.Newmoon = time.Unix(int64(m.NewMoon()), 0).UTC()
	solar.Nextnewmoon = time.Unix(int64(m.NextNewMoon()), 0).UTC()
	solar.Fullmoon = time.Unix(int64(m.FullMoon()), 0)
	solar.Phase = m.PhaseName()
	solar.Illuminated = math.Round(m.Illumination() * float64(100))
	solar.Age = math.Round(m.Age())

	return &solar
}

func (w Weather) getCurrent() (map[string]BoxProps, TemplateData, error) {
	query := `select id,mac,recorded,baromabsin,baromrelin,battout,Batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,
				maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,battlightning,lightningday,lightninghour,lightningtime,lightningdistance,  aqipm25, aqipm2524h
				from records order by recorded desc limit 1`
	rec := getRecord(w.DB, query)

	hourlyRain := getHourlyRain(w.DB)
	rec.Hourlyrain = hourlyRain

	minmax := w.Minmax()
	w.checkForecast()
	data := TemplateData{
		Units:    units,
		Record:   rec,
		Minmax:   minmax,
		Alerts:   w.Alerts(),
		Forecast: *w.Forecast,
		Wind:     w.getWind(),
		Astro:    *w.Astro,
		tTrend:   w.trend("tempf"),
		bTrend:   w.trend("baromabsin"),
	}
	box := buildBoxProps(units)

	return box, data, nil
}

func buildBoxProps(units string) map[string]BoxProps {
	box := map[string]BoxProps{}

	box["temperature"] = BoxProps{
		Icon:  "fa-temperature-three-quarters",
		Title: "Temperature",
		Unit:  tempLabel(units),
		Style: map[string]string{},
	}
	box["forecast"] = BoxProps{
		Icon:  "fa-cloud-sun-rain",
		Title: "Forecast",
		Unit:  tempLabel(units),
		Style: map[string]string{"forecast": "width: 570px"},
	}
	box["alerts"] = BoxProps{
		Icon:  "fa-triangle-exclamation",
		Title: "Alerts",
		Unit:  "",
		Style: map[string]string{},
	}
	box["wind"] = BoxProps{
		Icon:  "fa-wind",
		Title: "Wind",
		Unit:  windLabel(units),
		Style: map[string]string{},
	}
	box["rain"] = BoxProps{
		Icon:  "fa-cloud-showers-heavy",
		Title: "Rain",
		Unit:  rainLabel(units),
		Style: map[string]string{},
	}
	box["lightning"] = BoxProps{
		Icon:  "fa-bolt-lightning",
		Title: "Lightning",
		Unit:  "",
		Style: map[string]string{},
	}
	box["humidity"] = BoxProps{
		Icon:  "fa-droplet",
		Title: "Humidity",
		Unit:  "%",
		Style: map[string]string{},
	}
	box["barometer"] = BoxProps{
		Icon:  "fa-temperature-high",
		Title: "Barometer",
		Unit:  baroLabel(units),
		Style: map[string]string{},
	}
	box["sun"] = BoxProps{
		Icon:  "fa-sun",
		Title: "Sun",
		Unit:  "",
		Style: map[string]string{},
	}
	box["uv"] = BoxProps{
		Icon:  "fa-cloud-sun",
		Title: "UV | Solar",
		Unit:  "",
		Style: map[string]string{},
	}
	box["aqi"] = BoxProps{
		Icon:  "fa-lungs",
		Title: "Air Quality Index",
		Unit:  "",
		Style: map[string]string{},
	}
	box["tempin"] = BoxProps{
		Icon:  "fa-temperature-half",
		Title: "Living",
		Unit:  "",
		Style: map[string]string{},
	}
	box["temp1"] = BoxProps{
		Icon:  "fa-temperature-half",
		Title: "Basement",
		Unit:  "",
		Style: map[string]string{},
	}
	box["temp2"] = BoxProps{
		Icon:  "fa-temperature-half",
		Title: "Master Bedroom",
		Unit:  "",
		Style: map[string]string{},
	}
	box["temp3"] = BoxProps{
		Icon:  "fa-temperature-half",
		Title: "Office",
		Unit:  "",
		Style: map[string]string{},
	}
	return box
}
func (w *Weather) checkForecast() {
	dur := time.Since(w.Updated)
	if dur.Minutes() > 5 {
		f, err := getForecast()
		if err != nil {
			logger.Error(err)
		}
		w.Forecast = f
		w.Updated = time.Now()
		w.Astro = astro()
	}
}

func getHourlyRain(db *sql.DB) float64 {
	start := time.Now()
	end := start.Add(-60 * time.Minute)
	var maxrain float64
	query := fmt.Sprintf("select dailyrainin from records where recorded BETWEEN '%s' AND '%s' order by dailyrainin desc limit 1", formatDate(end), formatDate(start))
	logger.Debug(query)
	crows := db.QueryRow(query)
	err := crows.Scan(&maxrain)
	sqlError(err, query)
	return maxrain
}
