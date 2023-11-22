package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/IvanMenshykov/MoonPhase"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"os"
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

func apiin(w http.ResponseWriter, r *http.Request) {
	rec := Record{}
	logger.Info("Incoming Ambient Data Processing ....")
	err := json.NewDecoder(r.Body).Decode(&rec)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	inserted := insertRecord(rec)
	if inserted {
		go calculateStats()
	}

	i, err := w.Write([]byte("Success"))
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ambientin(w http.ResponseWriter, r *http.Request) {
	output := map[string]interface{}{}

	in := map[string]string{}
	in["baromabsin"] = "float"
	in["baromrelin"] = "float"
	in["batt_co2"] = "int"
	in["battlightning"] = "int"
	in["batt1"] = "int"
	in["batt2"] = "int"
	in["batt3"] = "int"
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

	values := r.URL.Query()
	if len(values) == 0 {
		return
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
		case "float":
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Printf("%s - %s\n", err, val)
			}
			output[k] = f
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

		}
	}
	output["date"] = time.Now()
	lastrainquery := "select recorded from records r where dailyrainin > 0 order by recorded desc limit 1"

	logger.Debug(lastrainquery)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	rec := Record{}
	err = json.Unmarshal(b, &rec)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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

	inserted := insertRecord(rec)
	if inserted {
		go calculateStats()
	}

	w.WriteHeader(200)
}

func alerts(w http.ResponseWriter, _ *http.Request) {
	now := time.Now()
	alertsSql := fmt.Sprintf("select * from alerts where onset <= '%s' and ends >= '%s'", formatDate(now), formatDate(now))
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

func astro(w http.ResponseWriter, _ *http.Request) {
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

	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
		if strings.Contains(v.ID, f) {
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

func getCurrentApp(w http.ResponseWriter, _ *http.Request) {
	query := `select id,mac,recorded,baromabsin,baromrelin,battout,Batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,
				maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,battlightning,lightningday,lightninghour,lightningtime,lightningdistance, aqipm25, aqipm2524h 
				from records order by recorded desc limit 1`
	rows := db.QueryRow(query)
	r := RecordApp{}
	err := rows.Scan(&r.ID, &r.Mac, &r.Recorded, &r.Baromabsin, &r.Baromrelin, &r.Battout, &r.Batt1, &r.Batt2, &r.Batt3, &r.Batt4, &r.Batt5, &r.Batt6, &r.Batt7, &r.Batt8, &r.Batt9, &r.Batt10, &r.Co2, &r.Dailyrainin, &r.Dewpoint, &r.Eventrainin, &r.Feelslike, &r.Hourlyrainin, &r.Hourlyrain, &r.Humidity, &r.Humidity1, &r.Humidity2, &r.Humidity3, &r.Humidity4, &r.Humidity5, &r.Humidity6, &r.Humidity7, &r.Humidity8, &r.Humidity9, &r.Humidity10, &r.Humidityin, &r.Lastrain, &r.Maxdailygust, &r.Relay1, &r.Relay2, &r.Relay3, &r.Relay4, &r.Relay5, &r.Relay6, &r.Relay7, &r.Relay8, &r.Relay9, &r.Relay10, &r.Monthlyrainin, &r.Solarradiation, &r.Tempf, &r.Temp1f, &r.Temp2f, &r.Temp3f, &r.Temp4f, &r.Temp5f, &r.Temp6f, &r.Temp7f, &r.Temp8f, &r.Temp9f, &r.Temp10f, &r.Tempinf, &r.Totalrainin, &r.Uv, &r.Weeklyrainin, &r.Winddir, &r.Windgustmph, &r.Windgustdir, &r.Windspeedmph, &r.Yearlyrainin, &r.Battlightning, &r.Lightningday, &r.Lightninghour, &r.Lightningtime, &r.Lightningdistance)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found ", query)
		} else {
			logger.Error("Scan:", err)
		}
	}

	hourlyRain := getHourlyRain()
	r.Hourlyrain = hourlyRain
	//f := getForecast()
	stats := getStats()

	r.Sunrise = "6:11"     //f.Days[0].Sunrise
	r.Sunset = "6:32"      //f.Days[0].Sunset
	r.Conditions = "Sunny" //f.Days[0].Conditions
	r.Visibility = 10      //f.Days[0].Visibility
	for _, v := range stats {
		if v.ID == "day_max_tempf" {
			fmt.Printf("%v", v)
			r.MaxTemp = v.Value
		}
		if v.ID == "day_min_tempf" {
			fmt.Printf("%v", v)
			r.MinTemp = v.Value
		}
		if v.ID == "day_avg_windspeedmph" {
			fmt.Printf("%v", v)
			r.AvgWind = v.Value
		}
	}

	b, err := json.Marshal(r)
	if err != nil {
		log.Println(err)
	}
	i, err := w.Write(b)
	if err != nil {
		logger.Error(i, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
