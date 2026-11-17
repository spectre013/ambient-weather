package fineoffset

// Firmware sentinels.
const (
	Invalid16      int16  = 0x7fff     // temperature / direction absent
	InvalidU16     uint16 = 0x7fff     // wind, gust, pressure, distance absent
	InvalidFFFF    uint16 = 0xffff     // PM / CO2 / aux fields absent
	InvalidHum     uint8  = 0x7f       // humidity / moisture absent
	InvalidByte    uint8  = 0xff       // UV / distance (km) absent
	InvalidLight   uint32 = 0x7fffffff // light absent
	InvalidCounter uint32 = 0xffffffff // WH54 24-bit field absent
)

// Measurement is one decoded sensor packet.
type Measurement interface {
	Family() Family
	// SensorID is the value the firmware registers in the slot (+0x08).
	SensorID() uint32
	// Channel is the channel nibble carried in the packet, or -1.
	Channel() int
	// LowBattery applies the per-family threshold the firmware stores in
	// slot byte +0x0d.
	LowBattery() bool
}

// Header is shared by all decoded types.
type Header struct {
	Raw0 byte   // first payload byte
	ID   uint32 // sensor ID as registered
}

func (h Header) Family() Family   { return FamilyOf(h.Raw0) }
func (h Header) SensorID() uint32 { return h.ID }
func (h Header) Channel() int     { return -1 }

// THP is family 0xEx (indoor T/H/baro, writes 0x1f10128+0x28/+0x02/+0x2a)
// and 0xDx (outdoor T/H, writes 0x1f10128+0x2e/+0x03). ID is 8 bits:
// (b1&0xf0)|(b0&0x0f).
type THP struct {
	Header
	Temp     int16  // 0.1 °C
	Hum      uint8  // %
	Pressure uint16 // 0.1 hPa absolute; InvalidU16 for 0xDx
	BattLow  bool   // b1 bit3
}

func (m THP) LowBattery() bool { return m.BattLow }

// WH24 is family 0x24 (fn_00da74). Wind/gust/UV are raw; use ScaleWH24*.
type WH24 struct {
	Header
	Extended bool   // 25-byte form (24-bit ID, barometer)
	WindDir  int16  // degrees, Invalid16
	Temp     int16  // 0.1 °C
	Hum      uint8  // %
	WindRaw  uint16 // 11-bit raw speed (see ScaleWH24Wind)
	GustRaw  uint16 // raw gust, InvalidU16
	Rain     uint16 // rain tip counter (b8:b9)
	UVRaw    uint16 // raw UV (µW/cm² class), InvalidU16
	Light    uint32 // 0.1 lux, InvalidLight
	BattLow  bool   // b3 bit3
	Pressure uint16 // 0.1 hPa, extended form only, else InvalidU16
}

func (m WH24) LowBattery() bool { return m.BattLow }

// Anemometer is family 0x68 (WS68) or 0x80 (WS80). Wind is 0.1 m/s before
// calibration (ScaleWind), UV is 0.1 UVI (ScaleUVI).
type Anemometer struct {
	Header
	Nibble  uint8  // b1>>4 (meaning unknown)
	Light   uint32 // 0.1 lux (raw*100), InvalidLight
	BattRaw uint8  // b6, 0.02 V steps (guess)
	Temp    int16
	Hum     uint8
	Wind    uint16 // InvalidU16, capped 500
	Dir     int16  // Invalid16
	Gust    uint16 // InvalidU16, capped 500
	UV      uint8  // InvalidByte if > 160
	// WS80 only
	Aux       uint16 // b14&0x3f:b15, validated 3000..12000, InvalidFFFF
	AuxHi     uint8  // b14>>6 (-> 0x1f123de)
	BattFlags uint8  // b7 bit3 | b7 bit7<<1 (-> 0x1f114a8)
}

func (m Anemometer) LowBattery() bool {
	if m.Family() == FamilyWS80 {
		return m.BattRaw < 0x79
	}
	return m.BattRaw < 0x3d
}

// WS90 is family 0x90 (fn_00d4ec).
type WS90 struct {
	Anemometer
	Length      int    // 32 or 39 (stored at 0x1f11eb6)
	F20         uint8  // b16>>4
	RainCounter uint32 // (b16&0xf):b17:b18, 20-bit
	Rain2       uint16 // b19:b20
	F22         uint8  // b21>>6
	F23         uint8  // b22
	F24         uint8  // b23 (third input to fn_00bd6c)
	F1A         uint16 // b24<<4 | b25>>4
	F1C         uint16 // (b25&0xf)<<8 | b26
	F1E, F1F    uint8  // b27, b28
	Version     uint8  // b29; must repeat before the packet is accepted
}

func (m WS90) LowBattery() bool { return m.BattRaw < 0x79 }

// WS85 is family 0x85.
type WS85 struct {
	Header
	BattRaw     uint8 // b4
	Flags       uint8 // b5
	Wind        uint16
	Dir         int16
	Gust        uint16
	RainCounter uint32 // (b12&0xf):b13:b14
	Rain2       uint16 // b15:b16
}

func (m WS85) LowBattery() bool { return m.BattRaw < 0x79 }

