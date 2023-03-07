package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
	_ "github.com/lib/pq"
)

var db *sql.DB
var logger = logrus.New()
var loc *time.Location
var client *http.Client
var tc *gocron.Job
var tf *gocron.Job
var aj *gocron.Job
var fj *gocron.Job
var LastModified time.Time
var days = []string{"SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"}

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

	creds := Credentials{
		AccessToken:       os.Getenv("ACCESS_TOKEN"),
		AccessTokenSecret: os.Getenv("ACCESS_TOKEN_SECRET"),
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
	}

	logLevel := logrus.InfoLevel
	if os.Getenv("LOGLEVEL") == "Debug" {
		logLevel = logrus.DebugLevel
	}
	logger.Info("Setting Debug Level to ", logLevel)
	logger.SetLevel(logLevel)

	if os.Getenv("LOGLEVEL") == "Debug" {
		fmt.Printf("user=%s password=%s host=%s port=5432 dbname=%s sslmode=disable\n", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_DATABASE"))
	}

	dburi := fmt.Sprintf("user=%s password=%s host=%s port=5432 dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_DATABASE"))
	db, err = sql.Open("postgres", dburi)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		fmt.Println(dburi)
		panic(err)
	}

	loc, err = time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Println(err)
		return
	}

	//client, err = getClient(&creds)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)
	client = config.Client(oauth1.NoContext, token)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Printf("Credentials: %v\n", creds)
		log.Println(err)
	}
	// Dont tweet if we are dev

	conditions := 60
	forecast := 300
	alert := 900
	s := gocron.NewScheduler(loc)

	if os.Getenv("LOGLEVEL") == "Debug" {
		conditions = 1
		forecast = 1
		tc, err = s.Every(conditions).Minute().Do(twitterConditions)
		if err != nil {
			logger.Error(err)
		}
		tf, err = s.Every(forecast).Minute().Do(twitterForecast)
		if err != nil {
			logger.Error(err)
		}
		aj, err = s.Every(alert).Minute().Do(updateAlerts)
		if err != nil {
			logger.Error(err)
		}
		fj, err = s.Every(forecast).Minute().Do(twitterForecastImage)
		if err != nil {
			logger.Error(err)
		}
	} else {
		tc, err = s.Cron(os.Getenv("TC_CRON")).Do(twitterConditions)
		if err != nil {
			logger.Error(err)
		}
		tf, err = s.Cron(os.Getenv("TF_CRON")).Do(twitterForecast)
		if err != nil {
			logger.Error(err)
		}
		aj, err = s.Cron(os.Getenv("ALERT_CRON")).Do(updateAlerts)
		if err != nil {
			logger.Error(err)
		}
		fj, err = s.Cron(os.Getenv("FI_CRON")).Do(twitterForecastImage)
		if err != nil {
			logger.Error(err)
		}
	}

	if err != nil {
		logger.Error(tc, tf, err)
	}

	s.StartAsync()

	r := mux.NewRouter()
	r.HandleFunc("/job/{jobid}", loggingMiddleware(getJob))
	r.HandleFunc("/force/{type}", loggingMiddleware(force))
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

