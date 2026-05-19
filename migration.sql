-- Migration: new stats design
-- ----------------------------------------------------------------------------
-- Drops the old stats table and replaces it with period_stats + supporting
-- objects. Run after deploying the new Go binary; the binary will backfill
-- when started with STATS_BACKFILL=true.
--
-- IMPORTANT: this drops the old `stats` table. Take a backup first.

BEGIN;

-- ----------------------------------------------------------------------------
-- 1. period_stats: one row per (period, period_start, field)
-- ----------------------------------------------------------------------------
DROP TABLE IF EXISTS stats CASCADE;

CREATE TABLE period_stats (
                              period       TEXT        NOT NULL,
                              period_start DATE        NOT NULL,
                              field        TEXT        NOT NULL,

                              min_value    NUMERIC(10,2),
                              min_at       TIMESTAMPTZ,
                              max_value    NUMERIC(10,2),
                              max_at       TIMESTAMPTZ,

    -- sum_value carries two meanings depending on the field's kind:
    --   normal fields:           running sum of sample values; avg = sum/count
    --   daily_cumulative fields: sum of per-day MAX values (e.g. monthly rain)
                              sum_value    NUMERIC(14,4),
                              sample_count BIGINT      NOT NULL DEFAULT 0,

                              updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

                              PRIMARY KEY (period, period_start, field),
                              CONSTRAINT period_stats_period_chk
                                  CHECK (period IN ('day','month','year','alltime'))
);

CREATE INDEX idx_period_stats_period_start
    ON period_stats (period, period_start DESC);

CREATE INDEX idx_period_stats_field
    ON period_stats (field, period, period_start DESC);

-- ----------------------------------------------------------------------------
-- 2. period_stats_v: convenience view that exposes avg_value
-- ----------------------------------------------------------------------------
CREATE OR REPLACE VIEW period_stats_v AS
SELECT
    period,
    period_start,
    field,
    min_value,
    min_at,
    max_value,
    max_at,
    CASE WHEN sample_count > 0
             THEN ROUND(sum_value / sample_count, 2)
         ELSE NULL
        END AS avg_value,
    sum_value,
    sample_count,
    updated_at
FROM period_stats;

-- ----------------------------------------------------------------------------
-- 3. stat_ingest_log: idempotency for record folding
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS stat_ingest_log (
                                               recorded  TIMESTAMPTZ PRIMARY KEY,
                                               folded_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Optional retention: keep last 90 days of ingest log entries. Older entries
-- can be pruned because the day they belong to has long since been finalized.
-- Run this from a separate maintenance cron if desired:
--
--   DELETE FROM stat_ingest_log WHERE recorded < now() - interval '90 days';

COMMIT;