package main

import (
	_ "database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (w *Weather) index(wr http.ResponseWriter, r *http.Request) {
	forecast, err := getForecast()
	if err != nil {
		logger.Error(err)
	}
	w.Forecast = forecast
	_ = Index().Render(r.Context(), wr)
}
func (w Weather) current(wr http.ResponseWriter, r *http.Request) {
	props, res, err := w.getCurrent()
	if err != nil {
		log.Println(err)
	}

	_ = Main(props, res).Render(r.Context(), wr)
}

func (w Weather) temperature(wr http.ResponseWriter, r *http.Request) {
	_ = Almanac("temp").Render(r.Context(), wr)
}
func (w Weather) forecast(wr http.ResponseWriter, r *http.Request) {

	_ = forecastDetail(w.Forecast, units).Render(r.Context(), wr)
}

func (w Weather) temp(wr http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t := vars["time"]
	a := w.Almanac(t)
	chart := w.Chart(t)
	_ = tempAlmanac(a, chart, t).Render(r.Context(), wr)
}

func (w Weather) Alert(wr http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	alertsSql := fmt.Sprintf("select * from alerts where id = '%s'", id)

	a := Alert{}
	row := w.DB.QueryRow(alertsSql)
	if err := row.Scan(&a.ID, &a.Alertid, &a.Wxtype, &a.Areadesc,
		&a.Sent, &a.Effective, &a.Onset, &a.Expires,
		&a.Ends, &a.Status, &a.Messagetype, &a.Category,
		&a.Severity, &a.Certainty, &a.Urgency, &a.Event,
		&a.Sender, &a.SenderName, &a.Headline, &a.Description,
		&a.Instruction, &a.Response); err != nil {
		sqlError(err, alertsSql)
	}
	a.Sent = a.Sent.In(loc)
	a.Effective = a.Effective.In(loc)
	a.Onset = a.Onset.In(loc)
	a.Ends = a.Ends.In(loc)

	e := AlertDetail(a, w.Alerts()).Render(r.Context(), wr)
	if e != nil {
		logger.Error(e)
	}
}

func (w Weather) Wind(wr http.ResponseWriter, r *http.Request) {

	_ = windDetail().Render(r.Context(), wr)
}

func (w Weather) Alerts() []Alert {
	now := time.Now().Local()

	alertsSql := fmt.Sprintf("select * from alerts where ends >= '%s' order by ends desc", formatDate(now))

	rows, err := w.DB.Query(alertsSql)
	if err != nil {
		logger.Error(err)
	}
	alerts := make([]Alert, 0)
	for rows.Next() {
		a := Alert{}
		err := rows.Scan(&a.ID, &a.Alertid, &a.Wxtype, &a.Areadesc, &a.Sent, &a.Effective, &a.Onset, &a.Expires, &a.Ends, &a.Status, &a.Messagetype, &a.Category, &a.Severity, &a.Certainty, &a.Urgency, &a.Event, &a.Sender, &a.SenderName, &a.Headline, &a.Description, &a.Instruction, &a.Response)
		sqlError(err, alertsSql)

		a.Sent = a.Sent.In(loc)
		a.Effective = a.Effective.In(loc)
		a.Onset = a.Onset.In(loc)
		a.Ends = a.Ends.In(loc)

		alerts = append(alerts, a)
	}
	return alerts
}

