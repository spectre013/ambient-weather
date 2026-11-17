package fineoffset

import (
	"errors"
	"fmt"
)

// Decode errors. Rejections mirror firmware paths that drop a packet that
// passed its checksums.
var (
	ErrEmpty         = errors.New("rf: empty payload")
	ErrUnknownFamily = errors.New("rf: unknown family")
	ErrShort         = errors.New("rf: payload too short")
	ErrCRC           = errors.New("rf: CRC-8 mismatch")
	ErrChecksum      = errors.New("rf: sum mismatch")
	ErrRejected      = errors.New("rf: packet rejected")
)

func reject(why string) error { return fmt.Errorf("%w: %s", ErrRejected, why) }

// Layout returns the CRC and sum check lengths (fn_012de4/fn_012e34
// arguments) for a payload. crcLen 0 means sum only. The packet length is
// sumLen.
func Layout(p []byte) (crcLen, sumLen int, ok bool) {
	if len(p) == 0 {
		return 0, 0, false
	}
	switch f := FamilyOf(p[0]); f {
	case FamilyIndoorTHP, FamilyOutdoorTH:
		return 0, 7, true
	case FamilyWH24:
		if crcOK(p, 0x18) && sumOK(p, 0x19) {
			return 0x18, 0x19, true
		}
		return 0x10, 0x11, true
	case FamilyWS68:
		return 0x0f, 0x10, true
	case FamilyWS80:
		return 0x11, 0x12, true
	case FamilyWH40, FamilyWH57, FamilyWN34:
		return 8, 9, true
	case FamilyWH31, FamilyWH31Alt:
		return 6, 7, true
	case FamilyPM41, FamilyPM42:
		return 7, 8, true
	case FamilyPM43, FamilyPM44:
		return 9, 10, true
	case FamilyWH51, FamilyWH48:
		return 0x0d, 0x0e, true
	case FamilyWH55:
		return 9, 10, true
	case FamilyWH45:
		return 0x0e, 0x0f, true
	case FamilyWH46:
		return 0x14, 0x15, true
	case FamilyWN35:
		return 10, 11, true
	case FamilyWS90:
		if len(p) > 0x1d && p[0x1d] >= 0x94 && crcOK(p, 0x26) && sumOK(p, 0x27) {
			return 0x26, 0x27, true
		}
		return 0x1f, 0x20, true
	case FamilyWS85:
		return 0x1b, 0x1c, true
	case FamilyWH54:
		return 0x12, 0x13, true
	}
	return 0, 0, false
}

// Decode verifies and decodes one payload as docTask does. The payload may
// be longer than the packet (the radio always delivers 0x20/0x27 bytes).
func Decode(p []byte) (Measurement, error) {
	if len(p) == 0 {
		return nil, ErrEmpty
	}
	crcLen, sumLen, ok := Layout(p)
	if !ok {
		return nil, fmt.Errorf("%w 0x%02x", ErrUnknownFamily, p[0])
	}
	if len(p) < sumLen {
		return nil, ErrShort
	}
	// WH24 and WS90 extended forms must also pass the base checks first.
	switch FamilyOf(p[0]) {
	case FamilyWH24:
		if !crcOK(p, 0x10) {
			return nil, ErrCRC
		}
		if !sumOK(p, 0x11) {
			return nil, ErrChecksum
		}
	case FamilyWS90:
		if !crcOK(p, 0x1f) {
			return nil, ErrCRC
		}
		if !sumOK(p, 0x20) {
			return nil, ErrChecksum
		}
	}
	if crcLen > 0 && !crcOK(p, crcLen) {
		return nil, ErrCRC
	}
	if !sumOK(p, sumLen) {
		return nil, ErrChecksum
	}
	switch f := FamilyOf(p[0]); f {
	case FamilyIndoorTHP, FamilyOutdoorTH:
		return decodeTHP(p)
	case FamilyWH24:
		return decodeWH24(p, sumLen == 0x19)
	case FamilyWS68, FamilyWS80:
		return decodeAnemometer(p)
	case FamilyWH40:
		return decodeWH40(p), nil
	case FamilyWH31, FamilyWH31Alt:
		return decodeWH31(p), nil
	case FamilyPM41, FamilyPM42, FamilyPM43, FamilyPM44:
		return decodePM(p), nil
	case FamilyWH51:
		return decodeWH51(p), nil
	case FamilyWH55:
		return decodeWH55(p), nil
	case FamilyWH57:
		return decodeWH57(p), nil
	case FamilyWH45, FamilyWH46:
		return decodeWH45(p), nil
	case FamilyWN34:
		return decodeWN34(p), nil
	case FamilyWH48:
		return decodeWH48(p), nil
	case FamilyWN35:
		return decodeWN35(p), nil
	case FamilyWS90:
		return decodeWS90(p, sumLen)
	case FamilyWS85:
		return decodeWS85(p), nil
	case FamilyWH54:
		return decodeWH54(p), nil
	}
	return nil, ErrUnknownFamily
}

