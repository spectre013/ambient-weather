# sdr

Receives Ambient Weather / Fine Offset / Ecowitt sensors directly with an
RTL-SDR dongle and publishes each reading as JSON on stdout and MQTT. It
replaces `rtl_433` in "Option B" of the weather README: the messages have the
same topics, model names and field names, so the processor consumes them
unchanged.

```
rtl_sdr ──IQ──▶ demod ──39-byte payload──▶ fineoffset ──▶ JSON ──▶ stdout / MQTT
 (dongle)   FSK, sync word, framing     console firmware's       ambient/<model>/<id>
            (what the CMT2300A does)    checks and decoders
```

| Part | What it is |
|---|---|
| `fineoffset/` | The sensor decoders from the HP2550 console firmware reconstruction (package `hp2550/rf`): CRC-8 / sum-8 checks and one decoder per sensor family, in the firmware's native units. Copied, not imported, so this module stands alone. |
| `demod/` | Software replacement for the console's radio chip: channel filter, FM discriminator, preamble + sync word (`AA AA 2D D4`) search, bit-clock recovery, 39-byte framing. Framing parameters come from the firmware's CMT2300A register table. |
| `message.go` | Maps a decoded measurement to rtl_433-compatible JSON. |
| `main.go`, `source.go`, `mqtt.go` | Configuration, the `rtl_sdr` subprocess, MQTT publishing. |

## Run

```sh
docker compose up -d --build     # in this directory: builds and prints readings
docker compose logs -f
```

In production it is the `sdr` service of `../docker-compose.yml`, enabled
with `COMPOSE_PROFILES=sdr`.

The container needs the dongle (`/dev/bus/usb`, see `docker-compose.yml`).
`rtl_sdr` detaches the kernel's DVB driver on start and hands the dongle back
on a clean stop, so nothing has to be blacklisted on the host. Only one
program can use a dongle at a time.

Every setting is a flag with an environment-variable default:

