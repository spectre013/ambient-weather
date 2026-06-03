-- Migration: observed current conditions table
-- ----------------------------------------------------------------------------
-- Adds public.conditions, populated by the stats service from the NWS
-- /stations/{STATION_ID}/observations/latest endpoint every 15 minutes.
-- The server reads the latest row per station to populate the /api/app
-- `conditions` (text) and `icon` (Visual-Crossing-style kind) fields.
--
-- Safe to run repeatedly (IF NOT EXISTS). Idempotent.

BEGIN;

CREATE TABLE IF NOT EXISTS public.conditions
(
    id               bigserial PRIMARY KEY,
    station          text        NOT NULL,
    observed_at      timestamptz NOT NULL,
    conditions       text,                 -- human text (NWS textDescription, refined)
    icon             text,                 -- derived Visual-Crossing-style icon kind
    text_description text,                 -- raw NWS textDescription
    present_weather  text,                 -- JSON-encoded raw presentWeather[]
    cloud_layers     text,                 -- JSON-encoded raw cloudLayers[]
    raw_icon         text,                 -- original NWS icon URL
    temperature      numeric(6,2),         -- °F (converted from observed °C)
    humidity         numeric(6,2),         -- %
    created_at       timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT conditions_station_observed_uk UNIQUE (station, observed_at)
);

ALTER TABLE IF EXISTS public.conditions
    OWNER to ambient;

CREATE INDEX IF NOT EXISTS idx_conditions_observed
    ON public.conditions USING btree
        (station, observed_at DESC);

COMMIT;
