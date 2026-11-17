// Package demod turns raw RTL-SDR IQ samples into Fine Offset radio payloads.
//
// It does in software what the CMT2300A does in the HP2550 console: 2-FSK
// demodulation, preamble / sync-word detection and packet framing. The
// framing comes from the console's radio register table (hp2550 rf/tables.go,
// regs 0x38..0x46): preamble byte 0xAA, sync word 0x2D 0xD4, fixed 39-byte
// (0x27) payload read from the FIFO. Checksums and field layouts are left to
// package fineoffset, exactly as the console leaves them to docTask.
package demod

import "math"

const (
	// DefaultBitRate is the Fine Offset symbol rate (58 µs per bit).
	DefaultBitRate = 1e6 / 58.0
	// DefaultPayloadLen is the CMT2300A FIFO read length (fn_043c34).
	DefaultPayloadLen = 0x27

	// syncPattern is 16 bits of preamble followed by the sync word.
	syncPattern  uint32 = 0xAAAA2DD4
	syncBits            = 32
	preambleBits        = 16
)

// Config describes the sample stream.
type Config struct {
	SampleRate float64 // samples per second (complex)
	BitRate    float64 // 0 = DefaultBitRate
	PayloadLen int     // bytes after the sync word; 0 = DefaultPayloadLen
}

// Packet is one framed payload. The payload has not been checksummed.
type Packet struct {
	Payload []byte
	// Sample is the index (from the start of the stream) of the last sync bit.
	Sample int64
	// RSSI, Noise and SNR are in dB relative to full scale.
	RSSI, Noise, SNR float64
	// FreqHi and FreqLo are the mark and space frequencies in Hz relative to
	// the tuned frequency.
	FreqHi, FreqLo float64
}

// Demod is a streaming demodulator. It is not safe for concurrent use.
type Demod struct {
	fs   float64
	sps  float64 // samples per bit
	nbit int     // payload bits

	// channel low-pass filter
	taps   []float32
	hi, hq []float32
	hpos   int

	prevI, prevQ float32

	// post-discriminator moving average
	ma    []float32
	maPos int
	maSum float32

	// filtered discriminator output (rad/sample) and power; freq[0] is
	// stream sample `base`.
	freq, pow []float32
	base      int64
	scan      int // next index in freq to test for the sync pattern

	syncOff [syncBits]int // sample offsets of each sync bit before the last one

	noise  float64 // noise floor estimate (linear power)
	blkSum float64
	blkN   int

	carry     byte // odd byte left over from the previous Write
	haveCarry bool

	lookback  int // samples needed before a sync candidate
	lookahead int // samples needed after it to slice the payload

	// sync candidate being refined
	candidate bool
	candBest  int
	candScore float32
	candThr   float32
	candEnd   int
}

// New returns a demodulator for the given stream.
func New(cfg Config) *Demod {
	if cfg.BitRate == 0 {
		cfg.BitRate = DefaultBitRate
	}
	if cfg.PayloadLen == 0 {
		cfg.PayloadLen = DefaultPayloadLen
	}
	d := &Demod{fs: cfg.SampleRate, sps: cfg.SampleRate / cfg.BitRate, nbit: cfg.PayloadLen * 8}

	// Channel filter: windowed-sinc low-pass. The signal is +-33 kHz
	// deviation at 17.2 kbit/s, so about +-50 kHz wide; allow +-15 kHz of
	// transmitter / dongle frequency error on top.
	cutoff := 65e3 / cfg.SampleRate
	if cutoff > 0.45 {
		cutoff = 0.45
	}
	n := 2*int(math.Ceil(2/cutoff)) + 1
	d.taps = lowpass(n, cutoff)
	d.hi = make([]float32, n)
	d.hq = make([]float32, n)

	// Half a bit of averaging after the discriminator.
	d.ma = make([]float32, max(1, int(d.sps/2)))

	for k := range syncBits {
		d.syncOff[k] = int(math.Round(float64(syncBits-1-k) * d.sps))
	}
	d.lookback = d.syncOff[0] + 2
	d.lookahead = int(float64(d.nbit+3)*d.sps) + 4
	d.noise = -1
	return d
}

func lowpass(n int, cutoff float64) []float32 {
	t := make([]float32, n)
	var sum float64
	m := float64(n-1) / 2
	for i := range t {
		x := float64(i) - m
		v := 2 * cutoff
		if x != 0 {
			v = math.Sin(2*math.Pi*cutoff*x) / (math.Pi * x)
		}
		v *= 0.54 - 0.46*math.Cos(2*math.Pi*float64(i)/float64(n-1))
		t[i] = float32(v)
		sum += v
	}
	for i := range t {
		t[i] = float32(float64(t[i]) / sum)
	}
	return t
}

// Write feeds unsigned 8-bit interleaved IQ (the rtl_sdr / .cu8 format) and
// returns the packets completed by this chunk.
func (d *Demod) Write(iq []byte) []Packet {
	if d.haveCarry && len(iq) > 0 {
		d.sample(d.carry, iq[0])
		iq = iq[1:]
		d.haveCarry = false
	}
	for i := 0; i+1 < len(iq); i += 2 {
		d.sample(iq[i], iq[i+1])
	}
	if len(iq)&1 == 1 {
		d.carry, d.haveCarry = iq[len(iq)-1], true
	}
	return d.search()
}

