package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
)

func getCSS() string {
	css := "/css/index.css"
	if os.Getenv("LOGLEVEL") == "Debug" {
		css = fmt.Sprintf("/css/index.css?v=%d", time.Now().Unix())
	}
	return css
}

func toString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
func getConditions(forecast ForecastImage) string {
	//return templ.EscapeString(forecast.Days[0].Conditions)
	return templ.EscapeString("clear")
}

func getIcon(icon string) string {
	return templ.EscapeString("/images/icons/" + icon + ".png")
}

func cssToString(s string, d string) string {
	return templ.EscapeString(fmt.Sprintf("%s %s", s, d))
}

func dateFormat(layout string, format string, d string) string {
	t, _ := time.Parse(layout, d)
	return t.Format(format)
}

func floatDisplay(val float64) string {
	return toString(toFixed(val, 0))
}

func timeFormat(t time.Time) string {
	t = t.Local()
	return t.Format("15:04")
}

func year(date time.Time) string {
	return toString(date.Year())
}
func month(date time.Time) string {
	return date.Format("Jan")
}
func full(date time.Time) string {
	date = date.Local()
	return date.Format(time.DateTime)
}

func baroLabel(units string) string {
	if units == "metric" {
		return "hPA"
	} else {
		return "inHG"
	}
}

func baroDisplay(baro float64, units string) string {
	if units == "metric" {
		return toString(toFixed(baro*33.86, 2))
	}
	return toString(toFixed(baro, 2))
}

func dewPointClass(dewpoint float64) string {
	if dewpoint > 69.8 {
		return "tempmodulehome25-30c"
	} else if dewpoint >= 68 {
		return "tempmodulehome20-25c"
	} else if dewpoint >= 59 {
		return "tempmodulehome15-20c"
	} else if dewpoint >= 50 {
		return "tempmodulehome10-15c"
	} else if dewpoint > 41 {
		return "tempmodulehome5-10c"
	} else if dewpoint >= 32 {
		return "tempmodulehome0-5c"
	} else if dewpoint > 14 {
		return "tempmodulehome-10-0c"
	} else if dewpoint >= -50 {
		return "tempmodulehome-50-10c"
	}
	return "tempmodulehome0-5c"
}
func humidityClass(humidity float64) string {
	if humidity > 90 {
		return "temphumcircle80-100"
	} else if humidity > 70 {
		return "temphumcircle60-80"
	} else if humidity > 35 {
		return "temphumcircle35-60"
	} else if humidity > 25 {
		return "temphumcircle25-35"
	} else if humidity <= 25 {
		return "temphumcircle0-25"
	}
	return ""
}

func getDay(dateString string) string {
	// Parse the date string
	date, _ := time.Parse("2006-01-02", dateString)
	return strings.ToUpper(date.Format("Mon"))
}
func tempLabel(units string) string {
	if units == "metric" {
		return "&deg;C"
	} else {
		return "&deg;F"
	}
}

func tempDisplay(temp float64, units string) string {
	t := temp
	if units == "metric" {
		t = ((temp - 32) * 5) / 9
	}
	return fmt.Sprintf("%.0f", toFixed(t, 0))
}
func tempColor(temp float64) string {

	if temp <= -5 {
		return "tempcolorminus10"
	} else if temp <= 5 {
		return "tempcolorminus5"
	} else if temp <= 14 {
		return "tempcolorminus"
	} else if temp <= 23 {
		return "tempcolor0-5"
	} else if temp <= 32 {
		return "tempcolorzero"
	} else if temp <= 41 {
		return "tempcolor0-5"
	} else if temp < 50 {
		return "tempcolor6-10"
	} else if temp < 59 {
		return "tempcolor11-15"
	} else if temp < 68 {
		return "tempcolor16-20"
	} else if temp < 77 {
		return "tempcolor21-25"
	} else if temp < 86 {
		return "tempcolor26-30"
	} else if temp < 95 {
		return "tempcolor31-35"
	} else if temp < 104 {
		return "tempcolor36-40"
	} else if temp < 113 {
		return "tempcolor41-45"
	} else if temp < 212 {
		return "tempcolor50"
	}
	return ""
}

func alertColor(alert Alert) string {

	if strings.HasPrefix(alert.Event, "911") {
		return strings.Replace("Telephone Outage 911", " ", " -", -1)
	}
	return strings.Replace(strings.ToLower(alert.Event), " ", "-", -1)
}

