// Package fineoffset decodes Fine Offset / Ambient Weather / Ecowitt sensor
// payloads the way the HP2550 display console does.
//
// The code is the sensor-decoding half of the hp2550 firmware reconstruction
// (package rf there: the per-family decoders that run in docTask fn_00e1c4,
// fn_00da74 and fn_00d4ec, plus the CRC-8 / sum-8 checks fn_012de4 and
// fn_012e34). The radio driver, packet ring and sensor-slot bookkeeping were
// left behind: package demod replaces the radio chip.
//
// Values are kept in the firmware's native units and sentinels:
//
//	temperature  int16, 0.1 °C, Invalid16 (0x7fff) when absent
//	humidity     uint8, %, InvalidHum (0x7f) when absent
//	pressure     0.1 hPa
//	wind         raw units per family (see the type docs and Scale* helpers)
//	light        0.1 lux, InvalidLight (0x7fffffff) when absent
package fineoffset