func be16(hi, lo byte) uint16 { return uint16(hi)<<8 | uint16(lo) }

// id20 is (b1&0xf)<<16 | b2<<8 | b3.
func id20(p []byte) uint32 { return uint32(p[1]&0x0f)<<16 | uint32(p[2])<<8 | uint32(p[3]) }

// id24 is b1<<16 | b2<<8 | b3.
func id24(p []byte) uint32 { return uint32(p[1])<<16 | uint32(p[2])<<8 | uint32(p[3]) }

// temp11 converts an 11-bit raw temperature (offset 400) with the
// "raw+200 < 0x835" validity test used by most families.
func temp11(raw uint16) int16 {
	if raw+200 < 0x835 {
		return int16(raw) - 400
	}
	return Invalid16
}

// temp1000 is the WS68/WS80/WS90 variant: valid only for raw <= 1000.
func temp1000(raw uint16) int16 {
	if raw < 0x3e9 {
		return int16(raw) - 400
	}
	return Invalid16
}

func hum99(b byte) uint8 {
	if b-1 > 0x62 { // also catches 0 via wraparound
		return InvalidHum
	}
	return b
}

// wind9 reproduces the 9-bit speed/gust decode: low byte plus an MSB flag;
// 0x1ff means absent, else capped at 500.
func wind9(lo byte, msb bool) uint16 {
	v := uint16(lo)
	if msb {
		v |= 0x100
		if v == 0x1ff {
			return InvalidU16
		}
	}
	if v > 500 {
		v = 500
	}
	return v
}

// dir9: direction low byte plus MSB flag; > 359 is absent (only reachable
// with the flag set).
func dir9(lo byte, msb bool) int16 {
	v := int16(lo)
	if msb {
		v |= 0x100
		if v > 0x167 {
			return Invalid16
		}
	}
	return v
}

func light16(hi, lo byte) uint32 {
	if v := be16(hi, lo); v <= 30000 {
		return uint32(v) * 100
	}
	return InvalidLight
}

// docTask 0xDx / 0xEx blocks.
func decodeTHP(p []byte) (Measurement, error) {
	m := THP{
		Header:   Header{Raw0: p[0], ID: uint32(p[1]&0xf0 | p[0]&0x0f)},
		Temp:     temp11(uint16(p[2]) | uint16(p[1]&7)<<8),
		Hum:      hum99(p[3]),
		BattLow:  p[1]&0x08 != 0,
		Pressure: InvalidU16,
	}
	if m.Temp == Invalid16 || m.Hum == InvalidHum {
		return nil, reject("temperature/humidity out of range")
	}
	pr := be16(p[4], p[5])
	if p[0]&0xF0 == 0xE0 {
		if pr-3000 >= 0x2329 {
			return nil, reject("pressure out of range")
		}
		m.Pressure = pr
	} else if pr != 0xffff {
		return nil, reject("0xDx packet carries pressure bytes")
	}
	return m, nil
}

