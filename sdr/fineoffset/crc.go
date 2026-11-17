package fineoffset

// crc8Table is the table at 0xa1450 used by fn_012de4: CRC-8, poly 0x31,
// init 0x00, MSB first, no reflection, no final xor.
var crc8Table = func() (t [256]byte) {
	for i := range t {
		c := byte(i)
		for range 8 {
			if c&0x80 != 0 {
				c = c<<1 ^ 0x31
			} else {
				c <<= 1
			}
		}
		t[i] = c
	}
	return t
}()

// CRC8 computes the Fine Offset CRC-8 (poly 0x31, init 0).
func CRC8(b []byte) byte {
	var c byte
	for _, x := range b {
		c = crc8Table[x^c]
	}
	return c
}

// Sum8 is the 8-bit additive checksum (fn_012e34).
func Sum8(b []byte) byte {
	var s byte
	for _, x := range b {
		s += x
	}
	return s
}

// crcOK mirrors fn_012de4(buf, n): CRC8(buf[0:n-1]) == buf[n-1].
func crcOK(b []byte, n int) bool { return len(b) >= n && CRC8(b[:n-1]) == b[n-1] }

// sumOK mirrors fn_012e34(buf, n): Sum8(buf[0:n-1]) == buf[n-1].
func sumOK(b []byte, n int) bool { return len(b) >= n && Sum8(b[:n-1]) == b[n-1] }

// Seal fills in the CRC byte at crcLen-1 (if crcLen > 0) and then the sum
// byte at sumLen-1, producing a packet that passes the firmware checks.
// It is used to build test vectors and simulator packets.
func Seal(b []byte, crcLen, sumLen int) []byte {
	if crcLen > 0 {
		b[crcLen-1] = CRC8(b[:crcLen-1])
	}
	b[sumLen-1] = Sum8(b[:sumLen-1])
	return b
}
