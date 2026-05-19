package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
)

// ----------------------------------------------------------------------------
// Field configuration
// ----------------------------------------------------------------------------

// fieldKind controls how a field's stats are aggregated.
//
//	normal           -- track running min/max/sum/count; avg = sum/count.
//	daily_cumulative -- the raw column is already a running total within a day
//	                    (e.g. dailyrainin, lightningday). The "day" period
//	                    tracks the running MAX. For month/year/alltime, the
//	                    incremental math also tracks the running MAX (which is
//	                    the cumulative total over the period since the column
//	                    accumulates monotonically until the station resets).
type fieldKind int

const (
	fieldNormal fieldKind = iota
	fieldDailyCumulative
)

type fieldDef struct {
	column string // SQL column on records
	kind   fieldKind
}

// statFieldDefs is the canonical list of fields stats are computed for.
var statFieldDefs = []fieldDef{
	{"tempf", fieldNormal},
	{"tempinf", fieldNormal},
	{"temp1f", fieldNormal},
	{"temp2f", fieldNormal},
	{"temp3f", fieldNormal},
	{"temp4f", fieldNormal},
	{"baromrelin", fieldNormal},
	{"uv", fieldNormal},
	{"humidity", fieldNormal},
	{"humidityin", fieldNormal},
	{"humidity1", fieldNormal},
	{"humidity2", fieldNormal},
	{"humidity3", fieldNormal},
	{"humidity4", fieldNormal},
	{"windspeedmph", fieldNormal},
	{"windgustmph", fieldNormal},
	{"dewpoint", fieldNormal},
	{"aqipm25", fieldNormal},
	{"dailyrainin", fieldDailyCumulative},
	{"lightningday", fieldDailyCumulative},
}

const tzName = "America/Denver"

// periodTrunc returns the SQL fragment that buckets a TIMESTAMPTZ into the
// given period in local time. Returns empty string for "alltime" (never
// resets).
func periodTrunc(period, ts string) string {
	switch period {
	case "day":
		return fmt.Sprintf("date_trunc('day', %s AT TIME ZONE '%s')", ts, tzName)
	case "month":
		return fmt.Sprintf("date_trunc('month', %s AT TIME ZONE '%s')", ts, tzName)
	case "year":
		return fmt.Sprintf("date_trunc('year', %s AT TIME ZONE '%s')", ts, tzName)
	}
	return ""
}

// ----------------------------------------------------------------------------
// MQTT hot path: fold one record into all 4 periods
// ----------------------------------------------------------------------------