// WH40 is family 0x40 (rain gauge).
type WH40 struct {
	Header
	BattRaw      uint8  // b4&0x7f, 0.1 V (guess)
	BattReported bool   // b4 bit7; when clear the low flag is forced off
	Rain         uint16 // b5:b6 counter
}

func (m WH40) LowBattery() bool { return m.BattReported && m.BattRaw < 0x0d }

// WH31 is family 0x30/0x37.
type WH31 struct {
	Header
	Ch      uint8 // 0..7
	Temp    int16
	HasHum  bool
	Hum     uint8
	BattLow bool
}

func (m WH31) Channel() int     { return int(m.Ch) }
func (m WH31) LowBattery() bool { return m.BattLow }

// PM25 is family 0x41..0x44.
type PM25 struct {
	Header
	Long      bool   // 0x43/0x44: 20-bit ID
	BattLevel uint8  // 4-bit level
	PM25      uint16 // 0.1 µg/m³ rounded to whole µg (x10), InvalidFFFF
	PM10      uint16 // second value (PM10 by analogy with rtl_433; unused by fw), InvalidFFFF
}

func (m PM25) LowBattery() bool { return m.BattLevel < 2 }

// WH51 is family 0x51 (soil moisture).
type WH51 struct {
	Header
	Nibble   uint8  // b1>>4
	Moisture uint8  // %, InvalidHum if > 100
	BattRaw  uint8  // b4&0x1f, 0.1 V (guess)
	AD       uint16 // (b7&7):b8 raw ADC for user calibration
}

func (m WH51) LowBattery() bool { return m.BattRaw < 0x0d }

// WH55 is family 0x55 (leak).
type WH55 struct {
	Header
	Ch        uint8 // b1>>4, direct slot index
	BattLevel uint8 // b4
	Value     uint16
	Leak      bool // b7 bit7
	Flag6     bool // b7 bit6
}

func (m WH55) Channel() int     { return int(m.Ch) }
func (m WH55) LowBattery() bool { return m.BattLevel < 2 }

// WH57 is family 0x57 (lightning).
type WH57 struct {
	Header
	State      uint8 // 3 if b1 bit7, 2 if bit6, else b1 bit4
	BattLevel  uint8 // b4
	DistanceKm uint8 // b5&0x3f, InvalidByte if > 40
	Count      uint8 // b6, rolling strike counter
}

func (m WH57) LowBattery() bool { return m.BattLevel < 2 }

// WH45 is family 0x45 (15 bytes) / 0x46 (21 bytes).
type WH45 struct {
	Header
	Temp      int16 // firmware sentinel is 0x7ff; mapped to Invalid16
	Hum       uint8 // 0xff = absent
	BattLevel uint8 // (b9>>6)|((b7&0xc0)>>4); 6 sets slot+0x1c
	PM25      uint16
	PM10      uint16
	CO2       uint16 // ppm
	Ext10     uint16 // 0x46: b13:b14 (PM1.0? guess), else InvalidFFFF
	Ext12     uint16 // 0x46: b15:b16 (PM4.0? guess), else InvalidFFFF
	Ext14     uint16 // 0x46: b17:b18 raw, else InvalidFFFF
}

func (m WH45) LowBattery() bool { return m.BattLevel < 2 }

// WN34 is family 0x34.
type WN34 struct {
	Header
	Temp    int16 // 0.1 °C
	Signed  bool  // b4 bit6: raw is signed 12-bit without the 400 offset
	BattRaw uint8 // b6, 0.02 V (guess)
}

func (m WN34) LowBattery() bool { return m.BattRaw < 0x3d }

// WH48 is family 0x48 (accepted only when Variant == 0).
type WH48 struct {
	Header
	Ch      uint8
	Temp    int16
	BattLow bool  // b4 bit4
	Flag3   bool  // b4 bit3
	Hum     uint8 // b6
	CO2     uint16
	Value2  uint16 // b9:b10
}

func (m WH48) Channel() int     { return int(m.Ch) }
func (m WH48) LowBattery() bool { return m.BattLow }

// WN35 is family 0x35 (leaf wetness).
type WN35 struct {
	Header
	HasTemp bool
	Temp    int16
	Wetness uint8 // %, InvalidHum
	BattRaw uint8 // b8
}

func (m WN35) LowBattery() bool { return m.BattRaw < 0x3d }

// WH54 is family 0x54 ("LDS").
type WH54 struct {
	Header
	Nibble   uint8  // b1>>4 (fw replaces with slot index)
	BattRaw  uint8  // b4
	Flag11   bool   // b5 bit7
	F13      uint8  // b5&0x7f
	Distance uint16 // (b6&0x3f):b7, mm (guess), InvalidU16
	Flag12   bool   // b6 bit6
	Type     uint8  // b6 bit7; selects battery threshold
	F6       uint16 // (b8&0x3f):b9
	FA       uint16 // b10:b11
	FC       uint32 // b12:b13:b14, InvalidCounter
	F14      uint16 // (b8&0xc0)<<2 | b15, 10-bit accumulating counter
	F17      uint8  // b17 (the CRC byte position)
}

func (m WH54) LowBattery() bool {
	if m.Type == 1 {
		return m.BattRaw < 0x3d
	}
	return m.BattRaw < 0x79
}
