package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// other imports
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	_ "github.com/lib/pq"
)

var db *sql.DB
var logger = logrus.New()
var loc *time.Location
var client *twitter.Client
var tc *gocron.Job
var tf *gocron.Job
var aj *gocron.Job
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

	client, err = getClient(&creds)
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
	}

	if err != nil {
		logger.Error(tc, tf, err)
	}

	s.StartAsync()

	r := mux.NewRouter()
	r.HandleFunc("/job/{jobid}", loggingMiddleware(getJob))
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


	checkSql := `select id from alerts where alertid = '%s'`


	iAlerts := getAlerts()
	for _,v := range iAlerts {
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
		//if _, ok := resp.Header["Last-Modified"]; ok {
		//	lModified := ""
		//	if len(resp.Header["Last-Modified"]) > 0 {
		//		lModified = resp.Header["Last-Modified"][0]
		//	}
		//	l, err := time.Parse("Mon, 2 Jan 2006 15:04:05 GMT", lModified)
		//	if err != nil {
		//		return []byte(""), err
		//	}
		//	logger.Info(l.Before(LastModified), l, LastModified)
		//	if (l.Before(LastModified) || l.Equal(LastModified)) && !LastModified.IsZero() {
		//		logger.Debug("No Updates")
		//		logger.Debug(l, LastModified)
		//		err = errors.New("cached")
		//		return body, err
		//	}
		//}
	logger.Debug("Alert Updates")
	if t != "head" {
		body, err = ioutil.ReadAll(resp.Body)
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
	}
}

func getClient(creds *Credentials) (*twitter.Client, error) {
	// Pass in your consumer key (API Key) and your Consumer Secret (API Secret)
	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	// Pass in your Access Token and your Access Token Secret
	token := oauth1.NewToken(creds.AccessToken, creds.AccessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)

	// Verify Credentials
	verifyParams := &twitter.AccountVerifyParams{
		SkipStatus:   twitter.Bool(true),
		IncludeEmail: twitter.Bool(true),
	}

	// we can retrieve the user and verify if the credentials
	// we have used successfully allow us to log in!
	user, _, err := client.Accounts.VerifyCredentials(verifyParams)
	if err != nil {
		return nil, err
	}
	logger.Info("Logged ", user.Name, " into Twitter")
	return client, nil
}

func twitterConditions() {
	query := `select id,mac,recorded,baromabsin,baromrelin,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidityin,lastrain,
				maxdailygust,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,
				battlightning, lightningday,lightninghour,lightningtime,lightningdistance 
				from records order by recorded desc limit 1`
	rec := getRecord(query)
	t := buildMessage(rec)
	tweet(t)

}
func twitterForecast() {
	f, err := getForecast()
	if err != nil {
		logger.Error(err)
	}
	name := f.Properties.Periods[0].Name
	details := f.Properties.Periods[0].Detailedforecast
	t := fmt.Sprintf("%s ~ %s #COwx #KCOCOLOR663", name, details)
	tweet(t)
}

func getForecast() (Forecast, error) {
	url := "https://api.weather.gov/zones/forecast/COZ085/forecast"
	header := map[string]string{}
	res, err := makeRequest(url, header)
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

func tweet(message string) {
	if logger.Level == logrus.DebugLevel {
		logger.Info(message)
	}

	if logger.Level != logrus.DebugLevel {
		_, _, err := client.Statuses.Update(message, nil)
		if err != nil {
			logger.Println(err)
		}
	}
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

func makeRequest(url string, header map[string]string) ([]byte, error) {
	logger.Debug(url)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return body, nil
}