func updateAlerts() {
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
				logger.Error("No Alerts Found")
				id = -1
			} else {
				logger.Error("Scan: %v", err)
			}
		}
		if id == -1 {
			_, err := db.Exec(insertSql, v.IDURI, v.Type, v.AreaDesc, v.Sent.UTC(), v.Effective.UTC(), v.Onset.UTC(),
				v.Expires.UTC(), v.Ends.UTC(), v.Status,
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

func getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := vars["jobid"]
	switch t {
	case "tf":
		r := fmt.Sprintf("Last Run: %s, Next Run: %s", tf.LastRun(), tf.NextRun())
		w.Write([]byte(r))
	case "tc":
		r := fmt.Sprintf("Last Run: %s, Next Run: %s", tc.LastRun(), tc.NextRun())
		w.Write([]byte(r))
	case "aj":
		r := fmt.Sprintf("Last Run: %s, Next Run: %s", aj.LastRun(), aj.NextRun())
		w.Write([]byte(r))
	case "fj":
		r := fmt.Sprintf("Last Run: %s, Next Run: %s", fj.LastRun(), fj.NextRun())
		w.Write([]byte(r))
	}
}

func force(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := vars["type"]
	switch t {
	case "conditions":
		twitterConditions()
	case "forecast":
		twitterForecast()
	case "forecastImage":
		twitterForecastImage()
	}
}

//func getClient(creds *Credentials) (*twitter.Client, error) {
//	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
//	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
//	// Pass in your Access Token and your Access Token Secret
//	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)
//
//	httpClient := config.Client(oauth1.NoContext, token)
//	client := twitter.NewClient(httpClient)
//
//	// Verify Credentials
//	verifyParams := &twitter.AccountVerifyParams{
//		SkipStatus:   twitter.Bool(true),
//		IncludeEmail: twitter.Bool(true),
//	}
//
//	// we can retrieve the user and verify if the credentials
//	// we have used successfully allow us to log in!
//	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
//	if err != nil {
//		return nil, err
//	}
//	logger.Info("Logged ", user.Name, " into Twitter")
//	return client, nil
//}

func twitterConditions() {
	query := `select id,mac,recorded,baromabsin,baromrelin,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidityin,lastrain,
				maxdailygust,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,
				battlightning, lightningday,lightninghour,lightningtime,lightningdistance 
				from records order by recorded desc limit 1`
	rec := getRecord(query)
	t := buildMessage(rec)
	tweet(t, false)

}
func twitterForecast() {
	f, err := getForecast()
	if err != nil {
		logger.Error(err)
	}
	name := f.Properties.Periods[0].Name
	details := f.Properties.Periods[0].Detailedforecast
	t := fmt.Sprintf("%s ~ %s #COwx #KCOCOLOR663", name, details)
	tweet(t, false)
}

func twitterForecastImage() {
	tweet("Colorado Springs Forecast #COwx #KCOCOLOR663", true)
}

func getForecast() (Forecast, error) {
	url := "https://api.weather.gov/zones/forecast/COZ085/forecast"
	header := map[string]string{}
	res, err := makeRequest(url, "GET", nil, header)
	if err != nil {
		logger.Error(err)
	}
	f := Forecast{}
	err = json.Unmarshal(res, &f)
	if err != nil {
		logger.Error(err)
		return Forecast{}, err
	}
	return f, nil

}

func buildMessage(rec Record) string {
	t := ""
	if rec.Tempf < 50 {
		t = fmt.Sprintf("Temp: %.2f°F Sustained winds at %.2fmph, gusts to: %.2fmph, Feels like %.2f #COwx #KCOCOLOR663", rec.Tempf, rec.Windspeedmph, rec.Windgustmph, rec.Feelslike)
	} else {
		t = fmt.Sprintf("Temp: %.2f°F Sustained winds at %.2fmph, gusts to: %.2fmph, Rain: %.2fin Lightning: %d Today #COwx #KCOCOLOR663", rec.Tempf, rec.Windspeedmph, rec.Windgustmph, rec.Dailyrainin, rec.Lightningday)
	}
	return t
}

func tweet(message string, includeImage bool) {

	if logger.Level == logrus.DebugLevel {
		logger.Info(message)
		if includeImage {
			res := uploadImage()
			logger.Info(res.MediaID)
		}
	}
	values := url.Values{"status": {message}}
	if logger.Level != logrus.DebugLevel {
		if includeImage {
			res := uploadImage()
			values = url.Values{"status": {message}, "media_ids": {fmt.Sprintf("%d", res.MediaID)}}
		}

		// post status with media id
		resp, err := client.PostForm("https://api.twitter.com/1.1/statuses/update.json", values)
		if err != nil {
			logger.Error("Error: %s\n", err)
		}
		// parse response
		_, err = io.ReadAll(resp.Body)
		if err != nil {
			logger.Error("Error: %s\n", err)
		}
	}
}

func Open(f string) *os.File {
	r, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	return r
}

func uploadImage() *MediaResponse {
	var err error
	createImage()

	// create body form
	b := &bytes.Buffer{}
	form := multipart.NewWriter(b)

	// create media paramater
	fw, err := form.CreateFormFile("media", "image.png")
	if err != nil {
		panic(err)
	}
	// open file
	opened, err := os.Open("image.png")
	if err != nil {
		panic(err)
	}
	// copy to form
	_, err = io.Copy(fw, opened)
	if err != nil {
		panic(err)
	}
	// close form
	form.Close()

	// upload media
	resp, err := client.Post("https://upload.twitter.com/1.1/media/upload.json?media_category=tweet_image", form.FormDataContentType(), bytes.NewReader(b.Bytes()))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	}
	defer resp.Body.Close()

	// decode response and get media id
	m := &MediaResponse{}
	_ = json.NewDecoder(resp.Body).Decode(m)

	return m

}