func (w Weather) Alertview(wr http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ids := vars["id"]
	id, err := strconv.Atoi(ids)
	errorHandler(err, "Error converting id to int")
	alerts := w.Alerts()
	alert := Alert{}
	prev := 0
	next := 0

	if id == 0 && len(alerts) == 1 {
		alert = alerts[0]
	} else if id == 0 && len(alerts) > 1 {
		alert = alerts[0]
		next = alerts[1].ID
	} else if id != 0 && len(alerts) > 0 {
		for i, a := range alerts {
			if a.ID == id {
				if i != 0 {
					prev = alerts[i-1].ID
				}
				alert = a
				if i != len(alerts)-1 {
					next = alerts[i+1].ID
				}
			}
		}
	}

	e := ShowAlert(alert, prev, next).Render(r.Context(), wr)
	errorHandler(e, "Error rendering alert")
}
func buildSQL(t string, interval string) string {
	where := fmt.Sprintf(" where recorded >= NOW() - interval '%s' AND recorded <= NOW() ", interval)
	if t == "at" {
		where = " "
	}

	parts := map[string]string{"tempmax": "tempf", "tempmin": "tempf", "feelmax": "feelslike", "feelmin": "feelslike",
		"humiditymax": "humidity", "humiditymin": "humidity", "baromax": "Baromrelin", "baromin": "Baromrelin",
		"dewmax": "dewpoint", "dewmin": "dewpoint", "windmax": "windspeedmph", "gustmax": "windgustmph"}
	base := "(select '%s' as label, CAST(COALESCE(%s,0.0) AS decimal(10,2)) as value, recorded from records%sorder by %s %s limit 1)"
	sqlparts := make([]string, 0)
	for k, v := range parts {
		order := "desc"
		if strings.Contains(k, "min") {
			order = "asc"
		}
		sqlparts = append(sqlparts, fmt.Sprintf(base, k, v, where, v, order))
	}
	return strings.Join(sqlparts, " union all ")
}
func (w Weather) Almanac(t string) AlmanacInfo {
	p := getPill(t)
	almSQL := buildSQL(t, p.Interval)

	rows, err := w.DB.Query(almSQL)
	if err != nil {
		logger.Error(err)
	}
	a := AlmanacInfo{}
	for rows.Next() {
		r := AlmanacData{}
		err = rows.Scan(&r.Label, &r.Value, &r.Recorded)
		if err != nil {
			logger.Error("Scan:", err)
		}

		switch r.Label {
		case "tempmax":
			a.TempMax = r
		case "tempmin":
			a.TempMin = r
		case "feelmax":
			a.FeelMax = r
		case "feelmin":
			a.FeelMin = r
		case "humiditymax":
			a.HumidityMax = r
		case "humiditymin":
			a.HumidityMin = r
		case "baromax":
			a.BaroMax = r
		case "baromin":
			a.BaroMin = r
		case "dewmax":
			a.DewpointMax = r
		case "dewmin":
			a.DewpointMin = r
		}
	}
	return a
}

type ChartValue struct {
	Ts  time.Time
	Max float64
	Min float64
}