//func DegreeToDirection(deg float64) string {
//	// Normalize the degree value
//	normalizedDeg := math.Mod(deg, 360.0)
//
//	switch {
//	case normalizedDeg >= 0 && normalizedDeg < 11.25:
//		return "N"
//	case normalizedDeg >= 11.25 && normalizedDeg < 33.75:
//		return "NNE"
//	case normalizedDeg >= 33.75 && normalizedDeg < 56.25:
//		return "NE"
//	case normalizedDeg >= 56.25 && normalizedDeg < 78.75:
//		return "ENE"
//	case normalizedDeg >= 78.75 && normalizedDeg < 101.25:
//		return "E"
//	case normalizedDeg >= 101.25 && normalizedDeg < 123.75:
//		return "ESE"
//	case normalizedDeg >= 123.75 && normalizedDeg < 146.25:
//		return "SE"
//	case normalizedDeg >= 146.25 && normalizedDeg < 168.75:
//		return "SSE"
//	case normalizedDeg >= 168.75 && normalizedDeg < 191.25:
//		return "S"
//	case normalizedDeg >= 191.25 && normalizedDeg < 213.75:
//		return "SSW"
//	case normalizedDeg >= 213.75 && normalizedDeg < 236.25:
//		return "SW"
//	case normalizedDeg >= 236.25 && normalizedDeg < 258.75:
//		return "WSW"
//	case normalizedDeg >= 258.75 && normalizedDeg < 281.25:
//		return "W"
//	case normalizedDeg >= 281.25 && normalizedDeg < 303.75:
//		return "WNW"
//	case normalizedDeg >= 303.75 && normalizedDeg < 326.25:
//		return "NW"
//	case normalizedDeg >= 326.25 && normalizedDeg < 348.75:
//		return "NNW"
//	default:
//		return "N"
//	}
//}

func windRun(wind float64) float64 {
	return wind * float64(time.Now().Hour())
}

func getBeaufort(windspeed float64) int {
	speed := windspeed
	if speed > 0 {
		speed = windspeed / 1.151
	}
	switch {
	case speed < 1:
		return 0
	case speed < 4:
		return 1
	case speed < 7:
		return 2
	case speed < 11:
		return 3
	case speed < 17:
		return 4
	case speed < 22:
		return 5
	case speed < 28:
		return 6
	case speed < 34:
		return 7
	case speed < 41:
		return 8
	case speed < 48:
		return 9
	case speed < 56:
		return 10
	case speed < 64:
		return 11
	default:
		return 12
	}
}

