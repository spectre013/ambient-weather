package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/dyindude/moonphase"
	"github.com/gorilla/mux"
	"github.com/nathan-osman/go-sunrise"
)

func home(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Zoms.net weather API visit https://weather.zoms.net for more information",
	})
}

// alerts now caches results via alertsCache.
func alerts() []Alert {
	v := alertsCache.Get(func() interface{} {
		return loadAlerts()
	})
	if v == nil {
		return []Alert{}
	}
	return v.([]Alert)
}

func loadAlerts() []Alert {
	now := time.Now()
	loc := config.Location

	// Date is now parameterized.
	const alertsSQL = "select * from alerts where ends >= $1"
	logger.Debug(alertsSQL)

	rows, err := db.Query(alertsSQL, now)
	if err != nil {
		logger.WithError(err).Error("alerts query")
		return []Alert{}
	}
	defer rows.Close()

	alerts := make([]Alert, 0)
	for rows.Next() {
		a := Alert{}
		if err := rows.Scan(&a.AlertID, &a.Wxtype, &a.Areadesc, &a.Sent, &a.Effective, &a.Onset, &a.Expires, &a.Ends, &a.Status, &a.Messagetype, &a.Category, &a.Severity, &a.Certainty, &a.Urgency, &a.Event, &a.Sender, &a.SenderName, &a.Headline, &a.Description, &a.Instruction, &a.Response, &a.Summary); err != nil {
			if err == sql.ErrNoRows {
				continue
			}
			logger.WithError(err).Error("alerts scan")
			continue
		}
		a.Sent = a.Sent.In(loc)
		a.Effective = a.Effective.In(loc)
		a.Onset = a.Onset.In(loc)
		a.Ends = a.Ends.In(loc)
		alerts = append(alerts, a)
	}
	if err := rows.Err(); err != nil {
		logger.WithError(err).Error("alerts rows.Err")
	}
	return alerts
}

// astro now uses config.Latitude / config.Longitude rather than reading
// and parsing env vars on every call.
func astro() Astro {
	lat := config.Latitude
	lng := config.Longitude

	rise, set := sunrise.SunriseSunset(lat, lng, time.Now().Year(), time.Now().Month(), time.Now().Day())
	t := time.Now().Add(24 * time.Hour)
	riset, sett := sunrise.SunriseSunset(lat, lng, t.Year(), t.Month(), t.Day())
	elevation := sunrise.Elevation(lat, lng, time.Now())

	hasSunSet := elevation <= 0
	phase, illumination, name := getMoonPhase()

	return Astro{
		Sunrise:         rise,
		Sunset:          set,
		SunriseTomorrow: riset,
		SunsetTomorrow:  sett,
		Darkness:        riset.Sub(set),
		Daylight:        set.Sub(rise),
		Elevation:       correctSunElevation(elevation, time.Now(), rise, set),
		HasSunset:       hasSunSet,
		MoonIlluminance: illumination,
		MoonPhase:       phase,
		MoonPhaseName:   name,
	}
}

func getMoonPhase() (float64, float64, string) {
	m := moonphase.New(time.Now())
	return m.Phase(), m.Illumination() * 100, m.PhaseName()
}

// trend uses a column allowlist so the interpolated identifier can never
// be user-controlled. Dates are parameterized.
func trend(t string) Trend {
	if !allowedTrendColumns[t] {
		logger.WithField("column", t).Error("trend: column not in allowlist")
		return Trend{}
	}

	end := time.Now()
	start := end.Add(-30 * time.Minute)

	avgQuery := fmt.Sprintf("select AVG(%s) from records where recorded BETWEEN $1 AND $2", t)
	logger.Debug(avgQuery)

	var avg float64
	if err := db.QueryRow(avgQuery, start, end).Scan(&avg); err != nil {
		sqlError("Trend1", err, avgQuery)
	}

	const currentQuery = "select id,baromrelin,tempf from records order by recorded desc limit 1"
	logger.Debug(currentQuery)

	cr := Record{}
	if err := db.QueryRow(currentQuery).Scan(&cr.ID, &cr.Baromrelin, &cr.Tempf); err != nil {
		sqlError("Trend2", err, currentQuery)
	}

	trend := Trend{}
	if strings.Contains(t, "temp") {
		if cr.Tempf > avg {
			trend.Trend = "up"
			trend.By = toFixed(cr.Tempf-avg, 2)
		} else {
			trend.Trend = "down"
			trend.By = toFixed(avg-cr.Tempf, 2)
		}
	} else {
		if cr.Baromrelin > avg {
			trend.Trend = "Steady"
			if (cr.Baromrelin - avg) > .5 {
				trend.Trend = "Rising"
			}
		} else {
			trend.Trend = "Steady"
			if (avg - cr.Baromrelin) > .5 {
				trend.Trend = "Falling"
			}
		}
	}
	return trend
}

