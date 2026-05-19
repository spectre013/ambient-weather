package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from a .env file if one exists. Real env vars
	// already set in the process environment take precedence (godotenv.Load
	// does not overwrite existing values). The path can be overridden with
	// ENV_FILE — useful for containers that mount config as a secret.
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}
	if err := godotenv.Load(envFile); err != nil {
		// Not fatal — missing .env is the normal case in production where
		// vars come straight from the orchestrator's environment.
		log.Printf("dotenv: not loading %s (%v)", envFile, err)
	} else {
		log.Printf("dotenv: loaded %s", envFile)
	}

	broker := envOr("MQTT_BROKER", "tcp://localhost:1883")
	clientID := envOr("MQTT_CLIENT_ID", "weather-processor-test")
	username := os.Getenv("MQTT_USER")
	password := os.Getenv("MQTT_PASS")
	topic := envOr("MQTT_TOPIC", "ambient/#")
	statsTopic := envOr("MQTT_STATS_TOPIC", "weather/stats")
	stationID := envOr("STATION_MAC", "WS-5000")
	pgDSN := envOr("DATABASE_URL", "postgres://weather:weather@localhost:5432/weather?sslmode=disable")
	interval := envDuration("SNAPSHOT_INTERVAL", time.Minute)
	rainEventGap := envDuration("RAIN_EVENT_GAP", time.Hour)
	elevationM := envFloat("ELEVATION", 0)

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.Printf("weather-mqtt starting (broker=%s topic=%s interval=%s elevation=%gm)",
		broker, topic, interval, elevationM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to postgres first — no point reading messages if we can't store them.
	store, err := NewStore(ctx, pgDSN)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer store.Close()

	state := NewSensorState()
	state.rainEventGap = rainEventGap
	state.elevationM = elevationM

	// Seed running totals from the most recent DB row so rain/lightning
	// period counters stay continuous across restarts.
	if last, err := store.LatestRecord(ctx); err != nil {
		log.Printf("seed: no previous record (%v)", err)
	} else {
		log.Printf("seed: found last record at %s", last.Recorded.Format(time.RFC3339))
		state.Seed(last, time.Now())
	}

	// Connect to MQTT and route every message into state.Apply.
	client, err := StartMQTT(MQTTConfig{
		Broker:   broker,
		ClientID: clientID,
		Username: username,
		Password: password,
		Topic:    topic,
	}, func(m *SensorMessage) {
		state.Apply(m, time.Now())
	})
	if err != nil {
		log.Fatalf("mqtt: %v", err)
	}
	defer client.Disconnect(500)

	// Tick every minute and write a snapshot. The first tick is aligned to the
	// next minute boundary so all the rows have :00 timestamps.
	go runSnapshotLoop(ctx, state, store, stationID, interval, client, statsTopic)

	// Wait for SIGINT/SIGTERM.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh
	log.Printf("received %s, shutting down", sig)
	cancel()

	// One final snapshot before exit — but only if new data has arrived since
	// the last write.
	if final, ok := state.TrySnapshot(stationID, time.Now()); ok {
		if err := store.Insert(context.Background(), final); err != nil {
			log.Printf("final insert: %v", err)
		} else {
			PublishStats(client, statsTopic, final.Mac, final.Recorded)
		}
	}
}

func runSnapshotLoop(ctx context.Context, state *SensorState, store *Store, mac string, interval time.Duration, mqttClient mqtt.Client, statsTopic string) {
	// Align first tick to the top of the next interval (e.g. next minute).
	now := time.Now()
	next := now.Truncate(interval).Add(interval)
	first := time.NewTimer(next.Sub(now))
	defer first.Stop()

	select {
	case <-ctx.Done():
		return
	case <-first.C:
	}

	snapshotOnce(ctx, state, store, mac, mqttClient, statsTopic)

	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			snapshotOnce(ctx, state, store, mac, mqttClient, statsTopic)
		}
	}
}

func snapshotOnce(ctx context.Context, state *SensorState, store *Store, mac string, mqttClient mqtt.Client, statsTopic string) {
	r, ok := state.TrySnapshot(mac, time.Now())
	if !ok {
		log.Printf("snapshot: no new sensor data since last write, skipping insert")
		return
	}
	if err := store.Insert(ctx, r); err != nil {
		log.Printf("insert: %v", err)
		return
	}
	log.Printf("snapshot: tempf=%.1f hum=%d wind=%.1f gust=%.1f rain_today=%.3fin lightning_today=%d aqi=%d",
		r.Tempf, r.Humidity, r.Windspeedmph, r.Windgustmph,
		r.Dailyrainin, r.Lightningday, r.Aqipm25)
	// Notify downstream stats service that a new row landed.
	PublishStats(mqttClient, statsTopic, r.Mac, r.Recorded)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envDuration(key string, def time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("invalid %s=%q, using default %s: %v", key, v, def, err)
		return def
	}
	return d
}

func envFloat(key string, def float64) float64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.Printf("invalid %s=%q, using default %g: %v", key, v, def, err)
		return def
	}
	return f
}
