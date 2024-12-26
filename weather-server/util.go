package main

import (
	"encoding/json"
	"log"
	"math"
	"regexp"
	"time"
)

func currentToJson(record Record) ([]byte, error) {
	b, err := json.Marshal(record)
	if err != nil {
		log.Println(err)
		return []byte(""), err
	}
	return b, nil
}

func conditionsToJson(condition Conditions) ([]byte, error) {
	b, err := json.Marshal(condition)
	if err != nil {
		log.Println(err)
		return []byte(""), err
	}
	return b, nil
}

func heatIndex(T float64, humidity int) float64 {
	RH := float64(humidity)
	feelsLike := -42.379 + 2.04901523*T + 10.14333127*RH - .22475541*T*RH - .00683783*T*T - .05481717*RH*RH + .00122874*T*T*RH + .00085282*T*RH*RH - .00000199*T*T*RH*RH
	if RH < 13 && (T >= 80 && T <= 112) {
		feelsLike = feelsLike - ((13-RH)/4)*math.Sqrt((17-math.Abs(T-95.))/17)
		if RH > 85 && (T >= 80 && T <= 87) {
			feelsLike = feelsLike + ((RH-85)/10)*((87-RH)/5)
		}
	}
	return toFixed(feelsLike, 2)
}

func windChill(temperature float64, windSpeed float64) float64 {
	if windSpeed < 3 || temperature > 50 {
		return temperature
	}

	windChill := 35.74 + 0.6215*temperature - 35.75*math.Pow(windSpeed, 0.16) + 0.4275*temperature*math.Pow(windSpeed, 0.16)
	return toFixed(windChill, 2)
}

func dewpoint(temp float64, humidity int) float64 {
	tc := (temp - 32) * 5 / 9
	L := math.Log(float64(humidity) / 100)
	M := 17.27 * tc
	N := 237.3 + tc
	B := (L + (M / N)) / 17.27
	dp := (237.3 * B) / (1 - B)
	return toFixed((dp*9/5)+32, 2)
}

func getTimeframe(timeframe string) []time.Time {
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		logger.Error(err)
	}
	var dates []time.Time
	now := time.Now()

	if timeframe == "yesterday" {
		dates = append(dates, time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, loc))
		dates = append(dates, time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, loc))
	} else if timeframe == "day" {
		dates = append(dates, time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc))
		dates = append(dates, time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc))
	} else if timeframe == "month" {
		dates = append(dates, time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc))
		dates = append(dates, time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, loc))
	} else if timeframe == "year" {
		dates = append(dates, time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc))
		dates = append(dates, time.Date(now.Year(), 12, 31, 23, 59, 59, 0, loc))
	}
	return dates
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func cleanString(s string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9 -]")
	replaceStr := reg.ReplaceAllString(s, "")
	return replaceStr
}

func formatDate(date time.Time) string {
	format := "2006-01-02 15:04:05 -0700"
	return date.Format(format)
}
