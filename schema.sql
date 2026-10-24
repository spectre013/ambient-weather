--
-- PostgreSQL database dump
--

\restrict q2kQrjCANIw2QU2KNRFymfKk7yFc0jsxkSkBwN6D9X60hG0IqsjACRNW0ZDScv4

-- Dumped from database version 16.15 (Ubuntu 16.15-0ubuntu0.24.04.1)
-- Dumped by pg_dump version 16.15 (Ubuntu 16.15-0ubuntu0.24.04.1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: pg_stat_statements; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS pg_stat_statements WITH SCHEMA public;


--
-- Name: EXTENSION pg_stat_statements; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION pg_stat_statements IS 'track planning and execution statistics of all SQL statements executed';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: alerts; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.alerts (
    alertid character varying(255) DEFAULT ''::character varying NOT NULL,
    wxtype character varying(255) DEFAULT ''::character varying NOT NULL,
    areadesc text DEFAULT ''::character varying NOT NULL,
    sent timestamp(0) without time zone DEFAULT NULL::timestamp without time zone NOT NULL,
    effective timestamp(0) without time zone DEFAULT NULL::timestamp without time zone NOT NULL,
    onset timestamp(0) without time zone DEFAULT NULL::timestamp without time zone NOT NULL,
    expires timestamp(0) without time zone DEFAULT NULL::timestamp without time zone NOT NULL,
    ends timestamp(0) without time zone DEFAULT NULL::timestamp without time zone NOT NULL,
    status character varying(255) DEFAULT ''::character varying NOT NULL,
    messagetype character varying(255) DEFAULT ''::character varying NOT NULL,
    category character varying(255) DEFAULT ''::character varying NOT NULL,
    severity character varying(255) DEFAULT ''::character varying NOT NULL,
    certainty character varying(255) DEFAULT ''::character varying NOT NULL,
    urgency character varying(255) DEFAULT ''::character varying NOT NULL,
    event character varying(255) DEFAULT ''::character varying NOT NULL,
    sender character varying(255) DEFAULT ''::character varying NOT NULL,
    sendername character varying(255) DEFAULT ''::character varying NOT NULL,
    headline text DEFAULT ''::character varying NOT NULL,
    description text DEFAULT ''::text NOT NULL,
    instruction text DEFAULT ''::text NOT NULL,
    response character varying(255) DEFAULT ''::character varying NOT NULL,
    summary text
);


ALTER TABLE public.alerts OWNER TO ambient;

--
-- Name: conditions; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.conditions (
    id bigint NOT NULL,
    station text NOT NULL,
    observed_at timestamp with time zone NOT NULL,
    conditions text,
    icon text,
    text_description text,
    present_weather text,
    cloud_layers text,
    raw_icon text,
    temperature numeric(6,2),
    humidity numeric(6,2),
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.conditions OWNER TO ambient;

--
-- Name: conditions_id_seq; Type: SEQUENCE; Schema: public; Owner: ambient
--

CREATE SEQUENCE public.conditions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.conditions_id_seq OWNER TO ambient;

--
-- Name: conditions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ambient
--

ALTER SEQUENCE public.conditions_id_seq OWNED BY public.conditions.id;


--
-- Name: forecast; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.forecast (
    datetime timestamp with time zone NOT NULL,
    datetime_epoch bigint,
    tempmax numeric(5,2),
    tempmin numeric(5,2),
    temp numeric(5,2),
    feelslikemax numeric(5,2),
    feelslikemin numeric(5,2),
    feelslike numeric(5,2),
    dew numeric(5,2),
    humidity numeric(5,2),
    precip numeric(6,3),
    precipprob numeric(5,2),
    precipcover numeric(5,2),
    preciptype text,
    snow numeric(5,2),
    snowdepth numeric(5,2),
    windgust numeric(5,2),
    windspeed numeric(5,2),
    winddir numeric(5,2),
    pressure numeric(6,2),
    cloudcover numeric(5,2),
    visibility numeric(5,2),
    solarradiation numeric(6,2),
    solarenergy numeric(5,2),
    uvindex integer,
    severerisk integer,
    sunrise time without time zone,
    sunrise_epoch bigint,
    sunset time without time zone,
    sunset_epoch bigint,
    moonphase numeric(3,2),
    conditions text,
    description text,
    icon text,
    stations text,
    source text,
    hours jsonb,
    summary text,
    model_blend text,
    based_on_runs text,
    transformed_at timestamp with time zone DEFAULT now() NOT NULL,
    content_hash text,
    summary_generated_at timestamp with time zone,
    summary_hash text,
    location_id bigint NOT NULL
);


ALTER TABLE public.forecast OWNER TO ambient;

--
-- Name: forecast_data; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.forecast_data (
    id bigint NOT NULL,
    model_type text NOT NULL,
    requested_lat double precision NOT NULL,
    requested_lon double precision NOT NULL,
    grid_lat double precision,
    grid_lon double precision,
    model_run_time timestamp with time zone NOT NULL,
    forecast_hour integer NOT NULL,
    valid_time timestamp with time zone NOT NULL,
    temp_2m double precision,
    dpt_2m double precision,
    rh_2m double precision,
    pot_2m double precision,
    tmax_2m double precision,
    tmin_2m double precision,
    ugrd_10m double precision,
    vgrd_10m double precision,
    wind_10m double precision,
    maxvw_10m double precision,
    temp_surface double precision,
    vis_surface double precision,
    gust_surface double precision,
    pres_surface double precision,
    snowc_surface double precision,
    snod_surface double precision,
    frozr_surface double precision,
    frzr_surface double precision,
    prate_surface double precision,
    apcp_surface double precision,
    crain_surface double precision,
    csnow_surface double precision,
    cfrzr_surface double precision,
    cpofp_surface double precision,
    fricv_surface double precision,
    hail_atmosphere double precision,
    ltng_atmosphere double precision,
    tcdc_atmosphere double precision,
    inserted_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    dswrf_surface double precision,
    asnow_surface double precision,
    weasd_surface double precision,
    cape_surface double precision,
    pop_surface double precision,
    location_id bigint NOT NULL
);


ALTER TABLE public.forecast_data OWNER TO ambient;

--
-- Name: forecast_blend; Type: VIEW; Schema: public; Owner: ambient
--

CREATE VIEW public.forecast_blend AS
 WITH hrrr_pick AS (
         SELECT DISTINCT ON (forecast_data.location_id, forecast_data.valid_time) forecast_data.location_id,
            forecast_data.valid_time,
            forecast_data.model_type,
            forecast_data.model_run_time,
            forecast_data.forecast_hour,
            forecast_data.temp_2m,
            forecast_data.dpt_2m,
            forecast_data.rh_2m,
            forecast_data.tmax_2m,
            forecast_data.tmin_2m,
            forecast_data.ugrd_10m,
            forecast_data.vgrd_10m,
            forecast_data.wind_10m,
            forecast_data.temp_surface,
            forecast_data.vis_surface,
            forecast_data.gust_surface,
            forecast_data.pres_surface,
            forecast_data.snod_surface,
            forecast_data.prate_surface,
            forecast_data.apcp_surface,
            forecast_data.crain_surface,
            forecast_data.csnow_surface,
            forecast_data.cfrzr_surface,
            forecast_data.asnow_surface,
            forecast_data.weasd_surface,
            forecast_data.dswrf_surface,
            forecast_data.cape_surface,
            forecast_data.tcdc_atmosphere
           FROM public.forecast_data
          WHERE ((forecast_data.model_type = 'hrrr'::text) AND (forecast_data.temp_2m IS NOT NULL))
          ORDER BY forecast_data.location_id, forecast_data.valid_time, forecast_data.model_run_time DESC, forecast_data.forecast_hour
        ), gfs_pick AS (
         SELECT DISTINCT ON (forecast_data.location_id, forecast_data.valid_time) forecast_data.location_id,
            forecast_data.valid_time,
            forecast_data.model_type,
            forecast_data.model_run_time,
            forecast_data.forecast_hour,
            forecast_data.temp_2m,
            forecast_data.dpt_2m,
            forecast_data.rh_2m,
            forecast_data.tmax_2m,
            forecast_data.tmin_2m,
            forecast_data.ugrd_10m,
            forecast_data.vgrd_10m,
            NULL::double precision AS wind_10m,
            NULL::double precision AS temp_surface,
            forecast_data.vis_surface,
            forecast_data.gust_surface,
            forecast_data.pres_surface,
            forecast_data.snod_surface,
            forecast_data.prate_surface,
            forecast_data.apcp_surface,
            forecast_data.crain_surface,
            forecast_data.csnow_surface,
            forecast_data.cfrzr_surface,
            NULL::double precision AS asnow_surface,
            forecast_data.weasd_surface,
            forecast_data.dswrf_surface,
            forecast_data.cape_surface,
            forecast_data.tcdc_atmosphere
           FROM public.forecast_data
          WHERE ((forecast_data.model_type = 'gfs'::text) AND (forecast_data.temp_2m IS NOT NULL))
          ORDER BY forecast_data.location_id, forecast_data.valid_time, forecast_data.model_run_time DESC, forecast_data.forecast_hour
        ), deterministic AS (
         SELECT hrrr_pick.location_id,
            hrrr_pick.valid_time,
            hrrr_pick.model_type,
            hrrr_pick.model_run_time,
            hrrr_pick.forecast_hour,
            hrrr_pick.temp_2m,
            hrrr_pick.dpt_2m,
            hrrr_pick.rh_2m,
            hrrr_pick.tmax_2m,
            hrrr_pick.tmin_2m,
            hrrr_pick.ugrd_10m,
            hrrr_pick.vgrd_10m,
            hrrr_pick.wind_10m,
            hrrr_pick.temp_surface,
            hrrr_pick.vis_surface,
            hrrr_pick.gust_surface,
            hrrr_pick.pres_surface,
            hrrr_pick.snod_surface,
            hrrr_pick.prate_surface,
            hrrr_pick.apcp_surface,
            hrrr_pick.crain_surface,
            hrrr_pick.csnow_surface,
            hrrr_pick.cfrzr_surface,
            hrrr_pick.asnow_surface,
            hrrr_pick.weasd_surface,
            hrrr_pick.dswrf_surface,
            hrrr_pick.cape_surface,
            hrrr_pick.tcdc_atmosphere
           FROM hrrr_pick
        UNION ALL
         SELECT g.location_id,
            g.valid_time,
            g.model_type,
            g.model_run_time,
            g.forecast_hour,
            g.temp_2m,
            g.dpt_2m,
            g.rh_2m,
            g.tmax_2m,
            g.tmin_2m,
            g.ugrd_10m,
            g.vgrd_10m,
            g.wind_10m,
            g.temp_surface,
            g.vis_surface,
            g.gust_surface,
            g.pres_surface,
            g.snod_surface,
            g.prate_surface,
            g.apcp_surface,
            g.crain_surface,
            g.csnow_surface,
            g.cfrzr_surface,
            g.asnow_surface,
            g.weasd_surface,
            g.dswrf_surface,
            g.cape_surface,
            g.tcdc_atmosphere
           FROM (gfs_pick g
             LEFT JOIN hrrr_pick h ON (((h.location_id = g.location_id) AND (h.valid_time = g.valid_time))))
          WHERE (h.location_id IS NULL)
        ), nbm_pick AS (
         SELECT DISTINCT ON (forecast_data.location_id, forecast_data.valid_time) forecast_data.location_id,
            forecast_data.valid_time,
            forecast_data.pop_surface
           FROM public.forecast_data
          WHERE ((forecast_data.model_type = 'nbm'::text) AND (forecast_data.pop_surface IS NOT NULL))
          ORDER BY forecast_data.location_id, forecast_data.valid_time, forecast_data.model_run_time DESC, forecast_data.forecast_hour
        )
 SELECT d.location_id,
    d.valid_time,
    d.model_type,
    d.model_run_time,
    d.forecast_hour,
    d.temp_2m,
    d.dpt_2m,
    d.rh_2m,
    d.tmax_2m,
    d.tmin_2m,
    d.ugrd_10m,
    d.vgrd_10m,
    d.wind_10m,
    d.temp_surface,
    d.vis_surface,
    d.gust_surface,
    d.pres_surface,
    d.snod_surface,
    d.prate_surface,
    d.apcp_surface,
    d.crain_surface,
    d.csnow_surface,
    d.cfrzr_surface,
    d.asnow_surface,
    d.weasd_surface,
    d.dswrf_surface,
    d.cape_surface,
    d.tcdc_atmosphere,
    n.pop_surface
   FROM (deterministic d
     LEFT JOIN nbm_pick n ON (((n.location_id = d.location_id) AND (n.valid_time = d.valid_time))));


ALTER VIEW public.forecast_blend OWNER TO ambient;

--
-- Name: forecast_data_id_seq; Type: SEQUENCE; Schema: public; Owner: ambient
--

CREATE SEQUENCE public.forecast_data_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.forecast_data_id_seq OWNER TO ambient;

--
-- Name: forecast_data_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ambient
--

ALTER SEQUENCE public.forecast_data_id_seq OWNED BY public.forecast_data.id;


--
-- Name: ingest_log; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.ingest_log (
    id bigint NOT NULL,
    model_type text NOT NULL,
    model_run_time timestamp with time zone NOT NULL,
    forecast_hour integer NOT NULL,
    file_path text NOT NULL,
    status text NOT NULL,
    error_message text,
    completed_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.ingest_log OWNER TO ambient;

--
-- Name: ingest_log_id_seq; Type: SEQUENCE; Schema: public; Owner: ambient
--

CREATE SEQUENCE public.ingest_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.ingest_log_id_seq OWNER TO ambient;

--
-- Name: ingest_log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ambient
--

ALTER SEQUENCE public.ingest_log_id_seq OWNED BY public.ingest_log.id;


--
-- Name: map_locations; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.map_locations (
    id bigint NOT NULL,
    name text NOT NULL,
    lat double precision NOT NULL,
    lon double precision NOT NULL,
    note text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    timezone text,
    forecast_enabled boolean DEFAULT false NOT NULL,
    last_viewed_at timestamp with time zone,
    CONSTRAINT map_locations_forecast_needs_tz CHECK (((NOT forecast_enabled) OR (timezone IS NOT NULL))),
    CONSTRAINT map_locations_lat_range CHECK (((lat >= ('-90'::integer)::double precision) AND (lat <= (90)::double precision))),
    CONSTRAINT map_locations_lon_range CHECK (((lon >= ('-180'::integer)::double precision) AND (lon <= (180)::double precision))),
    CONSTRAINT map_locations_name_len CHECK (((char_length(name) >= 1) AND (char_length(name) <= 80)))
);


ALTER TABLE public.map_locations OWNER TO ambient;

--
-- Name: map_locations_id_seq; Type: SEQUENCE; Schema: public; Owner: ambient
--

CREATE SEQUENCE public.map_locations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.map_locations_id_seq OWNER TO ambient;

--
-- Name: map_locations_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ambient
--

ALTER SEQUENCE public.map_locations_id_seq OWNED BY public.map_locations.id;


--
-- Name: period_stats; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.period_stats (
    period text NOT NULL,
    field text NOT NULL,
    min_value numeric(10,2),
    min_at timestamp with time zone,
    max_value numeric(10,2),
    max_at timestamp with time zone,
    sum_value numeric(14,4),
    sample_count bigint DEFAULT 0 NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.period_stats OWNER TO ambient;

--
-- Name: period_stats_v; Type: VIEW; Schema: public; Owner: postgres
--

CREATE VIEW public.period_stats_v AS
 SELECT period,
    field,
    min_value,
    min_at,
    max_value,
    max_at,
        CASE
            WHEN (sample_count > 0) THEN (sum_value / (sample_count)::numeric)
            ELSE NULL::numeric
        END AS avg_value,
    sum_value,
    sample_count,
    updated_at
   FROM public.period_stats;


ALTER VIEW public.period_stats_v OWNER TO postgres;

--
-- Name: records_seq; Type: SEQUENCE; Schema: public; Owner: ambient
--

CREATE SEQUENCE public.records_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.records_seq OWNER TO ambient;

--
-- Name: records; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.records (
    id integer DEFAULT nextval('public.records_seq'::regclass) NOT NULL,
    mac character varying(255) DEFAULT ''::character varying NOT NULL,
    recorded timestamp(0) with time zone DEFAULT (now())::timestamp without time zone NOT NULL,
    baromabsin real DEFAULT 0.0 NOT NULL,
    baromrelin real DEFAULT 0.0 NOT NULL,
    battout integer DEFAULT 0 NOT NULL,
    batt1 integer DEFAULT 0 NOT NULL,
    batt2 integer DEFAULT 0 NOT NULL,
    batt3 integer DEFAULT 0 NOT NULL,
    batt4 integer DEFAULT 0 NOT NULL,
    batt5 integer DEFAULT 0 NOT NULL,
    batt6 integer DEFAULT 0 NOT NULL,
    batt7 integer DEFAULT 0 NOT NULL,
    batt8 integer DEFAULT 0 NOT NULL,
    batt9 integer DEFAULT 0 NOT NULL,
    batt10 integer DEFAULT 0 NOT NULL,
    co2 real DEFAULT 0.0 NOT NULL,
    dailyrainin real DEFAULT 0.0 NOT NULL,
    dewpoint real DEFAULT 0.0 NOT NULL,
    eventrainin real DEFAULT 0.0 NOT NULL,
    feelslike real DEFAULT 0.0 NOT NULL,
    hourlyrainin real DEFAULT 0.0 NOT NULL,
    hourlyrain real DEFAULT 0.0 NOT NULL,
    humidity integer DEFAULT 0 NOT NULL,
    humidity1 integer DEFAULT 0 NOT NULL,
    humidity2 integer DEFAULT 0 NOT NULL,
    humidity3 integer DEFAULT 0 NOT NULL,
    humidity4 integer DEFAULT 0 NOT NULL,
    humidity5 integer DEFAULT 0 NOT NULL,
    humidity6 integer DEFAULT 0 NOT NULL,
    humidity7 integer DEFAULT 0 NOT NULL,
    humidity8 integer DEFAULT 0 NOT NULL,
    humidity9 integer DEFAULT 0 NOT NULL,
    humidity10 integer DEFAULT 0 NOT NULL,
    humidityin integer DEFAULT 0 NOT NULL,
    lastrain timestamp(0) without time zone DEFAULT NULL::timestamp without time zone NOT NULL,
    maxdailygust real DEFAULT 0.0 NOT NULL,
    relay1 integer DEFAULT 0 NOT NULL,
    relay2 integer DEFAULT 0 NOT NULL,
    relay3 integer DEFAULT 0 NOT NULL,
    relay4 integer DEFAULT 0 NOT NULL,
    relay5 integer DEFAULT 0 NOT NULL,
    relay6 integer DEFAULT 0 NOT NULL,
    relay7 integer DEFAULT 0 NOT NULL,
    relay8 integer DEFAULT 0 NOT NULL,
    relay9 integer DEFAULT 0 NOT NULL,
    relay10 integer DEFAULT 0 NOT NULL,
    monthlyrainin real DEFAULT 0.0 NOT NULL,
    solarradiation real DEFAULT 0.0 NOT NULL,
    tempf real DEFAULT 0.0 NOT NULL,
    temp1f real DEFAULT 0.0 NOT NULL,
    temp2f real DEFAULT 0.0 NOT NULL,
    temp3f real DEFAULT 0.0 NOT NULL,
    temp4f real DEFAULT 0.0 NOT NULL,
    temp5f real DEFAULT 0.0 NOT NULL,
    temp6f real DEFAULT 0.0 NOT NULL,
    temp7f real DEFAULT 0.0 NOT NULL,
    temp8f real DEFAULT 0.0 NOT NULL,
    temp9f real DEFAULT 0.0 NOT NULL,
    temp10f real DEFAULT 0.0 NOT NULL,
    tempinf real DEFAULT 0.0 NOT NULL,
    totalrainin real DEFAULT 0.0 NOT NULL,
    uv real DEFAULT 0.0 NOT NULL,
    weeklyrainin real DEFAULT 0.0 NOT NULL,
    winddir integer DEFAULT 0 NOT NULL,
    windgustmph real DEFAULT 0.0 NOT NULL,
    windgustdir integer DEFAULT 0 NOT NULL,
    windspeedmph real DEFAULT 0.0 NOT NULL,
    yearlyrainin real DEFAULT 0.0 NOT NULL,
    battlightning integer DEFAULT 0 NOT NULL,
    lightningday integer DEFAULT 0 NOT NULL,
    lightninghour integer DEFAULT 0 NOT NULL,
    lightningtime timestamp(0) without time zone DEFAULT NULL::timestamp without time zone NOT NULL,
    lightningdistance real DEFAULT 0.0 NOT NULL,
    aqipm25 integer DEFAULT 0,
    aqipm2524h integer DEFAULT 0,
    lightningmonth integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.records OWNER TO ambient;

--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.schema_migrations (
    name text NOT NULL,
    applied_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.schema_migrations OWNER TO ambient;

--
-- Name: spc_outlook_history; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.spc_outlook_history (
    id bigint NOT NULL,
    location text NOT NULL,
    outlook_day integer NOT NULL,
    layer text NOT NULL,
    in_area boolean NOT NULL,
    label text,
    severity double precision DEFAULT 0 NOT NULL,
    checked_at timestamp with time zone NOT NULL
);


ALTER TABLE public.spc_outlook_history OWNER TO ambient;

--
-- Name: spc_outlook_history_id_seq; Type: SEQUENCE; Schema: public; Owner: ambient
--

CREATE SEQUENCE public.spc_outlook_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.spc_outlook_history_id_seq OWNER TO ambient;

--
-- Name: spc_outlook_history_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: ambient
--

ALTER SEQUENCE public.spc_outlook_history_id_seq OWNED BY public.spc_outlook_history.id;


--
-- Name: spc_outlook_status; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.spc_outlook_status (
    location text NOT NULL,
    lat double precision NOT NULL,
    lon double precision NOT NULL,
    outlook_day integer NOT NULL,
    layer text NOT NULL,
    in_area boolean NOT NULL,
    label text,
    severity double precision DEFAULT 0 NOT NULL,
    checked_at timestamp with time zone NOT NULL
);


ALTER TABLE public.spc_outlook_status OWNER TO ambient;

--
-- Name: stat_ingest_log; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.stat_ingest_log (
    recorded timestamp with time zone NOT NULL,
    folded_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.stat_ingest_log OWNER TO ambient;

--
-- Name: vcforecast; Type: TABLE; Schema: public; Owner: ambient
--

CREATE TABLE public.vcforecast (
    datetime timestamp with time zone NOT NULL,
    datetime_epoch bigint,
    tempmax numeric(5,2),
    tempmin numeric(5,2),
    temp numeric(5,2),
    feelslikemax numeric(5,2),
    feelslikemin numeric(5,2),
    feelslike numeric(5,2),
    dew numeric(5,2),
    humidity numeric(5,2),
    precip numeric(5,3),
    precipprob numeric(5,2),
    precipcover numeric(5,2),
    preciptype text,
    snow numeric(5,2),
    snowdepth numeric(5,2),
    windgust numeric(5,2),
    windspeed numeric(5,2),
    winddir numeric(5,2),
    pressure numeric(6,2),
    cloudcover numeric(5,2),
    visibility numeric(5,2),
    solarradiation numeric(6,2),
    solarenergy numeric(5,2),
    uvindex integer,
    severerisk integer,
    sunrise time without time zone,
    sunrise_epoch bigint,
    sunset time without time zone,
    sunset_epoch bigint,
    moonphase numeric(3,2),
    conditions text,
    description text,
    icon text,
    stations text,
    source text,
    hours jsonb,
    summary text
);


ALTER TABLE public.vcforecast OWNER TO ambient;

--
-- Name: conditions id; Type: DEFAULT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.conditions ALTER COLUMN id SET DEFAULT nextval('public.conditions_id_seq'::regclass);


--
-- Name: forecast_data id; Type: DEFAULT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.forecast_data ALTER COLUMN id SET DEFAULT nextval('public.forecast_data_id_seq'::regclass);


--
-- Name: ingest_log id; Type: DEFAULT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.ingest_log ALTER COLUMN id SET DEFAULT nextval('public.ingest_log_id_seq'::regclass);


--
-- Name: map_locations id; Type: DEFAULT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.map_locations ALTER COLUMN id SET DEFAULT nextval('public.map_locations_id_seq'::regclass);


--
-- Name: spc_outlook_history id; Type: DEFAULT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.spc_outlook_history ALTER COLUMN id SET DEFAULT nextval('public.spc_outlook_history_id_seq'::regclass);


--
-- Name: alerts alerts_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.alerts
    ADD CONSTRAINT alerts_pkey PRIMARY KEY (alertid);


--
-- Name: conditions conditions_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.conditions
    ADD CONSTRAINT conditions_pkey PRIMARY KEY (id);


--
-- Name: conditions conditions_station_observed_uk; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.conditions
    ADD CONSTRAINT conditions_station_observed_uk UNIQUE (station, observed_at);


--
-- Name: forecast_data forecast_data_ident_key; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.forecast_data
    ADD CONSTRAINT forecast_data_ident_key UNIQUE (model_type, location_id, model_run_time, forecast_hour);


--
-- Name: forecast_data forecast_data_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.forecast_data
    ADD CONSTRAINT forecast_data_pkey PRIMARY KEY (id);


--
-- Name: forecast forecast_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.forecast
    ADD CONSTRAINT forecast_pkey PRIMARY KEY (location_id, datetime);


--
-- Name: ingest_log ingest_log_model_type_model_run_time_forecast_hour_key; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.ingest_log
    ADD CONSTRAINT ingest_log_model_type_model_run_time_forecast_hour_key UNIQUE (model_type, model_run_time, forecast_hour);


--
-- Name: ingest_log ingest_log_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.ingest_log
    ADD CONSTRAINT ingest_log_pkey PRIMARY KEY (id);


--
-- Name: map_locations map_locations_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.map_locations
    ADD CONSTRAINT map_locations_pkey PRIMARY KEY (id);


--
-- Name: period_stats period_stats_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.period_stats
    ADD CONSTRAINT period_stats_pkey PRIMARY KEY (period, field);


--
-- Name: records records_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.records
    ADD CONSTRAINT records_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (name);


--
-- Name: spc_outlook_history spc_outlook_history_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.spc_outlook_history
    ADD CONSTRAINT spc_outlook_history_pkey PRIMARY KEY (id);


--
-- Name: spc_outlook_status spc_outlook_status_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.spc_outlook_status
    ADD CONSTRAINT spc_outlook_status_pkey PRIMARY KEY (location, outlook_day, layer);


--
-- Name: stat_ingest_log stat_ingest_log_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.stat_ingest_log
    ADD CONSTRAINT stat_ingest_log_pkey PRIMARY KEY (recorded);


--
-- Name: vcforecast vcforecast_pkey; Type: CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.vcforecast
    ADD CONSTRAINT vcforecast_pkey PRIMARY KEY (datetime);


--
-- Name: baromrelinidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX baromrelinidx ON public.records USING btree (baromrelin, recorded);


--
-- Name: dailyraininidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX dailyraininidx ON public.records USING btree (dailyrainin, recorded);


--
-- Name: date; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX date ON public.records USING btree (recorded);


--
-- Name: datetempfidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX datetempfidx ON public.records USING btree (recorded, tempf);


--
-- Name: dewpointidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX dewpointidx ON public.records USING btree (dewpoint, recorded);


--
-- Name: feelsidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX feelsidx ON public.records USING btree (feelslike, recorded);


--
-- Name: humidity1idx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX humidity1idx ON public.records USING btree (humidity1, recorded);


--
-- Name: humidity2idx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX humidity2idx ON public.records USING btree (humidity2, recorded);


--
-- Name: humidity3idx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX humidity3idx ON public.records USING btree (humidity3, recorded);


--
-- Name: humidity4idx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX humidity4idx ON public.records USING btree (humidity4, recorded);


--
-- Name: humidityidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX humidityidx ON public.records USING btree (humidity, recorded);


--
-- Name: humidityinidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX humidityinidx ON public.records USING btree (humidityin, recorded);


--
-- Name: idx_conditions_observed; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX idx_conditions_observed ON public.conditions USING btree (station, observed_at DESC);


--
-- Name: idx_forecast_data_lookup; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX idx_forecast_data_lookup ON public.forecast_data USING btree (location_id, model_type, valid_time);


--
-- Name: idx_forecast_data_run_hour_loc; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX idx_forecast_data_run_hour_loc ON public.forecast_data USING btree (model_type, model_run_time, forecast_hour, location_id);


--
-- Name: idx_forecast_datetime; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX idx_forecast_datetime ON public.forecast USING btree (datetime);


--
-- Name: idx_forecast_run; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX idx_forecast_run ON public.forecast_data USING btree (model_type, model_run_time);


--
-- Name: idx_forecast_stale_summary; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX idx_forecast_stale_summary ON public.forecast USING btree (location_id, datetime) WHERE ((summary_hash IS NULL) OR (summary_hash <> content_hash));


--
-- Name: idx_period_field; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX idx_period_field ON public.period_stats USING btree (period, field);


--
-- Name: idx_spc_hist_time; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX idx_spc_hist_time ON public.spc_outlook_history USING btree (checked_at);


--
-- Name: idx_spc_status_inarea; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX idx_spc_status_inarea ON public.spc_outlook_status USING btree (in_area) WHERE in_area;


--
-- Name: lightningdayidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX lightningdayidx ON public.records USING btree (lightningday, recorded);


--
-- Name: map_locations_attention; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX map_locations_attention ON public.map_locations USING btree (last_viewed_at DESC) WHERE forecast_enabled;


--
-- Name: map_locations_name_key; Type: INDEX; Schema: public; Owner: ambient
--

CREATE UNIQUE INDEX map_locations_name_key ON public.map_locations USING btree (lower(name));


--
-- Name: recorded; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX recorded ON public.records USING btree (recorded DESC NULLS LAST);


--
-- Name: temp1fidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX temp1fidx ON public.records USING btree (temp1f, recorded);


--
-- Name: temp2fidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX temp2fidx ON public.records USING btree (temp2f, recorded);


--
-- Name: temp3fidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX temp3fidx ON public.records USING btree (temp3f, recorded);


--
-- Name: temp4fidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX temp4fidx ON public.records USING btree (temp4f, recorded);


--
-- Name: tempfidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX tempfidx ON public.records USING btree (tempf, recorded);


--
-- Name: tempinfidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX tempinfidx ON public.records USING btree (tempinf, recorded);


--
-- Name: uvidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX uvidx ON public.records USING btree (uv, recorded);


--
-- Name: windgustidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX windgustidx ON public.records USING btree (windgustmph, recorded);


--
-- Name: windgustmphidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX windgustmphidx ON public.records USING btree (windgustmph, recorded);


--
-- Name: windidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX windidx ON public.records USING btree (windspeedmph, recorded);


--
-- Name: windspeedmphidx; Type: INDEX; Schema: public; Owner: ambient
--

CREATE INDEX windspeedmphidx ON public.records USING btree (windspeedmph, recorded);


--
-- Name: forecast_data forecast_data_location_fk; Type: FK CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.forecast_data
    ADD CONSTRAINT forecast_data_location_fk FOREIGN KEY (location_id) REFERENCES public.map_locations(id) ON DELETE CASCADE;


--
-- Name: forecast forecast_location_fk; Type: FK CONSTRAINT; Schema: public; Owner: ambient
--

ALTER TABLE ONLY public.forecast
    ADD CONSTRAINT forecast_location_fk FOREIGN KEY (location_id) REFERENCES public.map_locations(id) ON DELETE CASCADE;


