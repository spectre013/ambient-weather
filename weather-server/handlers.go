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

func alerts() []Alert {
	now := time.Now()
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Println(err)
		return []Alert{}
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
	return alerts
}

func astro() Astro {
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
	return astro
}

func trend(t string) Trend {
	sel := fmt.Sprintf("AVG(%s)", t)

	start := time.Now()
	end := start.Add(-30 * time.Minute)
	var avg float64
	cr := Record{}

	avgQuery := fmt.Sprintf("select %s from records where recorded BETWEEN '%s' AND '%s'", sel, formatDate(end), formatDate(start))
	logger.Debug(avgQuery)
	rows := db.QueryRow(avgQuery)
	err := rows.Scan(&avg)
	sqlError("Trend1", err, avgQuery)

	currentQuery := "select id,baromrelin,tempf from records order by recorded desc limit 1"
	logger.Debug(currentQuery)
	crows := db.QueryRow(currentQuery)
	err = crows.Scan(&cr.ID, &cr.Baromrelin, &cr.Tempf)
	sqlError("Trend2", err, currentQuery)

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
	return trend
}

func lightningMonth() int {
	d := getTimeframe("month")
	start := formatDate(d[0])
	end := formatDate(d[1])
	s := fmt.Sprintf(`SELECT SUM(A.value) as value
			FROM (SELECT TO_CHAR(recorded,'YYY-MM-DD') as ldate, 
				  MAX(lightningday) as value 
				  FROM records 
				  where recorded between '%s' and '%s' 
			GROUP BY ldate) A`, start, end)
	rows := db.QueryRow(s)
	lightningMonth := 0
	err := rows.Scan(&lightningMonth)
	sqlError("Lightning Month", err, "Error Getting Record: ")

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
	errorHandler(err, "Error getting chart records")
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
	sqlError("alltime", err, query)

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

func wind() Wind {
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
	sqlError("wind", err, maxSpeed)

	logger.Debug(maxGust)
	mgrows := db.QueryRow(maxGust)
	err = mgrows.Scan(&maxgust.Value, &maxgust.Reccorded)
	sqlError("wind", err, maxGust)

	logger.Debug(avgSpeed)
	asrows := db.QueryRow(avgSpeed)
	err = asrows.Scan(&avg)
	sqlError("wind", err, avgSpeed)

	logger.Debug(avgDir)
	crows := db.QueryRow(avgDir)
	err = crows.Scan(&avgdir)
	sqlError("wind", err, avgDir)

	wind := Wind{
		Maxdailygust: maxgust.Value,
		Avg:          int(avg),
	}
	return wind
}
func minmax(f string) map[string]map[string]StatValue {
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
	return minmax
}

func buildConditions(record Record) Conditions {
	conditions := Conditions{
		Mac:      record.Mac,
		Recorded: record.Recorded,
	}

	conditions.Barometer = Barometer{
		Baromabsin: record.Baromabsin,
		Baromrelin: record.Baromrelin,
		Trend:      trend("baromrelin"),
		MinMax:     minmax("baromrelin"),
	}

	conditions.Humidity = Humidity{
		Humidity: record.Humidity,
		Dewpoint: record.Dewpoint,
		MinMax:   minmax("humidity"),
	}

	conditions.Temp = Temp{
		Temp:      record.Tempf,
		Humidity:  record.Humidity,
		Feelslike: record.Feelslike,
		Dewpoint:  record.Dewpoint,
		MinMax:    minmax("tempf"),
	}

	conditions.Tempin = Tempin{
		Temp:     record.Tempinf,
		Humidity: record.Humidityin,
		MinMax:   minmax("tempinf"),
	}

	conditions.Temp1 = Tempin{
		Temp:     record.Temp1f,
		Humidity: record.Humidity1,
		MinMax:   minmax("temp1f"),
	}

	conditions.Temp2 = Tempin{
		Temp:     record.Temp2f,
		Humidity: record.Humidity2,
		MinMax:   minmax("temp2f"),
	}

	conditions.Temp3 = Tempin{
		Temp:     record.Temp3f,
		Humidity: record.Humidity3,
		MinMax:   minmax("temp3f"),
	}

	w := wind()
	conditions.Wind = Wind{
		Dir:          record.Winddir,
		Gustmph:      record.Windgustmph,
		Gustdir:      record.Windgustdir,
		Speedmph:     record.Windspeedmph,
		Maxdailygust: w.Maxdailygust,
		Avg:          w.Avg,
		MinMax:       minmax("windspeedmph"),
	}

	conditions.UV = UV{
		Uv:             record.Uv,
		Solarradiation: record.Solarradiation,
		MinMax:         minmax("uv"),
	}

	conditions.Rain = Rain{
		Daily:    record.Dailyrainin,
		Weekly:   record.Weeklyrainin,
		Monthly:  record.Monthlyrainin,
		Yearly:   record.Yearlyrainin,
		Total:    record.Totalrainin,
		Lastrain: record.Lastrain,
	}

	conditions.Alert = alerts()

	conditions.Lightning = Lightning{
		Day:      record.Lightningday,
		Hour:     record.Lightninghour,
		Distance: record.Lightningdistance,
		Time:     record.Lightningtime,
		Month:    record.LightningMonth,
		Minmax:   minmax("lightning"),
	}

	conditions.AQI = AQI{
		Pm25:    record.Aqipm25,
		Pm2524h: record.Aqipm2524h,
		MinMax:  minmax("aqipm25"),
	}

	conditions.Astro = astro()

	return conditions

}

func getConditions() Conditions {
	res := getCurrent()
	cond := buildConditions(res)
	return cond
}

func current(w http.ResponseWriter, _ *http.Request) {
	cond := getConditions()

	b, err := json.Marshal(cond)
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
