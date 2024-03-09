package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var topics = make(map[string]Topic)

func (w Weather) Notify() {
	r := w.getWeatherData()

	for k, v := range topics {
		t := v
		switch t.Name {
		case "tempf":
			if t.Value > r.Tempf {
				v.Last = time.Now()
				topics[k] = v
				SendNotification(fmt.Sprintf(t.Message, t.Value, formatDate(v.Last)))
			}
		}
	}
}

func SendNotification(message string) {
	ntfyRequest("https://ntfy.zoms.net/weather", message)
}
func buildTopics() {
	topics["temp"] = Topic{
		Name:     "tempf",
		Last:     time.Now().Add(-6 * time.Hour),
		Value:    25.0,
		Duration: time.Duration(5 * time.Hour),
		Message:  "Temperature below %.2f degrees at %s",
	}

}

func (w Weather) getWeatherData() Record {
	query := `select id,mac,recorded,baromabsin,baromrelin,battout,Batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,
				maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,battlightning,lightningday,lightninghour,lightningtime,lightningdistance,  aqipm25, aqipm2524h
				from records order by recorded desc limit 1`
	rec := getRecord(w.DB, query)
	return rec
}

func ntfyRequest(url string, message string) {
	logger.Debug(url)
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(message))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		logger.Error(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(err)
	}
	logger.Debug(resp)
}

type Topic struct {
	Name     string
	Last     time.Time
	Value    float64
	Duration time.Duration
	Message  string
}