// fn_00da74.
func decodeWH24(p []byte, ext bool) (Measurement, error) {
	m := WH24{Header: Header{Raw0: p[0], ID: uint32(p[1])}, Extended: ext, Pressure: InvalidU16}
	if ext {
		m.ID = uint32(p[0x15])<<16 | uint32(p[0x16])<<8 | uint32(p[1])
	}
	m.WindDir = dir9(p[2], p[3]&0x80 != 0)
	if p[3]&7 == 7 && p[4] == 0xff {
		m.Temp, m.Hum = Invalid16, InvalidHum
	} else {
		m.Temp = temp11(be16(p[3], p[4]) & 0x7ff)
		if p[5] == 0xff {
			m.Hum = InvalidHum
		} else {
			m.Hum = hum99(p[5])
		}
	}
	var hi uint16
	if p[3]&0x40 != 0 {
		hi = uint16(p[3]&0x10) << 4
	} else {
		hi = uint16(p[3]&0x30) << 4
	}
	m.WindRaw = hi | uint16(p[6])
	if ext && p[0x11]&0x80 != 0 {
		m.WindRaw |= 0x400
	}
	m.GustRaw = uint16(p[7])
	if p[7] == 0xff {
		m.GustRaw = InvalidU16
	}
	m.Rain = be16(p[8], p[9])
	m.UVRaw = be16(p[10], p[11])
	if m.UVRaw == 0xffff {
		m.UVRaw = InvalidU16
	}
	m.Light = uint32(p[12])<<16 | uint32(p[13])<<8 | uint32(p[14])
	if m.Light == 0xffffff {
		m.Light = InvalidLight
	}
	m.BattLow = p[3]&0x08 != 0
	if ext {
		if p[0x14] != p[0x11]+p[0x12]+p[0x13] {
			return nil, reject("WH24 barometer sum")
		}
		raw := uint32(p[0x11]&1)<<16 | uint32(p[0x12])<<8 | uint32(p[0x13])
		if raw != 0x1ffff {
			v := raw / 10
			if v-3000 > 9000 {
				return nil, reject("WH24 pressure out of range")
			}
			m.Pressure = uint16(v)
		}
	}
	return m, nil
}

// fillAnemometer decodes the b1..b13 layout shared by 0x68, 0x80 and 0x90.
func fillAnemometer(m *Anemometer, p []byte) {
	m.Raw0 = p[0]
	m.Nibble = p[1] >> 4
	m.Light = light16(p[4], p[5])
	m.BattRaw = p[6]
	f := p[7]
	m.Temp = temp1000(uint16(p[8]) | uint16(f&7)<<8)
	m.Hum = hum99(p[9])
	m.Wind = wind9(p[10], f&0x10 != 0)
	m.Dir = dir9(p[11], f&0x20 != 0)
	m.Gust = wind9(p[12], f&0x40 != 0)
	m.UV = p[13]
	if m.UV > 0xa0 {
		m.UV = InvalidByte
	}
	m.Aux = InvalidFFFF
}

// docTask 0x68 / 0x80 blocks.
func decodeAnemometer(p []byte) (Measurement, error) {
	var m Anemometer
	fillAnemometer(&m, p)
	m.ID = id20(p)
	if p[0] == byte(FamilyWS80) {
		aux := uint16(p[15]) | uint16(p[14]&0x3f)<<8
		if aux == 0x3fff {
			aux = InvalidFFFF
		} else if aux-3000 > 9000 {
			return nil, reject("WS80 aux out of range")
		}
		m.Aux = aux
		m.AuxHi = p[14] >> 6
		m.BattFlags = (p[7]>>3)&1 | (p[7]>>6)&2
	}
	return m, nil
}

// fn_00d4ec.
func decodeWS90(p []byte, n int) (Measurement, error) {
	m := WS90{Length: n}
	fillAnemometer(&m.Anemometer, p)
	m.ID = id24(p)
	f := p[7]
	if d := uint16(p[11]) | uint16(f&0x20)<<3; d == 0x1ff || d > 0x167 {
		return nil, reject("WS90 direction")
	}
	if m.Gust < m.Wind {
		return nil, reject("WS90 gust < speed")
	}
	aux := uint16(p[15]) | uint16(p[14]&0x3f)<<8
	if aux == 0x3fff {
		aux = InvalidFFFF
	} else if aux-3000 > 9000 {
		return nil, reject("WS90 aux out of range")
	}
	m.Aux = aux
	m.BattFlags = (f>>3)&1 | (f>>6)&2
	m.AuxHi = p[14] >> 6
	m.F20 = p[16] >> 4
	m.RainCounter = uint32(p[16]&0x0f)<<16 | uint32(p[17])<<8 | uint32(p[18])
	m.Rain2 = be16(p[19], p[20])
	m.F22 = p[21] >> 6
	m.F23 = p[22]
	m.F24 = p[23]
	m.F1A = uint16(p[24])<<4 | uint16(p[25]>>4)
	m.F1C = uint16(p[26]) | uint16(p[25]&0x0f)<<8
	m.F1E, m.F1F = p[27], p[28]
	m.Version = p[29]
	return m, nil
}

