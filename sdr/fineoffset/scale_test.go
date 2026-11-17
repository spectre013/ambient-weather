package fineoffset

import "testing"

func TestScale(t *testing.T) {
	if got := ScaleWH24Wind(0x110, 100); got != 173 { // 272*0.6375 = 173.4
		t.Errorf("WH24 wind %d", got)
	}
	if got := ScaleWH24Gust(20, 100); got != 102 {
		t.Errorf("WH24 gust %d", got)
	}
	if got := ScaleWind(261, 100); got != 261 {
		t.Errorf("wind %d", got)
	}
	if got := ScaleWind(500, 200); got != 800 {
		t.Errorf("wind cap %d", got)
	}
	if ScaleWind(InvalidU16, 100) != InvalidU16 {
		t.Error("wind invalid")
	}
	for _, tc := range []struct{ d, off, want int16 }{{350, 20, 10}, {10, -20, 350}, {0, -360, 360}, {Invalid16, 5, Invalid16}} {
		if got := OffsetDir(tc.d, tc.off); got != tc.want {
			t.Errorf("OffsetDir(%d,%d)=%d", tc.d, tc.off, got)
		}
	}
	if ScaleUVI(35, 100) != 3 || ScaleUVI(200, 100) != 15 || ScaleUVI(InvalidByte, 100) != InvalidByte {
		t.Error("UVI")
	}
	if ScaleLight(123456, 100) != 123456 || ScaleLight(4000000, 100) != 3000000 {
		t.Error("light")
	}
	if ClampTemp(1490, 20) != 1500 || ClampTemp(-590, -20) != -600 {
		t.Error("temp clamp")
	}
	if ClampHum(98, 5) != 99 || ClampHum(2, -5) != 1 || ClampHum(InvalidHum, 3) != InvalidHum {
		t.Error("hum clamp")
	}
}

func TestWH24UVIndex(t *testing.T) {
	for _, tc := range []struct {
		uv   uint16
		want uint8
	}{{0, 0}, {0x62, 0}, {0x63, 1}, {0xaad, 6}, {0xaae, 7}, {0x1262, 10}, {0x1263, 11}, {0x1883, 13}, {0x1884, 14}, {0x1a7a, 15}, {20000, 15}, {20001, 0xff}} {
		if got := WH24UVIndex(tc.uv); got != tc.want {
			t.Errorf("UVIndex(%#x) = %d want %d", tc.uv, got, tc.want)
		}
	}
}