func (w Weather) Chart(t string) string {
	chart := BuildChart()
	chartSQL := chartQueries(t)

	rows, err := w.DB.Query(chartSQL)
	if err != nil {
		logger.Error(err)
	}

	ChartValues := make([]ChartValue, 0)
	for rows.Next() {
		a := ChartValue{}
		err := rows.Scan(&a.Ts, &a.Max, &a.Min)
		sqlError(err, chartSQL)
		ChartValues = append(ChartValues, a)
	}
	seriesMax := Series{
		Name: "Max Temperature",
		Data: make([]float64, 0),
	}
	seriesMin := Series{
		Name: "Min Temperature",
		Data: make([]float64, 0),
	}
	xcat := make([]string, 0)
	for _, v := range ChartValues {
		seriesMax.Data = append(seriesMax.Data, v.Max)
		seriesMin.Data = append(seriesMin.Data, v.Min)
		tformat := v.Ts.Format("15:04")
		if t == "1m" || t == "1y" || t == "at" {
			tformat = v.Ts.Format("Jan 02, 06")
		}
		xcat = append(xcat, tformat)
	}
	chart.Series = append(chart.Series, seriesMax)
	chart.Series = append(chart.Series, seriesMin)
	chart.Xaxis.Categories = xcat
	b, err := json.Marshal(chart)
	if err != nil {
		log.Println(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

func BuildChart() ChartData {
	chart := ChartData{}
	chart.Legend = Legend{
		Show:            true,
		Position:        "top",
		HorizontalAlign: "left",
		Labels: legendLabels{
			Colors:          []string{},
			UseSeriesColors: true,
		},
	}

	chart.Colors = []string{"#ff7c39cc", "#3b9cac"}
	chart.Chart = ChartDef{
		FontFamily: "Satoshi, sans-serif",
		Height:     600,
		Width:      1200,
		Type:       "area",
		DropShadow: DropShadow{
			Enabled: true,
			Color:   "#623CEA14",
			Top:     10,
			Blur:    4,
			Left:    0,
			Opacity: 0.1,
		},
		Toolbar: Toolbar{Show: false},
	}

	chart.Responsive = make([]Responsive, 0)
	optSM := Options{ChartOptions{Height: 300}}
	optMD := Options{ChartOptions{Height: 350}}
	chart.Responsive = append(chart.Responsive, Responsive{
		Breakpoint: 1024,
		Options:    optSM,
	})
	chart.Responsive = append(chart.Responsive, Responsive{
		Breakpoint: 1366,
		Options:    optMD,
	})
	chart.Stroke = Stroke{
		Width: []int{2, 2},
		Curve: "straight",
	}
	chart.Labels = Labels{
		Show:     false,
		Position: "top",
	}
	chart.Grid = Grid{
		Xaxis: Gridaxis{Lines: Lines{Show: true}},
		Yaxis: Gridaxis{Lines: Lines{Show: true}},
		Row: GridRowColumn{
			Colors:  "#C0C0C0",
			Opacity: 0.5,
		},
		Column: GridRowColumn{
			Colors:  "#C0C0C0",
			Opacity: 0.5,
		},
	}
	chart.DataLabels = DataLabels{Enabled: false}
	chart.Markers = Markers{
		Size:            4,
		Colors:          "#fff",
		StrokeColors:    []string{"#3056D3", "#80CAEE"},
		StrokeWidth:     3,
		StrokeOpacity:   0.9,
		StrokeDashArray: 0,
		FillOpacity:     1,
		Discrete:        []interface{}{},
		Hover:           Hover{SizeOffset: 5, Size: 0},
	}

	chart.Xaxis = Xaxis{
		Type:       "category",
		Categories: createAxis(1, 60),
		AxisBorder: AxisBorder{Show: false},
		AxisTicks:  AxisTicks{Show: false},
		Labels: AxisLabel{
			Show:  true,
			Align: "right",
			Style: LabelStyle{
				Colors:   "#fff",
				FontSize: "12px",
			},
		},
	}

	chart.Yaxis = Yaxis{
		Title: Title{Style: Style{FontSize: "20px"}},
		Min:   -20,
		Max:   100,
		Labels: AxisLabel{
			Show:  true,
			Align: "right",
			Style: LabelStyle{
				Colors:   "#fff",
				FontSize: "12px",
			},
		},
	}

	return chart
}

func createAxis(start int, end int) []string {
	var result []string
	for i := start; i <= end; i++ {
		s := strconv.Itoa(i)
		if i < 10 {
			s = "0" + s
		}
		result = append(result, s)
	}

	return result
}

func (w Weather) trend(t string) Trend {
	sel := fmt.Sprintf("AVG(%s)", t)

	start := time.Now()
	end := start.Add(-30 * time.Minute)
	var avg float64
	cr := Record{}

	avgQuery := fmt.Sprintf("select %s from records where recorded BETWEEN '%s' AND '%s'", sel, formatDate(end), formatDate(start))
	logger.Debug(avgQuery)
	rows := w.DB.QueryRow(avgQuery)
	err := rows.Scan(&avg)
	sqlError(err, avgQuery)

	currentQuery := "select id,baromrelin,tempf from records order by recorded desc limit 1"
	logger.Debug(currentQuery)
	crows := w.DB.QueryRow(currentQuery)
	err = crows.Scan(&cr.ID, &cr.Baromrelin, &cr.Tempf)
	sqlError(err, currentQuery)

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
func (w Weather) getWind() map[string]StatValue {

	dates := getTimeframe("day")
	start := dates[1]
	end := dates[0]

	maxwind := StatValue{}
	maxgust := StatValue{}
	var avg float64
	var avgdir float64

	maxSpeed := fmt.Sprintf("select windspeedmph as value, recorded from records where recorded BETWEEN '%s' AND '%s' order by windspeedmph desc limit 1", formatDate(end), formatDate(start))
	maxGust := fmt.Sprintf("select windgustmph as value, recorded from records where recorded BETWEEN '%s' AND '%s' order by windgustmph desc limit 1", formatDate(end), formatDate(start))
	avgSpeed := fmt.Sprintf("select AVG(windspeedmph) as value from records where recorded BETWEEN '%s' AND '%s'", formatDate(end), formatDate(start))
	avgDir := fmt.Sprintf("select AVG(winddir) as value from records where recorded BETWEEN '%s' AND '%s'", formatDate(end), formatDate(start))

	logger.Debug(maxSpeed)
	mrows := w.DB.QueryRow(maxSpeed)
	err := mrows.Scan(&maxwind.Value, &maxwind.Recorded)
	sqlError(err, maxSpeed)

	logger.Debug(maxGust)
	mgrows := w.DB.QueryRow(maxGust)
	err = mgrows.Scan(&maxgust.Value, &maxgust.Recorded)
	sqlError(err, maxGust)

	logger.Debug(avgSpeed)
	asrows := w.DB.QueryRow(avgSpeed)
	err = asrows.Scan(&avg)
	sqlError(err, avgSpeed)

	logger.Debug(avgDir)
	crows := w.DB.QueryRow(avgDir)
	err = crows.Scan(&avgdir)
	sqlError(err, avgDir)

	res := map[string]StatValue{}
	res["dir"] = StatValue{Value: avgdir}
	res["wind"] = maxwind
	res["gust"] = maxgust
	res["avg"] = StatValue{Value: avg}

	return res

}
