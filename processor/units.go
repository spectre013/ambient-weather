package main

import "math"

// cToF converts Celsius to Fahrenheit.
func cToF(c float64) float64 {
	return c*9.0/5.0 + 32.0
}

// msToMph converts meters/second to miles/hour.
func msToMph(ms float64) float64 {
	return ms * 2.236936
}

// mmToIn converts millimeters to inches.
func mmToIn(mm float64) float64 {
	return mm / 25.4
}

// kmToMi converts kilometers to miles.
func kmToMi(km float64) float64 {
	return km * 0.621371
}

// miToKm converts miles back to kilometers (used when seeding state from the
// imperial-unit DB column).
func miToKm(mi float64) float64 {
	return mi / 0.621371
}

// hpaToInHg converts hectopascals to inches of mercury.
func hpaToInHg(hpa float64) float64 {
	return hpa * 0.02953
}

// stationToSeaLevelHPa converts absolute (station) pressure to sea-level
// (relative) pressure using the hypsometric approximation.
//
//	pStation: absolute pressure in hPa from the barometer
//	elevationM: station elevation above sea level in meters
//	tempC: outdoor temperature in Celsius (use 15°C if unknown)
//
// Accurate to ~0.1 hPa for typical home-station elevations.
func stationToSeaLevelHPa(pStation, elevationM, tempC float64) float64 {
	tK := tempC + 273.15
	return pStation * math.Exp(elevationM/(29.3*tK))
}

// luxToWm2 is a rough conversion of lux to W/m^2 of solar radiation.
// The Davis convention is lux / 126.7.
func luxToWm2(lux float64) float64 {
	return lux / 126.7
}

// dewPointF computes dew point in F given temperature in F and relative humidity %.
// Uses the Magnus formula.
func dewPointF(tempF float64, rh float64) float64 {
	if rh <= 0 {
		return tempF
	}
	tC := (tempF - 32.0) * 5.0 / 9.0
	a := 17.625
	b := 243.04
	alpha := math.Log(rh/100.0) + (a*tC)/(b+tC)
	dpC := (b * alpha) / (a - alpha)
	return cToF(dpC)
}

// heatIndexF returns the NOAA heat index given T(F) and RH(%).
// Below 80F it returns the raw temperature.
func heatIndexF(tempF, rh float64) float64 {
	if tempF < 80.0 {
		return tempF
	}
	t := tempF
	r := rh
	hi := -42.379 +
		2.04901523*t +
		10.14333127*r -
		0.22475541*t*r -
		0.00683783*t*t -
		0.05481717*r*r +
		0.00122874*t*t*r +
		0.00085282*t*r*r -
		0.00000199*t*t*r*r
	return hi
}

// windChillF returns the NWS wind-chill given T(F) and wind speed(mph).
// Only valid when T<=50 and wind>3.
func windChillF(tempF, windMph float64) float64 {
	if tempF > 50.0 || windMph <= 3.0 {
		return tempF
	}
	return 35.74 + 0.6215*tempF - 35.75*math.Pow(windMph, 0.16) +
		0.4275*tempF*math.Pow(windMph, 0.16)
}

// feelsLikeF picks heat index above 80, wind chill below 50, otherwise temp.
func feelsLikeF(tempF, rh, windMph float64) float64 {
	switch {
	case tempF >= 80.0:
		return heatIndexF(tempF, rh)
	case tempF <= 50.0:
		return windChillF(tempF, windMph)
	default:
		return tempF
	}
}

// aqiFromPM25 converts a PM2.5 concentration to a US AQI value.
// Standard EPA breakpoints.
func aqiFromPM25(pm float64) int {
	type bp struct {
		cLo, cHi, iLo, iHi float64
	}
	bps := []bp{
		{0.0, 12.0, 0, 50},
		{12.1, 35.4, 51, 100},
		{35.5, 55.4, 101, 150},
		{55.5, 150.4, 151, 200},
		{150.5, 250.4, 201, 300},
		{250.5, 350.4, 301, 400},
		{350.5, 500.4, 401, 500},
	}
	for _, b := range bps {
		if pm >= b.cLo && pm <= b.cHi {
			return int(math.Round(((b.iHi-b.iLo)/(b.cHi-b.cLo))*(pm-b.cLo) + b.iLo))
		}
	}
	if pm > 500.4 {
		return 500
	}
	return 0
}
