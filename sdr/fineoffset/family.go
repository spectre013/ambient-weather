package fineoffset

import "fmt"

// Family is the first payload byte that docTask (fn_00e1c4) dispatches on.
// For FamilyIndoorTHP and FamilyOutdoorTH only the high nibble selects the
// family; the low nibble is part of the sensor ID.
type Family byte

const (
	FamilyWH24       Family = 0x24 // 7-in-1 array (WH24/WH65), fn_00da74
	FamilyWH31       Family = 0x30 // T&H CHn
	FamilyWH31Alt    Family = 0x37 // same layout as 0x30; model unconfirmed
	FamilyWN34       Family = 0x34 // WN34 temperature probe
	FamilyWN35       Family = 0x35 // WN35 leaf wetness
	FamilyWH40       Family = 0x40 // rain gauge
	FamilyPM41       Family = 0x41 // PM2.5, short form
	FamilyPM42       Family = 0x42 // PM2.5, short form
	FamilyPM43       Family = 0x43 // PM2.5, long form (20-bit ID)
	FamilyPM44       Family = 0x44 // PM2.5, long form (20-bit ID)
	FamilyWH45       Family = 0x45 // AQIN CO2/PM/T/H
	FamilyWH46       Family = 0x46 // AQIN, extended
	FamilyWH48       Family = 0x48 // "CO2IN"/"WH48" T/H/CO2 per channel (name from UI strings: guess)
	FamilyWH51       Family = 0x51 // soil moisture
	FamilyWH54       Family = 0x54 // "LDS" laser distance
	FamilyWH55       Family = 0x55 // leak
	FamilyWH57       Family = 0x57 // lightning ("Thunder")
	FamilyWS68       Family = 0x68 // anemometer + light/UV
	FamilyWS80       Family = 0x80 // ultrasonic array
	FamilyWS85       Family = 0x85 // ultrasonic wind + piezo rain
	FamilyWS90       Family = 0x90 // ultrasonic wind + piezo rain + T/H/light/UV
	FamilyOutdoorTH  Family = 0xD0 // high nibble 0xD: outdoor T/H (WH32 class)
	FamilyIndoorTHP  Family = 0xE0 // high nibble 0xE: indoor T/H/baro (WH25/WH32B class)
	familyNibbleMask Family = 0xF0
)

// FamilyOf returns the dispatch family of a payload's first byte, folding
// the 0xDx/0xEx nibble families.
func FamilyOf(b0 byte) Family {
	switch b0 & 0xF0 {
	case 0xD0, 0xE0:
		return Family(b0 & 0xF0)
	}
	return Family(b0)
}

var familyNames = map[Family]string{
	FamilyWH24: "WH24/WH65", FamilyWH31: "WH31", FamilyWH31Alt: "WH31(0x37)",
	FamilyWN34: "WN34", FamilyWN35: "WN35", FamilyWH40: "WH40",
	FamilyPM41: "WH41(0x41)", FamilyPM42: "WH41(0x42)", FamilyPM43: "WH41(0x43)", FamilyPM44: "WH41(0x44)",
	FamilyWH45: "WH45", FamilyWH46: "WH46", FamilyWH48: "WH48", FamilyWH51: "WH51",
	FamilyWH54: "WH54", FamilyWH55: "WH55", FamilyWH57: "WH57", FamilyWS68: "WS68",
	FamilyWS80: "WS80", FamilyWS85: "WS85", FamilyWS90: "WS90",
	FamilyOutdoorTH: "WH32", FamilyIndoorTHP: "WH25/WH32B",
}

func (f Family) String() string {
	if n, ok := familyNames[FamilyOf(byte(f))]; ok {
		return n
	}
	return fmt.Sprintf("Family(0x%02x)", byte(f))
}

// Interval is the per-family "expected transmit period" in seconds. The
// firmware passes the same constant both as the learning window to
// fn_00d48c/fn_00e004 and as the missed-packet period to fn_00e180.
func (f Family) Interval() int {
	switch FamilyOf(byte(f)) {
	case FamilyWH24:
		return 0x11
	case FamilyWS68:
		return 0x12
	case FamilyWS80:
		return 6
	case FamilyWS85, FamilyWS90:
		return 10
	case FamilyWH40:
		return 0x32
	case FamilyIndoorTHP, FamilyOutdoorTH:
		return 0x41
	case FamilyWH31, FamilyWH31Alt:
		return 0x46
	case FamilyPM41, FamilyPM42, FamilyPM43, FamilyPM44, FamilyWH45, FamilyWH46:
		return 0x262
	case FamilyWH51, FamilyWH57, FamilyWH48:
		return 0x50
	case FamilyWH55:
		return 0x42
	case FamilyWN34:
		return 0x4e
	case FamilyWN35:
		return 0x51
	case FamilyWH54:
		return 0x55
	}
	return 0
}

// Variant is the product-variant byte at RAM 0x40038f (initialised to 2 from
// .data). Its meaning is not confirmed; it changes which families rfTask
// forwards (see Accept), PM2.5 slot allocation and the frequency-trim path.
type Variant byte

// DefaultVariant is the .data initial value of 0x40038f.
const DefaultVariant Variant = 2

// Accept reproduces the rfTask (fn_0437ac) filter that runs before a payload
// is pushed into the ring buffer.
func Accept(v Variant, p []byte) bool {
	if len(p) == 0 {
		return false
	}
	f := p[0]
	bit := func(i int, m byte) bool { return len(p) > i && p[i]&m != 0 }
	if v == 0 {
		switch {
		case f == 0x24, f&0xF0 == 0xE0, f == 0x30, f == 0x42, f == 0x44:
			return true
		case f == 0x51, f == 0x40, f == 0x90:
			return bit(1, 0x80)
		case f&0xF0 == 0xD0, f == 0x80, f == 0x48:
			return true
		case f == 0x55:
			return bit(7, 0x01)
		case f == 0x57:
			return bit(5, 0x80)
		case f == 0x45, f == 0x35, f == 0x46:
			return bit(4, 0x80)
		}
		return false
	}
	switch {
	case f == 0x24, f&0xF0 == 0xE0, f&0xF0 == 0xD0,
		f == 0x30, f == 0x37, f >= 0x41 && f <= 0x44, f == 0x51,
		f == 0x80, f == 0x40, f == 0x57, f == 0x45, f == 0x55, f == 0x35, f == 0x46:
		return true
	case f == 0x68, f == 0x34, f == 0x90, f == 0x85, f == 0x54:
		return v != 1
	}
	return false
}
