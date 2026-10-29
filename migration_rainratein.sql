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


ALTER TABLE public.records
    -- WS80 (outdoor)
    ALTER COLUMN tempf DROP NOT NULL,
    ALTER COLUMN humidity DROP NOT NULL,
    ALTER COLUMN dewpoint DROP NOT NULL,
    ALTER COLUMN feelslike DROP NOT NULL,
    ALTER COLUMN winddir DROP NOT NULL,
    ALTER COLUMN windspeedmph DROP NOT NULL,
    ALTER COLUMN windgustmph DROP NOT NULL,
    ALTER COLUMN windgustdir DROP NOT NULL,
    ALTER COLUMN uv DROP NOT NULL,
    ALTER COLUMN solarradiation DROP NOT NULL,
    -- WH32B (indoor / barometer)
    ALTER COLUMN tempinf DROP NOT NULL,
    ALTER COLUMN humidityin DROP NOT NULL,
    ALTER COLUMN baromabsin DROP NOT NULL,
    ALTER COLUMN baromrelin DROP NOT NULL,
    -- WH31E channels
    ALTER COLUMN temp1f DROP NOT NULL,
    ALTER COLUMN temp2f DROP NOT NULL,
    ALTER COLUMN temp3f DROP NOT NULL,
    ALTER COLUMN temp4f DROP NOT NULL,
    ALTER COLUMN temp5f DROP NOT NULL,
    ALTER COLUMN temp6f DROP NOT NULL,
    ALTER COLUMN temp7f DROP NOT NULL,
    ALTER COLUMN temp8f DROP NOT NULL,
    ALTER COLUMN temp9f DROP NOT NULL,
    ALTER COLUMN temp10f DROP NOT NULL,
    ALTER COLUMN humidity1 DROP NOT NULL,
    ALTER COLUMN humidity2 DROP NOT NULL,
    ALTER COLUMN humidity3 DROP NOT NULL,
    ALTER COLUMN humidity4 DROP NOT NULL,
    ALTER COLUMN humidity5 DROP NOT NULL,
    ALTER COLUMN humidity6 DROP NOT NULL,
    ALTER COLUMN humidity7 DROP NOT NULL,
    ALTER COLUMN humidity8 DROP NOT NULL,
    ALTER COLUMN humidity9 DROP NOT NULL,
    ALTER COLUMN humidity10 DROP NOT NULL,
    -- WH0290 (air quality)
    ALTER COLUMN aqipm25 DROP NOT NULL,
    ALTER COLUMN aqipm2524h DROP NOT NULL,
    -- Timestamps that are NULL until the first event
    ALTER COLUMN lastrain DROP NOT NULL,
    ALTER COLUMN lightningtime DROP NOT NULL;
