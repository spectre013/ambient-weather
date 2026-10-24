package main

import (
	"testing"
	"time"
)

func f64(v float64) *float64 { return &v }
func intp(v int) *int        { return &v }

// t0 is a Wednesday afternoon, well away from any period boundary.
var t0 = time.Date(2025, 3, 26, 14, 10, 0, 0, time.UTC)

func rainMsg(mm float64) *SensorMessage {
	return &SensorMessage{Model: "EcoWitt-WH40", RainMM: f64(mm)}
}

func strikeMsg(count int, state string) *SensorMessage {
	return &SensorMessage{Model: "FineOffset-WH31L", StrikeCount: intp(count), State: state,
		StormDistance: f64(12), BatteryOK: f64(1)}
}

func pmMsg(pm float64) *SensorMessage {
	return &SensorMessage{Model: "Fineoffset-WH0290", PM25: f64(pm), PM10: f64(pm), BatteryOK: f64(1)}
}

// feedRain applies readings one minute apart starting at start and returns
// the time of the last one.
func feedRain(s *SensorState, start time.Time, mm ...float64) time.Time {
	now := start
	for i, v := range mm {
		if i > 0 {
			now = now.Add(time.Minute)
		}
		s.Apply(rainMsg(v), now)
	}
	return now
}

func TestRainAccumulation(t *testing.T) {
	tests := []struct {
		name     string
		readings []float64
		wantDay  float64 // inches
	}{
		{"steady increase", []float64{100.0, 100.3, 100.5}, 0.020},
		// 6553.0 -> 6553.5 is 0.5 mm; 6553.5 -> 0.3 is 0.4 mm (wrap at 65536 counts).
		{"counter wraps past 6553.5 mm", []float64{6553.0, 6553.5, 0.3}, 0.035},
		{"battery swap keeps totals", []float64{100.0, 105.0, 0.0, 0.3}, 0.209},
		{"implausible jump is dropped", []float64{100.0, 300.0, 300.2}, 0.008},
		{"no change", []float64{42.0, 42.0, 42.0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSensorStateAt(t0)
			now := feedRain(s, t0, tt.readings...)
			r := s.Snapshot("test", now)
			if r.Dailyrainin != tt.wantDay {
				t.Errorf("Dailyrainin = %v, want %v", r.Dailyrainin, tt.wantDay)
			}
			for name, v := range map[string]float64{"weekly": r.Weeklyrainin, "monthly": r.Monthlyrainin, "yearly": r.Yearlyrainin} {
				if v != tt.wantDay {
					t.Errorf("%s = %v, want %v", name, v, tt.wantDay)
				}
			}
		})
	}
}

func TestRainJumpLimitScalesWithElapsedTime(t *testing.T) {
	// 20 mm is 200 counts: too much after 10 s (limit 44), fine after an hour
	// (limit capped at 1023 counts).
	s := newSensorStateAt(t0)
	s.Apply(rainMsg(10), t0)
	s.Apply(rainMsg(30), t0.Add(10*time.Second))
	if got := s.rain.dayMM; got != 0 {
		t.Fatalf("after 10 s: dayMM = %v, want 0", got)
	}
	s.Apply(rainMsg(50), t0.Add(time.Hour))
	if got := s.rain.dayMM; got != 20 {
		t.Fatalf("after 1 h: dayMM = %v, want 20", got)
	}
}

func TestRainPeriodRollsWithoutNewMessage(t *testing.T) {
	evening := time.Date(2025, 3, 26, 23, 50, 0, 0, time.UTC) // Wednesday
	s := newSensorStateAt(evening)
	feedRain(s, evening, 10.0, 15.0) // 5.0 mm = 0.197 in at 23:51

	r := s.Snapshot("test", evening.Add(15*time.Minute)) // Thursday 00:05
	if r.Dailyrainin != 0 {
		t.Errorf("after midnight: daily=%v, want 0", r.Dailyrainin)
	}
	if r.Weeklyrainin != 0.197 || r.Monthlyrainin != 0.197 {
		t.Errorf("after midnight: weekly=%v monthly=%v, want 0.197", r.Weeklyrainin, r.Monthlyrainin)
	}
	// Hourly is a rolling 60 minutes, so it does not reset at midnight.
	if r.Hourlyrainin != 0.197 {
		t.Errorf("after midnight: hourly=%v, want 0.197", r.Hourlyrainin)
	}
}

