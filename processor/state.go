package main

import (
	"fmt"
	"log"
	"math"
	"sync"
	"time"
)

// Counter limits, taken from the HP2550 console firmware (fn_00bb40 for rain,
// the WH57 handler in docTask for lightning).
const (
	// rainCounterMod: the WH40 reports a 16-bit counter of 0.1 mm tips, so
	// rain_mm wraps at 6553.6 mm.
	rainCounterMod = 1 << 16
	// rainMaxPerSec / rainMaxDelta: the console accepts at most 4 counts
	// (0.4 mm) per elapsed second, capped at 1023 counts per packet. Larger
	// jumps are treated as a counter reset (battery swap) or a bad packet.
	rainMaxPerSec = 4
	rainMaxDelta  = 1023
	// strikeCounterMod: strike_count is a single byte and wraps at 256.
	strikeCounterMod = 1 << 8
	// strikeMaxDelta: the console only accepts increments below 100.
	strikeMaxDelta = 100
	// pmNoReading is the WH0290's raw 0x3fff "no reading" value after
	// rtl_433's integer division by 10.
	pmNoReading = 1638
	// rainEventResetMM: the console ends a rain event only once the last
	// 24 hours hold less than 1.0 mm (and the last hour none).
	rainEventResetMM = 1.0
	// rainEpsilonMM absorbs float rounding when testing a window for "no rain".
	rainEpsilonMM = 0.01
	// WH32B glitch filter limits (console fn_00c214): a packet is accepted only
	// if it is this close to the previous packet.
	wh32BMaxDeltaC   = 0.5
	wh32BMaxDeltaRH  = 3
	wh32BMaxDeltaHPa = 1.0
)

// pm25History keeps a sliding window of PM2.5 readings so we can compute a 24h AQI.
type pm25Sample struct {
	t  time.Time
	pm float64
}

// rainTotals holds the rain accumulated in each period. Like the console,
// every WH40 message adds the counter increment since the previous message
// rather than subtracting a period baseline, so a counter wrap or a
// battery-swap reset never loses the totals already collected.
type rainTotals struct {
	dayMM   float64
	weekMM  float64
	monthMM float64
	yearMM  float64

	hourStart  time.Time
	dayStart   time.Time
	weekStart  time.Time
	monthStart time.Time
	yearStart  time.Time

	// eventMM follows the console's rule: it keeps accumulating until an hour
	// boundary at which the last 60 minutes were dry and the last 24 hours
	// had less than rainEventResetMM (see rollRain).
	eventMM      float64
	lastRainTime time.Time

	// window holds per-minute rain for the rolling hourly / rate / 24 h values.
	window rainWindow

	// lastCount is the previous raw counter in 0.1 mm units and lastMsg the
	// time it arrived (for the per-second plausibility limit).
	lastCount int
	lastMsg   time.Time
}

// lightningTotals mirrors rainTotals for strike counts.
type lightningTotals struct {
	hourCount  int
	dayCount   int
	monthCount int

	hourStart  time.Time
	dayStart   time.Time
	monthStart time.Time

	lastStrikeTime time.Time
	lastDistanceKM float64
}

// SensorState is the in-memory "current world" — every minute we snapshot it
// into a Record and push to postgres.
type SensorState struct {
	mu sync.Mutex

	// Outdoor WS80
	outdoorTempC    float64
	outdoorHumidity int
	windDirDeg      int
	windAvgMS       float64
	windMaxMS       float64
	windGustDir     int
	maxDailyGustMS  float64
	maxGustDay      time.Time
	uvi             float64
	lightLux        float64
	ws80BatteryMV   float64
	ws80HaveData    bool

	// Indoor WH32B
	indoorTempC    float64
	indoorHumidity int
	pressureHPa    float64
	wh32BBattery   int
	wh32BHaveData  bool
	// wh32BPrev is the previous WH32B packet, for the console's glitch
	// filter; it is updated on every packet, accepted or not.
	wh32BPrev     wh32BReading
	wh32BPrevHave bool

	// WH31E indoor channels (1=basement, 2=hannah, 3=master, 4=garage)
	channelTempC    [11]float64 // index 1..10 used
	channelHumidity [11]int
	channelBattery  [11]int
	channelHave     [11]bool

	// Air quality
	pm25         float64
	pm10         float64
	pm25History  []pm25Sample
	aqHaveData   bool
	aqBatteryRaw float64

	// Rain (from EcoWitt-WH40)
	totalRainMM     float64
	totalRainHaveIt bool
	rain            rainTotals

	// Lightning (from FineOffset-WH31L)
	totalStrikes     int
	totalStrikesHave bool
	lightningBattery float64
	lightning        lightningTotals

	// Seeded values from the most recent DB record. These describe how much
	// rain/lightning has *already* accumulated in the current period before
	// this process started. They become the starting period totals on the
	// first WH40/WH31L message after startup so that period totals stay
	// continuous across restarts.
	seededRain     bool
	seedMaxGustMph float64
	seedMaxGustDay time.Time

	seededLightning bool

	seededAQ       bool
	seedAQIPM25    int
	seedAQIPM2524h int

	// dirty is set by Apply whenever a sensor message updates the state and
	// cleared by Snapshot. It lets the caller skip an insert when nothing new
	// has arrived since the last snapshot (e.g. the MQTT feed has stalled).
	dirty      bool
	lastUpdate time.Time

	// elevationM is the station's elevation above sea level in meters, used to
	// convert the WH32B's absolute pressure to the sea-level-corrected
	// relative pressure (Baromrelin).
	elevationM float64
}