func beaufortScale(windSpeed float64) Beaufort {
	bft := getBeaufort(windSpeed)

	if bft >= 12 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft12" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 29.96 C 5.55 36.68 11.21 43.31 16.80 50.00 C 18.26 49.99 19.73 49.99 21.19 49.99 C 18.93 47.26 16.67 44.53 14.40 41.79 C 15.94 40.54 17.47 39.27 19.00 38.00 C 22.34 42.00 25.66 46.01 29.01 50.00 C 42.72 49.98 56.43 50.03 70.14 49.98 C 71.17 47.82 72.07 45.50 73.83 43.81 C 77.91 39.62 84.85 39.15 89.85 41.94 C 93.15 43.97 95.29 47.56 96.00 51.33 L 96.00 54.56 C 95.35 58.38 93.17 62.01 89.84 64.06 C 85.44 66.52 79.67 66.42 75.46 63.60 C 72.81 61.81 71.37 58.87 70.15 56.02 C 46.76 55.98 23.38 56.01 0.00 56.00 L 0.00 29.96 Z" /></svg>`,
			Text:  "Hurricane",
			Class: "beaufort6",
		}
	} else if bft >= 11 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft11" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 29.97 C 5.57 36.68 11.20 43.33 16.81 50.00 C 34.60 49.99 52.38 50.02 70.16 49.99 C 71.98 43.63 78.44 39.00 85.10 40.36 C 90.77 40.90 95.07 45.87 96.00 51.29 L 96.00 54.67 C 95.15 59.33 91.95 63.89 87.21 65.17 C 82.45 66.67 76.62 65.56 73.32 61.64 C 71.87 60.01 71.03 57.98 70.16 56.01 C 46.77 55.99 23.39 56.01 0.00 56.00 L 0.00 29.97 Z" /></svg>`,
			Text:  "Violent Storm",
			Class: "beaufort6",
		}
	} else if bft >= 10 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft10"  width="12pt" height="12pt" viewBox="0 0 96 96" version="1.1" ><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 29.97 C 5.30 36.33 10.66 42.65 15.99 48.98 C 16.01 42.66 15.99 36.34 16.00 30.02 C 21.62 36.67 27.19 43.35 32.81 50.00 C 34.20 50.00 35.60 49.99 36.99 50.00 C 33.74 45.99 30.46 42.01 27.21 38.00 C 28.66 36.67 30.12 35.34 31.58 34.01 C 36.02 39.32 40.38 44.69 44.81 50.00 C 53.27 49.99 61.72 50.02 70.18 49.99 C 71.39 46.85 73.14 43.69 76.15 41.96 C 80.11 39.71 85.11 39.63 89.20 41.59 C 92.87 43.50 95.27 47.34 96.00 51.35 L 96.00 54.56 C 95.18 60.08 90.75 65.14 85.02 65.65 C 78.40 66.97 71.95 62.35 70.18 56.01 C 46.79 55.98 23.39 56.01 0.00 56.00 L 0.00 29.97 Z" /></svg>`,
			Text:  "Storm",
			Class: "beaufort6",
		}
	} else if bft >= 9 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft9" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 29.97 C 5.29 36.34 10.66 42.65 15.99 48.99 C 16.01 42.66 15.99 36.34 16.00 30.01 C 21.61 36.67 27.19 43.34 32.80 50.00 C 45.26 49.99 57.71 50.02 70.16 49.98 C 71.97 43.66 78.38 39.03 85.02 40.35 C 90.73 40.87 95.12 45.87 96.00 51.36 L 96.00 54.55 C 95.18 60.08 90.75 65.14 85.00 65.66 C 78.37 66.96 71.98 62.34 70.16 56.02 C 46.77 55.98 23.39 56.01 0.00 56.00 L 0.00 29.97 Z" /></svg>`,
			Text:  "Strong Gale",
			Class: "beaufort6",
		}
	} else if bft >= 8 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft8" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 4.64 30.07 C 10.05 36.70 15.41 43.37 20.82 50.01 C 22.21 50.00 23.60 50.00 25.00 49.99 C 20.66 44.66 16.33 39.33 12.00 34.00 C 13.54 32.67 15.07 31.34 16.60 30.00 C 22.01 36.67 27.40 43.35 32.81 50.01 C 34.21 50.00 35.60 49.99 37.00 49.99 C 32.66 44.66 28.33 39.33 24.00 34.00 C 25.54 32.67 27.07 31.34 28.60 30.00 C 34.01 36.67 39.40 43.35 44.82 50.01 C 46.21 50.00 47.60 50.00 49.00 49.99 C 44.66 44.66 40.33 39.33 36.00 34.00 C 37.54 32.67 39.07 31.34 40.60 30.00 C 46.01 36.67 51.40 43.35 56.81 50.01 C 58.34 50.00 59.86 50.00 61.39 49.99 C 58.60 46.59 55.80 43.20 53.00 39.80 C 54.54 38.53 56.07 37.27 57.61 36.01 C 61.73 40.79 65.44 45.94 69.89 50.42 C 71.21 47.70 72.41 44.73 74.89 42.83 C 79.11 39.58 85.30 39.36 89.89 41.99 C 93.19 43.96 95.20 47.55 96.00 51.23 L 96.00 54.77 C 95.21 58.43 93.21 62.00 89.94 63.98 C 85.52 66.55 79.63 66.43 75.40 63.55 C 72.77 61.77 71.38 58.81 70.11 56.01 C 52.74 55.99 35.38 56.01 18.01 56.00 C 11.92 48.95 6.57 41.23 0.00 34.64 L 0.00 33.40 C 1.68 32.49 3.18 31.30 4.64 30.07 Z" /></svg>`,
			Text:  "Gale",
			Class: "beaufort6",
		}
	} else if bft >= 7 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft7" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 34.01 C 1.53 32.68 3.03 31.30 4.61 30.03 C 10.06 36.65 15.35 43.40 20.85 49.98 C 22.23 50.02 23.60 50.02 24.98 49.97 C 20.67 44.64 16.31 39.36 12.04 34.00 C 13.53 32.64 15.05 31.31 16.61 30.03 C 22.05 36.65 27.35 43.39 32.84 49.98 C 34.22 50.02 35.60 50.02 36.98 49.98 C 32.69 44.64 28.30 39.37 24.04 34.00 C 25.53 32.64 27.05 31.31 28.61 30.03 C 34.05 36.65 39.36 43.39 44.83 49.98 C 46.35 50.02 47.86 50.02 49.38 49.99 C 46.62 46.57 43.78 43.22 41.03 39.80 C 42.53 38.52 44.05 37.24 45.61 36.03 C 49.51 40.65 53.29 45.38 57.22 49.98 C 61.55 50.03 65.88 50.00 70.21 49.99 C 71.17 47.29 72.62 44.67 74.86 42.84 C 78.91 39.72 84.66 39.43 89.20 41.60 C 92.85 43.49 95.26 47.32 96.00 51.30 L 96.00 54.66 C 95.11 60.04 90.82 65.13 85.16 65.58 C 78.59 67.06 71.90 62.43 70.21 56.01 C 52.82 55.97 35.43 56.04 18.04 55.98 C 11.96 48.71 6.04 41.31 0.00 34.01 L 0.00 34.01 Z" /></svg>`,
			Text:  "Near Gale",
			Class: "beaufort6",
		}
	} else if bft >= 6 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft6" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 4.55 30.01 C 10.03 36.62 15.37 43.35 20.81 50.00 C 22.20 50.00 23.60 50.00 24.99 49.99 C 20.67 44.65 16.33 39.34 12.01 34.00 C 13.53 32.66 15.07 31.33 16.60 30.00 C 22.02 36.67 27.39 43.38 32.84 50.02 C 34.22 50.01 35.60 49.99 36.98 49.98 C 32.67 44.64 28.31 39.34 24.01 33.99 C 25.54 32.66 27.07 31.33 28.60 30.01 C 34.01 36.67 39.39 43.35 44.81 50.00 C 53.26 49.99 61.71 50.01 70.15 49.99 C 71.04 48.00 71.89 45.95 73.36 44.31 C 76.67 40.43 82.45 39.34 87.19 40.83 C 91.91 42.08 95.07 46.60 96.00 51.22 L 96.00 54.75 C 95.20 58.73 92.83 62.57 89.13 64.44 C 84.81 66.48 79.42 66.27 75.43 63.58 C 72.80 61.79 71.34 58.86 70.15 56.01 C 52.77 55.99 35.39 56.01 18.01 56.00 C 11.92 48.94 6.51 41.22 0.00 34.56 L 0.00 33.45 C 1.83 32.80 3.11 31.23 4.55 30.01 Z" /></svg>`,
			Text:  "Strong Breeze",
			Class: "beaufort4-5",
		}
	} else if bft >= 5 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft5" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 4.55 30.01 C 10.04 36.62 15.37 43.37 20.82 50.01 C 22.21 50.00 23.60 49.99 25.00 49.99 C 20.67 44.66 16.33 39.33 12.00 34.00 C 13.53 32.67 15.07 31.34 16.60 30.01 C 22.01 36.67 27.39 43.35 32.82 50.01 C 45.26 49.98 57.71 50.02 70.15 49.99 C 71.41 46.91 73.07 43.77 76.03 42.02 C 79.40 40.12 83.56 39.63 87.24 40.85 C 91.95 42.11 95.08 46.63 96.00 51.23 L 96.00 55.03 C 95.11 58.56 93.16 61.97 90.02 63.95 C 85.60 66.53 79.71 66.45 75.44 63.58 C 72.80 61.79 71.34 58.86 70.15 56.01 C 52.77 55.99 35.39 56.00 18.02 56.00 C 11.93 48.90 6.44 41.24 0.00 34.48 L 0.00 33.53 C 1.72 32.64 3.15 31.32 4.55 30.01 Z" /></svg>`,
			Text:  "Fresh Breeze",
			Class: "beaufort4-5",
		}
	} else if bft >= 4 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft4" width="12pt" height="12pt" viewBox="0 0 96 96" ><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 33.39 C 1.62 32.38 3.17 31.27 4.69 30.10 C 10.05 36.74 15.43 43.37 20.80 50.01 C 22.27 49.99 23.73 49.99 25.20 49.99 C 22.39 46.60 19.61 43.19 16.80 39.80 C 18.34 38.53 19.87 37.27 21.40 36.00 C 25.26 40.67 29.13 45.33 33.00 50.00 C 45.36 49.99 57.72 50.02 70.08 49.98 C 71.35 47.43 72.52 44.67 74.84 42.87 C 79.08 39.57 85.34 39.34 89.94 42.02 C 93.23 44.01 95.21 47.59 96.00 51.27 L 96.00 54.84 C 95.16 58.45 93.23 61.98 89.99 63.95 C 85.38 66.65 79.11 66.44 74.86 63.15 C 72.54 61.35 71.34 58.58 70.08 56.02 C 52.72 55.98 35.37 56.01 18.01 56.00 C 11.92 48.97 6.60 41.23 0.00 34.67 L 0.00 33.39 Z" /></svg>`,
			Text:  "Moderate Breeze",
			Class: "beaufort3-4",
		}
	} else if bft >= 3 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft3" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 33.44 C 1.67 32.50 3.17 31.28 4.64 30.06 C 10.04 36.70 15.41 43.36 20.80 50.00 C 37.24 49.99 53.68 50.02 70.12 49.98 C 71.39 47.19 72.76 44.24 75.38 42.46 C 79.66 39.55 85.61 39.46 90.05 42.08 C 93.25 44.09 95.22 47.60 96.00 51.23 L 96.00 54.90 C 95.16 58.48 93.20 61.96 90.01 63.95 C 85.59 66.53 79.71 66.44 75.44 63.58 C 72.79 61.80 71.39 58.83 70.12 56.02 C 52.75 55.98 35.38 56.01 18.01 56.00 C 11.92 48.94 6.53 41.24 0.00 34.58 L 0.00 33.44 Z" /></svg>`,
			Text:  "Gentle Breeze",
			Class: "beaufort1-3",
		}
	} else if bft >= 2 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft2" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 33.38 C 1.68 32.46 3.19 31.28 4.67 30.09 C 10.04 36.72 15.42 43.36 20.80 50.00 C 37.23 49.99 53.66 50.03 70.09 49.98 C 71.41 47.21 72.76 44.23 75.39 42.45 C 79.66 39.54 85.60 39.45 90.03 42.07 C 93.26 44.08 95.26 47.64 96.00 51.31 L 96.00 54.79 C 95.15 58.68 92.92 62.47 89.30 64.34 C 84.74 66.62 78.83 66.29 74.79 63.09 C 72.52 61.29 71.31 58.57 70.10 56.02 C 46.73 55.97 23.37 56.02 0.00 56.00 L 0.00 49.94 C 4.33 50.04 8.66 50.00 13.00 49.99 C 8.62 44.92 4.88 39.28 0.00 34.68 L 0.00 33.38 Z" /></svg>`,
			Text:  "Light Breeze",
			Class: "beaufort1-3",
		}
	} else if bft >= 1 {
		return Beaufort{
			SVG:   `<svg id="weather34 bft1" width="12pt" height="12pt" viewBox="0 0 96 96" ><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 73.92 43.89 C 77.12 40.10 82.80 39.45 87.34 40.81 C 91.48 42.01 93.99 45.85 96.00 49.39 L 96.00 56.58 C 94.00 60.14 91.49 63.99 87.34 65.19 C 82.80 66.55 77.13 65.90 73.92 62.11 C 72.32 60.28 71.03 58.19 69.69 56.16 C 46.47 55.76 23.23 56.12 0.00 56.00 L 0.00 50.00 C 23.23 49.88 46.47 50.24 69.69 49.84 C 71.03 47.81 72.31 45.73 73.92 43.89 Z" /></svg>`,
			Text:  "Light Air",
			Class: "beaufort1-3",
		}
	}

	return Beaufort{
		SVG:   `<svg id="weather34 bft1" width="12pt" height="12pt" viewBox="0 0 96 96" ><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 73.92 43.89 C 77.12 40.10 82.80 39.45 87.34 40.81 C 91.48 42.01 93.99 45.85 96.00 49.39 L 96.00 56.58 C 94.00 60.14 91.49 63.99 87.34 65.19 C 82.80 66.55 77.13 65.90 73.92 62.11 C 72.32 60.28 71.03 58.19 69.69 56.16 C 46.47 55.76 23.23 56.12 0.00 56.00 L 0.00 50.00 C 23.23 49.88 46.47 50.24 69.69 49.84 C 71.03 47.81 72.31 45.73 73.92 43.89 Z" /></svg>`,
		Text:  "Light Air",
		Class: "beaufort1-3",
	} // Default return if no condition matches
}

