package main

import "time"

// rainWindowMinutes covers the longest rolling window we need: the 24 h
// total used by the rain-event reset rule.
const rainWindowMinutes = 24 * 60

// rainWindow is a ring of one-minute rain buckets. It backs the console-style
// rolling values: hourly rain (last 60 min), rain rate (last 10 min x 6) and
// the 24 h total that decides when a rain event ends.
type rainWindow struct {
	mm     [rainWindowMinutes]float64
	minute [rainWindowMinutes]int64 // unix minute each bucket belongs to
}

func unixMinute(t time.Time) int64 { return t.Unix() / 60 }

// add records mm of rain in the minute containing t.
func (w *rainWindow) add(t time.Time, mm float64) {
	m := unixMinute(t)
	i := m % rainWindowMinutes
	if w.minute[i] != m {
		w.minute[i], w.mm[i] = m, 0
	}
	w.mm[i] += mm
}

// sum returns the rain in the n minutes ending with the minute containing
// now (the current, partial minute included).
func (w *rainWindow) sum(now time.Time, n int) float64 {
	end := unixMinute(now)
	var total float64
	for m := end - int64(n) + 1; m <= end; m++ {
		if i := m % rainWindowMinutes; w.minute[i] == m {
			total += w.mm[i]
		}
	}
	return total
}

// hourlyMM is the rolling 60-minute total.
func (w *rainWindow) hourlyMM(now time.Time) float64 { return w.sum(now, 60) }

// rateMMPerHour is the console's rain rate: the last 10 minutes x 6.
func (w *rainWindow) rateMMPerHour(now time.Time) float64 { return w.sum(now, 10) * 6 }

// last24hMM is the rolling 24-hour total.
func (w *rainWindow) last24hMM(now time.Time) float64 { return w.sum(now, rainWindowMinutes) }

// RainSample is one stored record's daily rain total, used to rebuild the
// rolling window after a restart.
type RainSample struct {
	Recorded    time.Time
	DailyRainIn float64
}

// SeedRainHistory refills the rolling window from recent records, oldest
// first. Each increase in dailyrainin is rain that fell in that row's minute;
// a decrease is the midnight reset, after which the new value is that day's
// rain so far.
func (s *SensorState) SeedRainHistory(rows []RainSample) {
	s.mu.Lock()
	defer s.mu.Unlock()
	const inToMM = 25.4
	var prev float64
	for i, r := range rows {
		d := r.DailyRainIn
		if i > 0 {
			if r.DailyRainIn >= prev {
				d = r.DailyRainIn - prev
			}
			if d > 0 {
				s.rain.window.add(r.Recorded, d*inToMM)
			}
		}
		prev = r.DailyRainIn
	}
}