// sample runs one IQ pair through the channel filter, FM discriminator and
// post filter and appends the result to the work buffers.
func (d *Demod) sample(bi, bq byte) {
	n := len(d.taps)
	d.hi[d.hpos] = float32(bi) - 127.5
	d.hq[d.hpos] = float32(bq) - 127.5
	d.hpos++
	if d.hpos == n {
		d.hpos = 0
	}
	var i, q float32
	j := d.hpos
	for _, t := range d.taps {
		i += t * d.hi[j]
		q += t * d.hq[j]
		j++
		if j == n {
			j = 0
		}
	}

	// Polar discriminator: angle of x[n] * conj(x[n-1]).
	re := i*d.prevI + q*d.prevQ
	im := q*d.prevI - i*d.prevQ
	d.prevI, d.prevQ = i, q
	f := float32(math.Atan2(float64(im), float64(re)))

	d.maSum += f - d.ma[d.maPos]
	d.ma[d.maPos] = f
	d.maPos++
	if d.maPos == len(d.ma) {
		d.maPos = 0
		// Re-sum periodically so float error cannot accumulate.
		var s float32
		for _, v := range d.ma {
			s += v
		}
		d.maSum = s
	}
	p := i*i + q*q
	d.freq = append(d.freq, d.maSum/float32(len(d.ma)))
	d.pow = append(d.pow, p)

	// Noise floor: follows quiet blocks quickly and loud ones very slowly.
	d.blkSum += float64(p)
	d.blkN++
	if d.blkN == 2048 {
		b := d.blkSum / float64(d.blkN)
		switch {
		case d.noise < 0:
			d.noise = b
		case b < d.noise:
			d.noise += (b - d.noise) * 0.25
		default:
			d.noise += (b - d.noise) * 0.002
		}
		d.blkSum, d.blkN = 0, 0
	}
}

// matchSync tests whether the 32 bits ending with a bit centred on sample n
// are the preamble + sync word. The slicing threshold is the mean of the
// preamble bits, which is the carrier centre whatever the frequency offset.
// score is the smallest distance of any bit from the threshold.
func (d *Demod) matchSync(n int) (thr, score float32, ok bool) {
	f := d.freq
	for k := range preambleBits {
		thr += f[n-d.syncOff[k]]
	}
	thr /= preambleBits
	score = math.MaxFloat32
	for k := range syncBits {
		v := f[n-d.syncOff[k]] - thr
		if (v > 0) != (syncPattern>>(syncBits-1-k)&1 == 1) {
			return 0, 0, false
		}
		if v < 0 {
			v = -v
		}
		if v < score {
			score = v
		}
	}
	return thr, score, true
}

func (d *Demod) search() []Packet {
	var out []Packet
	for d.scan+d.lookahead < len(d.freq) {
		n := d.scan
		if n < d.lookback {
			d.scan++
			continue
		}
		thr, score, ok := d.matchSync(n)
		if d.candidate {
			// The pattern matches over a few adjacent samples; keep the
			// one with the widest eye as the bit-centre estimate.
			if ok && score > d.candScore {
				d.candBest, d.candScore, d.candThr = n, score, thr
			}
			if n >= d.candEnd {
				out = append(out, d.decode(d.candBest, d.candThr))
				d.candidate = false
				d.scan = d.candBest + int(2*d.sps)
				continue
			}
		} else if ok {
			d.candidate = true
			d.candBest, d.candScore, d.candThr = n, score, thr
			d.candEnd = n + int(d.sps)
		}
		d.scan++
	}
	// Drop history that can no longer be needed.
	if keep := d.scan - d.lookback - int(d.sps); keep > 1<<16 && !d.candidate {
		d.freq = append(d.freq[:0], d.freq[keep:]...)
		d.pow = append(d.pow[:0], d.pow[keep:]...)
		d.base += int64(keep)
		d.scan -= keep
	}
	return out
}

func (d *Demod) at(t float64) float32 {
	i := int(t)
	a := float32(t - float64(i))
	return d.freq[i]*(1-a) + d.freq[i+1]*a
}

// decode slices the payload that follows a sync word whose last bit is
// centred on sample n. Bit timing follows the transitions in the data, so a
// transmitter whose clock is a little off stays in step for the whole packet.
func (d *Demod) decode(n int, thr float32) Packet {
	p := Packet{Payload: make([]byte, d.nbit/8), Sample: d.base + int64(n)}

	t := float64(n)
	prev := false // the sync word ends in a 0 bit
	for i := range d.nbit {
		t += d.sps
		bit := d.at(t) > thr
		if bit {
			p.Payload[i/8] |= 0x80 >> (i % 8)
		}
		if bit != prev {
			// Locate the threshold crossing between the two bit centres.
			lo, hi := int(t-d.sps), int(t)+1
			for j := lo; j < hi; j++ {
				a, b := d.freq[j]-thr, d.freq[j+1]-thr
				if (a > 0) != (b > 0) {
					cross := float64(j) + float64(a/(a-b))
					err := cross - (t - d.sps/2)
					err = max(min(err, d.sps/4), -d.sps/4)
					t += 0.3 * err
					break
				}
			}
		}
		prev = bit
	}

	// Mark / space frequency from the sync bits.
	var hi, lo float64
	var nh, nl int
	for k := range syncBits {
		v := float64(d.freq[n-d.syncOff[k]])
		if syncPattern>>(syncBits-1-k)&1 == 1 {
			hi += v
			nh++
		} else {
			lo += v
			nl++
		}
	}
	toHz := d.fs / (2 * math.Pi)
	p.FreqHi, p.FreqLo = hi/float64(nh)*toHz, lo/float64(nl)*toHz

	// Signal level over the sync word and payload.
	var sum float64
	a, b := n-d.syncOff[0], min(int(t), len(d.pow)-1)
	for _, v := range d.pow[a:b] {
		sum += float64(v)
	}
	const fullScale = 127.5 * 127.5
	db := func(x float64) float64 { return 10 * math.Log10(max(x, 1e-3)/fullScale) }
	p.RSSI = db(sum / float64(b-a))
	p.Noise = db(d.noise)
	p.SNR = p.RSSI - p.Noise
	return p
}
