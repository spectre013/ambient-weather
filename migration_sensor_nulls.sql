-- Migration: NULL for sensors that are not reporting
-- ----------------------------------------------------------------------------
-- The processor now writes NULL to a sensor's reading columns once that
-- sensor has been silent for its stale period (10 minutes; 30 for the
-- WH0290), and for sensors it has never heard from, instead of repeating the
-- last reading or writing 0. Stats and the server skip NULLs.
--
-- Run before deploying the processor that writes NULLs; its inserts fail on
-- the NOT NULL constraints otherwise. The server must be deployed alongside
-- it (it reads these columns as nullable). Existing rows are left as they are.
--
-- lastrain and lightningtime are included because the processor already
-- sends NULL for them until the first rain / strike has been seen, which a
-- fresh database rejects.
--
-- Idempotent: DROP NOT NULL on a nullable column is a no-op.

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