func TestRainRollingHourlyAndRate(t *testing.T) {
	s := newSensorStateAt(t0)
	s.Apply(rainMsg(10.0), t0)
	s.Apply(rainMsg(12.5), t0.Add(time.Minute))    // 2.5 mm at t0+1m
	s.Apply(rainMsg(15.0), t0.Add(30*time.Minute)) // 2.5 mm at t0+30m
	s.Apply(rainMsg(15.0), t0.Add(55*time.Minute)) // nothing new
	for _, tt := range []struct {
		at             time.Duration
		hourly, rateIn float64
	}{
		{1 * time.Minute, 0.098, 0.591},  // 2.5 mm in the last 10 min -> 15 mm/h
		{30 * time.Minute, 0.197, 0.591}, // both falls in the hour; latest drives the rate
		{39 * time.Minute, 0.197, 0.591}, // 10-minute rate window still holds t0+30m
		{40 * time.Minute, 0.197, 0},     // rate window has passed it
		{60 * time.Minute, 0.197, 0},     // t0+1m is the 60th minute back
		{61 * time.Minute, 0.098, 0},     // first fall has left the hour
		{90 * time.Minute, 0, 0},
	} {
		r := s.Snapshot("test", t0.Add(tt.at))
		if r.Hourlyrainin != tt.hourly || r.Hourlyrain != tt.hourly || r.Rainratein != tt.rateIn {
			t.Errorf("t0+%v: hourly=%v/%v rate=%v, want hourly=%v rate=%v",
				tt.at, r.Hourlyrainin, r.Hourlyrain, r.Rainratein, tt.hourly, tt.rateIn)
		}
	}
}

func TestSeedRainHistory(t *testing.T) {
	// Rows one minute apart across midnight: 0.10 in at 23:58, +0.02 at 23:59,
	// then the daily reset with 0.03 in by 00:00, +0.01 at 00:01.
	base := time.Date(2025, 3, 26, 23, 58, 0, 0, time.UTC)
	rows := []RainSample{
		{base, 0.10},
		{base.Add(1 * time.Minute), 0.12},
		{base.Add(2 * time.Minute), 0.03},
		{base.Add(3 * time.Minute), 0.04},
	}
	s := newSensorStateAt(base.Add(3 * time.Minute))
	s.SeedRainHistory(rows)
	r := s.Snapshot("test", base.Add(3*time.Minute))
	if r.Hourlyrainin != 0.06 {
		t.Errorf("hourly after seeding = %v, want 0.06", r.Hourlyrainin)
	}
}

func TestRainEventConsoleRule(t *testing.T) {
	s := newSensorStateAt(t0)          // 14:10
	now := feedRain(s, t0, 10.0, 15.0) // 5.0 mm = 0.197 in at 14:11
	if r := s.Snapshot("test", now); r.Eventrainin != 0.197 {
		t.Fatalf("event = %v, want 0.197", r.Eventrainin)
	}
	// The event keeps showing after the rain stops: at 16:00 the last hour
	// was dry but the last 24 h still hold 5 mm (>= 1.0 mm).
	at16 := time.Date(2025, 3, 26, 16, 0, 0, 0, time.UTC)
	if r := s.Snapshot("test", at16); r.Eventrainin != 0.197 {
		t.Fatalf("event at 16:00 = %v, want 0.197", r.Eventrainin)
	}
	// More rain the next morning adds to the same event.
	next := time.Date(2025, 3, 27, 9, 0, 0, 0, time.UTC)
	s.Apply(rainMsg(15.0), next)
	s.Apply(rainMsg(17.0), next.Add(time.Minute)) // 2.0 mm
	if r := s.Snapshot("test", next.Add(time.Minute)); r.Eventrainin != 0.276 {
		t.Fatalf("event next morning = %v, want 0.276", r.Eventrainin)
	}
	// At 10:00 the last hour had rain: no reset. Once 24 h pass with < 1 mm,
	// the next hour boundary clears it: by 10:00 on the 28th only 0 mm remain.
	if r := s.Snapshot("test", next.Add(time.Hour)); r.Eventrainin != 0.276 {
		t.Fatalf("event at 10:00 = %v, want 0.276", r.Eventrainin)
	}
	end := time.Date(2025, 3, 28, 10, 0, 0, 0, time.UTC)
	if r := s.Snapshot("test", end); r.Eventrainin != 0 {
		t.Fatalf("event after a dry day = %v, want 0", r.Eventrainin)
	}
}

