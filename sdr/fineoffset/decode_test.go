package fineoffset

import (
	"errors"
	"reflect"
	"testing"
)

// pkt builds an n-byte packet from the given prefix and seals it.
func pkt(n, crcLen int, b ...byte) []byte {
	p := make([]byte, n)
	copy(p, b)
	return Seal(p, crcLen, n)
}

func TestCRC8TableMatchesFirmware(t *testing.T) {
	// first 16 bytes of the table at 0xa1450
	want := []byte{0x00, 0x31, 0x62, 0x53, 0xc4, 0xf5, 0xa6, 0x97, 0xb9, 0x88, 0xdb, 0xea, 0x7d, 0x4c, 0x1f, 0x2e}
	if !reflect.DeepEqual(crc8Table[:16], want) {
		t.Fatalf("table = % x", crc8Table[:16])
	}
	if CRC8([]byte{0x01}) != 0x31 || Sum8([]byte{0xff, 0x02}) != 0x01 {
		t.Fatal("CRC8/Sum8")
	}
}

func wh24Base() []byte {
	return []byte{0x24, 0x7b, 0x2c, 0x92, 0x8a, 55, 0x10, 20, 0x01, 0x23, 0x04, 0x56, 0x01, 0xe2, 0x40}
}

func wh24Ext() []byte {
	p := make([]byte, 0x19)
	copy(p, wh24Base())
	Seal(p, 0x10, 0x11)
	p[0x11], p[0x12], p[0x13] = 0x81, 0x8b, 0xcd // 101325 -> 10132
	p[0x14] = p[0x11] + p[0x12] + p[0x13]
	p[0x15], p[0x16] = 0x0a, 0x0b
	return Seal(p, 0x18, 0x19)
}

