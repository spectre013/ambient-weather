package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
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

func TestMeanDirection(t *testing.T) {
	rad := func(d float64) float64 { return d * math.Pi / 180 }
	vec := func(dirs ...float64) (float64, float64) {
		var s, c float64
		for _, d := range dirs {
			s += math.Sin(rad(d))
			c += math.Cos(rad(d))
		}
		return s, c
	}
	tests := []struct {
		name string
		dirs []float64
		want int
	}{
		{"across north", []float64{350, 10}, 0},
		{"east", []float64{80, 100}, 90},
		{"south-west", []float64{200, 250}, 225},
		{"north-west", []float64{300, 340}, 320},
	}
	for _, tt := range tests {
		s, c := vec(tt.dirs...)
		if got := meanDirection(s, c, -1); got != tt.want {
			t.Errorf("%s: meanDirection = %d, want %d", tt.name, got, tt.want)
		}
	}
	if got := meanDirection(0, 0, 123); got != 123 {
		t.Errorf("calm: meanDirection = %d, want fallback 123", got)
	}
}

func TestWallClockIn(t *testing.T) {
	loc, err := time.LoadLocation("America/Denver")
	if err != nil {
		t.Skip(err)
	}
	// 14:00 stored without a zone comes back from lib/pq labelled UTC.
	got := wallClockIn(time.Date(2026, 7, 4, 14, 0, 0, 0, time.UTC), loc)
	if got.Hour() != 14 || got.Location() != loc || got.UTC().Hour() != 20 {
		t.Errorf("wallClockIn = %v, want 14:00 MDT", got)
	}
	if !wallClockIn(time.Time{}, loc).IsZero() {
		t.Error("zero time should stay zero")
	}
}
