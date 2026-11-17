package main

import "testing"

func TestSensorFilter(t *testing.T) {
	all, err := parseSensors("")
	if err != nil || !all.Allow(&Message{Model: "Fineoffset-WS80", ID: 1}) {
		t.Fatal("an empty list must accept everything")
	}

	f, err := parseSensors("Fineoffset-WS80/34306, AmbientWeather-WH31E/156\nAmbientWeather-WH31E/130")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		model string
		id    int
		want  bool
	}{
		{"Fineoffset-WS80", 34306, true},
		{"Fineoffset-WS80", 34307, false},
		{"AmbientWeather-WH31E", 130, true},
		{"AmbientWeather-WH31E", 34306, false}, // right ID, wrong model
		{"Fineoffset-WH32B", 156, false},
	}
	for _, tc := range tests {
		if got := f.Allow(&Message{Model: tc.model, ID: tc.id}); got != tc.want {
			t.Errorf("%s/%d: allowed = %v, want %v", tc.model, tc.id, got, tc.want)
		}
	}

	for _, bad := range []string{"34306", "Fineoffset-WS80/", "/12", "Fineoffset-WS80/abc"} {
		if _, err := parseSensors(bad); err == nil {
			t.Errorf("%q should be rejected", bad)
		}
	}
}