func ws90Packet(n int, version byte) []byte {
	p := make([]byte, 0x27)
	copy(p, []byte{0x90, 0x01, 0x02, 0x03, 0x03, 0xe8, 0x9b, 0x02, 0x8a, 50, 30, 90, 40, 12, 0x3f, 0xff,
		0x51, 0x23, 0x45, 0x00, 0x10, 0xc0, 0x22, 0x33, 0xab, 0xcd, 0xef, 0x44, 0x55, version})
	Seal(p, 0x1f, 0x20)
	if n == 0x27 {
		Seal(p, 0x26, 0x27)
	}
	return p
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		want Measurement
	}{
		{"indoor THP 0xE5", pkt(7, 0, 0xe5, 0x32, 0x8a, 45, 0x27, 0x94),
			THP{Header: Header{0xe5, 0x35}, Temp: 250, Hum: 45, Pressure: 10132}},
		{"outdoor TH 0xD1", pkt(7, 0, 0xd1, 0x4a, 0x58, 60, 0xff, 0xff),
			THP{Header: Header{0xd1, 0x41}, Temp: 200, Hum: 60, Pressure: InvalidU16, BattLow: true}},
		{"WH24 17-byte", pkt(0x11, 0x10, wh24Base()...),
			WH24{Header: Header{0x24, 0x7b}, WindDir: 300, Temp: 250, Hum: 55, WindRaw: 0x110, GustRaw: 20,
				Rain: 0x0123, UVRaw: 0x0456, Light: 123456, Pressure: InvalidU16}},
		{"WH24 25-byte", wh24Ext(),
			WH24{Header: Header{0x24, 0x0a0b7b}, Extended: true, WindDir: 300, Temp: 250, Hum: 55, WindRaw: 0x510,
				GustRaw: 20, Rain: 0x0123, UVRaw: 0x0456, Light: 123456, Pressure: 10132}},
		{"WH24 no temp, gust/uv/light absent", pkt(0x11, 0x10, 0x24, 0x01, 0x10, 0x47, 0xff, 30, 0x05, 0xff, 0, 0, 0xff, 0xff, 0xff, 0xff, 0xff),
			WH24{Header: Header{0x24, 1}, WindDir: 0x10, Temp: Invalid16, Hum: InvalidHum, WindRaw: 0x005, GustRaw: InvalidU16,
				UVRaw: InvalidU16, Light: InvalidLight, Pressure: InvalidU16}},
		{"WS68", pkt(0x10, 0x0f, 0x68, 0x2a, 0x12, 0x34, 0x03, 0xe8, 0x50, 0x32, 0x8a, 70, 0x05, 0x0e, 0x30, 35),
			Anemometer{Header: Header{0x68, 0xa1234}, Nibble: 2, Light: 100000, BattRaw: 0x50, Temp: 250, Hum: 70,
				Wind: 0x105, Dir: 270, Gust: 0x30, UV: 35, Aux: InvalidFFFF}},
		{"WS80", pkt(0x12, 0x11, 0x80, 0x01, 0x00, 0x42, 0x75, 0x31, 0x80, 0xfb, 0xff, 0, 0xff, 0x80, 0xff, 200, 0x3f, 0xff),
			Anemometer{Header: Header{0x80, 0x10042}, Light: InvalidLight, BattRaw: 0x80, Temp: Invalid16, Hum: InvalidHum,
				Wind: InvalidU16, Dir: Invalid16, Gust: InvalidU16, UV: InvalidByte, Aux: InvalidFFFF, BattFlags: 3}},
		{"WH40", pkt(9, 8, 0x40, 0x05, 0x67, 0x89, 0x8e, 0x01, 0x02),
			WH40{Header: Header{0x40, 0x56789}, BattRaw: 14, BattReported: true, Rain: 0x0102}},
		{"WH31 ch2", pkt(7, 6, 0x30, 0xc7, 0xa2, 0x8a, 48),
			WH31{Header: Header{0x30, 0xc7}, Ch: 2, Temp: 250, HasHum: true, Hum: 48}},
		{"WH31 0x37 no hum, batt low", pkt(7, 6, 0x37, 0x10, 0x09, 0x00, 99),
			WH31{Header: Header{0x37, 0x10}, Ch: 0, Temp: -144, Hum: InvalidHum, BattLow: true}},
		{"PM 0x42 short", pkt(8, 7, 0x42, 0x33, 0x40, 123, 0x40, 200),
			PM25{Header: Header{0x42, 0x33}, BattLevel: 5, PM25: 120, PM10: 200}},
		{"PM 0x44 long", pkt(10, 9, 0x44, 0x0f, 0xff, 0xff, 0xff, 0xff, 0xc0, 0x05),
			PM25{Header: Header{0x44, 0xfffff}, Long: true, BattLevel: 0x0f, PM25: InvalidFFFF, PM10: 10}},
		{"WH51", pkt(14, 13, 0x51, 0x81, 0x22, 0x33, 0x0e, 0, 42, 0x03, 0x21),
			WH51{Header: Header{0x51, 0x12233}, Nibble: 8, Moisture: 42, BattRaw: 14, AD: 0x321}},
		{"WH55 leak", pkt(10, 9, 0x55, 0x21, 0x02, 0x03, 4, 0x12, 0x34, 0xc1),
			WH55{Header: Header{0x55, 0x10203}, Ch: 2, BattLevel: 4, Value: 0x1234, Leak: true, Flag6: true}},
		{"WH57", pkt(9, 8, 0x57, 0x45, 0x55, 0x66, 3, 0xd5, 7),
			WH57{Header: Header{0x57, 0x55566}, State: 2, BattLevel: 3, DistanceKm: 0x15, Count: 7}},
		{"WH57 far", pkt(9, 8, 0x57, 0x10, 0, 1, 1, 0x3f, 0),
			WH57{Header: Header{0x57, 1}, State: 1, BattLevel: 1, DistanceKm: InvalidByte}},
		{"WH45", pkt(15, 14, 0x45, 0x11, 0x22, 0x33, 0x02, 0x8a, 50, 0x80, 0x64, 0x40, 0xc8, 0x01, 0xf4),
			WH45{Header: Header{0x45, 0x112233}, Temp: 250, Hum: 50, BattLevel: 9, PM25: 100, PM10: 200, CO2: 500,
				Ext10: InvalidFFFF, Ext12: InvalidFFFF, Ext14: InvalidFFFF}},
		{"WH46", pkt(21, 20, 0x46, 0x11, 0x22, 0x33, 0x07, 0xff, 0xff, 0x3f, 0xff, 0x3f, 0xff, 0x01, 0xf4, 0x00, 0x0a, 0x3f, 0xff, 0x12, 0x34),
			wh46Expect()},
		{"WN34 offset format", pkt(9, 8, 0x34, 0x00, 0x01, 0x02, 0x02, 0x8a, 0x50),
			WN34{Header: Header{0x34, 0x102}, Temp: 250, BattRaw: 0x50}},
		{"WN34 signed negative", pkt(9, 8, 0x34, 0x00, 0x01, 0x02, 0x4f, 0x9c, 0x50),
			WN34{Header: Header{0x34, 0x102}, Temp: -100, Signed: true, BattRaw: 0x50}},
		{"WN34 absent", pkt(9, 8, 0x34, 0x00, 0x01, 0x02, 0x07, 0xff, 0x50),
			WN34{Header: Header{0x34, 0x102}, Temp: Invalid16, BattRaw: 0x50}},
		{"WH48", pkt(14, 13, 0x48, 0x31, 0x00, 0x05, 0x1a, 0x8a, 40, 0x01, 0xf4, 0x00, 0x07),
			WH48{Header: Header{0x48, 0x10005}, Ch: 3, Temp: 250, BattLow: true, Flag3: true, Hum: 40, CO2: 500, Value2: 7}},
		{"WN35 with temp", pkt(11, 10, 0x35, 0x0a, 0x0b, 0x0c, 0x42, 0x8a, 0, 77, 0x40),
			WN35{Header: Header{0x35, 0x0a0b0c}, HasTemp: true, Temp: 250, Wetness: 77, BattRaw: 0x40}},
		{"WN35 no temp", pkt(11, 10, 0x35, 0x0a, 0x0b, 0x0c, 0x00, 0x8a, 0, 100, 0x40),
			WN35{Header: Header{0x35, 0x0a0b0c}, Temp: Invalid16, Wetness: InvalidHum, BattRaw: 0x40}},
		{"WS85", pkt(0x1c, 0x1b, 0x85, 0x01, 0x02, 0x03, 0x9b, 0x70, 0, 0x2c, 0x0e, 0xff, 0, 0, 0x01, 0x23, 0x45, 0x00, 0x99),
			WS85{Header: Header{0x85, 0x010203}, BattRaw: 0x9b, Flags: 0x70, Wind: 0x12c, Dir: 270, Gust: InvalidU16,
				RainCounter: 0x12345, Rain2: 0x0099}},
		{"WH54", wh54Packet(), wh54Expect()},
		{"WS90 32-byte", ws90Packet(0x20, 0x90), ws90Expect(0x20, 0x90)},
		{"WS90 39-byte", ws90Packet(0x27, 0x95), ws90Expect(0x27, 0x95)},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Decode(tc.in)
			if err != nil {
				t.Fatalf("Decode: %v", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got  %+v\nwant %+v", got, tc.want)
			}
		})
	}
}