| Flag | Env | Default | |
|---|---|---|---|
| `-input` | `SDR_INPUT` | `rtl` | `rtl` runs `rtl_sdr`; `-` reads cu8 samples on stdin; anything else is a `.cu8` file to replay |
| `-freq` | `SDR_FREQ` | `915` | a console band (`433`, `868`, `915`, `920`: the firmware's 433.925 / 868.35 / 915.0 / 920.9 MHz) or e.g. `915M` |
| `-rate` | `SDR_SAMPLE_RATE` | `250k` | sample rate |
| `-gain` | `SDR_GAIN` | `0` | tuner gain in dB, `0` = automatic |
| `-device` | `SDR_DEVICE` | `0` | dongle index or serial |
| `-ppm` | `SDR_PPM` | `0` | dongle frequency correction |
| `-watchdog` | `SDR_WATCHDOG` | `5m` | exit (so Docker restarts it) after this long without a valid packet; `0` = off |
| `-sensors` | `SDR_SENSORS` | *(all)* | allow-list, comma-separated `<model>/<id>` (the MQTT topic pair), e.g. `Fineoffset-WS80/34306,AmbientWeather-WH31E/156` |
| `-repeats` | `SDR_REPEATS` | off | also output the second copy of each packet |
| `-raw` | `SDR_RAW` | off | add the packet bytes (`payload`, hex) to each message |
| `-v` | `SDR_VERBOSE` | off | log packets that fail their checks |
| `-mqtt` | `MQTT_BROKER` | *(unset)* | e.g. `tcp://host:1883`; unset = stdout only |
| `-topic` | `MQTT_TOPIC_PREFIX` | `ambient` | messages go to `<prefix>/<model>/<id>` |
| | `MQTT_USER`, `MQTT_PASS`, `MQTT_CLIENT_ID` | | client id defaults to `weather-sdr` |

Each sensor is logged the first time it is heard (`sensor
Fineoffset-WS80/34306: receiving`, or `not in SDR_SENSORS, ignoring it`),
which is how to find the IDs for the allow-list. A sensor picks a new ID
when its batteries are replaced, so it then has to be re-listed. This
replaces the console's ID learning, which registers the first sensor of each
type it hears.

Logs go to stderr, readings to stdout. A summary line is logged every 10
minutes (sync words found, readings, repeats, failed checks).

## Output

```json
{"time":"2026-11-16 16:45:34","model":"Fineoffset-WS80","id":34306,"battery_ok":1.1,"battery_mV":3160,"temperature_C":7.9,"humidity":67,"wind_dir_deg":302,"wind_avg_m_s":0.6,"wind_max_m_s":0.8,"uvi":0,"light_lux":490,"flags":161,"battery_low":0,"mic":"CRC","mod":"FSK","freq1":915.035,"freq2":914.964,"rssi":-1.4,"snr":20.4,"noise":-21.8}
```

Differences from rtl_433, all additions or corrections:

- `battery_low` (0/1) on every message is the console's own low-battery
  decision, using its per-sensor thresholds.
- A value the console treats as "absent" is left out instead of being
  published as a number. In particular a WH0290 with no reading has no
  `pm2_5_ug_m3` key rather than rtl_433's `1638`.
- WH0290 PM values are rounded the way the console rounds them
  (`(raw+5)/10`) rather than truncated.
- Sensors transmit each packet twice, about 25 ms apart. Only the first copy
  is output unless `-repeats` is set.
- The WH40 always includes `battery_V` / `battery_ok` when the sensor reports
  them.

`rssi`, `snr` and `noise` are dB relative to full scale. With automatic gain a
nearby sensor clips the dongle and reads slightly above 0; that is harmless
for FSK.

## Supported sensors

Confirmed off the air against rtl_433 (identical values on every shared
field, 150 s capture):

| Model in JSON | Sensor | Firmware family byte |
|---|---|---|
| `Fineoffset-WS80` | ultrasonic array | `0x80` |
| `Fineoffset-WH32B` | indoor T/H/barometer | `0xE_` |
| `AmbientWeather-WH31E` | T/H channels | `0x30` |
| `EcoWitt-WH40` | rain gauge | `0x40` |
| `FineOffset-WH31L` | lightning (WH57) | `0x57` |

Received live and passing its checks, but not compared field by field (it
was not in the capture, and it is currently sending "no reading" with a low
battery): `Fineoffset-WH0290` (PM2.5, WH41, family `0x41`–`0x44`).

Decoded from the firmware's layouts but **not checked against a real sensor**
(none in range): `Fineoffset-WH24` / `WH65B`, `EcoWitt-WS68`,
`Fineoffset-WS90`, `-WS85`, `-WH32` (outdoor, `0xD_`), `AmbientWeather-WH31B`
(`0x37`), `Fineoffset-WH51`, `-WH45`, `-WH46`, `-WN34`, `-WN35`, `-WH55`,
`-WH48`, `-WH54`. Their field names are a best effort; see
the firmware notes (`re/rf.md` in the hp2550 reconstruction) for which fields are still unidentified.

Left out on purpose: the console's signal-quality bars and its indoor-sensor "close to the previous reading" filter. The processor
already applies that filter to the WH32B. Sensor-ID learning is replaced by
the `SDR_SENSORS` allow-list.

## Development

Nothing needs to be installed on the host; `run-dev.sh` runs a command in a
container with the Go toolchain, `rtl_sdr` and `rtl_433`:

```sh
docker build -t sdr-dev -f Dockerfile.dev .
./run-dev.sh go test ./...
./run-dev.sh sh -c 'go build -o sdr . && ./sdr -input captures/cap1_915M_250k.cu8 -raw'

# record 60 s from the dongle (stop the receiver first)
DOCKER_ARGS="--device=/dev/bus/usb/001/006" ./run-dev.sh \
    rtl_sdr -f 915000000 -s 250000 -g 0 -n 15000000 captures/new_915M_250k.cu8
```

Tests: `fineoffset/` carries the firmware-derived vectors; `message_test.go`
has packets received off the air with the values rtl_433 printed for them;
`demod/testdata/*.cu8` are real bursts cut from a capture; the rest of the
demodulator tests synthesise FSK with frequency offset, clock error, noise
and clipping.

Sensitivity, measured by attenuating the capture and adding noise (66
transmissions in the file):

| added noise (sigma) | rtl_433 | sdr |
|---|---|---|
| 12 | 61–63 | 66 |
| 15 | 29–33 | 63 |
| 18 | 1 | 58 |
| 21 | 0 | 46 |

CPU: about 1.6 % of one core at 250 kS/s.