func windLabel(units string) string {
	if units == "imperial" {
		return "MPH"
	} else {
		return "M/S"
	}
}

func windDisplay(wind float64, units string) string {
	if units == "ms" {
		return fmt.Sprintf("%.f", toFixed(wind/2.237, 0))
	}
	return fmt.Sprintf("%.f", toFixed(wind, 0))
}

func degToCompass(num float64) string {
	val := math.Floor(num/22.5 + 0.5)
	arr := []string{
		"North",
		"NNE",
		"NE",
		"ENE",
		"East",
		"ESE",
		"SE",
		"SSE",
		"South",
		"SSW",
		"SW",
		"WSW",
		"West",
		"WNW",
		"NW",
		"NNW",
	}
	dir := int(math.Mod(val, 16.0))
	return arr[dir]
}

func rainLabel(units string) string {
	if units == "metric" {
		return "mm"
	} else {
		return "in"
	}
}

func rainDisplay(rn float64, units string) string {
	if units == "metric" {
		return fmt.Sprintf("%.2f", toFixed(rn*25.4, 2))
	}
	return fmt.Sprintf("%.2f", toFixed(rn, 2))
}

func sunTime(luna Astro) map[string]time.Time {
	times := map[string]time.Time{}
	times["sunrise"] = luna.Sunrise
	times["sunset"] = luna.Sunset
	times["tomorrow"] = luna.SunriseTomorrow
	return times
}

