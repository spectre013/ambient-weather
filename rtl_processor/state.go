package main

import (
	"log"
	"sync"
	"time"
)

// pm25History keeps a sliding window of PM2.5 readings so we can compute a 24h AQI.
type pm25Sample struct {
	t  time.Time
	pm float64
}

// rainBaselines tracks the cumulative rain reading at the start of each period.
// Because the WH40 reports a monotonically-increasing total rain_mm, the rainfall
// during a period is simply (current - baseline).
type rainBaselines struct {
	hourTotalMM  float64
	dayTotalMM   float64
	weekTotalMM  float64
	monthTotalMM float64
	yearTotalMM  float64

	hourStart  time.Time
	dayStart   time.Time
	weekStart  time.Time
	monthStart time.Time
	yearStart  time.Time

	// eventStart is reset whenever rain stops for an hour or more.
	eventTotalMM float64
	lastRainTime time.Time
}

// lightningBaselines mirrors rainBaselines for strike counts.
type lightningBaselines struct {
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
	rain            rainBaselines

	// Lightning (from FineOffset-WH31L)
	totalStrikes     int
	totalStrikesHave bool
	lightningBattery float64
	lightning        lightningBaselines

	// Seeded values from the most recent DB record. These describe how much
	// rain/lightning has *already* accumulated in the current period before
	// this process started. They are consumed (folded into the baselines) on
	// the first WH40/WH31L message after startup so that period totals stay
	// continuous across restarts.
	seededRain      bool
	seedHourRainIn  float64
	seedDayRainIn   float64
	seedWeekRainIn  float64
	seedMonthRainIn float64
	seedYearRainIn  float64
	seedEventRainIn float64
	seedLastRain    time.Time
	seedMaxGustMph  float64
	seedMaxGustDay  time.Time

	seededLightning     bool
	seedLightningHour   int
	seedLightningDay    int
	seedLightningMonth  int
	seedLightningTime   time.Time
	seedLightningDistKM float64

	seededAQ       bool
	seedAQIPM25    int
	seedAQIPM2524h int

	// dirty is set by Apply whenever a sensor message updates the state and
	// cleared by Snapshot. It lets the caller skip an insert when nothing new
	// has arrived since the last snapshot (e.g. the MQTT feed has stalled).
	dirty      bool
	lastUpdate time.Time

	// rainEventGap is the idle interval after which a rain event is considered
	// ended. Once an event has ended, Eventrainin reports 0 until the next
	// rainfall starts a new event. NOAA convention is 1 hour.
	rainEventGap time.Duration

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
//  2. Rain and lightning period totals are stashed as carry-over so the
//     baselines computed on the first WH40/WH31L reading reflect the
//     accumulation that already happened this period. Period seeds that fall
//     in a previous period (e.g. last record was yesterday) are discarded so
//     the new period correctly starts at zero.
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
	// were down and the seed for it should be 0.
	if startOfHour(last.Recorded).Equal(startOfHour(now)) {
		s.seedHourRainIn = last.Hourlyrainin
	}
	if startOfDay(last.Recorded).Equal(startOfDay(now)) {
		s.seedDayRainIn = last.Dailyrainin
		s.seedMaxGustMph = last.Maxdailygust
		s.seedMaxGustDay = startOfDay(now)
	}
	if startOfWeek(last.Recorded).Equal(startOfWeek(now)) {
		s.seedWeekRainIn = last.Weeklyrainin
	}
	if startOfMonth(last.Recorded).Equal(startOfMonth(now)) {
		s.seedMonthRainIn = last.Monthlyrainin
	}
	if startOfYear(last.Recorded).Equal(startOfYear(now)) {
		s.seedYearRainIn = last.Yearlyrainin
	}
	// Event total carries forward as long as it's been raining recently.
	if !last.Lastrain.IsZero() && now.Sub(last.Lastrain) < time.Hour {
		s.seedEventRainIn = last.Eventrainin
	}
	// lastrain timestamp always carries forward — it's a historical fact, not
	// a per-period accumulation. Set it directly on the rain baselines so that
	// Snapshot emits it even before the first WH40 message arrives.
	if !last.Lastrain.IsZero() {
		s.seedLastRain = last.Lastrain
		s.rain.lastRainTime = last.Lastrain
	}
	s.seededRain = true

	// Restore today's running max gust.
	if !s.seedMaxGustDay.IsZero() {
		s.maxDailyGustMS = s.seedMaxGustMph / 2.236936
		s.maxGustDay = s.seedMaxGustDay
	}

	// Lightning seeds.
	if startOfHour(last.Recorded).Equal(startOfHour(now)) {
		s.seedLightningHour = last.Lightninghour
	}
	if startOfDay(last.Recorded).Equal(startOfDay(now)) {
		s.seedLightningDay = last.Lightningday
	}
	if startOfMonth(last.Recorded).Equal(startOfMonth(now)) {
		s.seedLightningMonth = last.LightningMonth
	}
	// Same as lastrain: the lightning timestamp/distance are historical facts.
	// Set them directly so they survive into the next Snapshot without waiting
	// for the WH31L to chime in. DB column is miles; in-memory state is km.
	if !last.Lightningtime.IsZero() {
		s.seedLightningTime = last.Lightningtime
		s.seedLightningDistKM = miToKm(last.Lightningdistance)
		s.lightning.lastStrikeTime = last.Lightningtime
		s.lightning.lastDistanceKM = miToKm(last.Lightningdistance)
	}
	s.seededLightning = true

	log.Printf("seed: live tempf=%.1f hum=%d baromabsin=%.3f aqi=%d (ws80=%v wh32b=%v ch1..4=%v,%v,%v,%v)",
		last.Tempf, last.Humidity, last.Baromabsin, last.Aqipm25,
		s.ws80HaveData, s.wh32BHaveData,
		s.channelHave[1], s.channelHave[2], s.channelHave[3], s.channelHave[4])
	log.Printf("seed: rain hour=%.3fin day=%.3fin week=%.3fin month=%.3fin year=%.3fin event=%.3fin",
		s.seedHourRainIn, s.seedDayRainIn, s.seedWeekRainIn,
		s.seedMonthRainIn, s.seedYearRainIn, s.seedEventRainIn)
	log.Printf("seed: lightning hour=%d day=%d month=%d",
		s.seedLightningHour, s.seedLightningDay, s.seedLightningMonth)
}

