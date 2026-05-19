package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Store wraps the pgx connection pool. The schema (table `records`, indexes)
// is managed externally; this package only reads and writes.
type Store struct {
	pool *pgxpool.Pool
}

// NewStore opens a pgx connection pool against the supplied DSN. It does not
// create or alter any schema — callers must ensure the `records` table exists.
func NewStore(ctx context.Context, dsn string) (*Store, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}
	// Ping so we fail fast at startup rather than on first insert.
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}
	return &Store{pool: pool}, nil
}

func (s *Store) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

// Insert pushes one Record into the `records` table. Zero-valued time.Time
// fields are sent as NULL so we don't pollute the table with year-0001
// timestamps.
func (s *Store) Insert(ctx context.Context, r Record) error {
	var lastRain, lightningTime any
	if !r.Lastrain.IsZero() {
		lastRain = r.Lastrain
	}
	if !r.Lightningtime.IsZero() {
		lightningTime = r.Lightningtime
	}

	const sql = `
INSERT INTO records (
    mac, recorded,
    baromabsin, baromrelin, battout,
    batt1, batt2, batt3, batt4, batt5, batt6, batt7, batt8, batt9, batt10,
    battlightning, co2,
    dailyrainin, dewpoint, eventrainin, feelslike,
    hourlyrainin, hourlyrain, humidity,
    humidity1, humidity2, humidity3, humidity4, humidity5,
    humidity6, humidity7, humidity8, humidity9, humidity10,
    humidityin, lastrain, maxdailygust,
    monthlyrainin, solarradiation,
    tempf, temp1f, temp2f, temp3f, temp4f, temp5f, temp6f, temp7f, temp8f, temp9f, temp10f,
    tempinf, totalrainin, uv, weeklyrainin,
    winddir, windgustmph, windgustdir, windspeedmph, yearlyrainin,
    lightningday, lightninghour, lightningdistance, lightningtime, lightningmonth,
    aqipm25, aqipm2524h
) VALUES (
    $1, $2,
    $3, $4, $5,
    $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
    $16, $17,
    $18, $19, $20, $21,
    $22, $23, $24,
    $25, $26, $27, $28, $29,
    $30, $31, $32, $33, $34,
    $35, $36, $37,
    $38, $39,
    $40, $41, $42, $43, $44, $45, $46, $47, $48, $49, $50,
    $51, $52, $53, $54,
    $55, $56, $57, $58, $59,
    $60, $61, $62, $63, $64,
    $65, $66
)`
	_, err := s.pool.Exec(ctx, sql,
		r.Mac, r.Recorded,
		r.Baromabsin, r.Baromrelin, r.Battout,
		r.Batt1, r.Batt2, r.Batt3, r.Batt4, r.Batt5, r.Batt6, r.Batt7, r.Batt8, r.Batt9, r.Batt10,
		r.Battlightning, r.Co2,
		r.Dailyrainin, r.Dewpoint, r.Eventrainin, r.Feelslike,
		r.Hourlyrainin, r.Hourlyrain, r.Humidity,
		r.Humidity1, r.Humidity2, r.Humidity3, r.Humidity4, r.Humidity5,
		r.Humidity6, r.Humidity7, r.Humidity8, r.Humidity9, r.Humidity10,
		r.Humidityin, lastRain, r.Maxdailygust,
		r.Monthlyrainin, r.Solarradiation,
		r.Tempf, r.Temp1f, r.Temp2f, r.Temp3f, r.Temp4f, r.Temp5f, r.Temp6f, r.Temp7f, r.Temp8f, r.Temp9f, r.Temp10f,
		r.Tempinf, r.Totalrainin, r.Uv, r.Weeklyrainin,
		r.Winddir, r.Windgustmph, r.Windgustdir, r.Windspeedmph, r.Yearlyrainin,
		r.Lightningday, r.Lightninghour, r.Lightningdistance, lightningTime, r.LightningMonth,
		r.Aqipm25, r.Aqipm2524h,
	)
	return err
}

// LatestRecord returns the most recent row, used at startup to recover the
// rain/lightning baselines and seed the rest of the in-memory state so that
// snapshots taken before each sensor has reported still produce sensible rows.
// Returns the error from pgx (pgx.ErrNoRows if the table is empty).
func (s *Store) LatestRecord(ctx context.Context) (*Record, error) {
	const sql = `
SELECT recorded,
       baromabsin, baromrelin, battout,
       batt1, batt2, batt3, batt4, batt5, batt6, batt7, batt8, batt9, batt10,
       battlightning, co2,
       dailyrainin, dewpoint, eventrainin, feelslike,
       hourlyrainin, hourlyrain, humidity,
       humidity1, humidity2, humidity3, humidity4, humidity5,
       humidity6, humidity7, humidity8, humidity9, humidity10,
       humidityin, lastrain, maxdailygust,
       monthlyrainin, solarradiation,
       tempf, temp1f, temp2f, temp3f, temp4f, temp5f, temp6f, temp7f, temp8f, temp9f, temp10f,
       tempinf, totalrainin, uv, weeklyrainin,
       winddir, windgustmph, windgustdir, windspeedmph, yearlyrainin,
       lightningday, lightninghour, lightningdistance, lightningtime, lightningmonth,
       aqipm25, aqipm2524h
FROM records
ORDER BY recorded DESC
LIMIT 1`
	var r Record
	var lastRain, lightningTime *time.Time
	err := s.pool.QueryRow(ctx, sql).Scan(
		&r.Recorded,
		&r.Baromabsin, &r.Baromrelin, &r.Battout,
		&r.Batt1, &r.Batt2, &r.Batt3, &r.Batt4, &r.Batt5, &r.Batt6, &r.Batt7, &r.Batt8, &r.Batt9, &r.Batt10,
		&r.Battlightning, &r.Co2,
		&r.Dailyrainin, &r.Dewpoint, &r.Eventrainin, &r.Feelslike,
		&r.Hourlyrainin, &r.Hourlyrain, &r.Humidity,
		&r.Humidity1, &r.Humidity2, &r.Humidity3, &r.Humidity4, &r.Humidity5,
		&r.Humidity6, &r.Humidity7, &r.Humidity8, &r.Humidity9, &r.Humidity10,
		&r.Humidityin, &lastRain, &r.Maxdailygust,
		&r.Monthlyrainin, &r.Solarradiation,
		&r.Tempf, &r.Temp1f, &r.Temp2f, &r.Temp3f, &r.Temp4f, &r.Temp5f, &r.Temp6f, &r.Temp7f, &r.Temp8f, &r.Temp9f, &r.Temp10f,
		&r.Tempinf, &r.Totalrainin, &r.Uv, &r.Weeklyrainin,
		&r.Winddir, &r.Windgustmph, &r.Windgustdir, &r.Windspeedmph, &r.Yearlyrainin,
		&r.Lightningday, &r.Lightninghour, &r.Lightningdistance, &lightningTime, &r.LightningMonth,
		&r.Aqipm25, &r.Aqipm2524h,
	)
	if err != nil {
		return nil, err
	}
	if lastRain != nil {
		r.Lastrain = *lastRain
	}
	if lightningTime != nil {
		r.Lightningtime = *lightningTime
	}
	return &r, nil
}
