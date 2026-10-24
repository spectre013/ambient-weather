package main

import "testing"

func TestAQIFromPM25(t *testing.T) {
	tests := []struct {
		pm   float64
		want int
	}{
		{0, 0},
		{9.0, 38},
		{12.0, 50},
		{12.05, 50}, // between segments: truncated to 12.0
		{12.09, 50},
		{12.1, 51},
		{35.4, 100},
		{35.45, 100},
		{35.5, 101},
		{55.45, 150},
		{150.45, 200},
		{250.45, 300},
		{350.45, 400},
		{500.4, 500},
		{600, 500},
		{-1, 0},
	}
	for _, tt := range tests {
		if got := aqiFromPM25(tt.pm); got != tt.want {
			t.Errorf("aqiFromPM25(%v) = %d, want %d", tt.pm, got, tt.want)
		}
	}
}

// Every 0.01 µg/m³ step from 0 to 500 must land in a segment.
func TestAQIFromPM25HasNoGaps(t *testing.T) {
	for i := 1; i <= 50000; i++ {
		pm := float64(i) / 100
		if aqiFromPM25(pm) == 0 && pm >= 1 {
			t.Fatalf("aqiFromPM25(%v) = 0", pm)
		}
	}
}

func TestHeatIndexF(t *testing.T) {
	tests := []struct {
		t, rh, want float64
	}{
		{79, 90, 79},     // below 80F: temperature
		{82, 20, 79.9},   // simple formula averages < 80
		{85, 60, 89.3},   // Rothfusz, no adjustment
		{95, 10, 89.4},   // low-humidity adjustment
		{84, 90, 98.3},   // high-humidity adjustment
		{100, 40, 109.3}, // Rothfusz
	}
	for _, tt := range tests {
		if got := round1(heatIndexF(tt.t, tt.rh)); got != tt.want {
			t.Errorf("heatIndexF(%v, %v) = %v, want %v", tt.t, tt.rh, got, tt.want)
		}
	}
}

func TestRoundNegative(t *testing.T) {
	tests := []struct{ in, want float64 }{
		{-5.0, -5.0}, {-5.08, -5.1}, {-5.04, -5.0}, {5.05, 5.1}, {-0.04, 0},
	}
	for _, tt := range tests {
		if got := round1(tt.in); got != tt.want {
			t.Errorf("round1(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
	if got := round3(-0.0125); got != -0.013 {
		t.Errorf("round3(-0.0125) = %v, want -0.013", got)
	}
}
