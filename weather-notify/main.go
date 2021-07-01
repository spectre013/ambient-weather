package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	client, err = getClient(&creds)
	if err != nil {
		log.Println("Error getting Twitter Client")
		log.Printf("Credentials: %v\n", creds)
		log.Println(err)
	}
	// Dont tweet if we are dev

	conditions := 60
	forecast := 300
	s := gocron.NewScheduler(loc)

	if os.Getenv("LOGLEVEL") == "Debug" {
		conditions = 1
		forecast = 1
		tc, err = s.Every(conditions).Minute().Do(twitterConditions)
		tf, err = s.Every(forecast).Minute().Do(twitterForecast)
	} else {
		tc, err = s.Cron(os.Getenv("TC_CRON")).Do(twitterConditions)
		tf, err = s.Cron(os.Getenv("TF_CRON")).Do(twitterForecast)
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

func getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := vars["jobid"]
	switch t {
	case "tf":
		r := fmt.Sprintf("Last Run: %s, Next Run: %s",tf.LastRun(),tf.NextRun())
		w.Write([]byte(r))
		break
	case "tc":
		r := fmt.Sprintf("Last Run: %s, Next Run: %s",tc.LastRun(),tc.NextRun())
		w.Write([]byte(r))
		break
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
	f,err := getForecast()
	if err != nil {
		logger.Error(err)
	}
	name := f.Properties.Periods[0].Name
	details := f.Properties.Periods[0].Detailedforecast
	t := fmt.Sprintf("%s ~ %s #COwx #KCOCOLOR663",name,details)
	tweet(t)
}

func getForecast() (Forecast,error) {
 url := "https://api.weather.gov/zones/forecast/COZ085/forecast"
	header := map[string]string{}
	res, err := makeRequest(url, header)
	f:= Forecast{}
	err = json.Unmarshal(res,&f)
	if err != nil {
		logger.Error(err)
		return Forecast{},err
	}
	return f,nil

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