// NewSensorState constructs an empty state and seeds the rain/lightning period
// starts so the first minute after startup behaves sensibly.
func NewSensorState() *SensorState {
	now := time.Now()
	s := &SensorState{
		rainEventGap: time.Hour,
	}
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

func startOfHour(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, t.Location())
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
		// The WH0290 sends nonsense PM values (typically pegged at 1638) once
		// its battery drops below ~0.3 — treat those as "no signal" and force
		// the database row to 0 rather than polluting the history.
		batteryLow := m.BatteryOK != nil && *m.BatteryOK < 0.3
		if m.BatteryOK != nil {
			s.aqBatteryRaw = *m.BatteryOK
		}
		if batteryLow {
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

// applyRain assumes the WH40 reports a monotonically increasing total in mm.
// It rolls over period baselines when the clock crosses a boundary.
func (s *SensorState) applyRain(totalMM float64, now time.Time) {
	// First reading: seed baselines so this snapshot reads "zero" for the
	// period — unless we were given carry-over from a previous run, in which
	// case we offset the baselines so each period reads the seeded amount.
	if !s.totalRainHaveIt {
		// mmToIn is mm/25.4, so inToMM is in*25.4.
		const inToMM = 25.4
		s.rain.hourTotalMM = totalMM - s.seedHourRainIn*inToMM
		s.rain.dayTotalMM = totalMM - s.seedDayRainIn*inToMM
		s.rain.weekTotalMM = totalMM - s.seedWeekRainIn*inToMM
		s.rain.monthTotalMM = totalMM - s.seedMonthRainIn*inToMM
		s.rain.yearTotalMM = totalMM - s.seedYearRainIn*inToMM
		s.rain.eventTotalMM = totalMM - s.seedEventRainIn*inToMM
		if !s.seedLastRain.IsZero() {
			s.rain.lastRainTime = s.seedLastRain
		}
		s.totalRainMM = totalMM
		s.totalRainHaveIt = true
		if s.seededRain {
			log.Printf("rain: first reading %.2fmm, applied seed (day=%.3fin)",
				totalMM, s.seedDayRainIn)
		} else {
			log.Printf("rain: first reading, baseline=%.2fmm", totalMM)
		}
		return
	}

	// Detect a hardware reset (counter decreased): rebaseline everything.
	if totalMM < s.totalRainMM {
		log.Printf("rain: counter went backwards (%.2f -> %.2f), rebaselining",
			s.totalRainMM, totalMM)
		s.rain.hourTotalMM = totalMM
		s.rain.dayTotalMM = totalMM
		s.rain.weekTotalMM = totalMM
		s.rain.monthTotalMM = totalMM
		s.rain.yearTotalMM = totalMM
		s.rain.eventTotalMM = totalMM
		s.totalRainMM = totalMM
		return
	}

	// Roll over period baselines if the clock crossed a boundary.
	if startOfHour(now).After(s.rain.hourStart) {
		s.rain.hourStart = startOfHour(now)
		s.rain.hourTotalMM = s.totalRainMM
	}
	if startOfDay(now).After(s.rain.dayStart) {
		s.rain.dayStart = startOfDay(now)
		s.rain.dayTotalMM = s.totalRainMM
	}
	if startOfWeek(now).After(s.rain.weekStart) {
		s.rain.weekStart = startOfWeek(now)
		s.rain.weekTotalMM = s.totalRainMM
	}
	if startOfMonth(now).After(s.rain.monthStart) {
		s.rain.monthStart = startOfMonth(now)
		s.rain.monthTotalMM = s.totalRainMM
	}
	if startOfYear(now).After(s.rain.yearStart) {
		s.rain.yearStart = startOfYear(now)
		s.rain.yearTotalMM = s.totalRainMM
	}

	// If the gauge actually counted up, it just rained.
	if totalMM > s.totalRainMM {
		// If this raindrop arrives after the rain-event idle gap, the previous
		// event is over — start a new event at the current counter so that
		// Eventrainin counts only this new event's accumulation.
		if !s.rain.lastRainTime.IsZero() &&
			now.Sub(s.rain.lastRainTime) > s.rainEventGap {
			s.rain.eventTotalMM = s.totalRainMM
		}
		s.rain.lastRainTime = now
	}

	s.totalRainMM = totalMM
}

// applyLightning is the strike-count analogue of applyRain.
func (s *SensorState) applyLightning(total int, distance *float64, state string, now time.Time) {
	if !s.totalStrikesHave {
		s.lightning.hourCount = total - s.seedLightningHour
		s.lightning.dayCount = total - s.seedLightningDay
		s.lightning.monthCount = total - s.seedLightningMonth
		if !s.seedLightningTime.IsZero() {
			s.lightning.lastStrikeTime = s.seedLightningTime
			s.lightning.lastDistanceKM = s.seedLightningDistKM
		}
		s.totalStrikes = total
		s.totalStrikesHave = true
		if s.seededLightning {
			log.Printf("lightning: first reading %d, applied seed (day=%d)",
				total, s.seedLightningDay)
		} else {
			log.Printf("lightning: first reading, baseline=%d", total)
		}
		return
	}

	if total < s.totalStrikes {
		log.Printf("lightning: counter went backwards (%d -> %d), rebaselining",
			s.totalStrikes, total)
		s.lightning.hourCount = total
		s.lightning.dayCount = total
		s.lightning.monthCount = total
		s.totalStrikes = total
		return
	}

	if startOfHour(now).After(s.lightning.hourStart) {
		s.lightning.hourStart = startOfHour(now)
		s.lightning.hourCount = s.totalStrikes
	}
	if startOfDay(now).After(s.lightning.dayStart) {
		s.lightning.dayStart = startOfDay(now)
		s.lightning.dayCount = s.totalStrikes
	}
	if startOfMonth(now).After(s.lightning.monthStart) {
		s.lightning.monthStart = startOfMonth(now)
		s.lightning.monthCount = s.totalStrikes
	}

	// The WH31L's counter advances on both real strikes and noise events.
	// Only update the timestamp and distance when the sensor classifies the
	// event as a real strike. Other state values seen in the wild include
	// "interference" and "disturber".
	if total > s.totalStrikes && state == "strike" {
		s.lightning.lastStrikeTime = now
		if distance != nil {
			s.lightning.lastDistanceKM = *distance
		}
	}

	s.totalStrikes = total
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

	// Rain aggregations.
	if s.totalRainHaveIt {
		r.Totalrainin = round3(mmToIn(s.totalRainMM))
		r.Hourlyrainin = round3(mmToIn(s.totalRainMM - s.rain.hourTotalMM))
		r.Hourlyrain = r.Hourlyrainin
		r.Dailyrainin = round3(mmToIn(s.totalRainMM - s.rain.dayTotalMM))
		r.Weeklyrainin = round3(mmToIn(s.totalRainMM - s.rain.weekTotalMM))
		r.Monthlyrainin = round3(mmToIn(s.totalRainMM - s.rain.monthTotalMM))
		r.Yearlyrainin = round3(mmToIn(s.totalRainMM - s.rain.yearTotalMM))
		// Event rain is the accumulation from the start of the current rain
		// event up to the moment it stops. Once it's been more than
		// rainEventGap since the last rainfall, the event is over and
		// Eventrainin reads 0 until new rain starts the next event.
		if s.rain.lastRainTime.IsZero() ||
			now.Sub(s.rain.lastRainTime) > s.rainEventGap {
			r.Eventrainin = 0
		} else {
			r.Eventrainin = round3(mmToIn(s.totalRainMM - s.rain.eventTotalMM))
		}
	} else if s.seededRain {
		// No WH40 message yet this run, but we have values from the last DB
		// row — emit those so periods stay continuous.
		r.Hourlyrainin = s.seedHourRainIn
		r.Hourlyrain = s.seedHourRainIn
		r.Dailyrainin = s.seedDayRainIn
		r.Weeklyrainin = s.seedWeekRainIn
		r.Monthlyrainin = s.seedMonthRainIn
		r.Yearlyrainin = s.seedYearRainIn
		// Seed path: same end-of-event logic, using the last-known rain time.
		if s.rain.lastRainTime.IsZero() ||
			now.Sub(s.rain.lastRainTime) > s.rainEventGap {
			r.Eventrainin = 0
		} else {
			r.Eventrainin = s.seedEventRainIn
		}
	}
	if !s.rain.lastRainTime.IsZero() {
		r.Lastrain = s.rain.lastRainTime
	}

	// Lightning aggregations.
	if s.totalStrikesHave {
		r.Lightninghour = s.totalStrikes - s.lightning.hourCount
		r.Lightningday = s.totalStrikes - s.lightning.dayCount
		r.LightningMonth = s.totalStrikes - s.lightning.monthCount
		if s.lightningBattery >= 1.0 {
			r.Battlightning = 1
		} else {
			r.Battlightning = 0
		}
	} else if s.seededLightning {
		// No WH31L message yet this run — fall back to seeded period counts.
		r.Lightninghour = s.seedLightningHour
		r.Lightningday = s.seedLightningDay
		r.LightningMonth = s.seedLightningMonth
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

func round1(v float64) float64 { return float64(int(v*10+0.5)) / 10.0 }
func round2(v float64) float64 { return float64(int(v*100+0.5)) / 100.0 }
func round3(v float64) float64 { return float64(int(v*1000+0.5)) / 1000.0 }
