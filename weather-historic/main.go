package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"time"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var logger = logrus.New()
var loc *time.Location
var apiurl = "https://api.ambientweather.net/v1/devices/84:F3:EB:20:DA:9E?apiKey=%s&applicationKey=%s&endDate=%s&limit=%d"
var APIKEY = os.Getenv("AMBIENT_WEATHER_API_KEY")
var APPKEY = os.Getenv("AMBIENT_WEATHER_APPLICATION_KEY")

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

	dburi := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=true&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_DATABASE"))
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

	gaps := findGaps()
	//fmt.Println(gaps)
	buildURLS(gaps)
}

func buildURLS(gaps []Gap) []string {
	var result []string

	for _,v := range gaps {
		fmt.Println(v.From, v.duration.Minutes()/5, v.To)
		result = append(result,fmt.Sprintf(apiurl,APIKEY,APPKEY,v.To,10))
	}
	return result
}


func findGaps() []Gap {
	type Result struct {
		Date time.Time
	}
	var results []Result
	err := db.Raw("select date from records order by date desc").Scan(&results) // (*sql.Rows, error)
	if err != nil {
		log.Println(err)
	}

	gaps := make([]Gap,0)

	for k,v := range results {
		if k+1 < len(results) {
			dur := v.Date.Sub(results[k+1].Date)
			if dur.Minutes() > 5 {
				gaps = append(gaps, Gap{
					From: v.Date,
					To:   results[k+1].Date,
					duration: dur,
				})
				//fmt.Println(dur,v.Date, results[k+1].Date,dur)
			}
		}
	}

	return gaps
}

type Gap struct {
	From time.Time
	To  time.Time
	duration time.Duration
}

type Record struct {
	ID               int       `json:"id" db:"id"`
	Mac              string    `json:"mac" db:"mac"`
	Date             time.Time `json:"date" db:"date"`
	Baromabsin       float64   `json:"baroabsin" db:"baromabsin"`
	Baromrelin       float64   `json:"baromrelin" db:"baromrelin"`
	Battout          string    `json:"battout" db:"battout"`
	Batt1            string    `json:"batt1" db:"batt1"`
	Batt2            string    `json:"batt2" db:"batt2"`
	Batt3            string    `json:"batt3" db:"batt3"`
	Batt4            string    `json:"batt4" db:"batt4"`
	Batt5            string    `json:"batt5" db:"batt5"`
	Batt6            string    `json:"batt6" db:"batt6"`
	Batt7            string    `json:"batt7" db:"batt7"`
	Batt8            string    `json:"batt8" db:"batt8"`
	Batt9            string    `json:"batt9" db:"batt9"`
	Batt10           string    `json:"batt10" db:"batt10"`
	Co2              float64   `json:"co2" db:"co2"`
	Dailyrainin      float64   `json:"dailyrainin" db:"dailyrainin"`
	Dewpoint         float64   `json:"dewpoint" db:"dewpoint"`
	Eventrainin      float64   `json:"eventrain" db:"eventrainin"`
	Feelslike        float64   `json:"feelslike" db:"feelslike"`
	Hourlyrainin     float64   `json:"hourlyrainin" db:"hourlyrainin"`
	Hourlyrain       float64   `json:"hourlyrain"`
	Humidity         int       `json:"humidity" db:"humidity"`
	Humidity1        int       `json:"humidity1" db:"humidity1"`
	Humidity2        int       `json:"humidity2" db:"humidity2"`
	Humidity3        int       `json:"humidity3" db:"humidity3"`
	Humidity4        int       `json:"humidity4" db:"humidity4"`
	Humidity5        int       `json:"humidity5" db:"humidity5"`
	Humidity6        int       `json:"humidity6" db:"humidity6"`
	Humidity7        int       `json:"humidity7" db:"humidity7"`
	Humidity8        int       `json:"humidity8" db:"humidity8"`
	Humidity9        int       `json:"humidity9" db:"humidity9"`
	Humidity10       int       `json:"humidity10" db:"humidity10"`
	Humidityin       int       `json:"humidityin" db:"humidityin"`
	Lastrain         time.Time `json:"lastrain" db:"lastrain"`
	Maxdailygust     float64   `json:"maxdailygust" db:"maxdailygust"`
	Relay1           int       `json:"relay1" db:"relay1"`
	Relay2           int       `json:"relay2" db:"relay2"`
	Relay3           int       `json:"relay3" db:"relay3"`
	Relay4           int       `json:"relay4" db:"relay4"`
	Relay5           int       `json:"relay5" db:"relay5"`
	Relay6           int       `json:"relay6" db:"relay6"`
	Relay7           int       `json:"relay7" db:"relay7"`
	Relay8           int       `json:"relay8" db:"relay8"`
	Relay9           int       `json:"relay9" db:"relay9"`
	Relay10          int       `json:"relay10" db:"relay10"`
	Monthlyrainin    float64   `json:"monthlyrainin" db:"monthlyrainin"`
	Solarradiation   float64   `json:"solarradiation" db:"solarradiation"`
	Tempf            float64   `json:"tempf" db:"tempf"`
	Temp1f           float64   `json:"temp1f" db:"temp1f"`
	Temp2f           float64   `json:"temp2f" db:"temp2f"`
	Temp3f           float64   `json:"temp3f" db:"temp3f"`
	Temp4f           float64   `json:"temp4f" db:"temp4f"`
	Temp5f           float64   `json:"temp5f" db:"temp5f"`
	Temp6f           float64   `json:"temp6f" db:"temp6f"`
	Temp7f           float64   `json:"temp7f" db:"temp7f"`
	Temp8f           float64   `json:"temp8f" db:"temp8f"`
	Temp9f           float64   `json:"temp9f" db:"temp9f"`
	Temp10f          float64   `json:"temp10f" db:"temp10f"`
	Tempinf          float64   `json:"tempinf" db:"tempinf"`
	Totalrainin      float64   `json:"totalrainin" db:"totalrainin"`
	Uv               float64   `json:"uv" db:"uv"`
	Weeklyrainin     float64   `json:"weeklyrainin" db:"weeklyrainin"`
	Winddir          int       `json:"winddir" db:"winddir"`
	Windgustmph      float64   `json:"windgustmph" db:"windgustmph"`
	Windgustdir      int       `json:"windgustdir" db:"windgustdir"`
	Windspeedmph     float64   `json:"windspeedmph" db:"windspeedmph"`
	Winddiravg2m     int       `json:"winddiravg2m" db:"winddir_avg2m"`
	Windspdmphavg2m  float64   `json:"windspeedavg2m" db:"windspdmph_avg2m"`
	Winddiravg10m    int       `json:"winddiravg10m" db:"winddir_avg10m"`
	Windspdmphavg10m float64   `json:"windspeedavg10m" db:"windspdmph_avg10m"`
	Yearlyrainin     float64   `json:"yearlyrainin" db:"yearlyrainin"`
}