// Seed pre-populates state from the most recent DB record so that the very
// first snapshot after startup has plausible values for every column, instead
// of zeros for sensors that haven't reported yet. Two things happen:
//
//  1. Live sensor values (temp/humidity/wind/baro/channel temps/AQI) are
//     copied into the in-memory state and the corresponding *HaveData flag is
//     set. As real MQTT messages arrive they overwrite the seeded values.
//
//  2. Rain and lightning period totals are loaded into the running
//     accumulators, so increments from the first WH40/WH31L reading add to
//     what already fell this period, and the periods roll over at their
//     boundaries even if the sensor never reports again. Period totals that
//     fall in a previous period (e.g. last record was yesterday) are
//     discarded so the new period correctly starts at zero.
func (s *SensorState) Seed(last *Record, now time.Time) {
	if last == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	last.Recorded = last.Recorded.In(now.Location())

	// --- Pre-populate live sensor state from the last DB row so the very first
	// snapshot has plausible values for every column. As real MQTT messages
	// arrive they overwrite these seeded values. A field is only seeded if the
	// stored value looks "real" (non-zero), so we don't propagate genuine
	// zero-defaults from a fresh DB row.

	// Outdoor (WS80): tempf/humidity/wind are useless if the unit hasn't ever
	// reported, so gate on tempf being non-zero.
	if last.Tempf != 0 || last.Humidity != 0 {
		s.outdoorTempC = (last.Tempf - 32) * 5.0 / 9.0
		s.outdoorHumidity = last.Humidity
		s.windDirDeg = last.Winddir
		s.windAvgMS = last.Windspeedmph / 2.236936
		s.windMaxMS = last.Windgustmph / 2.236936
		s.windGustDir = last.Windgustdir
		s.uvi = last.Uv
		s.lightLux = last.Solarradiation * 126.7
		if last.Battout > 0 {
			s.ws80BatteryMV = 3000 // good-enough placeholder; live msg will overwrite
		}
		s.ws80HaveData = true
	}

	// Indoor barometer (WH32B).
	if last.Tempinf != 0 || last.Baromabsin != 0 {
		s.indoorTempC = (last.Tempinf - 32) * 5.0 / 9.0
		s.indoorHumidity = last.Humidityin
		s.pressureHPa = last.Baromabsin / 0.02953
		s.wh32BBattery = 1
		s.wh32BHaveData = true
	}

	// WH31E channels — temp/humidity per channel.
	type chSeed struct {
		idx   int
		tempF float64
		hum   int
		batt  int
	}
	chs := []chSeed{
		{1, last.Temp1f, last.Humidity1, last.Batt1},
		{2, last.Temp2f, last.Humidity2, last.Batt2},
		{3, last.Temp3f, last.Humidity3, last.Batt3},
		{4, last.Temp4f, last.Humidity4, last.Batt4},
		{5, last.Temp5f, last.Humidity5, last.Batt5},
		{6, last.Temp6f, last.Humidity6, last.Batt6},
		{7, last.Temp7f, last.Humidity7, last.Batt7},
		{8, last.Temp8f, last.Humidity8, last.Batt8},
		{9, last.Temp9f, last.Humidity9, last.Batt9},
		{10, last.Temp10f, last.Humidity10, last.Batt10},
	}
	for _, c := range chs {
		if c.tempF == 0 && c.hum == 0 {
			continue
		}
		s.channelTempC[c.idx] = (c.tempF - 32) * 5.0 / 9.0
		s.channelHumidity[c.idx] = c.hum
		s.channelBattery[c.idx] = c.batt
		s.channelHave[c.idx] = true
	}

	// Air quality (WH0290). We don't store raw pm2_5 in the DB, only AQI, so
	// just mark it as having data and let aqi/aqi24 carry forward via the
	// fallback path until the first real WH0290 message arrives.
	if last.Aqipm25 != 0 {
		s.seedAQIPM25 = last.Aqipm25
		s.seedAQIPM2524h = last.Aqipm2524h
		s.seededAQ = true
	}

	// --- Rain period seeds ----------------------------------------------------
	// Only carry a period's accumulation forward if the previous record was
	// inside the *same* period. Otherwise that period has rolled over while we
	// were down and the seed for it should be 0. The period starts are reset
	// to now's periods so rollRain zeroes each total at its next boundary.
	const inToMM = 25.4 // mmToIn is mm/25.4
	s.rain.hourStart = startOfHour(now)
	s.rain.dayStart = startOfDay(now)
	s.rain.weekStart = startOfWeek(now)
	s.rain.monthStart = startOfMonth(now)
	s.rain.yearStart = startOfYear(now)
	if startOfDay(last.Recorded).Equal(startOfDay(now)) {
		s.rain.dayMM = last.Dailyrainin * inToMM
		s.seedMaxGustMph = last.Maxdailygust
		s.seedMaxGustDay = startOfDay(now)
	}
	if startOfWeek(last.Recorded).Equal(startOfWeek(now)) {
		s.rain.weekMM = last.Weeklyrainin * inToMM
	}
	if startOfMonth(last.Recorded).Equal(startOfMonth(now)) {
		s.rain.monthMM = last.Monthlyrainin * inToMM
	}
	if startOfYear(last.Recorded).Equal(startOfYear(now)) {
		s.rain.yearMM = last.Yearlyrainin * inToMM
	}
	// The event total always carries forward; the next hourly check in
	// rollRain clears it if the last 24 hours were (nearly) dry.
	s.rain.eventMM = last.Eventrainin * inToMM
	// lastrain timestamp always carries forward — it's a historical fact, not
	// a per-period accumulation.
	if !last.Lastrain.IsZero() {
		s.rain.lastRainTime = last.Lastrain
	}
	s.seededRain = true

	// Restore today's running max gust.
	if !s.seedMaxGustDay.IsZero() {
		s.maxDailyGustMS = s.seedMaxGustMph / 2.236936
		s.maxGustDay = s.seedMaxGustDay
	}

	// Lightning seeds, loaded the same way as rain.
	s.lightning.hourStart = startOfHour(now)
	s.lightning.dayStart = startOfDay(now)
	s.lightning.monthStart = startOfMonth(now)
	if startOfHour(last.Recorded).Equal(startOfHour(now)) {
		s.lightning.hourCount = last.Lightninghour
	}
	if startOfDay(last.Recorded).Equal(startOfDay(now)) {
		s.lightning.dayCount = last.Lightningday
	}
	if startOfMonth(last.Recorded).Equal(startOfMonth(now)) {
		s.lightning.monthCount = last.LightningMonth
	}
	// Same as lastrain: the lightning timestamp/distance are historical facts.
	// DB column is miles; in-memory state is km.
	if !last.Lightningtime.IsZero() {
		s.lightning.lastStrikeTime = last.Lightningtime
		s.lightning.lastDistanceKM = miToKm(last.Lightningdistance)
	}
	s.seededLightning = true

	log.Printf("seed: live tempf=%.1f hum=%d baromabsin=%.3f aqi=%d (ws80=%v wh32b=%v ch1..4=%v,%v,%v,%v)",
		last.Tempf, last.Humidity, last.Baromabsin, last.Aqipm25,
		s.ws80HaveData, s.wh32BHaveData,
		s.channelHave[1], s.channelHave[2], s.channelHave[3], s.channelHave[4])
	log.Printf("seed: rain day=%.3fin week=%.3fin month=%.3fin year=%.3fin event=%.3fin",
		mmToIn(s.rain.dayMM), mmToIn(s.rain.weekMM),
		mmToIn(s.rain.monthMM), mmToIn(s.rain.yearMM), mmToIn(s.rain.eventMM))
	log.Printf("seed: lightning hour=%d day=%d month=%d",
		s.lightning.hourCount, s.lightning.dayCount, s.lightning.monthCount)
}

