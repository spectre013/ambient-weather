package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nathan-osman/go-sunrise"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func home(w http.ResponseWriter, _ *http.Request) {
	res := map[string]string{}
	res["message"] = "Zoms.net weather API visit https://weather.zoms.net for more information"
	b, _ := json.Marshal(res)
	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func alerts(w http.ResponseWriter, _ *http.Request) {
	now := time.Now()
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Println(err)
		return
	}
	alertsSql := fmt.Sprintf("select * from alerts where ends >= '%s'", formatDate(now))
	logger.Info(alertsSql)
	rows, err := db.Query(alertsSql)
	if err != nil {
		logger.Error(err)
	}
	alerts := make([]Alert, 0)
	for rows.Next() {
		a := Alert{}
		err := rows.Scan(&a.ID, &a.Alertid, &a.Wxtype, &a.Areadesc, &a.Sent, &a.Effective, &a.Onset, &a.Expires, &a.Ends, &a.Status, &a.Messagetype, &a.Category, &a.Severity, &a.Certainty, &a.Urgency, &a.Event, &a.Sender, &a.SenderName, &a.Headline, &a.Description, &a.Instruction, &a.Response)
		if err != nil {
			if err == sql.ErrNoRows {
				logger.Error("Zero Rows Found ", alertsSql)
			} else {
				logger.Error("Scan: ", err)
			}
		}

		a.Sent = a.Sent.In(loc)
		a.Effective = a.Effective.In(loc)
		a.Onset = a.Onset.In(loc)
		a.Ends = a.Ends.In(loc)

		alerts = append(alerts, a)
	}
	b, err := json.Marshal(&alerts)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func astro(w http.ResponseWriter, r *http.Request) {
	lat, err := strconv.ParseFloat(os.Getenv("LAT"), 64)
	errorHandler(err, "Error parsing LAT")
	lng, err := strconv.ParseFloat(os.Getenv("LON"), 64)
	errorHandler(err, "Error parsing LON")

	rise, set := sunrise.SunriseSunset(
		lat, lng, time.Now().Year(), time.Now().Month(), time.Now().Day(),
	)
	t := time.Now().Add(24 * time.Hour)
	riset, sett := sunrise.SunriseSunset(
		lat, lng, t.Year(), t.Month(), t.Day(),
	)
	elevation := sunrise.Elevation(lat, lng, time.Now())
	hasSunSet := false
	if elevation <= 0 {
		hasSunSet = true
	}

	astro := Astro{
		Sunrise:         rise,
		Sunset:          set,
		SunriseTomorrow: riset,
		SunsetTomorrow:  sett,
		Darkness:        riset.Sub(set),
		Daylight:        set.Sub(rise),
		Elevation:       elevation,
		HasSunset:       hasSunSet,
	}

	b, err := json.Marshal(astro)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func lightning() int {

	s := `SELECT SUM(A.value) as value
			FROM (SELECT TO_CHAR(recorded,'YYY-MM-DD') as ldate, 
				  MAX(lightningday) as value 
				  FROM records 
				  where recorded between '04-01-2024' and '04-30-2024' 
			GROUP BY ldate) A`
	rows := db.QueryRow(s)
	lightningMonth := 0
	err := rows.Scan(&lightningMonth)
	sqlError(err, "Error Getting Record: ")

	return lightningMonth
}

func chart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := vars["type"]
	p := vars["period"]
	t = cleanString(t)
	dates := getTimeframe(p)
	dateformat := "TO_CHAR(recorded,'YYYY-MM-DD')"
	if p == "day" || p == "yesterday" {
		dateformat = "TO_CHAR(recorded,'HH24:MI:SS')"
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
		where = fmt.Sprintf("%s AS mmdd, max(%s) AS max", dateformat, t)
	} else {
		where = fmt.Sprintf("%s AS mmdd, max(%s) AS max, min(%s) AS min", dateformat, t, t)
	}

	query := fmt.Sprintf("select %s from records where recorded BETWEEN '%s' AND '%s' group by mmdd order by mmdd", where, formatDate(dates[0]), formatDate(dates[1]))
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
			logger.Error("Scan: ", err)
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

	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

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
			logger.Error("Scan: ", err)
		}
	}

	b, err := json.Marshal(rt)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
			logger.Error("Scan: ", err)
		}
	}
	currentQuery := "select id,mac,recorded,baromabsin,baromrelin,battout,batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,hourlyrainin,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,hourlyrain,lightningday,lightninghour,lightningtime,lightningdistance,battlightning from records order by recorded desc limit 1"
	logger.Debug(currentQuery)
	crows := db.QueryRow(currentQuery)
	err = crows.Scan(&cr.ID, &cr.Mac, &cr.Recorded, &cr.Baromabsin, &cr.Baromrelin, &cr.Battout, &cr.Batt1, &cr.Batt2, &cr.Batt3, &cr.Batt4, &cr.Batt5, &cr.Batt6, &cr.Batt7, &cr.Batt8, &cr.Batt9, &cr.Batt10, &cr.Co2, &cr.Dailyrainin, &cr.Dewpoint, &cr.Eventrainin, &cr.Feelslike, &cr.Hourlyrainin, &cr.Humidity, &cr.Humidity1, &cr.Humidity2, &cr.Humidity3, &cr.Humidity4, &cr.Humidity5, &cr.Humidity6, &cr.Humidity7, &cr.Humidity8, &cr.Humidity9, &cr.Humidity10, &cr.Humidityin, &cr.Lastrain, &cr.Maxdailygust, &cr.Relay1, &cr.Relay2, &cr.Relay3, &cr.Relay4, &cr.Relay5, &cr.Relay6, &cr.Relay7, &cr.Relay8, &cr.Relay9, &cr.Relay10, &cr.Monthlyrainin, &cr.Solarradiation, &cr.Tempf, &cr.Temp1f, &cr.Temp2f, &cr.Temp3f, &cr.Temp4f, &cr.Temp5f, &cr.Temp6f, &cr.Temp7f, &cr.Temp8f, &cr.Temp9f, &cr.Temp10f, &cr.Tempinf, &cr.Totalrainin, &cr.Uv, &cr.Weeklyrainin, &cr.Winddir, &cr.Windgustmph, &cr.Windgustdir, &cr.Windspeedmph, &cr.Yearlyrainin, &cr.Hourlyrain, &cr.Lightningday, &cr.Lightninghour, &cr.Lightningtime, &cr.Lightningdistance, &cr.Battlightning)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", currentQuery)
		} else {
			logger.Error("Scan: ", err)
		}
	}

	trend := Trend{}
	if strings.Contains(t, "temp") {
		if cr.Tempf > avg {
			//trend up
			trend.Trend = "up"
			trend.By = toFixed(cr.Tempf-avg, 2)
		} else {
			//trend down
			trend.Trend = "down"
			trend.By = toFixed(avg-cr.Tempf, 2)
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

	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func wind(w http.ResponseWriter, _ *http.Request) {
	type Result struct {
		Reccorded time.Time `json:"date"`
		Value     float64   `json:"value"`
	}
	dates := getTimeframe("day")
	start := dates[1]
	end := dates[0]

	maxwind := Result{}
	maxgust := Result{}
	var avg float64
	var avgdir float64

	maxSpeed := fmt.Sprintf("select windspeedmph as value, recorded from records where recorded BETWEEN '%s' AND '%s' order by windspeedmph desc limit 1", formatDate(end), formatDate(start))
	maxGust := fmt.Sprintf("select windgustmph as value, recorded from records where recorded BETWEEN '%s' AND '%s' order by windgustmph desc limit 1", formatDate(end), formatDate(start))
	avgSpeed := fmt.Sprintf("select AVG(windspeedmph) as value from records where recorded BETWEEN '%s' AND '%s'", formatDate(end), formatDate(start))
	avgDir := fmt.Sprintf("select AVG(winddir) as value from records where recorded BETWEEN '%s' AND '%s'", formatDate(end), formatDate(start))

	logger.Debug(maxSpeed)
	mrows := db.QueryRow(maxSpeed)
	err := mrows.Scan(&maxwind.Value, &maxwind.Reccorded)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows ", maxSpeed)
		} else {
			logger.Error("Scan: ", err)
		}
	}

	logger.Debug(maxGust)
	mgrows := db.QueryRow(maxGust)
	err = mgrows.Scan(&maxgust.Value, &maxgust.Reccorded)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", maxGust)
		} else {
			logger.Error("Scan: ", err)
		}
	}

	logger.Debug(avgSpeed)
	asrows := db.QueryRow(avgSpeed)
	err = asrows.Scan(&avg)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", avgSpeed)
		} else {
			logger.Error("Scan: ", err)
		}
	}

	logger.Debug(avgDir)
	crows := db.QueryRow(avgDir)
	err = crows.Scan(&avgdir)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found", avgdir)
		} else {
			logger.Error("Scan: ", err)
		}
	}

	res := map[string]Result{}
	res["dir"] = Result{Value: avgdir}
	res["wind"] = maxwind
	res["gust"] = maxgust
	res["avg"] = Result{Value: avg}

	b, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func minmax(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	f := vars["field"]
	s := getStats()

	minmax := make(map[string]map[string]StatValue)

	for _, v := range s {
		matched, _ := regexp.MatchString(f+"$", v.ID)
		if matched {
			parts := strings.Split(v.ID, "_")

			if _, ok := minmax[parts[1]]; !ok {
				minmax[parts[1]] = make(map[string]StatValue)
			}
			minmax[parts[1]][strings.ToLower(parts[0])] = StatValue{Recorded: v.Recorded, Value: v.Value}
		}
	}
	b, err := json.Marshal(minmax)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
func current(w http.ResponseWriter, _ *http.Request) {
	res, err := getCurrent()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	i, err := w.Write(res)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getForecastHandler(w http.ResponseWriter, _ *http.Request) {
	res, err := getForecast()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	b, err := json.Marshal(res)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
