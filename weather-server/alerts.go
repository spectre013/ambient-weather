package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

var LastModified time.Time

func updateAlerts() {
	fmt.Println("Updating Alerts ... ")
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
				logger.Error(fmt.Sprintf("Scan: %v", err))
			}
		}
		fmt.Println(v.Onset, v.Expires)
		if id == -1 {
			_, err := db.Exec(insertSql, v.IDURI, v.Type, v.AreaDesc, v.Sent, v.Effective, v.Onset,
				v.Expires, v.Ends, v.Status,
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
