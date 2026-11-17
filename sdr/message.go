package main

import (
	"bytes"
	"encoding/json"
	"math"

	fo "git.zoms.net/apps/weather/sdr/fineoffset"
)

// Message is one sensor reading as an ordered JSON object. Field names,
// units and model names follow rtl_433, because that is what the weather
// processor already consumes on MQTT ("ambient/<model>/<id>").
type Message struct {
	Model  string
	ID     int
	fields []field
}

type field struct {
	key string
	val any
}

func (m *Message) add(key string, val any) *Message {
	m.fields = append(m.fields, field{key, val})
	return m
}

// MarshalJSON keeps the fields in insertion order.
func (m *Message) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, f := range m.fields {
		if i > 0 {
			b.WriteByte(',')
		}
		k, _ := json.Marshal(f.key)
		v, err := json.Marshal(f.val)
		if err != nil {
			return nil, err
		}
		b.Write(k)
		b.WriteByte(':')
		b.Write(v)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

func round(v float64, places int) float64 {
	p := math.Pow(10, float64(places))
	return math.Round(v*p) / p
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}

// Optional-field helpers: a firmware "absent" sentinel leaves the key out.
func (m *Message) temp(t int16) *Message {
	if t != fo.Invalid16 {
		m.add("temperature_C", round(float64(t)*0.1, 1))
	}
	return m
}

func (m *Message) hum(h uint8) *Message {
	if h != fo.InvalidHum && h != 0xff {
		m.add("humidity", int(h))
	}
	return m
}

// wind adds speed, direction and gust given in 0.1 m/s and degrees.
func (m *Message) wind(dir int16, avg, gust uint16) *Message {
	if dir != fo.Invalid16 {
		m.add("wind_dir_deg", int(dir))
	}
	if avg != fo.InvalidU16 {
		m.add("wind_avg_m_s", round(float64(avg)*0.1, 1))
	}
	if gust != fo.InvalidU16 {
		m.add("wind_max_m_s", round(float64(gust)*0.1, 1))
	}
	return m
}

func (m *Message) light(l uint32) *Message {
	if l != fo.InvalidLight {
		m.add("light_lux", round(float64(l)*0.1, 1))
	}
	return m
}

func (m *Message) pm(key string, v uint16) *Message {
	if v != fo.InvalidFFFF {
		m.add(key, round(float64(v)*0.1, 1))
	}
	return m
}

// battLevel is rtl_433's 0..1 "battery_ok" level for a voltage range.
func battLevel(mv, emptyMV, fullMV int) float64 {
	return round(max(min(float64(mv-emptyMV)/float64(fullMV-emptyMV), 1.5), 0), 2)
}

// NewMessage converts a decoded measurement. Wind is published with the
// console's default 100 % gain and no direction offset.
func NewMessage(meas fo.Measurement, payload []byte) *Message {
	m := &Message{ID: int(meas.SensorID())}
	head := func(model string) *Message {
		m.Model = model
		return m.add("model", model).add("id", m.ID)
	}
	low := meas.LowBattery()

	switch v := meas.(type) {
	case fo.THP:
		// rtl_433 builds this ID with the nibbles the other way round;
		// keep its value so existing topics and logs do not change.
		m.ID = int(payload[0]&0x0f)<<4 | int(payload[1]>>4)
		if v.Family() == fo.FamilyIndoorTHP {
			head("Fineoffset-WH32B")
		} else {
			head("Fineoffset-WH32")
		}
		m.add("battery_ok", b2i(!low)).temp(v.Temp).hum(v.Hum)
		if v.Pressure != fo.InvalidU16 {
			m.add("pressure_hPa", round(float64(v.Pressure)*0.1, 1))
		}

	case fo.WH31:
		if v.Family() == fo.FamilyWH31Alt {
			head("AmbientWeather-WH31B")
		} else {
			head("AmbientWeather-WH31E")
		}
		m.add("channel", int(v.Ch)+1).add("battery_ok", b2i(!low)).temp(v.Temp).hum(v.Hum)

	case fo.Anemometer:
		mv := int(v.BattRaw) * 20
		if v.Family() == fo.FamilyWS80 {
			// rtl_433 reports the 24-bit ID for the WS80.
			m.ID = int(payload[1])<<16 | int(payload[2])<<8 | int(payload[3])
			head("Fineoffset-WS80").add("battery_ok", battLevel(mv, 1400, 3000))
		} else {
			head("EcoWitt-WS68").add("battery_ok", battLevel(mv, 1800, 3000))
		}
		m.add("battery_mV", mv).temp(v.Temp).hum(v.Hum).
			wind(v.Dir, fo.ScaleWind(v.Wind, 100), fo.ScaleWind(v.Gust, 100))
		if v.UV != fo.InvalidByte {
			m.add("uvi", round(float64(v.UV)*0.1, 1))
		}
		m.light(v.Light).add("flags", int(payload[7]))

	case fo.WS90:
		mv := int(v.BattRaw) * 20
		head("Fineoffset-WS90").add("battery_ok", battLevel(mv, 1400, 3000)).add("battery_mV", mv).
			temp(v.Temp).hum(v.Hum).
			wind(v.Dir, fo.ScaleWind(v.Wind, 100), fo.ScaleWind(v.Gust, 100))
		if v.UV != fo.InvalidByte {
			m.add("uvi", round(float64(v.UV)*0.1, 1))
		}
		m.light(v.Light).add("flags", int(payload[7])).
			add("rain_mm", round(float64(v.RainCounter)*0.1, 1)).
			add("firmware", int(v.Version))

	case fo.WS85:
		mv := int(v.BattRaw) * 20
		head("Fineoffset-WS85").add("battery_ok", battLevel(mv, 1400, 3000)).add("battery_mV", mv).
			wind(v.Dir, fo.ScaleWind(v.Wind, 100), fo.ScaleWind(v.Gust, 100)).
			add("rain_mm", round(float64(v.RainCounter)*0.1, 1))

	case fo.WH24:
		if v.Extended {
			head("Fineoffset-WH65B")
		} else {
			head("Fineoffset-WH24")
		}
		m.add("battery_ok", b2i(!low)).temp(v.Temp).hum(v.Hum).
			wind(v.WindDir, fo.ScaleWH24Wind(v.WindRaw, 100), fo.ScaleWH24Gust(v.GustRaw, 100)).
			// The WH24 bucket tips in 0.3 mm steps (rtl_433 figure; the
			// console keeps the raw counter).
			add("rain_mm", round(float64(v.Rain)*0.3, 1))
		if v.UVRaw != fo.InvalidU16 {
			m.add("uv", int(v.UVRaw)).add("uvi", int(fo.WH24UVIndex(fo.ScaleWH24UV(v.UVRaw, 100))))
		}
		m.light(v.Light)
		if v.Pressure != fo.InvalidU16 {
			m.add("pressure_hPa", round(float64(v.Pressure)*0.1, 1))
		}

	case fo.WH40:
		head("EcoWitt-WH40")
		if v.BattReported {
			m.add("battery_V", round(float64(v.BattRaw)*0.1, 1)).
				add("battery_ok", battLevel(int(v.BattRaw)*100, 900, 1500))
		}
		m.add("rain_mm", round(float64(v.Rain)*0.1, 1))

	case fo.PM25:
		// battery_ok is rtl_433's "bars out of 5".
		head("Fineoffset-WH0290").add("battery_ok", round(float64(v.BattLevel&7)*0.2, 1)).
			pm("pm2_5_ug_m3", v.PM25).pm("estimated_pm10_0_ug_m3", v.PM10).
			add("family", int(payload[0]))

	case fo.WH57:
		state := [...]string{"reset", "interference", "noise", "strike"}[v.State]
		head("FineOffset-WH31L").
			add("battery_ok", round(float64(v.BattLevel&6>>1)*0.5, 1)).
			add("state", state).
			add("flags", int(payload[1]>>4)<<12|int(payload[4])<<4|int(payload[5]>>4))
		if v.DistanceKm != fo.InvalidByte {
			m.add("storm_dist_km", int(v.DistanceKm))
		}
		m.add("strike_count", int(v.Count))

	case fo.WH51:
		mv := int(v.BattRaw) * 100
		head("Fineoffset-WH51").add("battery_ok", battLevel(mv, 1000, 1700)).add("battery_mV", mv)
		if v.Moisture != fo.InvalidHum {
			m.add("moisture", int(v.Moisture))
		}
		m.add("ad_raw", int(v.AD))

	case fo.WH45:
		if v.Family() == fo.FamilyWH46 {
			head("Fineoffset-WH46")
		} else {
			head("Fineoffset-WH45")
		}
		m.add("battery_ok", round(min(float64(v.BattLevel)*0.2, 1), 1)).temp(v.Temp).hum(v.Hum).
			pm("pm2_5_ug_m3", v.PM25).pm("pm10_ug_m3", v.PM10).add("co2_ppm", int(v.CO2))
		m.pm("pm1_ug_m3", v.Ext10).pm("pm4_ug_m3", v.Ext12)

	case fo.WN34:
		mv := int(v.BattRaw) * 20
		head("Fineoffset-WN34").add("battery_ok", battLevel(mv, 1400, 3000)).add("battery_mV", mv).temp(v.Temp)

	case fo.WN35:
		mv := int(v.BattRaw) * 20
		head("Fineoffset-WN35").add("battery_ok", battLevel(mv, 1400, 3000)).add("battery_mV", mv)
		if v.HasTemp {
			m.temp(v.Temp)
		}
		if v.Wetness != fo.InvalidHum {
			m.add("leaf_wetness", int(v.Wetness))
		}

	case fo.WH55:
		head("Fineoffset-WH55").add("channel", int(v.Ch)).add("battery_ok", round(min(float64(v.BattLevel)*0.2, 1), 1)).
			add("leak", b2i(v.Leak)).add("raw_value", int(v.Value))

	case fo.WH48:
		head("Fineoffset-WH48").add("channel", int(v.Ch)).add("battery_ok", b2i(!low)).
			temp(v.Temp).hum(v.Hum).add("co2_ppm", int(v.CO2))

	case fo.WH54:
		head("Fineoffset-WH54").add("battery_mV", int(v.BattRaw)*20)
		if v.Distance != fo.InvalidU16 {
			m.add("distance_mm", int(v.Distance))
		}

	default:
		head(meas.Family().String())
	}

	// The console's own low-battery decision (per-family thresholds).
	return m.add("battery_low", b2i(low)).add("mic", "CRC")
}