func TestRainEventNotResetBySmallDrizzleDay(t *testing.T) {
	// 0.5 mm in 24 h is below the 1.0 mm threshold, so the event ends at the
	// first dry hour boundary.
	s := newSensorStateAt(t0)
	feedRain(s, t0, 10.0, 10.5)
	at16 := time.Date(2025, 3, 26, 16, 0, 0, 0, time.UTC)
	if r := s.Snapshot("test", at16); r.Eventrainin != 0 {
		t.Fatalf("event = %v, want 0", r.Eventrainin)
	}
}

func TestRainSeedCarriesOver(t *testing.T) {
	s := newSensorStateAt(t0)
	s.Seed(&Record{
		Recorded:    t0.Add(-5 * time.Minute),
		Dailyrainin: 0.1, Weeklyrainin: 0.5, Monthlyrainin: 1, Yearlyrainin: 2,
	}, t0)
	feedRain(s, t0, 500.0, 505.0) // +5.0 mm = 0.197 in
	r := s.Snapshot("test", t0.Add(time.Minute))
	want := map[string][2]float64{
		"daily": {r.Dailyrainin, 0.297}, "weekly": {r.Weeklyrainin, 0.697},
		"monthly": {r.Monthlyrainin, 1.197}, "yearly": {r.Yearlyrainin, 2.197},
	}
	for name, v := range want {
		if v[0] != v[1] {
			t.Errorf("%s = %v, want %v", name, v[0], v[1])
		}
	}
}

func TestLightningCounter(t *testing.T) {
	type step struct {
		count int
		state string
	}
	tests := []struct {
		name    string
		steps   []step
		wantDay int
	}{
		{"strikes", []step{{10, "strike"}, {12, "strike"}, {15, "strike"}}, 5},
		{"counter wraps past 255", []step{{250, "strike"}, {254, "strike"}, {3, "strike"}}, 9},
		{"sensor reset resyncs", []step{{40, "strike"}, {42, "strike"}, {0, "reset"}, {2, "strike"}}, 4},
		{"implausible jump dropped", []step{{10, "strike"}, {200, "strike"}, {201, "strike"}}, 1},
		{"single interference skipped", []step{{10, "strike"}, {11, "interference"}}, 0},
		{"single noise skipped", []step{{10, "strike"}, {11, "noise"}, {12, "strike"}}, 1},
		{"multi-count jump tagged noise still counts", []step{{10, "strike"}, {13, "noise"}}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := newSensorStateAt(t0)
			now := t0
			for _, st := range tt.steps {
				now = now.Add(time.Minute)
				s.Apply(strikeMsg(st.count, st.state), now)
			}
			r := s.Snapshot("test", now)
			if r.Lightningday != tt.wantDay || r.LightningMonth != tt.wantDay || r.Lightninghour != tt.wantDay {
				t.Errorf("hour/day/month = %d/%d/%d, want %d",
					r.Lightninghour, r.Lightningday, r.LightningMonth, tt.wantDay)
			}
		})
	}
}

func TestLightningStrikeTimeOnlyOnStrike(t *testing.T) {
	s := newSensorStateAt(t0)
	s.Apply(strikeMsg(10, "strike"), t0)
	s.Apply(strikeMsg(11, "interference"), t0.Add(time.Minute))
	if !s.lightning.lastStrikeTime.IsZero() {
		t.Fatalf("interference set strike time %v", s.lightning.lastStrikeTime)
	}
	s.Apply(strikeMsg(12, "strike"), t0.Add(2*time.Minute))
	if !s.lightning.lastStrikeTime.Equal(t0.Add(2 * time.Minute)) {
		t.Fatalf("strike time = %v", s.lightning.lastStrikeTime)
	}
}

func TestLightningDayRollsWithoutNewMessage(t *testing.T) {
	evening := time.Date(2025, 3, 26, 23, 50, 0, 0, time.UTC)
	s := newSensorStateAt(evening)
	s.Apply(strikeMsg(10, "strike"), evening)
	s.Apply(strikeMsg(13, "strike"), evening.Add(time.Minute))
	r := s.Snapshot("test", evening.Add(20*time.Minute))
	if r.Lightningday != 0 || r.LightningMonth != 3 {
		t.Errorf("day/month = %d/%d, want 0/3", r.Lightningday, r.LightningMonth)
	}
}

func TestLightningBattery(t *testing.T) {
	for _, tt := range []struct {
		ok   float64
		want int
	}{{1.0, 1}, {0.5, 1}, {0, 0}} {
		s := newSensorStateAt(t0)
		m := strikeMsg(5, "strike")
		m.BatteryOK = f64(tt.ok)
		s.Apply(m, t0)
		if r := s.Snapshot("test", t0); r.Battlightning != tt.want {
			t.Errorf("battery_ok %.1f: Battlightning = %d, want %d", tt.ok, r.Battlightning, tt.want)
		}
	}
}