func getClimate(w http.ResponseWriter, _ *http.Request) {
	climate := ConvertRawToClimateRecords(almanacQueries())
	writeJSON(w, http.StatusOK, climate)
}

func getFirstFreeze(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, firstFreeze())
}

// lightningMonth uses parameterized dates and the corrected YYYY-MM-DD
// date format.
func lightningMonth() int {
	d := getTimeframe("month")
	if len(d) < 2 {
		return 0
	}

	const query = `SELECT SUM(A.value) as value
		FROM (SELECT TO_CHAR(recorded, 'YYYY-MM-DD') as ldate,
		             MAX(lightningday) as value
		      FROM records
		      WHERE recorded BETWEEN $1 AND $2
		      GROUP BY ldate) A`

	var lightning sql.NullInt64
	err := db.QueryRow(query, d[0], d[1]).Scan(&lightning)
	if err != nil && err != sql.ErrNoRows {
		sqlError("Lightning Month", err, "query failed")
	}
	if lightning.Valid {
		return int(lightning.Int64)
	}
	return 0
}

type ChartValue struct {
	Ts    time.Time `json:"date"`
	Value float64   `json:"value"`
}

type ChartData struct {
	Values []ChartValue `json:"values"`
	Key    string       `json:"key"`
	Color  string       `json:"color"`
}

// ChartSensor describes one series in a chart bundle. Typed
// struct rather than nested map[string]map[string][]map[string]string.
type ChartSensor struct {
	Sensor string
	Color  string
	Title  string
}

// chartConfig replaces the deeply-nested map literal that was rebuilt
// on every request.
var chartConfig = map[string][]ChartSensor{
	"temperature": {
		{Sensor: "tempf", Color: "#EE4B2B", Title: "Temperature"},
		{Sensor: "dewpoint", Color: "blue", Title: "Dewpoint"},
	},
	"humidity": {
		{Sensor: "humidity", Color: "green", Title: "Humidity"},
	},
	"windspeed": {
		{Sensor: "windspeedmph", Color: "orange", Title: "Wind Speed"},
		{Sensor: "windgustmph", Color: "red", Title: "Wind Gust"},
	},
	"baromrelin": {
		{Sensor: "baromrelin", Color: "purple", Title: "Barometric Pressure"},
	},
	"lightning": {
		{Sensor: "lightningday", Color: "yellow", Title: "Lightning"},
	},
	"rain": {
		{Sensor: "dailyrainin", Color: "blue", Title: "Daily Rain"},
	},
}

func chart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sensor := cleanString(vars["sensor"])
	t := cleanString(vars["time"])

	sensors, ok := chartConfig[sensor]
	if !ok {
		writeJSONError(w, http.StatusNotFound, "unknown sensor")
		return
	}

	data := make([]ChartData, 0, len(sensors))
	for _, s := range sensors {
		data = append(data, ChartData{
			Values: chartData(t, s.Sensor),
			Key:    s.Title,
			Color:  s.Color,
		})
	}

	writeJSON(w, http.StatusOK, data)
}