// foldRecord reads the record near the given timestamp and folds it into the
// day, month, year, and alltime stats for every configured field. One
// transaction per record. Idempotent across restarts via stat_ingest_log.
func foldRecord(ctx context.Context, recordedTS time.Time) error {
	rec, found, err := loadRecordNear(ctx, recordedTS)
	if err != nil {
		return fmt.Errorf("load record: %w", err)
	}
	if !found {
		return errors.New("no record found near timestamp")
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Idempotency: skip if we've already folded this exact record.
	res, err := tx.ExecContext(ctx, `
		INSERT INTO stat_ingest_log (recorded)
		VALUES ($1)
		ON CONFLICT (recorded) DO NOTHING`, rec.recorded)
	if err != nil {
		return fmt.Errorf("ingest log: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		logger.Debug("record already folded into stats, skipping ", rec.recorded)
		return tx.Commit()
	}

	for _, fd := range statFieldDefs {
		val, ok := rec.values[fd.column]
		if !ok {
			continue
		}
		// For daily_cumulative fields, month/year/alltime stats are
		// SUM-of-daily-maxes, which can't be computed incrementally from a
		// single sample without knowing all of today's prior samples. The
		// daily rebuild handles those; the hot path only touches the day row.
		periods := []string{"day", "month", "year", "alltime"}
		if fd.kind == fieldDailyCumulative {
			periods = []string{"day"}
		}
		for _, period := range periods {
			if err := upsertPeriodStat(ctx, tx, period, fd, val, rec.recorded); err != nil {
				return fmt.Errorf("upsert %s/%s: %w", period, fd.column, err)
			}
		}
	}

	return tx.Commit()
}

// upsertPeriodStat is the heart of the hot path. One INSERT...ON CONFLICT
// statement, one row touched.
//
// The boundary detection works like this: we compare period_trunc(updated_at)
// against period_trunc($at). If they match, the existing row belongs to the
// current period and we accumulate. If they don't match, the row is stale
// (e.g. it's the start of a new day) and we replace its values with just this
// sample. "alltime" never resets, so the check is always "accumulate".
func upsertPeriodStat(
	ctx context.Context,
	tx *sql.Tx,
	period string,
	fd fieldDef,
	value float64,
	at time.Time,
) error {
	// "same period" predicate: TRUE if the existing row's updated_at is in
	// the same period as the new sample's timestamp. For alltime, always TRUE.
	var samePeriod string
	if period == "alltime" {
		samePeriod = "TRUE"
	} else {
		samePeriod = fmt.Sprintf("%s = %s",
			periodTrunc(period, "period_stats.updated_at"),
			periodTrunc(period, "EXCLUDED.updated_at"),
		)
	}

	if fd.kind == fieldDailyCumulative {
		// daily_cumulative: track running MAX. min is conceptually always 0
		// (the column can't go negative), so we store 0 with the row's
		// updated_at as min_at. sum/count are not meaningful for cumulative
		// running totals at the day level.
		q := fmt.Sprintf(`
			INSERT INTO period_stats
				(period, field,
				 min_value, min_at, max_value, max_at,
				 sample_count, updated_at)
			VALUES ($1, $2, 0, $4, $3, $4, 1, $4)
			ON CONFLICT (period, field) DO UPDATE SET
				min_value = 0,
				min_at = EXCLUDED.updated_at,
				max_value = CASE
				    WHEN %[1]s THEN GREATEST(period_stats.max_value, EXCLUDED.max_value)
				    ELSE EXCLUDED.max_value
				END,
				max_at = CASE
				    WHEN %[1]s AND EXCLUDED.max_value > period_stats.max_value
				        THEN EXCLUDED.max_at
				    WHEN %[1]s
				        THEN period_stats.max_at
				    ELSE EXCLUDED.max_at
				END,
				sample_count = CASE
				    WHEN %[1]s THEN period_stats.sample_count + 1
				    ELSE 1
				END,
				sum_value = NULL,
				updated_at = EXCLUDED.updated_at`,
			samePeriod,
		)
		_, err := tx.ExecContext(ctx, q, period, fd.column, value, at)
		return err
	}

	// normal field: full min/max/sum/count tracking
	q := fmt.Sprintf(`
		INSERT INTO period_stats
			(period, field,
			 min_value, min_at, max_value, max_at,
			 sum_value, sample_count, updated_at)
		VALUES ($1, $2, $3, $4, $3, $4, $3, 1, $4)
		ON CONFLICT (period, field) DO UPDATE SET
			min_value = CASE
			    WHEN %[1]s THEN LEAST(period_stats.min_value, EXCLUDED.min_value)
			    ELSE EXCLUDED.min_value
			END,
			min_at = CASE
			    WHEN %[1]s AND EXCLUDED.min_value < period_stats.min_value
			        THEN EXCLUDED.min_at
			    WHEN %[1]s
			        THEN period_stats.min_at
			    ELSE EXCLUDED.min_at
			END,
			max_value = CASE
			    WHEN %[1]s THEN GREATEST(period_stats.max_value, EXCLUDED.max_value)
			    ELSE EXCLUDED.max_value
			END,
			max_at = CASE
			    WHEN %[1]s AND EXCLUDED.max_value > period_stats.max_value
			        THEN EXCLUDED.max_at
			    WHEN %[1]s
			        THEN period_stats.max_at
			    ELSE EXCLUDED.max_at
			END,
			sum_value = CASE
			    WHEN %[1]s THEN period_stats.sum_value + EXCLUDED.sum_value
			    ELSE EXCLUDED.sum_value
			END,
			sample_count = CASE
			    WHEN %[1]s THEN period_stats.sample_count + 1
			    ELSE 1
			END,
			updated_at = EXCLUDED.updated_at`,
		samePeriod,
	)
	_, err := tx.ExecContext(ctx, q, period, fd.column, value, at)
	return err
}

// ----------------------------------------------------------------------------
// Paranoia rebuild: recompute all 4 periods from records
// ----------------------------------------------------------------------------

// rebuildStatsFromRecords runs aggregate queries directly against `records`
// for each (period, field) pair and replaces the period_stats rows. This is
// the daily cron job and the backfill path -- not the hot path.
//
// Cost: 4 periods * ~18 fields = ~72 aggregate queries. With an index on
// records.recorded these are fast even with millions of rows.
func rebuildStatsFromRecords(ctx context.Context) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().In(loc)
	type window struct {
		period     string
		start, end time.Time
	}
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
	yearStart := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, loc)
	farFuture := time.Date(9999, 12, 31, 23, 59, 59, 0, loc)
	farPast := time.Date(1970, 1, 1, 0, 0, 0, 0, loc)

	windows := []window{
		{"day", dayStart, dayStart.AddDate(0, 0, 1)},
		{"month", monthStart, monthStart.AddDate(0, 1, 0)},
		{"year", yearStart, yearStart.AddDate(1, 0, 0)},
		{"alltime", farPast, farFuture},
	}

	for _, w := range windows {
		for _, fd := range statFieldDefs {
			if err := rebuildOne(ctx, tx, w.period, w.start, w.end, fd); err != nil {
				return fmt.Errorf("rebuild %s/%s: %w", w.period, fd.column, err)
			}
		}
	}

	return tx.Commit()
}