// docTask 0x85 block.
func decodeWS85(p []byte) Measurement {
	f := p[5]
	m := WS85{Header: Header{Raw0: p[0], ID: id24(p)}, BattRaw: p[4], Flags: f}
	m.Wind = wind9(p[7], f&0x10 != 0)
	m.Dir = dir9(p[8], f&0x20 != 0)
	m.Gust = wind9(p[9], f&0x40 != 0)
	m.RainCounter = uint32(p[12]&0x0f)<<16 | uint32(p[13])<<8 | uint32(p[14])
	m.Rain2 = be16(p[15], p[16])
	return m
}

// docTask 0x40 block.
func decodeWH40(p []byte) Measurement {
	return WH40{
		Header:       Header{Raw0: p[0], ID: id20(p)},
		BattRaw:      p[4] & 0x7f,
		BattReported: p[4]&0x80 != 0,
		Rain:         be16(p[5], p[6]),
	}
}

// docTask 0x30/0x37 block.
func decodeWH31(p []byte) Measurement {
	b2 := p[2]
	m := WH31{
		Header:  Header{Raw0: p[0], ID: uint32(p[1])},
		Ch:      (b2 & 0x7f) >> 4,
		Temp:    temp11(uint16(p[3]) | uint16(b2&7)<<8),
		HasHum:  b2&0x80 != 0,
		Hum:     InvalidHum,
		BattLow: b2&0x08 != 0,
	}
	if m.HasHum {
		m.Hum = hum99(p[4])
	}
	return m
}

func pmRound(raw uint16) uint16 {
	if raw == 0x3fff {
		return InvalidFFFF
	}
	return (raw + 5) / 10 * 10
}

// docTask 0x41..0x44 block.
func decodePM(p []byte) Measurement {
	m := PM25{Header: Header{Raw0: p[0]}}
	var pm25, pm10 uint16
	if p[0] < 0x43 {
		m.ID = uint32(p[1])
		m.BattLevel = p[4]>>6 | (p[2]&0xc0)>>4
		pm25 = uint16(p[3]) | uint16(p[2]&0x3f)<<8
		pm10 = uint16(p[5]) | uint16(p[4]&0x3f)<<8
	} else {
		m.Long = true
		m.ID = id20(p)
		m.BattLevel = p[6]>>6 | (p[4]&0xc0)>>4
		pm25 = uint16(p[5]) | uint16(p[4]&0x3f)<<8
		pm10 = uint16(p[7]) | uint16(p[6]&0x3f)<<8
	}
	m.PM25, m.PM10 = pmRound(pm25), pmRound(pm10)
	return m
}

// docTask 0x51 block.
func decodeWH51(p []byte) Measurement {
	m := WH51{
		Header:  Header{Raw0: p[0], ID: id20(p)},
		Nibble:  p[1] >> 4,
		BattRaw: p[4] & 0x1f,
		AD:      uint16(p[8]) | uint16(p[7]&7)<<8,
		Moisture: func() uint8 {
			if p[6] <= 100 {
				return p[6]
			}
			return InvalidHum
		}(),
	}
	return m
}

// docTask 0x55 block.
func decodeWH55(p []byte) Measurement {
	return WH55{
		Header:    Header{Raw0: p[0], ID: id20(p)},
		Ch:        p[1] >> 4,
		BattLevel: p[4],
		Value:     be16(p[5], p[6]),
		Leak:      p[7]&0x80 != 0,
		Flag6:     p[7]&0x40 != 0,
	}
}

// docTask 0x57 block.
func decodeWH57(p []byte) Measurement {
	b1 := p[1]
	var st uint8
	switch {
	case b1&0x80 != 0:
		st = 3
	case b1&0x40 != 0:
		st = 2
	default:
		st = (b1 >> 4) & 1
	}
	d := p[5] & 0x3f
	if d > 0x28 {
		d = InvalidByte
	}
	return WH57{Header: Header{Raw0: p[0], ID: id20(p)}, State: st, BattLevel: p[4], DistanceKm: d, Count: p[6]}
}

