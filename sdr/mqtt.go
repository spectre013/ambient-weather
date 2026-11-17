package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttConfig struct {
	Broker   string // e.g. tcp://192.168.1.10:1883
	ClientID string
	Username string
	Password string
	Prefix   string // topic prefix; messages go to <prefix>/<model>/<id>
}

// publisher sends each message to <prefix>/<model>/<id>, the layout rtl_433
// produces with "events=ambient[/model][/id]".
type publisher struct {
	client mqtt.Client
	prefix string
}

func newPublisher(cfg mqttConfig) (*publisher, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker).
		SetClientID(cfg.ClientID).
		SetAutoReconnect(true).
		SetKeepAlive(30 * time.Second).
		SetPingTimeout(10 * time.Second).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second).
		SetCleanSession(true)
	if cfg.Username != "" {
		opts.SetUsername(cfg.Username)
	}
	if cfg.Password != "" {
		opts.SetPassword(cfg.Password)
	}
	opts.SetOnConnectHandler(func(mqtt.Client) { log.Printf("mqtt: connected to %s", cfg.Broker) })
	opts.SetConnectionLostHandler(func(_ mqtt.Client, err error) { log.Printf("mqtt: connection lost: %v", err) })

	c := mqtt.NewClient(opts)
	// With ConnectRetry the token only completes once connected; do not
	// hold up reception for a broker that is down.
	if t := c.Connect(); t.WaitTimeout(10*time.Second) && t.Error() != nil {
		return nil, fmt.Errorf("mqtt connect: %w", t.Error())
	}
	return &publisher{client: c, prefix: cfg.Prefix}, nil
}

// Publish is fire-and-forget: a reading that cannot be delivered is
// superseded by the next one within seconds.
func (p *publisher) Publish(m *Message) {
	if !p.client.IsConnected() {
		return
	}
	body, err := json.Marshal(m)
	if err != nil {
		log.Printf("mqtt: marshal: %v", err)
		return
	}
	topic := fmt.Sprintf("%s/%s/%d", p.prefix, m.Model, m.ID)
	t := p.client.Publish(topic, 0, false, body)
	go func() {
		if t.WaitTimeout(5*time.Second) && t.Error() != nil {
			log.Printf("mqtt: publish %s: %v", topic, t.Error())
		}
	}()
}

func (p *publisher) Close() { p.client.Disconnect(250) }