func TestPMNoReadingSentinel(t *testing.T) {
	s := newSensorStateAt(t0)
	s.Apply(pmMsg(8), t0)
	s.Apply(pmMsg(pmNoReading), t0.Add(time.Minute))
	if s.pm25 != 0 || s.pm10 != 0 {
		t.Fatalf("pm25/pm10 = %v/%v after sentinel, want 0/0", s.pm25, s.pm10)
	}
	if len(s.pm25History) != 1 {
		t.Fatalf("history has %d samples, want 1 (sentinel must not be added)", len(s.pm25History))
	}
	// PM10 alone carrying the sentinel is cleared without affecting PM2.5.
	m := pmMsg(9)
	m.PM10 = f64(pmNoReading)
	s.Apply(m, t0.Add(2*time.Minute))
	if s.pm25 != 9 || s.pm10 != 0 {
		t.Fatalf("pm25/pm10 = %v/%v, want 9/0", s.pm25, s.pm10)
	}
}

func TestAQI24hAverageBetweenBreakpoints(t *testing.T) {
	// rtl_433 reports whole µg/m³; 19 x 12 + 1 x 13 averages 12.05, which
	// used to fall between the 12.0 and 12.1 breakpoints and give AQI 0.
	s := newSensorStateAt(t0)
	now := t0
	for i := 0; i < 19; i++ {
		s.Apply(pmMsg(12), now)
		now = now.Add(time.Minute)
	}
	s.Apply(pmMsg(13), now)
	r := s.Snapshot("test", now)
	if r.Aqipm2524h != 50 {
		t.Errorf("aqi24 = %d, want 50", r.Aqipm2524h)
	}
	if r.Aqipm25 != aqiFromPM25(13) {
		t.Errorf("aqi = %d, want %d", r.Aqipm25, aqiFromPM25(13))
	}
}

func wh32bMsg(tempC float64, hum int, hpa float64) *SensorMessage {
	return &SensorMessage{Model: "Fineoffset-WH32B", TemperatureC: f64(tempC), Humidity: intp(hum),
		PressureHPa: f64(hpa), BatteryOK: f64(1)}
}

func TestWH32BGlitchFilter(t *testing.T) {
	s := newSensorStateAt(t0)
	steps := []struct {
		m        *SensorMessage
		wantTemp float64
	}{
		{wh32bMsg(21.0, 40, 840.0), 21.0}, // first packet always accepted
		{wh32bMsg(21.4, 42, 840.8), 21.4}, // within 0.5C / 3% / 1.0hPa
		{wh32bMsg(35.0, 42, 840.8), 21.4}, // spike dropped
		{wh32bMsg(21.5, 42, 840.8), 21.4}, // far from the spike packet, dropped too
		{wh32bMsg(21.5, 42, 840.9), 21.5}, // two agreeing packets: accepted
		{wh32bMsg(21.5, 47, 840.9), 21.5}, // humidity jump of 5% dropped
		{wh32bMsg(21.5, 47, 842.5), 21.5}, // pressure jump of 1.6 hPa dropped
		{wh32bMsg(21.5, 47, 842.4), 21.5}, // agrees with previous: accepted
	}
	for i, st := range steps {
		s.Apply(st.m, t0.Add(time.Duration(i)*time.Minute))
		if s.indoorTempC != st.wantTemp {
			t.Fatalf("step %d: indoor temp = %v, want %v", i, s.indoorTempC, st.wantTemp)
		}
	}
	if s.indoorHumidity != 47 || s.pressureHPa != 842.4 {
		t.Fatalf("final hum/press = %d/%v, want 47/842.4", s.indoorHumidity, s.pressureHPa)
	}
}

func TestWH32BMissingPressureNotCompared(t *testing.T) {
	s := newSensorStateAt(t0)
	s.Apply(wh32bMsg(21.0, 40, 840.0), t0)
	m := wh32bMsg(21.2, 41, 0)
	m.PressureHPa = nil // rtl_433 omits pressure when the raw value is 0xffff
	s.Apply(m, t0.Add(time.Minute))
	if s.indoorTempC != 21.2 || s.pressureHPa != 840.0 {
		t.Fatalf("temp/press = %v/%v, want 21.2/840", s.indoorTempC, s.pressureHPa)
	}
}
