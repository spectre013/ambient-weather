package main

// SQL-injection defense: column and chart-key allowlists.
//
// `cleanString` strips non-alphanumerics but is NOT a substitute for column
// validation -- "idORrecorded" survives cleaning and would let a caller
// pivot to arbitrary columns when interpolated into `ORDER BY`.
// We use strict allowlists for any identifier that must be interpolated
// into SQL (Postgres cannot parameterize column names).

// allowedAlltimeColumns lists the columns /api/alltime is permitted to query.
var allowedAlltimeColumns = map[string]bool{
	"tempf":          true,
	"tempinf":        true,
	"temp1f":         true,
	"temp2f":         true,
	"temp3f":         true,
	"temp4f":         true,
	"humidity":       true,
	"humidityin":     true,
	"baromrelin":     true,
	"baromabsin":     true,
	"dewpoint":       true,
	"feelslike":      true,
	"windspeedmph":   true,
	"windgustmph":    true,
	"dailyrainin":    true,
	"hourlyrainin":   true,
	"monthlyrainin":  true,
	"yearlyrainin":   true,
	"weeklyrainin":   true,
	"totalrainin":    true,
	"solarradiation": true,
	"uv":             true,
	"lightningday":   true,
	"aqipm25":        true,
	"aqipm2524h":     true,
	"co2":            true,
}

// allowedAlltimeCalcs lists the aggregation directions /api/alltime supports.
var allowedAlltimeCalcs = map[string]bool{
	"max": true,
	"min": true,
}

// allowedChartSensors lists columns the chart query is permitted to read.
var allowedChartSensors = map[string]bool{
	"tempf":          true,
	"dewpoint":       true,
	"humidity":       true,
	"windspeedmph":   true,
	"windgustmph":    true,
	"baromrelin":     true,
	"lightningday":   true,
	"dailyrainin":    true,
	"solarradiation": true,
	"uv":             true,
}

// allowedTrendColumns lists columns trend() may read.
var allowedTrendColumns = map[string]bool{
	"tempf":      true,
	"baromrelin": true,
	"humidity":   true,
}

// allowedChartTimeframes maps the URL timeframe slug to a Postgres interval
// literal. The interval is selected from this fixed table, never built from
// user input.
var allowedChartTimeframes = map[string]string{
	"1h":  "1 hour",
	"3h":  "3 hours",
	"6h":  "6 hours",
	"12h": "12 hours",
	"24h": "24 hours",
	"1m":  "1 month",
}
