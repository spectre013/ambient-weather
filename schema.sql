-- SCHEMA: public

-- DROP SCHEMA IF EXISTS public ;

CREATE SCHEMA IF NOT EXISTS public
    AUTHORIZATION postgres;

COMMENT ON SCHEMA public
    IS 'standard public schema';

GRANT ALL ON SCHEMA public TO PUBLIC;

GRANT ALL ON SCHEMA public TO postgres;

-- SEQUENCE: public.alerts_seq

-- DROP SEQUENCE IF EXISTS public.alerts_seq;

CREATE SEQUENCE IF NOT EXISTS public.alerts_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

ALTER SEQUENCE public.alerts_seq
    OWNER TO ambient;

-- SEQUENCE: public.records_seq

-- DROP SEQUENCE IF EXISTS public.records_seq;

CREATE SEQUENCE IF NOT EXISTS public.records_seq
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

ALTER SEQUENCE public.records_seq
    OWNER TO ambient;

-- Table: public.alerts

-- DROP TABLE IF EXISTS public.alerts;

CREATE TABLE IF NOT EXISTS public.alerts
(
    id integer NOT NULL DEFAULT nextval('alerts_seq'::regclass),
    alertid character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    wxtype character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    areadesc text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    sent timestamp(0) without time zone NOT NULL DEFAULT NULL::timestamp without time zone,
    effective timestamp(0) without time zone NOT NULL DEFAULT NULL::timestamp without time zone,
    onset timestamp(0) without time zone NOT NULL DEFAULT NULL::timestamp without time zone,
    expires timestamp(0) without time zone NOT NULL DEFAULT NULL::timestamp without time zone,
    ends timestamp(0) without time zone NOT NULL DEFAULT NULL::timestamp without time zone,
    status character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    messagetype character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    category character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    severity character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    certainty character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    urgency character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    event character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    sender character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    sendername character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    headline text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    description text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    instruction text COLLATE pg_catalog."default" NOT NULL DEFAULT ''::text,
    response character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    CONSTRAINT alerts_pkey PRIMARY KEY (id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.alerts
    OWNER to ambient;

-- Table: public.records

-- DROP TABLE IF EXISTS public.records;

CREATE TABLE IF NOT EXISTS public.records
(
    id integer NOT NULL DEFAULT nextval('records_seq'::regclass),
    mac character varying(255) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    recorded timestamp(0) with time zone NOT NULL DEFAULT (now())::timestamp without time zone,
    baromabsin real NOT NULL DEFAULT 0.0,
    baromrelin real NOT NULL DEFAULT 0.0,
    battout integer NOT NULL DEFAULT 0,
    batt1 integer NOT NULL DEFAULT 0,
    batt2 integer NOT NULL DEFAULT 0,
    batt3 integer NOT NULL DEFAULT 0,
    batt4 integer NOT NULL DEFAULT 0,
    batt5 integer NOT NULL DEFAULT 0,
    batt6 integer NOT NULL DEFAULT 0,
    batt7 integer NOT NULL DEFAULT 0,
    batt8 integer NOT NULL DEFAULT 0,
    batt9 integer NOT NULL DEFAULT 0,
    batt10 integer NOT NULL DEFAULT 0,
    co2 real NOT NULL DEFAULT 0.0,
    dailyrainin real NOT NULL DEFAULT 0.0,
    dewpoint real NOT NULL DEFAULT 0.0,
    eventrainin real NOT NULL DEFAULT 0.0,
    feelslike real NOT NULL DEFAULT 0.0,
    hourlyrainin real NOT NULL DEFAULT 0.0,
    hourlyrain real NOT NULL DEFAULT 0.0,
    humidity integer NOT NULL DEFAULT 0,
    humidity1 integer NOT NULL DEFAULT 0,
    humidity2 integer NOT NULL DEFAULT 0,
    humidity3 integer NOT NULL DEFAULT 0,
    humidity4 integer NOT NULL DEFAULT 0,
    humidity5 integer NOT NULL DEFAULT 0,
    humidity6 integer NOT NULL DEFAULT 0,
    humidity7 integer NOT NULL DEFAULT 0,
    humidity8 integer NOT NULL DEFAULT 0,
    humidity9 integer NOT NULL DEFAULT 0,
    humidity10 integer NOT NULL DEFAULT 0,
    humidityin integer NOT NULL DEFAULT 0,
    lastrain timestamp(0) without time zone NOT NULL DEFAULT NULL::timestamp without time zone,
    maxdailygust real NOT NULL DEFAULT 0.0,
    relay1 integer NOT NULL DEFAULT 0,
    relay2 integer NOT NULL DEFAULT 0,
    relay3 integer NOT NULL DEFAULT 0,
    relay4 integer NOT NULL DEFAULT 0,
    relay5 integer NOT NULL DEFAULT 0,
    relay6 integer NOT NULL DEFAULT 0,
    relay7 integer NOT NULL DEFAULT 0,
    relay8 integer NOT NULL DEFAULT 0,
    relay9 integer NOT NULL DEFAULT 0,
    relay10 integer NOT NULL DEFAULT 0,
    monthlyrainin real NOT NULL DEFAULT 0.0,
    solarradiation real NOT NULL DEFAULT 0.0,
    tempf real NOT NULL DEFAULT 0.0,
    temp1f real NOT NULL DEFAULT 0.0,
    temp2f real NOT NULL DEFAULT 0.0,
    temp3f real NOT NULL DEFAULT 0.0,
    temp4f real NOT NULL DEFAULT 0.0,
    temp5f real NOT NULL DEFAULT 0.0,
    temp6f real NOT NULL DEFAULT 0.0,
    temp7f real NOT NULL DEFAULT 0.0,
    temp8f real NOT NULL DEFAULT 0.0,
    temp9f real NOT NULL DEFAULT 0.0,
    temp10f real NOT NULL DEFAULT 0.0,
    tempinf real NOT NULL DEFAULT 0.0,
    totalrainin real NOT NULL DEFAULT 0.0,
    uv real NOT NULL DEFAULT 0.0,
    weeklyrainin real NOT NULL DEFAULT 0.0,
    winddir integer NOT NULL DEFAULT 0,
    windgustmph real NOT NULL DEFAULT 0.0,
    windgustdir integer NOT NULL DEFAULT 0,
    windspeedmph real NOT NULL DEFAULT 0.0,
    yearlyrainin real NOT NULL DEFAULT 0.0,
    battlightning integer NOT NULL DEFAULT 0,
    lightningday integer NOT NULL DEFAULT 0,
    lightninghour integer NOT NULL DEFAULT 0,
    lightningtime timestamp(0) without time zone NOT NULL DEFAULT NULL::timestamp without time zone,
    lightningdistance real NOT NULL DEFAULT 0.0,
    aqipm25 integer DEFAULT 0,
    aqipm2524h integer DEFAULT 0,
    CONSTRAINT records_pkey PRIMARY KEY (id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.records
    OWNER to ambient;
-- Index: baromrelinidx

-- DROP INDEX IF EXISTS public.baromrelinidx;

CREATE INDEX IF NOT EXISTS baromrelinidx
    ON public.records USING btree
        (baromrelin ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: dailyraininidx

-- DROP INDEX IF EXISTS public.dailyraininidx;

CREATE INDEX IF NOT EXISTS dailyraininidx
    ON public.records USING btree
        (dailyrainin ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: date

-- DROP INDEX IF EXISTS public.date;

CREATE INDEX IF NOT EXISTS date
    ON public.records USING btree
        (recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: datetempfidx

-- DROP INDEX IF EXISTS public.datetempfidx;

CREATE INDEX IF NOT EXISTS datetempfidx
    ON public.records USING btree
        (recorded ASC NULLS LAST, tempf ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: dewpointidx

-- DROP INDEX IF EXISTS public.dewpointidx;

CREATE INDEX IF NOT EXISTS dewpointidx
    ON public.records USING btree
        (dewpoint ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: feelsidx

-- DROP INDEX IF EXISTS public.feelsidx;

CREATE INDEX IF NOT EXISTS feelsidx
    ON public.records USING btree
        (feelslike ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: humidity1idx

-- DROP INDEX IF EXISTS public.humidity1idx;

CREATE INDEX IF NOT EXISTS humidity1idx
    ON public.records USING btree
        (humidity1 ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: humidity2idx

-- DROP INDEX IF EXISTS public.humidity2idx;

CREATE INDEX IF NOT EXISTS humidity2idx
    ON public.records USING btree
        (humidity2 ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: humidityidx

-- DROP INDEX IF EXISTS public.humidityidx;

CREATE INDEX IF NOT EXISTS humidityidx
    ON public.records USING btree
        (humidity ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: humidityinidx

-- DROP INDEX IF EXISTS public.humidityinidx;

CREATE INDEX IF NOT EXISTS humidityinidx
    ON public.records USING btree
        (humidityin ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: lightningdayidx

-- DROP INDEX IF EXISTS public.lightningdayidx;

CREATE INDEX IF NOT EXISTS lightningdayidx
    ON public.records USING btree
        (lightningday ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: recorded

-- DROP INDEX IF EXISTS public.recorded;

CREATE INDEX IF NOT EXISTS recorded
    ON public.records USING btree
        (recorded DESC NULLS LAST)
    TABLESPACE pg_default;
-- Index: temp1fidx

-- DROP INDEX IF EXISTS public.temp1fidx;

CREATE INDEX IF NOT EXISTS temp1fidx
    ON public.records USING btree
        (temp1f ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: temp2fidx

-- DROP INDEX IF EXISTS public.temp2fidx;

CREATE INDEX IF NOT EXISTS temp2fidx
    ON public.records USING btree
        (temp2f ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: tempfidx

-- DROP INDEX IF EXISTS public.tempfidx;

CREATE INDEX IF NOT EXISTS tempfidx
    ON public.records USING btree
        (tempf ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: tempinfidx

-- DROP INDEX IF EXISTS public.tempinfidx;

CREATE INDEX IF NOT EXISTS tempinfidx
    ON public.records USING btree
        (tempinf ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: uvidx

-- DROP INDEX IF EXISTS public.uvidx;

CREATE INDEX IF NOT EXISTS uvidx
    ON public.records USING btree
        (uv ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: windgustidx

-- DROP INDEX IF EXISTS public.windgustidx;

CREATE INDEX IF NOT EXISTS windgustidx
    ON public.records USING btree
        (windgustmph ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: windgustmphidx

-- DROP INDEX IF EXISTS public.windgustmphidx;

CREATE INDEX IF NOT EXISTS windgustmphidx
    ON public.records USING btree
        (windgustmph ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: windidx

-- DROP INDEX IF EXISTS public.windidx;

CREATE INDEX IF NOT EXISTS windidx
    ON public.records USING btree
        (windspeedmph ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;
-- Index: windspeedmphidx

-- DROP INDEX IF EXISTS public.windspeedmphidx;

CREATE INDEX IF NOT EXISTS windspeedmphidx
    ON public.records USING btree
        (windspeedmph ASC NULLS LAST, recorded ASC NULLS LAST)
    TABLESPACE pg_default;

-- Table: public.stats

-- DROP TABLE IF EXISTS public.stats;

CREATE TABLE IF NOT EXISTS public.stats
(
    id character varying(100) COLLATE pg_catalog."default" NOT NULL DEFAULT ''::character varying,
    recorded timestamp(0) without time zone DEFAULT (now())::timestamp without time zone,
    value numeric NOT NULL DEFAULT 0.0,
    CONSTRAINT stats_pkey PRIMARY KEY (id)
)

    TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.stats
    OWNER to ambient;