// NewSensorState constructs an empty state and seeds the rain/lightning period
// starts so the first minute after startup behaves sensibly.
func NewSensorState() *SensorState {
	return newSensorStateAt(time.Now())
}

// newSensorStateAt is NewSensorState with the period starts taken from now.
func newSensorStateAt(now time.Time) *SensorState {
	s := &SensorState{}
	s.rain.hourStart = startOfHour(now)
	s.rain.dayStart = startOfDay(now)
	s.rain.weekStart = startOfWeek(now)
	s.rain.monthStart = startOfMonth(now)
	s.rain.yearStart = startOfYear(now)
	s.lightning.hourStart = startOfHour(now)
	s.lightning.dayStart = startOfDay(now)
	s.lightning.monthStart = startOfMonth(now)
	return s
}

// startOfHour truncates the absolute instant rather than rebuilding it from
// wall-clock fields: during the repeated hour at the end of DST, time.Date
// resolves 01:xx to the first (MDT) 01:00, so the second 01:00 hour would
// never start. Truncation assumes a whole-hour UTC offset, which holds for
// America/Denver.
func startOfHour(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}
func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
func startOfWeek(t time.Time) time.Time {
	// ISO week: start on Monday.
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	d := startOfDay(t).AddDate(0, 0, -(wd - 1))
	return d
}
func startOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}
func startOfYear(t time.Time) time.Time {
	return time.Date(t.Year(), 1, 1, 0, 0, 0, 0, t.Location())
}

