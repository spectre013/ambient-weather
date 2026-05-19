package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTConfig holds the broker connection settings.
type MQTTConfig struct {
	Broker    string // e.g. tcp://192.168.1.10:1883 or ssl://host:8883
	ClientID  string
	Username  string
	Password  string
	Topic     string // e.g. "ambient/#"
	TLSConfig *tls.Config
}

// MessageHandler is invoked for every successfully-decoded sensor message.
type MessageHandler func(m *SensorMessage)

// StartMQTT connects to the broker and subscribes to the configured topic.
// Returns the client so callers can disconnect cleanly on shutdown.
func StartMQTT(cfg MQTTConfig, handler MessageHandler) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID).
		SetAutoReconnect(true).
		SetKeepAlive(30 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetCleanSession(true).
		SetOrderMatters(false)

	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}
	if cfg.TLSConfig != nil {
		opts.SetTLSConfig(cfg.TLSConfig)
	}

	opts.SetOnConnectHandler(func(c mqtt.Client) {
		log.Printf("mqtt: connected, subscribing to %q", cfg.Topic)
		if t := c.Subscribe(cfg.Topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
			var m SensorMessage
			if err := json.Unmarshal(msg.Payload(), &m); err != nil {
				log.Printf("mqtt: bad json on %s: %v", msg.Topic(), err)
				return
			}
			// Fall back to the topic if model wasn't in the body for some reason.
			if m.Model == "" {
				m.Model = topicToModel(msg.Topic())
			}
			handler(&m)
		}); t.Wait() && t.Error() != nil {
			log.Printf("mqtt: subscribe error: %v", t.Error())
		}
	})

	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) {
		log.Printf("mqtt: connection lost: %v", err)
	})

	client := mqtt.NewClient(opts)
	if t := client.Connect(); t.Wait() && t.Error() != nil {
		return nil, fmt.Errorf("mqtt connect: %w", t.Error())
	}
	return client, nil
}

// topicToModel pulls the model segment out of "ambient/<Model>/<id>" if needed.
func topicToModel(topic string) string {
	// Walk segments instead of importing path: we just want segment [1].
	seg := 0
	start := 0
	for i := 0; i <= len(topic); i++ {
		if i == len(topic) || topic[i] == '/' {
			if seg == 1 {
				return topic[start:i]
			}
			seg++
			start = i + 1
		}
	}
	return ""
}

// StatsNotification is the payload published to the stats topic after every
// successful database insert. Downstream services subscribe to this topic to
// trigger recomputation of derived sensor statistics for display.
type StatsNotification struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// PublishStats fires a notification on the configured stats topic. Errors are
// logged but not returned — a stats publish failure must not block the main
// ingest loop.
func PublishStats(client mqtt.Client, topic, id string, ts time.Time) {
	if client == nil || !client.IsConnected() {
		log.Printf("mqtt: stats publish skipped, client not connected")
		return
	}
	payload, err := json.Marshal(StatsNotification{ID: id, Timestamp: ts})
	if err != nil {
		log.Printf("mqtt: stats marshal: %v", err)
		return
	}
	// QoS 1 so the consuming service is guaranteed to receive it at least once
	// even if it's reconnecting at the moment we publish.
	t := client.Publish(topic, 1, false, payload)
	// Don't block the ingest loop — fire and forget, but log if the publish
	// fails immediately (broker rejects, queue full, etc.).
	go func() {
		if !t.WaitTimeout(5 * time.Second) {
			log.Printf("mqtt: stats publish to %s timed out", topic)
			return
		}
		if err := t.Error(); err != nil {
			log.Printf("mqtt: stats publish to %s: %v", topic, err)
		}
	}()
}
