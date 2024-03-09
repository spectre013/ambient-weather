package main

import (
	"math"
	"time"
)

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

func formatDate(date time.Time) string {
	format := "2006-01-02 15:04:05 -0700"
	return date.Format(format)
}