// Apply integrates a freshly-decoded sensor message into the state.
func (s *SensorState) Apply(m *SensorMessage, now time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch m.Model {

	case "Fineoffset-WS80":
		if m.TemperatureC != nil {
			s.outdoorTempC = *m.TemperatureC
		}
		if m.Humidity != nil {
			s.outdoorHumidity = *m.Humidity
		}
		if m.WindDirDeg != nil {
			s.windDirDeg = *m.WindDirDeg
		}
		if m.WindAvgMS != nil {
			s.windAvgMS = *m.WindAvgMS
		}
		if m.WindMaxMS != nil {
			s.windMaxMS = *m.WindMaxMS
			// Track today's maximum gust.
			today := startOfDay(now)
			if !s.maxGustDay.Equal(today) {
				s.maxDailyGustMS = 0
				s.maxGustDay = today
			}
			if *m.WindMaxMS > s.maxDailyGustMS {
				s.maxDailyGustMS = *m.WindMaxMS
				s.windGustDir = s.windDirDeg
			}
		}
		if m.UVI != nil {
			s.uvi = *m.UVI
		}
		if m.LightLux != nil {
			s.lightLux = *m.LightLux
		}
		if m.BatteryMV != nil {
			s.ws80BatteryMV = *m.BatteryMV
		}
		s.ws80HaveData = true

	case "Fineoffset-WH32B":
		// rtl_433 checks WH32B packets with an 8-bit sum only, so corrupt
		// packets get through. Like the console (fn_00c214), accept a packet
		// only if it is close to the previous one; the first packet after
		// startup is always accepted, and a genuine change shows up one
		// packet later.
		cur := wh32BReadingOf(m)
		accept := !s.wh32BPrevHave || s.wh32BPrev.close(cur)
		s.wh32BPrev, s.wh32BPrevHave = cur, true
		if !accept {
			log.Printf("wh32b: dropping reading too far from the previous packet (last accepted %s, got %s)",
				s.wh32BPrevString(), cur)
			return
		}
		if m.TemperatureC != nil {
			s.indoorTempC = *m.TemperatureC
		}
		if m.Humidity != nil {
			s.indoorHumidity = *m.Humidity
		}
		if m.PressureHPa != nil {
			s.pressureHPa = *m.PressureHPa
		}
		if m.BatteryOK != nil {
			s.wh32BBattery = int(*m.BatteryOK)
		}
		s.wh32BHaveData = true

	case "AmbientWeather-WH31E":
		if m.Channel == nil {
			return
		}
		ch := *m.Channel
		if ch < 1 || ch > 10 {
			return
		}
		if m.TemperatureC != nil {
			s.channelTempC[ch] = *m.TemperatureC
		}
		if m.Humidity != nil {
			s.channelHumidity[ch] = *m.Humidity
		}
		if m.BatteryOK != nil {
			s.channelBattery[ch] = int(*m.BatteryOK)
		}
		s.channelHave[ch] = true

	case "Fineoffset-WH0290":
		// A raw PM value of 0x3fff is the sensor's "no reading" marker (the
		// console firmware treats it as invalid); rtl_433 divides by 10 and
		// reports it as 1638. The sensor also sends it once its battery drops
		// below ~0.3 — treat either case as "no signal" and force the database
		// row to 0 rather than polluting the history.
		batteryLow := m.BatteryOK != nil && *m.BatteryOK < 0.3
		noReading := m.PM25 != nil && *m.PM25 == pmNoReading
		if m.BatteryOK != nil {
			s.aqBatteryRaw = *m.BatteryOK
		}
		if batteryLow || noReading {
			s.pm25 = 0
			s.pm10 = 0
			// Don't append the bogus reading to the 24h rolling buffer, but
			// do trim out any genuinely old samples so aqi24 ages out cleanly.
			cutoff := now.Add(-24 * time.Hour)
			i := 0
			for i < len(s.pm25History) && s.pm25History[i].t.Before(cutoff) {
				i++
			}
			s.pm25History = s.pm25History[i:]
		} else {
			if m.PM25 != nil {
				s.pm25 = *m.PM25
				s.pm25History = append(s.pm25History, pm25Sample{t: now, pm: *m.PM25})
				cutoff := now.Add(-24 * time.Hour)
				i := 0
				for i < len(s.pm25History) && s.pm25History[i].t.Before(cutoff) {
					i++
				}
				s.pm25History = s.pm25History[i:]
			}
			if m.PM10 != nil {
				s.pm10 = *m.PM10
				if s.pm10 == pmNoReading {
					s.pm10 = 0
				}
			}
		}
		s.aqHaveData = true

	case "EcoWitt-WH40":
		if m.RainMM == nil {
			return
		}
		s.applyRain(*m.RainMM, now)

	case "FineOffset-WH31L", "Fineoffset-WH31L":
		if m.BatteryOK != nil {
			s.lightningBattery = *m.BatteryOK
		}
		if m.StrikeCount != nil {
			s.applyLightning(*m.StrikeCount, m.StormDistance, m.State, now)
		}

	default:
		// Unknown / future sensors — ignore quietly without marking the state
		// as dirty, so a stream of garbage messages doesn't trigger snapshots.
		return
	}

	s.dirty = true
	s.lastUpdate = now
}

