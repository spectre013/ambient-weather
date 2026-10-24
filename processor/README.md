# weather-mqtt

A Go service that consumes weather sensor readings from an MQTT broker, maintains
a current in-memory snapshot of every sensor, and writes a single aggregated row
to PostgreSQL once per minute.

## Sensors

| MQTT model name        | Role                                  | Notes |
|------------------------|---------------------------------------|-------|
| `Fineoffset-WS80`      | Outdoor (temp, humidity, wind, UV, light) | ~4s cadence |
| `Fineoffset-WH32B`     | Indoor / barometer                    | ~30s cadence |
| `AmbientWeather-WH31E` | Indoor multi-channel temp/hum         | ch1=basement, ch2=Hannah's room, ch3=master, ch4=garage |
| `Fineoffset-WH0290`    | Air quality (PM2.5 / PM10)            | ~30s cadence |
| `EcoWitt-WH40`         | Rain (cumulative `rain_mm`)           | aggregated hourly/daily/weekly/monthly/yearly |
| `FineOffset-WH31L`     | Lightning (cumulative `strike_count`) | aggregated hourly/daily/monthly |

## Aggregation strategy

Both the rain gauge (`rain_mm`) and the lightning sensor (`strike_count`) report
running counters. As in the HP2550 console firmware, each message adds the
increment since the previous message to every period total (day / week /
month / year, and hour for lightning), and each total is zeroed when the clock
crosses its boundary. Increments are taken modulo the counter width, so the
rain counter wrapping at 6553.5 mm or the strike counter wrapping at 255 still
counts. An increment larger than the console's limit (rain: 0.4 mm per elapsed
second, at most 102.3 mm; lightning: 100 strikes) is treated as a battery-swap
reset or bad packet and skipped, leaving the totals intact.

Rain also feeds a per-minute window that gives the console-style rolling values:

- `hourlyrainin` / `hourlyrain`: rain in the last 60 minutes (rolling, not
  since the top of the hour).
- `rainratein`: rain rate in in/hr, the last 10 minutes x 6.
- `eventrainin`: accumulates across showers and keeps showing after the rain
  stops. At each hour boundary it resets to 0 if the last 60 minutes were dry
  and the last 24 hours had less than 1.0 mm.

At startup the window is rebuilt from the last 24 hours of `dailyrainin`
values, so the rolling values are right immediately after a restart.

The lightning sensor's `strike_count` advances on noise too. A single-count
increase tagged `interference` or `noise` is not counted; larger jumps are
counted because the state only describes the latest event.

`lightningtime` and `lightningdistance` are only updated on messages whose
state is `strike`.

The WH32B is only protected by an 8-bit sum in rtl_433, so a reading is
accepted only if it is within 0.5 °C, 3 % and 1.0 hPa of the previous packet
(the console's filter). A genuine change is accepted one packet later.

## Configuration

All settings come from environment variables. On startup the service also
loads a `.env` file from the working directory if one exists (override the
path with `ENV_FILE`). Variables already set in the process environment take
precedence over the `.env` file, so production deployments can keep their
secrets in the orchestrator's env and use `.env` only for local development.
A starter template is included as `.env.example` — copy it to `.env` and edit.

| Env var             | Default                                                              |
|---------------------|----------------------------------------------------------------------|
| `ENV_FILE`          | `.env` — alternate path for the dotenv file                          |
| `MQTT_BROKER`       | `tcp://localhost:1883`                                               |
| `MQTT_CLIENT_ID`    | `weather-ingest`                                                     |
| `MQTT_USER`         | *(unset)*                                                            |
| `MQTT_PASS`         | *(unset)*                                                            |
| `MQTT_TOPIC`        | `ambient/#`                                                          |
| `STATION_MAC`       | `WS-5000`                                                            |
| `DATABASE_URL`      | `postgres://weather:weather@localhost:5432/weather?sslmode=disable`  |
| `SNAPSHOT_INTERVAL` | `1m` (parsed with `time.ParseDuration`)                              |
| `ELEVATION`         | `0` — station elevation above sea level in meters (for `baromrelin`) |

## Build & run

```sh
go mod tidy
go build -o weather-mqtt .

# Local dev: edit .env then just run the binary.
cp .env.example .env
./weather-mqtt

# Or set vars directly:
MQTT_BROKER=tcp://mqtt.local:1883 \
  DATABASE_URL='postgres://weather:secret@db.local:5432/weather?sslmode=disable' \
  ./weather-mqtt
```

The `records` table must already exist before the service starts — schema and
indexing are managed externally.

## Notes & caveats

- `feelslike` uses the full NWS heat-index procedure (Steadman, then
  Rothfusz with its humidity adjustments) above 80°F and NWS wind-chill below
  50°F. In dry heat the heat index can be below the air temperature.
- `solarradiation` is a rough lux→W/m² conversion (`lux / 126.7`, Davis convention).
- `baromrelin` is computed from `baromabsin` using the hypsometric formula
  with `ELEVATION` (meters) and the current outdoor temperature. If
  `ELEVATION=0` (the default) the relative reading equals the absolute. If
  the outdoor sensor hasn't reported yet, the formula falls back to ICAO
  standard temperature (15°C) for the correction.
- `aqi24` is computed from a rolling 24-hour buffer of PM2.5 samples held in
  memory — it warms up across the first day of operation, but `aqi`/`aqi24`
  are seeded from the last DB row at startup so they aren't zero during warmup.
- At startup the entire in-memory state is seeded from the most recent
  `records` row via `Store.LatestRecord`: outdoor temp/humidity/wind,
  indoor barometer, every WH31E channel, AQI, today's max gust, and rain /
  lightning period totals. As real MQTT messages arrive they overwrite the
  seeded values. Rain/lightning period seeds that fall in a previous period
  (e.g. last record was yesterday) are correctly discarded so the new period
  starts at zero.
- The periodic snapshot only writes a row when new sensor data has arrived
  since the previous write. If the MQTT feed stalls (broker down, antenna
  unplugged), the service logs `no new sensor data since last write, skipping
  insert` rather than inserting duplicate rows — so the gap in `recorded`
  timestamps directly reflects the outage.