func chartData(timeframe string, sensor string) []ChartValue {
	chartSQL, ok := chartQuery(timeframe, sensor)
	if !ok {
		logger.WithField("sensor", sensor).Warn("chartData: sensor not in allowlist")
		return []ChartValue{}
	}
	logger.Debug(chartSQL)

	rows, err := db.Query(chartSQL)
	if err != nil {
		logger.WithError(err).Error("chartData query")
		return []ChartValue{}
	}
	defer rows.Close()

	chartValues := make([]ChartValue, 0)
	for rows.Next() {
		a := ChartValue{}
		if sensor == "dailyrainin" || sensor == "lightningday" {
			var ts int
			var value float64
			if err := rows.Scan(&ts, &value); err != nil {
				logger.WithError(err).Error("chartData scan")
				continue
			}
			a = ChartValue{
				Ts:    time.Date(time.Now().Year(), time.Now().Month(), ts, 0, 0, 0, 0, time.Local),
				Value: value,
			}
		} else {
			if err := rows.Scan(&a.Ts, &a.Value); err != nil {
				logger.WithError(err).Error("chartData scan")
				continue
			}
		}
		chartValues = append(chartValues, a)
	}
	if err := rows.Err(); err != nil {
		logger.WithError(err).Error("chartData rows.Err")
	}
	return chartValues
}

func alltime(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	col := cleanString(vars["type"])
	calc := cleanString(vars["calc"])

	if !allowedAlltimeColumns[col] {
		writeJSONError(w, http.StatusBadRequest, "unknown column")
		return
	}
	if !allowedAlltimeCalcs[calc] {
		writeJSONError(w, http.StatusBadRequest, "calc must be max or min")
		return
	}

	order := ""
	if calc == "max" {
		order = " desc"
	}

	type Result struct {
		Date  string  `json:"date"`
		Value float64 `json:"value"`
	}
	rt := Result{}

	// Both identifiers are now allowlisted, so this interpolation is
	// safe.
	query := fmt.Sprintf("select %s as value, recorded from records order by %s%s limit 1", col, col, order)
	logger.Info(query)

	err := db.QueryRow(query).Scan(&rt.Value, &rt.Date)
	if err != nil && err != sql.ErrNoRows {
		sqlError("alltime", err, query)
	}

	writeJSON(w, http.StatusOK, rt)
}

func wind(timeFrame string) float64 {
	avgSpeedQ := fmt.Sprintf("SELECT ROUND(AVG(windspeedmph)::numeric, 2) AS value FROM records WHERE recorded >= NOW() - INTERVAL '%s'", timeFrame)

	var avg sql.NullFloat64
	if err := db.QueryRow(avgSpeedQ).Scan(&avg); err != nil && err != sql.ErrNoRows {
		sqlError("wind", err, avgSpeedQ)
	}

	return avg.Float64
}

var minmaxRegexCache sync.Map // map[string]*regexp.Regexp

func minmaxRegex(f string) *regexp.Regexp {
	if v, ok := minmaxRegexCache.Load(f); ok {
		return v.(*regexp.Regexp)
	}
	re := regexp.MustCompile(f + "$")
	minmaxRegexCache.Store(f, re)
	return re
}

//func minmax(f string) map[string]map[string]StatValue {
//	s := getStats() // reads from cache
//	result := make(map[string]map[string]StatValue)
//	re := minmaxRegex(f)
//
//	for _, v := range s {
//		if !re.MatchString(v.ID) {
//			continue
//		}
//		parts := strings.Split(v.ID, "_")
//		if len(parts) < 2 {
//			continue
//		}
//		if _, ok := result[parts[1]]; !ok {
//			result[parts[1]] = make(map[string]StatValue)
//		}
//		result[parts[1]][strings.ToLower(parts[0])] = StatValue{Recorded: v.Recorded, Value: v.Value}
//	}
//	return result
//}

// ----------------------------------------------------------------------------
// Loaders
// ----------------------------------------------------------------------------

