package main

import (
	"testing"
)

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestHeatIndex(t *testing.T) {
	hi := heatIndex(85, 45)
	want := 85.31
	t.Log(hi, want)
	if hi != want {
		t.Errorf("got %.2f, wanted %.2f", hi, want)
	}
}

func TestHeatIndex_lowh(t *testing.T) {
	hi := heatIndex(85, 27)
	want := 82.54
	t.Log(hi, want)
	if hi != want {
		t.Errorf("got %.2f, wanted %.2f", hi, want)
	}

}

func TestHeatIndex_highh(t *testing.T) {
	hi := heatIndex(87, 87)
	want := 107.22
	t.Log(hi, want)
	if hi != want {
		t.Errorf("got %.2f, wanted %.2f", hi, want)
	}

}

func Test_WindChill(t *testing.T) {
	wc := windChill(20, 10)
	want := 8.85
	t.Log(wc, want)
	if wc != want {
		t.Errorf("got %.2f, wanted %.2f", wc, want)
	}
}

//func Test_TimeFrame_yesterday(t *testing.T) {
//	now := time.Now()
//	tf := getTimeframe("yesterday")
//	timeStart := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, loc)
//	timeEnd := time.Date(now.Year(), now.Month(), now.Day()-1, 23, 59, 59, 0, loc)
//
//	if tf[0] != timeStart {
//		t.Errorf("got %v, wanted %v", tf[0], timeStart)
//	}
//	if tf[1] != timeEnd {
//		t.Errorf("got %v, wanted %v", tf[1], timeEnd)
//	}
//}