// applyRain integrates a WH40 rain_mm reading (a 16-bit counter of 0.1 mm).
// Each period total grows by the counter increment since the previous
// message, as in the console firmware (fn_00bb40): the increment is taken
// modulo the counter width so a wrap past 6553.5 mm still counts, and an
// increment larger than the console's per-second limit is treated as a
// reset (battery swap) or bad packet and dropped without touching the totals.
func (s *SensorState) applyRain(totalMM float64, now time.Time) {
	count := int(math.Round(totalMM * 10))

	// First reading: only synchronise the counter. The period totals are
	// already zero, or hold the carry-over Seed loaded from the last record.
	if !s.totalRainHaveIt {
		s.rain.lastCount = count
		s.rain.lastMsg = now
		s.totalRainMM = totalMM
		s.totalRainHaveIt = true
		if s.seededRain {
			log.Printf("rain: first reading %.2fmm, continuing seed (day=%.3fin)",
				totalMM, mmToIn(s.rain.dayMM))
		} else {
			log.Printf("rain: first reading, counter=%.2fmm", totalMM)
		}
		return
	}

	s.rollRain(now)

	delta := ((count-s.rain.lastCount)%rainCounterMod + rainCounterMod) % rainCounterMod
	elapsed := int(now.Sub(s.rain.lastMsg) / time.Second)
	if elapsed < 0 {
		elapsed = 0
	}
	limit := (elapsed + 1) * rainMaxPerSec
	if limit > rainMaxDelta {
		limit = rainMaxDelta
	}
	prevMM := s.totalRainMM
	s.rain.lastCount = count
	s.rain.lastMsg = now
	s.totalRainMM = totalMM

	switch {
	case delta == 0:
	case delta > limit:
		log.Printf("rain: counter jumped %.1f -> %.1f mm (%d counts, limit %d), treating as a reset",
			prevMM, totalMM, delta, limit)
	default:
		mm := float64(delta) / 10
		s.rain.window.add(now, mm)
		s.rain.dayMM += mm
		s.rain.weekMM += mm
		s.rain.monthMM += mm
		s.rain.yearMM += mm
		s.rain.eventMM += mm
		s.rain.lastRainTime = now
	}
}