func rebuildOne(
	ctx context.Context,
	tx *sql.Tx,
	period string,
	winStart, winEnd time.Time,
	fd fieldDef,
) error {
	if fd.kind == fieldDailyCumulative {
		// For daily_cumulative fields, the per-day "value" is MAX(column)
		// over the day. For day, that's just MAX over the window. For
		// month/year/alltime, we want SUM-of-daily-maxes.
		var q string
		if period == "day" {
			q = fmt.Sprintf(`
				WITH agg AS (
				    SELECT MAX(%[1]s) AS mx,
				           (SELECT recorded FROM records
				              WHERE recorded >= $2 AND recorded < $3
				                AND %[1]s IS NOT NULL
				              ORDER BY %[1]s DESC NULLS LAST LIMIT 1) AS mx_at,
				           COUNT(*) AS n
				    FROM records
				    WHERE recorded >= $2 AND recorded < $3
				)
				INSERT INTO period_stats
				    (period, field,
				     min_value, min_at, max_value, max_at,
				     sample_count, updated_at)
				SELECT $1, $4, 0, mx_at, mx, mx_at, n, now() FROM agg WHERE n > 0
				ON CONFLICT (period, field) DO UPDATE SET
				    min_value = 0,
				    min_at    = EXCLUDED.min_at,
				    max_value = EXCLUDED.max_value,
				    max_at    = EXCLUDED.max_at,
				    sum_value = NULL,
				    sample_count = EXCLUDED.sample_count,
				    updated_at = now()`, fd.column)
		} else {
			// SUM of per-day MAX values
			q = fmt.Sprintf(`
				WITH daily AS (
				    SELECT date_trunc('day', recorded AT TIME ZONE '%[2]s') AS d,
				           MAX(%[1]s) AS dmx,
				           MAX(recorded) AS d_last_at
				    FROM records
				    WHERE recorded >= $2 AND recorded < $3
				    GROUP BY 1
				), agg AS (
				    SELECT MAX(dmx) AS mx,
				           (SELECT d_last_at FROM daily
				              ORDER BY dmx DESC NULLS LAST LIMIT 1) AS mx_at,
				           SUM(dmx) AS s,
				           COUNT(*) AS n
				    FROM daily
				)
				INSERT INTO period_stats
				    (period, field,
				     min_value, min_at, max_value, max_at,
				     sum_value, sample_count, updated_at)
				SELECT $1, $4, 0, mx_at, mx, mx_at, s, n, now() FROM agg WHERE n > 0
				ON CONFLICT (period, field) DO UPDATE SET
				    min_value = 0,
				    min_at    = EXCLUDED.min_at,
				    max_value = EXCLUDED.max_value,
				    max_at    = EXCLUDED.max_at,
				    sum_value = EXCLUDED.sum_value,
				    sample_count = EXCLUDED.sample_count,
				    updated_at = now()`, fd.column, tzName)
		}
		_, err := tx.ExecContext(ctx, q, period, winStart, winEnd, fd.column)
		return err
	}

	// normal field
	q := fmt.Sprintf(`
		WITH agg AS (
		    SELECT MIN(%[1]s) AS mn,
		           (SELECT recorded FROM records
		              WHERE recorded >= $2 AND recorded < $3
		                AND %[1]s IS NOT NULL
		              ORDER BY %[1]s ASC NULLS LAST LIMIT 1) AS mn_at,
		           MAX(%[1]s) AS mx,
		           (SELECT recorded FROM records
		              WHERE recorded >= $2 AND recorded < $3
		                AND %[1]s IS NOT NULL
		              ORDER BY %[1]s DESC NULLS LAST LIMIT 1) AS mx_at,
		           SUM(%[1]s) AS s,
		           COUNT(%[1]s) AS n
		    FROM records
		    WHERE recorded >= $2 AND recorded < $3
		)
		INSERT INTO period_stats
		    (period, field,
		     min_value, min_at, max_value, max_at,
		     sum_value, sample_count, updated_at)
		SELECT $1, $4, mn, mn_at, mx, mx_at, s, n, now() FROM agg WHERE n > 0
		ON CONFLICT (period, field) DO UPDATE SET
		    min_value = EXCLUDED.min_value,
		    min_at    = EXCLUDED.min_at,
		    max_value = EXCLUDED.max_value,
		    max_at    = EXCLUDED.max_at,
		    sum_value = EXCLUDED.sum_value,
		    sample_count = EXCLUDED.sample_count,
		    updated_at = now()`, fd.column)

	_, err := tx.ExecContext(ctx, q, period, winStart, winEnd, fd.column)
	return err
}

