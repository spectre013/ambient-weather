package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//func setupSuite(tb testing.TB) func(tb testing.TB) {
//	log.Println("setup suite")
//
//	// Return a function to teardown the test
//	return func(tb testing.TB) {
//		log.Println("teardown suite")
//	}
//}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHeatIndex(t *testing.T) {
	hi := heatIndex(85, 45)
	want := 85.31
	t.Log(hi, want)
	assert.Equal(t, hi, want)
}

func TestHeatIndex_lowh(t *testing.T) {
	hi := heatIndex(85, 27)
	want := 82.54
	t.Log(hi, want)
	assert.Equal(t, hi, want)

}

func TestHeatIndex_highh(t *testing.T) {
	hi := heatIndex(87, 87)
	want := 107.22
	t.Log(hi, want)
	assert.Equal(t, hi, want)
}

func Test_WindChill(t *testing.T) {
	wc := windChill(20, 10)
	want := 8.85
	t.Log(wc, want)
	assert.Equal(t, wc, want)
}

func Test_Dewpoint(t *testing.T) {
	dp := dewpoint(50, 30)
	want := 19.85
	t.Log(dp, want)
	assert.Equal(t, dp, want)
}

func Test_TimeFrame(t *testing.T) {
	now := time.Now()
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Println(err)
		return
	}
	ty := getTimeframe("yesterday")
	timeStart := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, loc)
	timeEnd := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, loc)
	assert.WithinDuration(t, timeStart, ty[0], 0)
	assert.WithinDuration(t, timeEnd, ty[1], 0)

	td := getTimeframe("day")
	timeDayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	timeDayEnd := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	assert.WithinDuration(t, timeDayStart, td[0], 0)
	assert.WithinDuration(t, timeDayEnd, td[1], 0)

	tm := getTimeframe("month")
	timeMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	timeMonthEnd := time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 0, loc)
	assert.WithinDuration(t, timeMonthStart, tm[0], 0)
	assert.WithinDuration(t, timeMonthEnd, tm[1], 0)

	tyr := getTimeframe("year")
	timeYearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
	timeYearEnd := time.Date(now.Year(), 12, 31, 23, 59, 59, 0, loc)
	assert.WithinDuration(t, timeYearStart, tyr[0], 0)
	assert.WithinDuration(t, timeYearEnd, tyr[1], 0)
}

func Test_Round(t *testing.T) {
	have := round(10.1)
	want := 10
	assert.Equal(t, have, want)
}

func Test_ToFixed(t *testing.T) {
	have := toFixed(10.1234, 2)
	want := 10.12
	assert.Equal(t, have, want)

}

func Test_CleanString(t *testing.T) {
	have := cleanString("This%is%wrong")
	want := "Thisiswrong"
	assert.Equal(t, have, want)
}

func Test_FormatDate(t *testing.T) {
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		fmt.Println(err)
		return
	}
	date := time.Date(2023, 2, 10, 0, 0, 0, 0, loc)
	have := formatDate(date)
	want := "2023-02-10 00:00:00 -0700"
	assert.Equal(t, have, want)
}