func wh46Expect() WH45 {
	return WH45{Header: Header{0x46, 0x112233}, Temp: Invalid16, Hum: 0xff, BattLevel: 0, PM25: InvalidFFFF,
		PM10: InvalidFFFF, CO2: 500, Ext10: 10, Ext12: InvalidFFFF, Ext14: 0x1234}
}

func wh54Packet() []byte {
	return pkt(0x13, 0x12, 0x54, 0x31, 0x02, 0x03, 0x40, 0x85, 0xc4, 0xd2, 0xc1, 0x02, 0x00, 0x10, 0xff, 0xff, 0xff, 0x05)
}

func wh54Expect() WH54 {
	return WH54{Header: Header{0x54, 0x10203}, Nibble: 3, BattRaw: 0x40, Flag11: true, F13: 5, Distance: 0x04d2,
		Flag12: true, Type: 1, F6: 0x0102, FA: 0x0010, FC: InvalidCounter, F14: 0x305, F17: wh54Packet()[17]}
}

func ws90Expect(n int, ver byte) WS90 {
	return WS90{
		Anemometer: Anemometer{Header: Header{0x90, 0x010203}, Light: 100000, BattRaw: 0x9b, Temp: 250, Hum: 50,
			Wind: 30, Dir: 90, Gust: 40, UV: 12, Aux: InvalidFFFF},
		Length: n, F20: 5, RainCounter: 0x12345, Rain2: 0x0010, F22: 3, F23: 0x22, F24: 0x33,
		F1A: 0xabc, F1C: 0xdef, F1E: 0x44, F1F: 0x55, Version: ver,
	}
}