// rollRain zeroes every rain period whose boundary the clock has crossed,
// and at each new hour applies the console's rain-event reset (FUN_0000ce80):
// the event ends when the last 60 minutes were dry and the last 24 hours
// had less than rainEventResetMM.
func (s *SensorState) rollRain(now time.Time) {
	if t := startOfHour(now); t.After(s.rain.hourStart) {
		s.rain.hourStart = t
		if s.rain.window.hourlyMM(now) < rainEpsilonMM &&
			s.rain.window.last24hMM(now) < rainEventResetMM {
			s.rain.eventMM = 0
		}
	}
	if t := startOfDay(now); t.After(s.rain.dayStart) {
		s.rain.dayStart, s.rain.dayMM = t, 0
	}
	if t := startOfWeek(now); t.After(s.rain.weekStart) {
		s.rain.weekStart, s.rain.weekMM = t, 0
	}
	if t := startOfMonth(now); t.After(s.rain.monthStart) {
		s.rain.monthStart, s.rain.monthMM = t, 0
	}
	if t := startOfYear(now); t.After(s.rain.yearStart) {
		s.rain.yearStart, s.rain.yearMM = t, 0
	}
}

// applyLightning is the strike-count analogue of applyRain. strike_count is
// one byte, so increments are taken modulo 256; like the console, only
// increments below strikeMaxDelta count. A "reset" state means the sensor has
// restarted its counter, so that message only re-synchronises.
func (s *SensorState) applyLightning(total int, distance *float64, state string, now time.Time) {
	if !s.totalStrikesHave {
		s.totalStrikes = total
		s.totalStrikesHave = true
		if s.seededLightning {
			log.Printf("lightning: first reading %d, continuing seed (day=%d)",
				total, s.lightning.dayCount)
		} else {
			log.Printf("lightning: first reading, counter=%d", total)
		}
		return
	}

	s.rollLightning(now)

	delta := ((total-s.totalStrikes)%strikeCounterMod + strikeCounterMod) % strikeCounterMod
	switch {
	case delta == 0:
	case state == "reset":
		log.Printf("lightning: sensor reset (%d -> %d), resyncing", s.totalStrikes, total)
	case delta >= strikeMaxDelta:
		log.Printf("lightning: counter jumped %d -> %d (%d), treating as a reset",
			s.totalStrikes, total, delta)
	case delta == 1 && (state == "interference" || state == "noise"):
		// A single increment tagged as interference or noise is
		// unambiguously not a strike. Larger jumps may mix real strikes
		// with noise (state only describes the latest event), so they
		// still count.
	default:
		s.lightning.hourCount += delta
		s.lightning.dayCount += delta
		s.lightning.monthCount += delta
		// The WH31L's counter advances on both real strikes and noise events.
		// Only update the timestamp and distance when the sensor classifies
		// the event as a real strike. Other state values seen in the wild
		// include "interference" and "disturber".
		if state == "strike" {
			s.lightning.lastStrikeTime = now
			if distance != nil {
				s.lightning.lastDistanceKM = *distance
			}
		}
	}

	s.totalStrikes = total
}

// rollLightning zeroes every lightning period whose boundary has passed.
func (s *SensorState) rollLightning(now time.Time) {
	if t := startOfHour(now); t.After(s.lightning.hourStart) {
		s.lightning.hourStart, s.lightning.hourCount = t, 0
	}
	if t := startOfDay(now); t.After(s.lightning.dayStart) {
		s.lightning.dayStart, s.lightning.dayCount = t, 0
	}
	if t := startOfMonth(now); t.After(s.lightning.monthStart) {
		s.lightning.monthStart, s.lightning.monthCount = t, 0
	}
}