// loadAllStats fetches every period_stats row in a single query and groups
// them by field. Call this once per request and look up each field as you
// build the response payload.
//
// Returns a map keyed by field name (e.g. "tempf", "humidity"). Fields with
// no rows in period_stats simply won't appear; callers should treat a
// missing key the same as an empty FieldStats.
func loadAllStats() (map[string]FieldStats, error) {
	const q = `
		SELECT period, field,
		       min_value, min_at,
		       max_value, max_at,
		       sum_value, sample_count, updated_at
		FROM period_stats`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make(map[string]FieldStats)

	for rows.Next() {
		var (
			period      string
			field       string
			minValue    sql.NullFloat64
			minAt       sql.NullTime
			maxValue    sql.NullFloat64
			maxAt       sql.NullTime
			sumValue    sql.NullFloat64
			sampleCount int64
			updatedAt   time.Time
		)
		if err := rows.Scan(&period, &field,
			&minValue, &minAt, &maxValue, &maxAt,
			&sumValue, &sampleCount, &updatedAt); err != nil {
			return nil, err
		}

		fs := out[field]

		if minValue.Valid {
			sv := &StatValue{
				Value:    minValue.Float64,
				Recorded: nullTimeOr(minAt, updatedAt),
			}
			assignPeriod(&fs.Min, period, sv)
		}
		if maxValue.Valid {
			sv := &StatValue{
				Value:    maxValue.Float64,
				Recorded: nullTimeOr(maxAt, updatedAt),
			}
			assignPeriod(&fs.Max, period, sv)
		}
		if sumValue.Valid && sampleCount > 0 {
			if fs.Avg == nil {
				fs.Avg = &Periods{}
			}
			sv := &StatValue{
				Value:    sumValue.Float64 / float64(sampleCount),
				Recorded: updatedAt,
			}
			assignPeriod(fs.Avg, period, sv)
		}

		out[field] = fs
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// minmax is the single-field convenience wrapper. Equivalent to looking up
// loadAllStats()[f], but does a targeted query when you only need one field.
func minmax(f string) FieldStats {
	const q = `
		SELECT period,
		       min_value, min_at,
		       max_value, max_at,
		       sum_value, sample_count, updated_at
		FROM period_stats
		WHERE field = $1`

	rows, err := db.Query(q, f)
	if err != nil {
		logger.Error("minmax query: ", err)
		return FieldStats{}
	}
	defer rows.Close()

	var fs FieldStats
	for rows.Next() {
		var (
			period      string
			minValue    sql.NullFloat64
			minAt       sql.NullTime
			maxValue    sql.NullFloat64
			maxAt       sql.NullTime
			sumValue    sql.NullFloat64
			sampleCount int64
			updatedAt   time.Time
		)
		if err := rows.Scan(&period, &minValue, &minAt, &maxValue, &maxAt,
			&sumValue, &sampleCount, &updatedAt); err != nil {
			logger.Error("minmax scan: ", err)
			continue
		}

		if minValue.Valid {
			assignPeriod(&fs.Min, period, &StatValue{
				Value:    minValue.Float64,
				Recorded: nullTimeOr(minAt, updatedAt),
			})
		}
		if maxValue.Valid {
			assignPeriod(&fs.Max, period, &StatValue{
				Value:    maxValue.Float64,
				Recorded: nullTimeOr(maxAt, updatedAt),
			})
		}
		if sumValue.Valid && sampleCount > 0 {
			if fs.Avg == nil {
				fs.Avg = &Periods{}
			}
			assignPeriod(fs.Avg, period, &StatValue{
				Value:    sumValue.Float64 / float64(sampleCount),
				Recorded: updatedAt,
			})
		}
	}
	if err := rows.Err(); err != nil {
		logger.Error("minmax rows: ", err)
	}
	return fs
}

// ----------------------------------------------------------------------------
// helpers
// ----------------------------------------------------------------------------

func assignPeriod(p *Periods, period string, sv *StatValue) {
	switch period {
	case "day":
		p.Day = sv
	case "month":
		p.Month = sv
	case "year":
		p.Year = sv
	case "alltime":
		p.AllTime = sv
	}
}

func nullTimeOr(t sql.NullTime, fallback time.Time) time.Time {
	if t.Valid {
		return t.Time
	}
	return fallback
}

func buildBaseConditions(record Record) Conditions {
	c := Conditions{
		Mac:      record.Mac,
		Recorded: record.Recorded,
	}

	c.Barometer = Barometer{
		Baromabsin: record.Baromabsin,
		Baromrelin: record.Baromrelin,
		Trend:      trend("baromrelin"),
		MinMax:     minmax("baromrelin"),
	}
	c.Humidity = Humidity{
		Humidity: record.Humidity,
		Dewpoint: record.Dewpoint,
		MinMax:   minmax("humidity"),
	}
	c.Temp = Temp{
		Temp:      record.Tempf,
		Humidity:  record.Humidity,
		Feelslike: record.Feelslike,
		Dewpoint:  record.Dewpoint,
		MinMax:    minmax("tempf"),
	}
	c.Tempin = Tempin{Temp: record.Tempinf, Humidity: record.Humidityin, MinMax: minmax("tempinf")}
	c.Temp1 = Tempin{Temp: record.Temp1f, Humidity: record.Humidity1, MinMax: minmax("temp1f")}
	c.Temp2 = Tempin{Temp: record.Temp2f, Humidity: record.Humidity2, MinMax: minmax("temp2f")}
	c.Temp3 = Tempin{Temp: record.Temp3f, Humidity: record.Humidity3, MinMax: minmax("temp3f")}
	c.Temp4 = Tempin{Temp: record.Temp4f, Humidity: record.Humidity4, MinMax: minmax("temp4f")}

	c.Wind = Wind{
		Dir:        record.Winddir,
		Gustmph:    record.Windgustmph,
		Gustdir:    record.Windgustdir,
		Speedmph:   record.Windspeedmph,
		Avg:        wind("10 minutes"),
		MinMax:     minmax("windspeedmph"),
		GustMinMax: minmax("windgustmph"),
	}

	c.UV = UV{Uv: record.Uv, Solarradiation: record.Solarradiation, MinMax: minmax("uv")}
	c.Rain = Rain{
		Hourly:   record.Hourlyrainin,
		Daily:    record.Dailyrainin,
		Weekly:   record.Weeklyrainin,
		Monthly:  record.Monthlyrainin,
		Yearly:   record.Yearlyrainin,
		Total:    record.Totalrainin,
		Lastrain: record.Lastrain,
	}
	c.Alert = alerts()
	c.Lightning = Lightning{
		Day:      record.Lightningday,
		Hour:     record.Lightninghour,
		Distance: record.Lightningdistance,
		Time:     record.Lightningtime,
		Month:    record.LightningMonth,
		Minmax:   minmax("lightning"),
	}
	c.AQI = AQI{Pm25: record.Aqipm25, Pm2524h: record.Aqipm2524h, MinMax: minmax("aqipm25")}
	c.Astro = astro()
	return c
}

func buildAppConditions(record Record, forecast []ForecastDB) AppConditions {
	base := buildBaseConditions(record)

	app := AppConditions{
		ID:        base.ID,
		Mac:       base.Mac,
		Recorded:  base.Recorded,
		Barometer: base.Barometer,
		Humidity:  base.Humidity,
		Temp:      base.Temp,
		Tempin:    base.Tempin,
		Temp1:     base.Temp1,
		Temp2:     base.Temp2,
		Temp3:     base.Temp3,
		Temp4:     base.Temp4,
		Rain:      base.Rain,
		Lightning: base.Lightning,
		AQI:       base.AQI,
		Wind:      base.Wind,
		UV:        base.UV,
		Astro:     base.Astro,
		Alert:     base.Alert,
	}

	if len(forecast) > 0 {
		app.Forecast = forecast[0]
		app.Conditions = forecast[0].Conditions
	} else {
		app.Forecast = ForecastDB{}
		app.Conditions = ""
	}
	return app
}

func getConditions() Conditions {
	return buildBaseConditions(getCurrent())
}

func getAppConditions() AppConditions {
	res := getCurrent()
	forecast, err := GetForecasts()
	if err != nil {
		logger.WithError(err).Error("getAppConditions: forecast load failed")
	}
	return buildAppConditions(res, forecast)
}

func current(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, getConditions())
}

func app(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, getAppConditions())
}

func about(w http.ResponseWriter, _ *http.Request) {
	const query = `select count(*) as total, max(tempf) as maxtemp,
		max(windgustmph) as wind,
		min(tempf) as mintemp
		from records`

	r := About{}
	err := db.QueryRow(query).Scan(&r.Records, &r.Maxtemp, &r.Maxgust, &r.Mintemp)
	if err != nil && err != sql.ErrNoRows {
		logger.WithError(err).Error("about query")
		writeJSONError(w, http.StatusInternalServerError, "query failed")
		return
	}

	writeJSON(w, http.StatusOK, r)
}

func getForecastHandler(w http.ResponseWriter, _ *http.Request) {
	res, err := GetForecasts()
	if err != nil {
		logger.WithError(err).Error("getForecastHandler")
		writeJSONError(w, http.StatusInternalServerError, "forecast load failed")
		return
	}
	writeJSON(w, http.StatusOK, res)
}
