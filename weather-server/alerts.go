package main

//
//import (
//	"encoding/json"
//	"errors"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//	"time"
//)
//
//var LastModified time.Time
//
//func parseJson(j []byte) {
//	alerts := Alerts{}
//	err := json.Unmarshal(j, &alerts)
//	if err != nil {
//		logger.Error(err)
//	}
//	for k, v := range alerts.Features {
//		fmt.Println(k, v.Properties.Event)
//	}
//}
//
//func startAlerts() {
//	logger.Info("Running Alert System")
//	alertCheck()
//	for range time.NewTicker(1 * time.Minute).C {
//		alertCheck()
//	}
//}
//
//func alertCheck() bool {
//	uri := "https://api.weather.gov/alerts/active?area=CO"
//
//	if !LastModified.IsZero() {
//		_, err := alertRequest("HEAD", uri)
//		if err != nil {
//			logger.Error(err)
//			return false
//		}
//	}
//
//	res, err := alertRequest("GET", uri)
//	if err != nil {
//		logger.Error()
//		return false
//	}
//	parseJson(res)
//	return true
//}
//
//func alertRequest(t string, url string) (body []byte, err error) {
//	body = []byte("")
//	client := &http.Client{}
//	req, err := http.NewRequest(t, url, nil)
//	req.Header.Add("User-Agent", `Zoms Weather, wxcos@zoms.net`)
//
//	resp, err := client.Do(req)
//	if err != nil || resp.StatusCode != 200 {
//		err = errors.New("server responded with an error")
//		return body, err
//	}
//	if _, ok := resp.Header["Last-Modified"]; ok {
//		lModified := ""
//		if len(resp.Header["Last-Modified"]) > 0 {
//			lModified = resp.Header["Last-Modified"][0]
//		}
//		l, err := time.Parse("Mon, 2 Jan 2006 15:04:05 GMT", lModified)
//		if err != nil {
//			return []byte(""), err
//		}
//		logger.Info(l.Before(LastModified), l, LastModified)
//		if (l.Before(LastModified) || l.Equal(LastModified)) && !LastModified.IsZero() {
//			logger.Debug("No Updates")
//			logger.Debug(l, LastModified)
//			err = errors.New("cached")
//			return body, err
//		}
//	}
//	logger.Debug("Alert Updates")
//	if t != "head" {
//		body, err = ioutil.ReadAll(resp.Body)
//		if err != nil {
//			logger.Error(err)
//			return
//		}
//		l, _ := time.Parse("Mon, 2 Jan 2006 15:04:05 GMT", resp.Header["Last-Modified"][0])
//		LastModified = l
//	}
//
//	return
//}