// Snapshot builds a Record from the current state. Caller supplies the Mac
// (station identifier) and timestamp. This always returns a record, even when
// no new sensor data has arrived since the last call — use TrySnapshot if you
// want to skip writes while the MQTT feed is stalled.
func (s *SensorState) Snapshot(mac string, now time.Time) Record {
	s.mu.Lock()
	defer s.mu.Unlock()

	var r Record
	r.Mac = mac
	r.Recorded = now

	// Outdoor temperature / humidity / wind.
	if s.ws80HaveData {
		r.Tempf = round1(cToF(s.outdoorTempC))
		r.Humidity = s.outdoorHumidity
		r.Winddir = s.windDirDeg
		r.Windspeedmph = round1(msToMph(s.windAvgMS))
		r.Windgustmph = round1(msToMph(s.windMaxMS))
		r.Windgustdir = s.windGustDir
		r.Maxdailygust = round1(msToMph(s.maxDailyGustMS))
		r.Uv = s.uvi
		r.Solarradiation = round2(luxToWm2(s.lightLux))
		// 1.5V alkaline cell: 1=good, 0=low. Use 1 if reading > 1.0V.
		if s.ws80BatteryMV >= 2400 {
			r.Battout = 1
		} else {
			r.Battout = 0
		}
		r.Feelslike = round1(feelsLikeF(r.Tempf, float64(r.Humidity), r.Windspeedmph))
		r.Dewpoint = round1(dewPointF(r.Tempf, float64(r.Humidity)))
	}

	// Indoor barometer.
	if s.wh32BHaveData {
		r.Tempinf = round1(cToF(s.indoorTempC))
		r.Humidityin = s.indoorHumidity
		r.Baromabsin = round3(hpaToInHg(s.pressureHPa))
		// Relative (sea-level) pressure: corrects the station reading using
		// the configured elevation and current outdoor temperature. If the
		// outdoor sensor hasn't reported yet, fall back to ICAO standard temp
		// (15°C). If no elevation is configured, relative == absolute.
		outdoorC := s.outdoorTempC
		if !s.ws80HaveData {
			outdoorC = 15.0
		}
		r.Baromrelin = round3(hpaToInHg(
			stationToSeaLevelHPa(s.pressureHPa, s.elevationM, outdoorC)))
	}

	// WH31E channels. Channel 1 = basement -> temp1f/humidity1/batt1, etc.
	if s.channelHave[1] {
		r.Temp1f = round1(cToF(s.channelTempC[1]))
		r.Humidity1 = s.channelHumidity[1]
		r.Batt1 = s.channelBattery[1]
	}
	if s.channelHave[2] {
		r.Temp2f = round1(cToF(s.channelTempC[2]))
		r.Humidity2 = s.channelHumidity[2]
		r.Batt2 = s.channelBattery[2]
	}
	if s.channelHave[3] {
		r.Temp3f = round1(cToF(s.channelTempC[3]))
		r.Humidity3 = s.channelHumidity[3]
		r.Batt3 = s.channelBattery[3]
	}
	if s.channelHave[4] {
		r.Temp4f = round1(cToF(s.channelTempC[4]))
		r.Humidity4 = s.channelHumidity[4]
		r.Batt4 = s.channelBattery[4]
	}
	for ch := 5; ch <= 10; ch++ {
		if !s.channelHave[ch] {
			continue
		}
		f := round1(cToF(s.channelTempC[ch]))
		h := s.channelHumidity[ch]
		b := s.channelBattery[ch]
		switch ch {
		case 5:
			r.Temp5f, r.Humidity5, r.Batt5 = f, h, b
		case 6:
			r.Temp6f, r.Humidity6, r.Batt6 = f, h, b
		case 7:
			r.Temp7f, r.Humidity7, r.Batt7 = f, h, b
		case 8:
			r.Temp8f, r.Humidity8, r.Batt8 = f, h, b
		case 9:
			r.Temp9f, r.Humidity9, r.Batt9 = f, h, b
		case 10:
			r.Temp10f, r.Humidity10, r.Batt10 = f, h, b
		}
	}

	// Rain aggregations. Before the first WH40 message of this run the
	// totals are the seeded carry-over, which rolls over the same way.
	if s.totalRainHaveIt || s.seededRain {
		// Roll periods here too, so a period that ended since the last WH40
		// message reads zero rather than carrying the previous period's total.
		s.rollRain(now)
		if s.totalRainHaveIt {
			r.Totalrainin = round3(mmToIn(s.totalRainMM))
		}
		r.Dailyrainin = round3(mmToIn(s.rain.dayMM))
		r.Weeklyrainin = round3(mmToIn(s.rain.weekMM))
		r.Monthlyrainin = round3(mmToIn(s.rain.monthMM))
		r.Yearlyrainin = round3(mmToIn(s.rain.yearMM))
		// Event rain keeps showing after the rain stops; rollRain clears it
		// at the first hour boundary that meets the console's reset rule.
		r.Eventrainin = round3(mmToIn(s.rain.eventMM))
	}
	// Rolling values come from the per-minute window, which SeedRainHistory
	// refills after a restart. Hourly is the last 60 minutes (the console's
	// "Hourly"), rate is the last 10 minutes x 6 in in/hr.
	r.Hourlyrainin = round3(mmToIn(s.rain.window.hourlyMM(now)))
	r.Hourlyrain = r.Hourlyrainin
	r.Rainratein = round3(mmToIn(s.rain.window.rateMMPerHour(now)))
	if !s.rain.lastRainTime.IsZero() {
		r.Lastrain = s.rain.lastRainTime
	}

	// Lightning aggregations; seeded counts roll over like live ones.
	if s.totalStrikesHave || s.seededLightning {
		s.rollLightning(now)
		r.Lightninghour = s.lightning.hourCount
		r.Lightningday = s.lightning.dayCount
		r.LightningMonth = s.lightning.monthCount
	}
	if s.totalStrikesHave {
		// rtl_433 reports battery_ok as 0, 0.5 or 1.0; the console only
		// flags the WH57 battery as low at level 0.
		if s.lightningBattery > 0 {
			r.Battlightning = 1
		} else {
			r.Battlightning = 0
		}
	}
	// lastStrikeTime / distance are historical facts — emit whenever known,
	// regardless of whether we've seen a live message. The DB column uses
	// miles to match the rest of the imperial schema.
	if !s.lightning.lastStrikeTime.IsZero() {
		r.Lightningtime = s.lightning.lastStrikeTime
		r.Lightningdistance = round1(kmToMi(s.lightning.lastDistanceKM))
	}

	// Air quality.
	if s.aqHaveData {
		r.Aqipm25 = aqiFromPM25(s.pm25)
		// 24-hour average AQI.
		if len(s.pm25History) > 0 {
			var sum float64
			for _, p := range s.pm25History {
				sum += p.pm
			}
			avg := sum / float64(len(s.pm25History))
			r.Aqipm2524h = aqiFromPM25(avg)
		} else {
			r.Aqipm2524h = r.Aqipm25
		}
	} else if s.seededAQ {
		r.Aqipm25 = s.seedAQIPM25
		r.Aqipm2524h = s.seedAQIPM2524h
	}

	// Snapshot consumes the dirty flag — anyone calling Snapshot is committing
	// to writing the row, so further "no new data" decisions should start over.
	s.dirty = false

	return r
}

