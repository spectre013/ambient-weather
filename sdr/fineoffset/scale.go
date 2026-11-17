package fineoffset

// Calibration helpers reproducing the integer arithmetic docTask applies
// before publishing to the live-data block at 0x1f11d94. Gains are the
// percentage words of the calibration block at 0x1f11460 (100 = 1.0):
// +0x0e wind, +0x12 UV, +0x14 light; +0x0c is the direction offset.

// ScaleWH24Wind converts WH24 raw speed to 0.1 m/s: (g*raw*51+4000)/8000,
// capped at 800 (fn_00da74).
func ScaleWH24Wind(raw, gainPct uint16) uint16 {
	if raw == InvalidU16 {
		return raw
	}
	return min(uint16((uint32(gainPct)*uint32(raw)*51+4000)/8000), 800)
}

// ScaleWH24Gust converts WH24 raw gust to 0.1 m/s: (g*raw*51+500)/1000.
func ScaleWH24Gust(raw, gainPct uint16) uint16 {
	if raw == InvalidU16 {
		return raw
	}
	return min(uint16((uint32(gainPct)*uint32(raw)*51+500)/1000), 800)
}

// ScaleWind applies the WS68/WS80/WS85/WS90 gain: (g*v*10+500)/1000 in
// 0.1 m/s, capped at 800.
func ScaleWind(v, gainPct uint16) uint16 {
	if v == InvalidU16 {
		return v
	}
	return min(uint16((uint32(gainPct)*uint32(v)*10+500)/1000), 800)
}

// OffsetDir adds the direction offset and wraps into 0..359.
func OffsetDir(dir, offset int16) int16 {
	if dir == Invalid16 {
		return dir
	}
	d := int(dir) + int(offset)
	if d < 0 {
		d = d - d/360*360 + 360 // firmware form; yields 360 for exact multiples
	} else if d > 359 {
		d -= d / 360 * 360
	}
	return int16(d)
}

// ScaleUVI applies the UV gain to a 0.1-UVI byte (WS68/WS80/WS90):
// raw*g/1000, capped at 15.
func ScaleUVI(raw uint8, gainPct uint16) uint8 {
	if raw == InvalidByte {
		return raw
	}
	return uint8(min(uint32(raw)*uint32(gainPct)/1000, 15))
}

// ScaleWH24UV applies the UV gain to WH24 raw UV: raw*g/100, capped 20000.
func ScaleWH24UV(raw, gainPct uint16) uint16 {
	if raw == InvalidU16 {
		return raw
	}
	return uint16(min(uint32(raw)*uint32(gainPct)/100, 20000))
}

// ScaleLight applies the light gain: lux10*g/100 (32-bit wrap as on the
// MCU), capped at 3,000,000 (0.1 lux).
func ScaleLight(lux10 uint32, gainPct uint16) uint32 {
	if lux10 == InvalidLight {
		return lux10
	}
	return min(lux10*uint32(gainPct)/100, 3000000)
}

// SolarWm2 is the lux->W/m² conversion used when setting 0x1f19090+0xda
// is 1: (lux10*1000+633)/1267000 (fn_00d4ec, docTask 0x80 block).
func SolarWm2(lux10 uint32) uint32 {
	if lux10 == InvalidLight {
		return 0
	}
	return (lux10*1000 + 0x279) / 0x135538
}

// uviThresholds are the inclusive-lower bounds of UVI 1..15 in fn_00c7a8.
var uviThresholds = [...]uint16{0x63, 0x21c, 0x3e8, 0x578, 0x733, 0x8f4, 0xaae, 0xc42, 0xe40, 0x1064, 0x1263, 0x1459, 0x1667, 0x1884, 0x1a7a}

// WH24UVIndex maps scaled WH24 UV to an index (fn_00c7a8); 0xff if > 20000.
func WH24UVIndex(uv uint16) uint8 {
	if uv > 20000 {
		return InvalidByte
	}
	var i uint8
	for _, t := range uviThresholds {
		if uv < t {
			break
		}
		i++
	}
	return i
}

// ClampTemp applies the calibration offset and the -600..1500 clamp used
// for every published temperature.
func ClampTemp(t, offset int16) int16 {
	if t == Invalid16 {
		return t
	}
	return max(min(t+offset, 1500), -600)
}

// ClampHum applies the humidity offset (signed-char add, as on the MCU)
// and the 1..99 clamp.
func ClampHum(h uint8, offset int8) uint8 {
	if h == InvalidHum {
		return h
	}
	return uint8(max(min(int8(h+uint8(offset)), 99), 1))
}