func getRecord(sqlStatement string) Record {

	rows := db.QueryRow(sqlStatement)
	r := Record{}
	err := rows.Scan(&r.ID, &r.Mac, &r.Recorded, &r.Baromabsin, &r.Baromrelin, &r.Co2, &r.Dailyrainin, &r.Dewpoint, &r.Eventrainin, &r.Feelslike, &r.Hourlyrainin, &r.Hourlyrain, &r.Humidity, &r.Humidity1, &r.Humidity2, &r.Humidity3, &r.Humidityin, &r.Lastrain, &r.Maxdailygust, &r.Monthlyrainin, &r.Solarradiation, &r.Tempf, &r.Temp1f, &r.Temp2f, &r.Temp3f, &r.Tempinf, &r.Totalrainin, &r.Uv, &r.Weeklyrainin, &r.Winddir, &r.Windgustmph, &r.Windgustdir, &r.Windspeedmph, &r.Yearlyrainin, &r.Battlightning, &r.Lightningday, &r.Lightninghour, &r.Lightningtime, &r.Lightningdistance)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found ", sqlStatement)
		} else {
			logger.Error("Scan: %v", err)
		}
	}

	return r
}

func createImage() {

	img := gg.NewContext(2300, 1000)
	img.SetHexColor("#142a4d")
	img.Clear()

	img.SetHexColor("#FFFFFF")
	img.SetLineWidth(1)
	img.DrawLine(0, 120, 2300, 120)
	img.Stroke()

	if err := img.LoadFontFace("/go/bin/fonts/Arial.ttf", 96); err != nil {
		panic(err)
	}

	img.SetHexColor("#FFFFFF")
	img.DrawStringAnchored("COLORADO SPRINGS", 525, 60, 0.5, 0.5)

	img = DrawVerticle(img)

	img.SetHexColor("#FFFFFF")
	img.DrawRectangle(0, 824, 2300, 176)
	img.Fill()
	img = DaysOfWeek(img)
	img = ForecastValues(img)

	img.SavePNG("image.png")

}

func ForecastValues(img *gg.Context) *gg.Context {
	if err := img.LoadFontFace("/go/bin/fonts/Arial Black.ttf", 108); err != nil {
		panic(err)
	}
	forecast, err := getForecastImage()
	if err != nil {
		fmt.Println(err)
	}
	offSet := 163.0
	for i, v := range forecast.Days {

		img.SetHexColor("#FFFFFF")
		img.DrawStringAnchored(fmt.Sprintf("%d", int(v.Tempmax)), 326*float64(i)+offSet, 700, 0.5, 0.5)
		im, err := gg.LoadImage(fmt.Sprintf("/go/bin/icons/%s.png", v.Icon))
		if err != nil {
			log.Fatal(err)
		}
		img.DrawImage(im, 326*i+50, 300)
	}
	if err := img.LoadFontFace("/go/bin/fonts/Arial.ttf", 88); err != nil {
		panic(err)
	}

	for i, v := range forecast.Days {
		if i < 6 {
			img.SetHexColor("#142a4d")
			img.DrawStringAnchored(fmt.Sprintf("%d", int(v.Tempmin)), 326*float64(i+1), 875, 0.5, 0.5)
		}
	}

	return img
}

func DaysOfWeek(img *gg.Context) *gg.Context {
	if err := img.LoadFontFace("/go/bin/fonts/Arial Black.ttf", 76); err != nil {
		panic(err)
	}
	today := int(time.Now().Weekday())
	daysLeft := 7
	col := 0.0
	offSet := 163.0
	for t := today; t < len(days); t++ {
		img.SetHexColor("#FFFFFF")
		img.DrawStringAnchored(days[t], 326*col+offSet, 200, 0.5, 0.5)
		daysLeft--
		col = col + 1.0
	}
	for w := 0; w < daysLeft; w++ {
		img.SetHexColor("#FFFFFF")
		img.DrawStringAnchored(days[w], 326*col+offSet, 200, 0.5, 0.5)
		col = col + 1.0
	}
	return img
}

func DrawVerticle(img *gg.Context) *gg.Context {
	img.SetHexColor("#FFFFFF")
	img.SetLineWidth(1)
	for i := 0.0; i < 7.0; i++ {
		img.DrawLine(326*i, 120, 326*i, 824)
		img.Stroke()
	}
	img.DrawLine(2299, 120, 2299, 824)
	img.Stroke()
	return img
}

func getForecastImage() (ForecastImage, error) {
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Colorado%%20Springs?unitGroup=us&iconSets=icon2&include=days&key=%s&contentType=json", os.Getenv("WEATHER_API"))
	header := map[string]string{}
	res, err := makeRequest(url, "GET", nil, header)
	if err != nil {
		logger.Error(err)
	}
	f := ForecastImage{}
	err = json.Unmarshal(res, &f)
	if err != nil {
		logger.Error(err)
		return ForecastImage{}, err
	}
	return f, nil

}

func makeRequest(url string, method string, inputbody io.Reader, header map[string]string) ([]byte, error) {
	logger.Debug(url)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, inputbody)
	if err != nil {
		logger.Error(err)
	}
	if _, ok := header["User-Agent"]; !ok {
		req.Header.Add("User-Agent", `weather.zoms.net, brian@brianpaulson.com`)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return body, nil
}