// docTask 0x45/0x46 block.
func decodeWH45(p []byte) Measurement {
	m := WH45{Header: Header{Raw0: p[0], ID: id24(p)}}
	if t := uint16(p[5]) | uint16(p[4]&7)<<8; t == 0x7ff {
		m.Temp = Invalid16
	} else {
		m.Temp = int16(t) - 400
	}
	m.Hum = p[6]
	m.BattLevel = p[9]>>6 | (p[7]&0xc0)>>4
	ff := func(v uint16) uint16 {
		if v == 0x3fff {
			return InvalidFFFF
		}
		return v
	}
	m.PM25 = ff(uint16(p[8]) | uint16(p[7]&0x3f)<<8)
	m.PM10 = ff(uint16(p[10]) | uint16(p[9]&0x3f)<<8)
	m.CO2 = be16(p[11], p[12])
	m.Ext10, m.Ext12, m.Ext14 = InvalidFFFF, InvalidFFFF, InvalidFFFF
	if p[0] == byte(FamilyWH46) {
		m.Ext10 = ff(be16(p[13], p[14]))
		m.Ext12 = ff(be16(p[15], p[16]))
		m.Ext14 = be16(p[17], p[18])
	}
	return m
}

// docTask 0x34 block.
func decodeWN34(p []byte) Measurement {
	b4 := p[4]
	raw := uint16(p[5]) | uint16(b4&0x0f)<<8
	m := WN34{Header: Header{Raw0: p[0], ID: id24(p)}, BattRaw: p[6]}
	switch {
	case raw == 0x7ff:
		m.Temp = Invalid16
	case b4&0x40 == 0:
		m.Temp = int16(raw) - 400
	default:
		m.Signed = true
		m.Temp = int16(raw)
		if b4&0x08 != 0 {
			m.Temp = int16(uint16(p[5]) | uint16(b4|0xf0)<<8)
		}
	}
	return m
}

// docTask 0x48 block.
func decodeWH48(p []byte) Measurement {
	b4 := p[4]
	m := WH48{
		Header:  Header{Raw0: p[0], ID: id20(p)},
		Ch:      p[1] >> 4,
		BattLow: b4&0x10 != 0,
		Flag3:   b4&0x08 != 0,
		Hum:     p[6],
		CO2:     be16(p[7], p[8]),
		Value2:  be16(p[9], p[10]),
	}
	if t := uint16(p[5]) | uint16(b4&7)<<8; t == 0x7ff {
		m.Temp = Invalid16
	} else {
		m.Temp = int16(t) - 400
	}
	return m
}

// docTask 0x35 block.
func decodeWN35(p []byte) Measurement {
	m := WN35{Header: Header{Raw0: p[0], ID: id24(p)}, BattRaw: p[8], Temp: Invalid16, Wetness: InvalidHum}
	if p[4]&0x40 != 0 {
		m.HasTemp = true
		m.Temp = temp11(uint16(p[5]) | uint16(p[4]&7)<<8)
	}
	if p[7] < 100 {
		m.Wetness = p[7]
	}
	return m
}

// docTask 0x54 block.
func decodeWH54(p []byte) Measurement {
	m := WH54{
		Header:  Header{Raw0: p[0], ID: id20(p)},
		Nibble:  p[1] >> 4,
		BattRaw: p[4],
		Flag11:  p[5]&0x80 != 0,
		F13:     p[5] & 0x7f,
		Flag12:  p[6]&0x40 != 0,
		Type:    p[6] >> 7,
		F6:      uint16(p[9]) | uint16(p[8]&0x3f)<<8,
		FA:      be16(p[10], p[11]),
		F14:     uint16(p[15]) | uint16(p[8]&0xc0)<<2,
		F17:     p[17],
	}
	m.Distance = uint16(p[7]) | uint16(p[6]&0x3f)<<8
	if m.Distance == 0x3fff {
		m.Distance = InvalidU16
	}
	m.FC = uint32(p[12])<<16 | uint32(p[13])<<8 | uint32(p[14])
	if m.FC == 0xffffff {
		m.FC = InvalidCounter
	}
	return m
}
