-- Migration: rain rate column
-- ----------------------------------------------------------------------------
-- Adds records.rainratein (rain rate in in/hr: the last 10 minutes x 6, as the
-- HP2550 console computes it). Run before deploying the processor that
-- writes it; the server reads it as well.
--
-- From this release hourlyrainin / hourlyrain hold the rolling 60-minute
-- total instead of rain since the top of the clock hour. Existing rows are
-- left as they are.

ALTER TABLE public.records
    ADD COLUMN IF NOT EXISTS rainratein real NOT NULL DEFAULT 0.0;