// TrySnapshot returns a Record only if new sensor data has arrived since the
// last snapshot. The second return value is false when nothing has changed, in
// which case the caller should skip the insert entirely to avoid hammering the
// database with duplicate rows while the MQTT feed is stalled.
func (s *SensorState) TrySnapshot(mac string, now time.Time) (Record, bool) {
	s.mu.Lock()
	dirty := s.dirty
	last := s.lastUpdate
	s.mu.Unlock()

	if !dirty {
		return Record{}, false
	}
	_ = last // available for callers that want to log staleness in future
	return s.Snapshot(mac, now), true
}

// math.Round rounds half away from zero, so negative values (sub-zero
// temperatures, dew points) round symmetrically with positive ones.
func round1(v float64) float64 { return math.Round(v*10) / 10.0 }
func round2(v float64) float64 { return math.Round(v*100) / 100.0 }
func round3(v float64) float64 { return math.Round(v*1000) / 1000.0 }

// wh32BReading is the part of a WH32B packet the glitch filter compares.
type wh32BReading struct {
	tempC            float64
	hum              int
	hPa              float64
	hasT, hasH, hasP bool
}

func wh32BReadingOf(m *SensorMessage) wh32BReading {
	var r wh32BReading
	if m.TemperatureC != nil {
		r.tempC, r.hasT = *m.TemperatureC, true
	}
	if m.Humidity != nil {
		r.hum, r.hasH = *m.Humidity, true
	}
	if m.PressureHPa != nil {
		r.hPa, r.hasP = *m.PressureHPa, true
	}
	return r
}

// close reports whether cur is within the console's limits of r. Fields
// missing from either packet are not compared.
func (r wh32BReading) close(cur wh32BReading) bool {
	const eps = 1e-9
	if r.hasT && cur.hasT && math.Abs(cur.tempC-r.tempC) > wh32BMaxDeltaC+eps {
		return false
	}
	if r.hasH && cur.hasH && absInt(cur.hum-r.hum) > wh32BMaxDeltaRH {
		return false
	}
	if r.hasP && cur.hasP && math.Abs(cur.hPa-r.hPa) > wh32BMaxDeltaHPa+eps {
		return false
	}
	return true
}

func (r wh32BReading) String() string {
	return fmt.Sprintf("%.1fC %d%% %.1fhPa", r.tempC, r.hum, r.hPa)
}

// wh32BPrevString formats the last accepted WH32B values for the drop log.
func (s *SensorState) wh32BPrevString() string {
	return fmt.Sprintf("%.1fC %d%% %.1fhPa", s.indoorTempC, s.indoorHumidity, s.pressureHPa)
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
