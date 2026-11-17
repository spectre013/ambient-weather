package demod

import (
	"bytes"
	"encoding/hex"
	"math"
	"math/rand/v2"
	"os"
	"testing"
)

// modulate builds cu8 IQ for one transmission: preamble, sync, payload,
// 2-FSK with continuous phase, plus noise and lead-in / lead-out silence.
func modulate(payload []byte, fs, bitRate, dev, offset, amp, noise float64, rng *rand.Rand) []byte {
	frame := append([]byte{0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0x2D, 0xD4}, payload...)
	sps := fs / bitRate
	lead := int(fs * 0.01)
	total := lead*2 + int(float64(len(frame)*8)*sps)
	out := make([]byte, 0, total*2)
	var phase float64
	q8 := func(v float64) byte { return byte(max(min(math.Round(v+127.5), 255), 0)) }
	for n := range total {
		var i, q float64
		if k := int(float64(n-lead) / sps); n >= lead && k < len(frame)*8 {
			f := offset - dev
			if frame[k/8]&(0x80>>(k%8)) != 0 {
				f = offset + dev
			}
			phase += 2 * math.Pi * f / fs
			i, q = amp*math.Cos(phase), amp*math.Sin(phase)
		}
		out = append(out, q8(i+rng.NormFloat64()*noise), q8(q+rng.NormFloat64()*noise))
	}
	return out
}

func testPayload() []byte {
	p := make([]byte, DefaultPayloadLen)
	for i := range p {
		p[i] = byte(i*37 + 11)
	}
	p[5], p[6], p[7] = 0x00, 0x00, 0xff // long runs without transitions
	return p
}

func TestDemod(t *testing.T) {
	const fs = 250000
	want := testPayload()
	tests := []struct {
		name                        string
		bitRate, offset, amp, noise float64
	}{
		{"clean", DefaultBitRate, 0, 100, 1},
		{"offset +15 kHz", DefaultBitRate, 15e3, 100, 1},
		{"offset -15 kHz", DefaultBitRate, -15e3, 100, 1},
		{"clock +1%", DefaultBitRate * 1.01, 0, 100, 1},
		{"clock -1%", DefaultBitRate * 0.99, 3e3, 100, 1},
		{"weak, 12 dB SNR", DefaultBitRate, -5e3, 12, 3},
		{"clipping", DefaultBitRate, 2e3, 400, 1},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rng := rand.New(rand.NewPCG(1, 2))
			iq := modulate(want, fs, tc.bitRate, 33e3, tc.offset, tc.amp, tc.noise, rng)
			// Pad so the packet is flushed out of the look-ahead window.
			for range 20000 {
				iq = append(iq, byte(127.5+rng.NormFloat64()*tc.noise))
			}
			d := New(Config{SampleRate: fs})
			var got []Packet
			// Feed in odd-sized chunks to exercise the stream handling.
			for len(iq) > 0 {
				n := min(4097, len(iq))
				got = append(got, d.Write(iq[:n])...)
				iq = iq[n:]
			}
			if len(got) != 1 {
				t.Fatalf("got %d packets, want 1", len(got))
			}
			if !bytes.Equal(got[0].Payload, want) {
				t.Fatalf("payload\n got % x\nwant % x", got[0].Payload, want)
			}
			if mid := (got[0].FreqHi + got[0].FreqLo) / 2; math.Abs(mid-tc.offset) > 2e3 {
				t.Errorf("centre %.0f Hz, want %.0f", mid, tc.offset)
			}
		})
	}
}

func TestNoiseOnly(t *testing.T) {
	rng := rand.New(rand.NewPCG(3, 4))
	iq := make([]byte, 2_000_000)
	for i := range iq {
		iq[i] = byte(127.5 + rng.NormFloat64()*4)
	}
	if got := New(Config{SampleRate: 250000}).Write(iq); len(got) != 0 {
		t.Fatalf("%d packets from noise", len(got))
	}
}

// Real bursts cut from a 915 MHz capture (NESDR SMArt, 250 kS/s, automatic
// gain, so the samples clip). Each expected payload was confirmed against
// rtl_433 decoding the same capture.
func TestRealCaptures(t *testing.T) {
	tests := []struct {
		file, payload string
		count         int // sensors send each packet twice
	}{
		{"ws80.cu8", "8000860200819ea1e641084e09003fff830f", 1},
		{"wh31e.cu8", "309c925d244726", 2},
		{"wh40.cu8", "4080cead8e06469db2", 2},
	}
	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			iq, err := os.ReadFile("testdata/" + tc.file)
			if err != nil {
				t.Fatal(err)
			}
			want, _ := hex.DecodeString(tc.payload)
			got := New(Config{SampleRate: 250000}).Write(iq)
			if len(got) != tc.count {
				t.Fatalf("got %d packets, want %d", len(got), tc.count)
			}
			for _, p := range got {
				if !bytes.HasPrefix(p.Payload, want) {
					t.Errorf("payload % x, want prefix % x", p.Payload, want)
				}
				if dev := (p.FreqHi - p.FreqLo) / 2; dev < 25e3 || dev > 45e3 {
					t.Errorf("deviation %.0f Hz", dev)
				}
			}
		})
	}
}