// ----------------------------------------------------------------------------
// Maintenance
// ----------------------------------------------------------------------------

// pruneIngestLog drops ingest log entries older than 24 hours. Anything
// older than that can no longer be a "duplicate" we care about.
func pruneIngestLog(ctx context.Context) error {
	_, err := db.ExecContext(ctx,
		`DELETE FROM stat_ingest_log WHERE folded_at < now() - interval '24 hours'`)
	return err
}

// ----------------------------------------------------------------------------
// Helpers
// ----------------------------------------------------------------------------

type recordRow struct {
	recorded time.Time
	values   map[string]float64
}

func loadRecordNear(ctx context.Context, ts time.Time) (recordRow, bool, error) {
	cols := make([]string, 0, len(statFieldDefs)+1)
	cols = append(cols, "recorded")
	for _, fd := range statFieldDefs {
		cols = append(cols, fd.column)
	}

	// +/-1s window to absorb sub-second precision drift between publisher and DB
	query := fmt.Sprintf(`
		SELECT %s
		FROM records
		WHERE recorded BETWEEN $1 AND $2
		ORDER BY ABS(EXTRACT(EPOCH FROM (recorded - $3)))
		LIMIT 1`,
		strings.Join(cols, ", "))

	scanTargets := make([]interface{}, len(cols))
	var recordedTS time.Time
	scanTargets[0] = &recordedTS
	values := make([]sql.NullFloat64, len(statFieldDefs))
	for i := range values {
		scanTargets[i+1] = &values[i]
	}

	err := db.QueryRowContext(ctx, query,
		ts.Add(-time.Second), ts.Add(time.Second), ts,
	).Scan(scanTargets...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return recordRow{}, false, nil
		}
		return recordRow{}, false, err
	}

	rec := recordRow{
		recorded: recordedTS,
		values:   make(map[string]float64, len(statFieldDefs)),
	}
	for i, fd := range statFieldDefs {
		if values[i].Valid {
			rec.values[fd.column] = values[i].Float64
		}
	}
	return rec, true, nil
}
