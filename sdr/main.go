// Command sdr receives Fine Offset / Ambient Weather sensors with an RTL-SDR
// dongle and publishes each reading as JSON on stdout and, optionally, MQTT.
//
// Pipeline: rtl_sdr (IQ samples) -> demod (FSK, sync word, framing, the job
// of the console's CMT2300A) -> fineoffset (the console firmware's checks
// and per-sensor decoders) -> rtl_433-compatible JSON.
package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"git.zoms.net/apps/weather/sdr/demod"
	fo "git.zoms.net/apps/weather/sdr/fineoffset"
)

// bandHz are the console's four frequency settings (firmware table 0xcbfdc,
// Init-file key "Frequency").
var bandHz = map[string]float64{"433": 433.925e6, "868": 868.35e6, "915": 915e6, "920": 920.9e6}

type config struct {
	input    string
	device   string
	freq     float64
	rate     float64
	gain     string
	ppm      int
	raw      bool
	repeats  bool
	sensors  *sensorFilter
	verbose  bool
	watchdog time.Duration
	mqtt     mqttConfig
}

func env(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return def
}

// parseHz accepts "915000000", "915M", "433.92M", "250k" or a band name
// from bandHz ("915").
func parseHz(s string) (float64, error) {
	if hz, ok := bandHz[s]; ok {
		return hz, nil
	}
	mult := 1.0
	switch {
	case strings.HasSuffix(s, "M"):
		mult, s = 1e6, strings.TrimSuffix(s, "M")
	case strings.HasSuffix(s, "k"):
		mult, s = 1e3, strings.TrimSuffix(s, "k")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil || v <= 0 {
		return 0, fmt.Errorf("bad frequency %q", s)
	}
	return v * mult, nil
}

func loadConfig() (config, error) {
	var c config
	var freq, rate, wd, sensors string
	flag.StringVar(&c.input, "input", env("SDR_INPUT", "rtl"), `sample source: "rtl" (run rtl_sdr), "-" (cu8 on stdin) or a .cu8 file`)
	flag.StringVar(&c.device, "device", env("SDR_DEVICE", "0"), "rtl_sdr device index or serial")
	flag.StringVar(&freq, "freq", env("SDR_FREQ", "915"), "centre frequency: a console band (433, 868, 915, 920) or e.g. 915M")
	flag.StringVar(&rate, "rate", env("SDR_SAMPLE_RATE", "250k"), "sample rate")
	flag.StringVar(&c.gain, "gain", env("SDR_GAIN", "0"), "tuner gain in dB, 0 = automatic")
	flag.IntVar(&c.ppm, "ppm", 0, "dongle frequency correction in ppm (env SDR_PPM)")
	flag.BoolVar(&c.raw, "raw", env("SDR_RAW", "") != "", "add the raw payload (hex) to each message")
	flag.BoolVar(&c.repeats, "repeats", env("SDR_REPEATS", "") != "", "also output the repeat transmission of each packet")
	flag.StringVar(&sensors, "sensors", env("SDR_SENSORS", ""), "only output these sensors: comma-separated <model>/<id>, e.g. Fineoffset-WS80/34306 (empty = all)")
	flag.BoolVar(&c.verbose, "v", env("SDR_VERBOSE", "") != "", "log packets that fail their checks")
	flag.StringVar(&wd, "watchdog", env("SDR_WATCHDOG", "5m"), "exit if no valid packet arrives for this long (live input only, 0 = off)")
	flag.StringVar(&c.mqtt.Broker, "mqtt", env("MQTT_BROKER", ""), "MQTT broker URL, e.g. tcp://host:1883 (empty = stdout only)")
	flag.StringVar(&c.mqtt.Prefix, "topic", env("MQTT_TOPIC_PREFIX", "ambient"), "MQTT topic prefix; messages go to <prefix>/<model>/<id>")
	flag.Parse()

	c.mqtt.ClientID = env("MQTT_CLIENT_ID", "weather-sdr")
	c.mqtt.Username = env("MQTT_USER", "")
	c.mqtt.Password = env("MQTT_PASS", "")

	var err error
	if c.freq, err = parseHz(freq); err != nil {
		return c, err
	}
	if c.rate, err = parseHz(rate); err != nil {
		return c, err
	}
	if c.sensors, err = parseSensors(sensors); err != nil {
		return c, err
	}
	if c.watchdog, err = time.ParseDuration(wd); err != nil {
		return c, fmt.Errorf("bad watchdog %q", wd)
	}
	if v := os.Getenv("SDR_PPM"); v != "" && c.ppm == 0 {
		if c.ppm, err = strconv.Atoi(v); err != nil {
			return c, fmt.Errorf("bad SDR_PPM %q", v)
		}
	}
	return c, nil
}

type stats struct {
	synced, valid, repeats, filtered, badCheck, unknown, rejected int
}

func (s stats) String() string {
	return fmt.Sprintf("%d sync words: %d readings, %d repeats, %d from unlisted sensors, %d failed checks, %d unknown type, %d rejected",
		s.synced, s.valid, s.repeats, s.filtered, s.badCheck, s.unknown, s.rejected)
}

// repeatWindow: sensors send every packet twice about 25 ms apart (the WH40
// 50 ms). An identical payload inside this window is the same reading.
const repeatWindow = 0.5 // seconds

func main() {
	log.SetFlags(log.LstdFlags)
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := run(ctx, cfg); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}

func run(ctx context.Context, cfg config) error {
	src, live, err := openSource(ctx, cfg)
	if err != nil {
		return err
	}
	defer src.Close()

	var pub *publisher
	if cfg.mqtt.Broker != "" {
		if pub, err = newPublisher(cfg.mqtt); err != nil {
			return err
		}
		defer pub.Close()
	}

	d := demod.New(demod.Config{SampleRate: cfg.rate})
	out := json.NewEncoder(os.Stdout)
	var st stats
	lastValid := time.Now()
	lastReport := time.Now()
	buf := make([]byte, 1<<16)
	var lastPayload []byte
	var lastSample int64

	for {
		n, err := src.Read(buf)
		for _, p := range d.Write(buf[:n]) {
			st.synced++
			msg, derr := decode(cfg, p, live)
			if derr != nil {
				switch {
				case errors.Is(derr, fo.ErrUnknownFamily):
					st.unknown++
				case errors.Is(derr, fo.ErrRejected):
					st.rejected++
				default:
					st.badCheck++
				}
				if cfg.verbose {
					log.Printf("dropped: %v: %s (rssi %.1f dB, snr %.1f dB)", derr, hex.EncodeToString(p.Payload), p.RSSI, p.SNR)
				}
				continue
			}
			lastValid = time.Now()
			_, plen, _ := fo.Layout(p.Payload)
			repeat := bytes.Equal(lastPayload, p.Payload[:plen]) &&
				float64(p.Sample-lastSample) < repeatWindow*cfg.rate
			lastPayload, lastSample = p.Payload[:plen], p.Sample
			if repeat && !cfg.repeats {
				st.repeats++
				continue
			}
			if !cfg.sensors.Allow(msg) {
				st.filtered++
				continue
			}
			st.valid++
			if err := out.Encode(msg); err != nil {
				return err
			}
			if pub != nil {
				pub.Publish(msg)
			}
		}
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			if errors.Is(err, io.EOF) && !live {
				log.Printf("end of input: %s", st)
				return nil
			}
			return fmt.Errorf("sample source: %w", err)
		}
		if live {
			if cfg.watchdog > 0 && time.Since(lastValid) > cfg.watchdog {
				// The console re-initialises its radio after 60 s of
				// silence; here the equivalent is to exit and let the
				// supervisor restart rtl_sdr.
				return fmt.Errorf("no valid packet for %s, giving up", cfg.watchdog)
			}
			if time.Since(lastReport) >= 10*time.Minute {
				log.Printf("last 10 min: %s", st)
				st, lastReport = stats{}, time.Now()
			}
		}
	}
}

// decode checks and decodes one framed payload and builds its message.
func decode(cfg config, p demod.Packet, live bool) (*Message, error) {
	meas, err := fo.Decode(p.Payload)
	if err != nil {
		return nil, err
	}
	msg := NewMessage(meas, p.Payload)
	ts := fmt.Sprintf("@%.6fs", float64(p.Sample)/cfg.rate)
	if live {
		ts = time.Now().Format("2006-01-02 15:04:05")
	}
	msg.fields = append([]field{{"time", ts}}, msg.fields...)
	if cfg.raw {
		_, n, _ := fo.Layout(p.Payload)
		msg.add("payload", hex.EncodeToString(p.Payload[:n]))
	}
	msg.add("mod", "FSK").
		add("freq1", round((cfg.freq+p.FreqHi)/1e6, 3)).
		add("freq2", round((cfg.freq+p.FreqLo)/1e6, 3)).
		add("rssi", round(p.RSSI, 1)).
		add("snr", round(p.SNR, 1)).
		add("noise", round(p.Noise, 1))
	return msg, nil
}
