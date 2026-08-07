package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration parsed once at startup.
// Centralizes os.Getenv calls and validates required values
// rather than scattering parsing across the codebase.
type Config struct {
	Port        string
	DBUser      string
	DBPassword  string
	DBHost      string
	DBDatabase  string
	Latitude    float64
	Longitude   float64
	Location    *time.Location
	LocationStr string
	LogLevel    string
	GoEnv       string
	Station     string

	// ForecastLocationID selects which location's forecast this station shows.
	//
	// The forecast table used to hold exactly one location and every query
	// simply took whatever was there. It is now keyed on (location_id,
	// datetime) and can hold many, so an unfiltered read returns rows for
	// several places interleaved by date — ten rows that look like ten days and
	// are three days of three locations. Plausible, and wrong.
	//
	// Required rather than defaulted, because there is no safe default: the
	// forecaster does not know which location a weather station stands at, and
	// guessing would produce exactly the plausible-wrong output above. Find the
	// id on the forecaster's /locations page.
	ForecastLocationID int64
}

// LoadConfig reads environment variables and validates required values.
// Returns an error rather than panicking so main() controls exit behavior.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Port:        os.Getenv("PORT"),
		DBUser:      os.Getenv("DB_USER"),
		DBPassword:  os.Getenv("DB_PASSWORD"),
		DBHost:      os.Getenv("DB_HOST"),
		DBDatabase:  os.Getenv("DB_DATABASE"),
		LogLevel:    os.Getenv("LOGLEVEL"),
		GoEnv:       os.Getenv("GO_ENV"),
		LocationStr: getEnvDefault("TZ", "America/Denver"),
		Station:     getEnvDefault("STATION_ID", "KCOS"),
	}

	if cfg.Port == "" {
		return nil, fmt.Errorf("PORT environment variable is required")
	}
	if cfg.DBUser == "" || cfg.DBHost == "" || cfg.DBDatabase == "" {
		return nil, fmt.Errorf("DB_USER, DB_HOST, and DB_DATABASE are required")
	}

	fidStr := os.Getenv("FORECAST_LOCATION_ID")
	if fidStr == "" {
		return nil, fmt.Errorf(
			"FORECAST_LOCATION_ID is required: the forecast table is keyed on location_id " +
				"and holds one row per location per day, so a query without it mixes locations. " +
				"Find the id on the forecaster's /locations page")
	}
	fid, err := strconv.ParseInt(fidStr, 10, 64)
	if err != nil || fid <= 0 {
		return nil, fmt.Errorf("invalid FORECAST_LOCATION_ID %q: want a positive integer", fidStr)
	}
	cfg.ForecastLocationID = fid

	latStr := os.Getenv("LAT")
	lonStr := os.Getenv("LON")
	if latStr == "" || lonStr == "" {
		return nil, fmt.Errorf("LAT and LON environment variables are required")
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid LAT value %q: %w", latStr, err)
	}
	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid LON value %q: %w", lonStr, err)
	}
	cfg.Latitude = lat
	cfg.Longitude = lon

	loc, err := time.LoadLocation(cfg.LocationStr)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone %q: %w", cfg.LocationStr, err)
	}
	cfg.Location = loc

	return cfg, nil
}

// DSN returns the Postgres connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"user=%s password=%s host=%s port=5432 dbname=%s sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBDatabase,
	)
}

func getEnvDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

// config is the package-wide configuration, populated at startup.
var config *Config
