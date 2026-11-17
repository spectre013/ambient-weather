package main

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	fo "git.zoms.net/apps/weather/sdr/fineoffset"
)

// Payloads received off the air. The expected values are what rtl_433 printed
// for the same transmissions, apart from battery_low and the WH40 battery
// fields, which that rtl_433 build does not output.
func TestRealPayloads(t *testing.T) {
	tests := []struct{ payload, want string }{
		{"8000860200819ea1e641084e09003fff830f",
			`{"model":"Fineoffset-WS80","id":34306,"battery_ok":1.1,"battery_mV":3160,"temperature_C":8.6,"humidity":65,"wind_dir_deg":334,"wind_avg_m_s":0.8,"wind_max_m_s":0.9,"uvi":0,"light_lux":1290,"flags":161,"battery_low":0,"mic":"CRC"}`},
		{"e6325d22205b12",
			`{"model":"Fineoffset-WH32B","id":99,"battery_ok":1,"temperature_C":20.5,"humidity":34,"pressure_hPa":828.3,"battery_low":0,"mic":"CRC"}`},
		{"309c925d244726",
			`{"model":"AmbientWeather-WH31E","id":156,"channel":2,"battery_ok":1,"temperature_C":20.5,"humidity":36,"battery_low":0,"mic":"CRC"}`},
		{"3082a25722ffcc",
			`{"model":"AmbientWeather-WH31E","id":130,"channel":3,"battery_ok":1,"temperature_C":19.9,"humidity":34,"battery_low":0,"mic":"CRC"}`},
		{"4080cead8e06469db2",
			`{"model":"EcoWitt-WH40","id":52909,"battery_V":1.4,"battery_ok":0.83,"rain_mm":160.6,"battery_low":0,"mic":"CRC"}`},
		{"5710d98901bfadc2f8",
			`{"model":"FineOffset-WH31L","id":55689,"battery_ok":0,"state":"interference","flags":4123,"strike_count":173,"battery_low":1,"mic":"CRC"}`},
	}
	for _, tc := range tests {
		p, _ := hex.DecodeString(tc.payload)
		// The radio always hands over 39 bytes; what follows the packet is noise.
		p = append(p, make([]byte, 39-len(p))...)
		meas, err := fo.Decode(p)
		if err != nil {
			t.Errorf("%s: %v", tc.payload, err)
			continue
		}
		got, err := json.Marshal(NewMessage(meas, p))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != tc.want {
			t.Errorf("%s:\n got %s\nwant %s", tc.payload, got, tc.want)
		}
	}
}

// One flipped bit anywhere must fail the checks.
func TestCorruptPayloadRejected(t *testing.T) {
	good, _ := hex.DecodeString("8000860200819ea1e641084e09003fff830f")
	for bit := range len(good) * 8 {
		p := append([]byte(nil), good...)
		p[bit/8] ^= 0x80 >> (bit % 8)
		if _, err := fo.Decode(p); err == nil {
			t.Fatalf("bit %d flipped but the packet still decoded", bit)
		}
	}
}