func ParseDuration(d string) map[string]string {
	var hours string
	var minutes string
	if strings.Contains(d, "h") {
		h := strings.Split(d, "h")
		m := strings.Split(h[1], "m")
		hours = h[0]
		minutes = m[0]
	} else {
		m := strings.Split(d, "m")
		hours = "0"
		minutes = m[0]
	}

	return map[string]string{
		"hour": hours,
		"min":  minutes,
	}
}

func lightDark(t time.Duration) string {
	r := ParseDuration(t.String())
	return fmt.Sprintf("%s hrs %s min", r["hour"], r["min"])
}

func sunHasSet(luna Astro) bool {
	ss := false
	times := sunTime(luna)
	if time.Now().After(times["sunset"]) && time.Now().Before(times["tomorrow"]) {
		ss = true
	}
	return ss
}

func isSunSet(luna Astro) string {
	if luna.HasSunset {
		return "Time til Sunrise"
	} else {
		return "Time til Sunset"
	}
}

func riseSetClass(luna Astro) string {
	if luna.HasSunset {
		return "riseclr"
	} else {
		return "setclr"
	}
}
func sunBelow(luna Astro) string {
	if luna.HasSunset {
		return "sunbelow"
	} else {
		return "sunabove"
	}
}

func sunTimes(luna Astro) map[string]string {
	var t string
	times := sunTime(luna)
	if luna.HasSunset {
		t = time.Until(times["tomorrow"]).String()
	} else {
		t = time.Until(times["sunset"]).String()
	}
	return ParseDuration(t)
}