func TestDecodeWH54CRCField(t *testing.T) {
	// F17 is the byte at the CRC position, which the firmware stores.
	p := pkt(0x13, 0x12, 0x54, 0x31, 0x02, 0x03)
	m, err := Decode(p)
	if err != nil {
		t.Fatal(err)
	}
	if m.(WH54).F17 != p[17] {
		t.Fatalf("F17 = %#x, want %#x", m.(WH54).F17, p[17])
	}
}

func TestDecodeErrors(t *testing.T) {
	good := pkt(9, 8, 0x40, 0x05, 0x67, 0x89, 0x8e, 0x01, 0x02)
	badCRC := append([]byte(nil), good...)
	badCRC[7] ^= 1
	badSum := append([]byte(nil), good...)
	badSum[8] ^= 1
	extBadSum := wh24Ext()
	extBadSum[0x14]++
	Seal(extBadSum, 0x18, 0x19)
	ws80aux := pkt(0x12, 0x11, 0x80, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x00, 0x10)
	ws90dir := ws90Packet(0x20, 1)
	ws90dir[7] |= 0x20
	ws90dir[11] = 0x80 // 0x180 > 359
	Seal(ws90dir, 0x1f, 0x20)
	ws90gust := ws90Packet(0x20, 1)
	ws90gust[12] = 5
	Seal(ws90gust, 0x1f, 0x20)
	tests := []struct {
		name string
		in   []byte
		want error
	}{
		{"empty", nil, ErrEmpty},
		{"unknown", []byte{0x99, 0, 0}, ErrUnknownFamily},
		{"short", good[:5], ErrShort},
		{"crc", badCRC, ErrCRC},
		{"sum", badSum, ErrChecksum},
		{"D with pressure", pkt(7, 0, 0xd1, 0x4a, 0x58, 60, 0x27, 0x94), ErrRejected},
		{"E bad hum", pkt(7, 0, 0xe5, 0x32, 0x8a, 0, 0x27, 0x94), ErrRejected},
		{"E bad pressure", pkt(7, 0, 0xe5, 0x32, 0x8a, 45, 0x01, 0x00), ErrRejected},
		{"WH24 ext baro sum", extBadSum, ErrRejected},
		{"WS80 aux range", ws80aux, ErrRejected},
		{"WS90 dir", ws90dir, ErrRejected},
		{"WS90 gust<speed", ws90gust, ErrRejected},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Decode(tc.in)
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
		})
	}
}

func TestLayoutWS90ExtendedNeedsVersion(t *testing.T) {
	p := ws90Packet(0x27, 0x93) // ext checks valid but version < 0x94
	if _, n, _ := Layout(p); n != 0x20 {
		t.Fatalf("len = %d", n)
	}
}

func TestLowBattery(t *testing.T) {
	tests := []struct {
		m    Measurement
		want bool
	}{
		{Anemometer{Header: Header{Raw0: 0x68}, BattRaw: 0x3c}, true},
		{Anemometer{Header: Header{Raw0: 0x68}, BattRaw: 0x3d}, false},
		{Anemometer{Header: Header{Raw0: 0x80}, BattRaw: 0x78}, true},
		{WH40{BattRaw: 5}, false},
		{WH40{BattRaw: 5, BattReported: true}, true},
		{PM25{BattLevel: 1}, true},
		{WH54{BattRaw: 0x50, Type: 0}, true},
		{WH54{BattRaw: 0x50, Type: 1}, false},
	}
	for i, tc := range tests {
		if got := tc.m.LowBattery(); got != tc.want {
			t.Errorf("%d: %T LowBattery = %v", i, tc.m, got)
		}
	}
}
