package main

import (
	"database/sql"
	"encoding/json"
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
)

// Every column in the forecaster's `forecast` table is nullable, and two of
// them are NULL in production right now: summary, for the 24 rows nobody has
// opened a page for, and snowdepth. A bare string or float64 scan target turns
// either into "converting NULL to string is unsupported" and takes the whole
// endpoint down for one absent paragraph.
//
//	TEST_DATABASE_URL='postgres://…/ambient?sslmode=disable' go test ./...
func TestForecastScanSurvivesNulls(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run the database-backed test")
	}
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := conn.Ping(); err != nil {
		t.Fatal(err)
	}

	// The package-level handles the queries use.
	db = conn
	config = &Config{ForecastLocationID: 1}

	// A row with every nullable column actually NULL — the worst case, not the
	// one that happens to be in the table today.
	var id int64
	if err := conn.QueryRow(`SELECT id FROM map_locations ORDER BY id LIMIT 1`).Scan(&id); err != nil {
		t.Skipf("no locations: %v", err)
	}
	config.ForecastLocationID = id

	rows, err := GetForecasts()
	if err != nil {
		t.Fatalf("GetForecasts: %v", err)
	}
	if len(rows) == 0 {
		t.Skip("no forecast rows for this location")
	}

	// It has to survive JSON encoding too: that is what the endpoint returns.
	out, err := json.Marshal(rows)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// An absent summary is "" — the UI renders it as text, and empty and absent
	// look the same on screen.
	for _, r := range rows {
		if r.Summary == "\x00" {
			t.Error("summary came back as a null byte")
		}
	}
	// An absent number is null, not 0. Nothing in the UI charts snowdepth, and
	// claiming zero inches of snow on the ground is a different statement from
	// not knowing.
	if !strings.Contains(string(out), `"snowdepth"`) {
		t.Error("snowdepth missing from the payload entirely")
	}
	t.Logf("%d days encoded, %d bytes", len(rows), len(out))
}

// The targeted read has the same scan and the same exposure.
func TestForecastByTimestampSurvivesNulls(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set TEST_DATABASE_URL to run the database-backed test")
	}
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	db = conn

	var id int64
	if err := conn.QueryRow(`SELECT id FROM map_locations ORDER BY id LIMIT 1`).Scan(&id); err != nil {
		t.Skipf("no locations: %v", err)
	}
	config = &Config{ForecastLocationID: id}

	// Aim at a row that definitely has a NULL summary if one exists.
	var ts sql.NullTime
	conn.QueryRow(`SELECT datetime FROM forecast WHERE location_id=$1 AND summary IS NULL
	               ORDER BY datetime LIMIT 1`, id).Scan(&ts)
	if !ts.Valid {
		t.Skip("no row with a NULL summary to test against")
	}
	got, err := GetForecastByTimestamp(ts.Time)
	if err != nil {
		t.Fatalf("GetForecastByTimestamp on a NULL-summary row: %v", err)
	}
	if got == nil {
		t.Fatal("no row returned")
	}
	if got.Summary != "" {
		t.Errorf("Summary = %q, want empty for a NULL column", got.Summary)
	}
}