func todayTomorrow(t string, luna Astro) string {
	event := map[string]time.Time{}
	event["sunrise"] = luna.Sunrise
	event["sunset"] = luna.Sunset

	if time.Now().After(event[t]) {
		return "Tomorrow"
	} else {
		return "Today"
	}
}

func uvToday(data Record) string {
	if data.Uv >= 10 {
		return "uvtoday11"
	} else if data.Uv >= 8 {
		return "uvtoday9-10"
	} else if data.Uv >= 5 {
		return "uvtoday6-8"
	} else if data.Uv >= 3 {
		return "uvtoday4-5"
	} else if data.Uv >= 0 {
		return "uvtoday1-3"
	}
	return ""
}

func uvCaution(data TemplateData) string {
	uv := data.Record.Uv
	if uv >= 10 {
		return "Extreme"
	} else if uv >= 8 {
		return "Very High"
	} else if uv >= 5 {
		return "High"
	} else if uv >= 3 {
		return "Moderate"
	} else if !sunHasSet(data.Astro) && uv >= 0 {
		return "Low"
	} else if sunHasSet(data.Astro) && uv <= 0 {
		return "Below Horizon"
	}
	return ""
}

func Linear(AQIhigh float64, AQIlow float64, Conchigh float64, Conclow float64, Concentration float64) float64 {
	Conc := Concentration
	a := ((Conc-Conclow)/(Conchigh-Conclow))*(AQIhigh-AQIlow) + AQIlow
	return math.Round(a)
}

