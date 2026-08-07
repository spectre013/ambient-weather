package main

import "testing"

func TestForecastLocationIDIsRequired(t *testing.T) {
	base := map[string]string{
		"PORT": "6000", "DB_USER": "u", "DB_HOST": "h", "DB_DATABASE": "d",
		"LAT": "38.7", "LON": "-104.6",
	}
	for k, v := range base {
		t.Setenv(k, v)
	}
	t.Setenv("FORECAST_LOCATION_ID", "")
	if _, err := LoadConfig(); err == nil {
		t.Fatal("LoadConfig accepted a missing FORECAST_LOCATION_ID")
	} else {
		t.Logf("refused: %v", err)
	}

	t.Setenv("FORECAST_LOCATION_ID", "nope")
	if _, err := LoadConfig(); err == nil {
		t.Error("LoadConfig accepted a non-numeric FORECAST_LOCATION_ID")
	}
	t.Setenv("FORECAST_LOCATION_ID", "0")
	if _, err := LoadConfig(); err == nil {
		t.Error("LoadConfig accepted location id 0")
	}
	t.Setenv("FORECAST_LOCATION_ID", "3")
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.ForecastLocationID != 3 {
		t.Errorf("ForecastLocationID = %d, want 3", cfg.ForecastLocationID)
	}
}
