package main

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var logger = logrus.New()
var client mqtt.Client

func init() {
	logger.Out = os.Stdout
	logger.SetLevel(logrus.InfoLevel)
}
func main() {
	var err error
	if os.Getenv("GO_ENV") != "production" {
		err = godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	logLevel := logrus.InfoLevel
	if os.Getenv("LOGLEVEL") == "Debug" {
		logLevel = logrus.DebugLevel
	}
	logger.Info("Setting Debug Level to ", logLevel)
	logger.SetLevel(logLevel)

	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_HOST"))
	opts.SetClientID("go_publisher")
	opts.SetUsername(os.Getenv("MQTT_USERNAME"))
	opts.SetPassword(os.Getenv("MQTT_PASSWORD"))

	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	e := echo.New()
	e.GET("/", index)
	e.GET("/api/receiver", dataIn)
	e.Logger.Fatal(e.Start(":8080"))
}
func index(c echo.Context) error {
	return c.String(http.StatusOK, "Ambient Weather Receiver ")
}

func dataIn(c echo.Context) error {
	logger.Info(c.Request().URL.String())
	output := map[string]interface{}{}

	in := map[string]string{}
	in["baromabsin"] = "float"
	in["baromrelin"] = "float"
	in["batt_co2"] = "int"
	in["battlightning"] = "int"
	in["batt1"] = "int"
	in["batt2"] = "int"
	in["batt3"] = "int"
	in["batt4"] = "int"
	in["battin"] = "int"
	in["battout"] = "int"
	in["dailyrainin"] = "float"
	in["dateutc"] = "string"
	in["eventrainin"] = "float"
	in["hourlyrainin"] = "float"
	in["humidity"] = "int"
	in["humidity1"] = "int"
	in["humidity2"] = "int"
	in["humidity3"] = "int"
	in["humidity4"] = "int"
	in["humidityin"] = "int"
	in["lightningday"] = "int"
	in["lightningdistance"] = "int"
	in["lightningtime"] = "string"
	in["maxdailygust"] = "float"
	in["monthlyrainin"] = "float"
	in["solarradiation"] = "float"
	in["temp1f"] = "float"
	in["temp2f"] = "float"
	in["temp3f"] = "float"
	in["temp4f"] = "float"
	in["tempf"] = "float"
	in["tempinf"] = "float"
	in["uv"] = "int"
	in["weeklyrainin"] = "float"
	in["winddir"] = "int"
	in["windgustmph"] = "float"
	in["windspeedmph"] = "float"
	in["yearlyrainin"] = "float"
	in["aqipm25"] = "int"
	in["aqipm2524h"] = "int"

	values := c.QueryParams()
	if len(values) == 0 {
		return errors.New("no values received")
	}
	for k, v := range values {
		k = strings.Replace(k, "_", "", -1)
		val := v[0]
		switch in[k] {
		case "int":
			i, err := strconv.Atoi(val)
			if err != nil {
				log.Printf("%s - %s\n", err, val)
			}
			output[k] = i
			logger.Debug(k, " - ", i)
		case "float":
			f, err := strconv.ParseFloat(val, 64)
			if err != nil {
				log.Printf("%s - %s\n", err, val)
			}
			output[k] = toFixed(f, 2)
			logger.Debug(k, " - ", toFixed(f, 2))
		default:
			if k == "PASSKEY" {
				k = "mac"
			}

			output[k] = v[0]

			if k == "lightningtime" {
				i, err := strconv.ParseInt(v[0], 10, 64)
				if err != nil {
					logger.Error(err)
				}
				output[k] = time.Unix(i, 0)
			}
			logger.Debug(k, " - ", v[0])

		}
	}
	output["date"] = time.Now()

	b, err := json.Marshal(output)
	if err != nil {
		log.Println(err)
	}

	topic := os.Getenv("MQTT_PUBLISH")
	token := client.Publish(topic, 0, false, b)
	token.Wait()

	return c.NoContent(http.StatusOK)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