func AqiCalc(Concentration int) float64 {
	Conc := Concentration
	AQI := 0.0
	c := (math.Floor(10 * float64(Conc))) / 10
	if c >= 0 && c < 12.1 {
		AQI = Linear(50, 0, 12, 0, c)
	} else if c >= 12.1 && c < 35.5 {
		AQI = Linear(100, 51, 35.4, 12.1, c)
	} else if c >= 35.5 && c < 55.5 {
		AQI = Linear(150, 101, 55.4, 35.5, c)
	} else if c >= 55.5 && c < 150.5 {
		AQI = Linear(200, 151, 150.4, 55.5, c)
	} else if c >= 150.5 && c < 250.5 {
		AQI = Linear(300, 201, 250.4, 150.5, c)
	} else if c >= 250.5 && c < 350.5 {
		AQI = Linear(400, 301, 350.4, 250.5, c)
	} else if c >= 350.5 {
		AQI = Linear(500, 401, 500.4, 350.5, c)
	}
	return AQI
}

func getDetails(aqi float64) AqiCategory {
	categories := make([]AqiCategory, 0)
	categories = append(categories, AqiCategory{Max: 50, Color: "green", Name: "Good"})
	categories = append(categories, AqiCategory{Max: 100, Color: "yellow", Name: "Moderate"})
	categories = append(categories, AqiCategory{Max: 150, Color: "orange", Name: "Unhealthy for sensitive groups"})
	categories = append(categories, AqiCategory{Max: 200, Color: "red", Name: "Unhealthy"})
	categories = append(categories, AqiCategory{Max: 300, Color: "purple", Name: "Very unhealth"})
	categories = append(categories, AqiCategory{Max: 500, Color: "maroon", Name: "Hazardous"})

	if aqi < 0 || (aqi >= 0 && aqi <= 50) {
		return categories[0]
	} else if aqi > 50 || aqi <= 100 {
		return categories[1]
	} else if aqi > 100 || aqi <= 150 {
		return categories[2]
	} else if aqi > 150 || aqi <= 200 {
		return categories[3]
	} else if aqi > 200 || aqi <= 300 {
		return categories[4]
	} else if aqi > 300 {
		return categories[5]
	}
	return categories[0]
}

func getSensor(s string, rec Record) float64 {

	switch s {
	case "tempinf":
		return rec.Tempinf
	case "temp1f":
		return rec.Temp1f
	case "temp2f":
		return rec.Temp2f
	case "temp3f":
		return rec.Temp3f
	case "humidityin":
		return float64(rec.Humidityin)
	case "humidity1":
		return float64(rec.Humidity1)
	case "humidity2":
		return float64(rec.Humidity2)
	case "humidity3":
		return float64(rec.Humidity3)
	}

	return 0.0
}

func PillData() map[string]Pill {

	pill := map[string]Pill{}
	pill["1h"] = Pill{Interval: "1 hour", Name: "1 Hour"}
	pill["6h"] = Pill{Interval: "6 hours", Name: "6 Hours"}
	pill["12h"] = Pill{Interval: "12 hours", Name: "12 Hours"}
	pill["1d"] = Pill{Interval: "24 hours", Name: "24 Hours"}
	pill["1m"] = Pill{Interval: "1 month", Name: "Month"}
	pill["1y"] = Pill{Interval: "1 year", Name: "Year"}

	return pill
}

func getPill(pill string) Pill {
	pd := PillData()
	return pd[pill]
}

func lightningClass(cnt int) string {
	if cnt == 0 || cnt < 50 {
		return "green"
	} else if cnt >= 50 && cnt < 250 {
		return "yellow"
	} else if cnt >= 250 && cnt < 500 {
		return "orange"
	} else if cnt >= 500 {
		return "red"
	}
	return "green"
}

func distanceClass(d int) string {
	if d == 0 {
		return "green"
	} else if d >= 1 && d < 5 {
		return "red"
	} else if d >= 5 && d < 10 {
		return "orange"
	} else if d >= 10 && d < 15 {
		return "yellow"
	} else if d >= 15 {
		return "green"
	}
	return "green"
}
