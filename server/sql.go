package main

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"time"
)

const sqlDateFmt = "YYYY-MM-DD"

func getCurrent() Record {
	const query = `select id,mac,recorded,baromabsin,baromrelin,battout,Batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,
				maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,battlightning,lightningday,lightninghour,lightningtime,lightningdistance, aqipm25, aqipm2524h
				from records order by recorded desc limit 1`
	rec := getRecord(query)

	rec.Hourlyrain = getHourlyRain()
	rec.LightningMonth = lightningMonth()
	return rec
}

func getHourlyRain() float64 {
	end := time.Now()
	start := end.Add(-60 * time.Minute)
	const query = `select dailyrainin from records where recorded BETWEEN $1 AND $2 order by dailyrainin desc limit 1`
	logger.Debug(query)

	var maxrain float64
	err := db.QueryRow(query, start, end).Scan(&maxrain)
	if err != nil && err != sql.ErrNoRows {
		logger.WithError(err).Error("getHourlyRain scan")
	}
	return maxrain
}

func getRecord(sqlStatement string) Record {
	r := Record{}
	err := db.QueryRow(sqlStatement).Scan(&r.ID, &r.Mac, &r.Recorded, &r.Baromabsin, &r.Baromrelin, &r.Battout, &r.Batt1, &r.Batt2, &r.Batt3, &r.Batt4, &r.Batt5, &r.Batt6, &r.Batt7, &r.Batt8, &r.Batt9, &r.Batt10, &r.Co2, &r.Dailyrainin, &r.Dewpoint, &r.Eventrainin, &r.Feelslike, &r.Hourlyrainin, &r.Hourlyrain, &r.Humidity, &r.Humidity1, &r.Humidity2, &r.Humidity3, &r.Humidity4, &r.Humidity5, &r.Humidity6, &r.Humidity7, &r.Humidity8, &r.Humidity9, &r.Humidity10, &r.Humidityin, &r.Lastrain, &r.Maxdailygust, &r.Relay1, &r.Relay2, &r.Relay3, &r.Relay4, &r.Relay5, &r.Relay6, &r.Relay7, &r.Relay8, &r.Relay9, &r.Relay10, &r.Monthlyrainin, &r.Solarradiation, &r.Tempf, &r.Temp1f, &r.Temp2f, &r.Temp3f, &r.Temp4f, &r.Temp5f, &r.Temp6f, &r.Temp7f, &r.Temp8f, &r.Temp9f, &r.Temp10f, &r.Tempinf, &r.Totalrainin, &r.Uv, &r.Weeklyrainin, &r.Winddir, &r.Windgustmph, &r.Windgustdir, &r.Windspeedmph, &r.Yearlyrainin, &r.Battlightning, &r.Lightningday, &r.Lightninghour, &r.Lightningtime, &r.Lightningdistance, &r.Aqipm25, &r.Aqipm2524h)
	sqlError("getRecord", err, sqlStatement)
	return r
}

func GetForecasts() ([]ForecastDB, error) {
	start := time.Now().Format("2006-01-02")
	const query = `
		SELECT
			datetime, datetime_epoch, tempmax, tempmin, temp, feelslikemax,
			feelslikemin, feelslike, dew, humidity, precip, precipprob,
			precipcover, preciptype, snow, snowdepth, windgust, windspeed,
			winddir, pressure, cloudcover, visibility, solarradiation,
			solarenergy, uvindex, severerisk, sunrise, sunrise_epoch,
			sunset, sunset_epoch, moonphase, conditions, description,
			icon, stations, source, hours, summary
		FROM forecast
		WHERE datetime >= $1
		ORDER BY datetime ASC LIMIT 10`

	rows, err := db.Query(query, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forecasts []ForecastDB
	for rows.Next() {
		var f ForecastDB
		err := rows.Scan(
			&f.Datetime, &f.DatetimeEpoch, &f.TempMax, &f.TempMin, &f.Temp, &f.FeelsLikeMax,
			&f.FeelsLikeMin, &f.FeelsLike, &f.Dew, &f.Humidity, &f.Precip, &f.PrecipProb,
			&f.PrecipCover, &f.PrecipType, &f.Snow, &f.SnowDepth, &f.WindGust, &f.WindSpeed,
			&f.WindDir, &f.Pressure, &f.CloudCover, &f.Visibility, &f.SolarRadiation,
			&f.SolarEnergy, &f.UVIndex, &f.SevereRisk, &f.Sunrise, &f.SunriseEpoch,
			&f.Sunset, &f.SunsetEpoch, &f.MoonPhase, &f.Conditions, &f.Description,
			&f.Icon, &f.Stations, &f.Source, &f.Hours, &f.Summary,
		)
		if err != nil {
			return nil, err
		}
		forecasts = append(forecasts, f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return forecasts, nil
}

func GetForecastByTimestamp(ts time.Time) (*ForecastDB, error) {
	f := &ForecastDB{}
	const query = `
		SELECT
			datetime, datetime_epoch, tempmax, tempmin, temp, feelslikemax,
			feelslikemin, feelslike, dew, humidity, precip, precipprob,
			precipcover, preciptype, snow, snowdepth, windgust, windspeed,
			winddir, pressure, cloudcover, visibility, solarradiation,
			solarenergy, uvindex, severerisk, sunrise, sunrise_epoch,
			sunset, sunset_epoch, moonphase, conditions, description,
			icon, stations, source, hours, summary
		FROM forecast
		WHERE datetime = $1`

	err := db.QueryRow(query, ts).Scan(
		&f.Datetime, &f.DatetimeEpoch, &f.TempMax, &f.TempMin, &f.Temp, &f.FeelsLikeMax,
		&f.FeelsLikeMin, &f.FeelsLike, &f.Dew, &f.Humidity, &f.Precip, &f.PrecipProb,
		&f.PrecipCover, &f.PrecipType, &f.Snow, &f.SnowDepth, &f.WindGust, &f.WindSpeed,
		&f.WindDir, &f.Pressure, &f.CloudCover, &f.Visibility, &f.SolarRadiation,
		&f.SolarEnergy, &f.UVIndex, &f.SevereRisk, &f.Sunrise, &f.SunriseEpoch,
		&f.Sunset, &f.SunsetEpoch, &f.MoonPhase, &f.Conditions, &f.Description,
		&f.Icon, &f.Stations, &f.Source, &f.Hours, &f.Summary,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return f, nil
}

// GetLatestConditions returns the most recent observed conditions row for the
// given station, or (nil, nil) when none exist yet.
func GetLatestConditions(station string) (*StationConditions, error) {
	const query = `
		SELECT conditions, icon, observed_at, temperature, humidity
		FROM conditions
		WHERE station = $1
		ORDER BY observed_at DESC
		LIMIT 1`

	var (
		c           StationConditions
		conditions  sql.NullString
		icon        sql.NullString
		temperature sql.NullFloat64
		humidity    sql.NullFloat64
	)
	err := db.QueryRow(query, station).Scan(
		&conditions, &icon, &c.ObservedAt, &temperature, &humidity,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	c.Conditions = conditions.String
	c.Icon = icon.String
	c.Temperature = temperature.Float64
	c.Humidity = humidity.Float64
	return &c, nil
}

func chartQuery(timeframe, sensor string) (string, bool) {
	if !allowedChartSensors[sensor] {
		return "", false
	}

	if sensor == "dailyrainin" || sensor == "lightningday" {
		// Aggregate by day within the current month.
		q := fmt.Sprintf(`SELECT EXTRACT(DAY FROM recorded) AS ts, max(%s) as value
			FROM records
			WHERE recorded >= DATE_TRUNC('month', NOW())
			  AND recorded <  DATE_TRUNC('month', NOW()) + INTERVAL '1 month'
			GROUP BY ts
			ORDER BY ts`, sensor)
		return q, true
	}

	interval, ok := allowedChartTimeframes[timeframe]
	if !ok {
		interval = "1 hour"
	}

	q := fmt.Sprintf(`SELECT recorded AS ts, %s as value
		FROM records
		WHERE recorded >= NOW() - interval '%s'
		  AND recorded <= NOW()
		ORDER BY ts ASC`, sensor, interval)
	return q, true
}

func almanacQueries() []ClimateRaw {
	const query = `SELECT
		EXTRACT(YEAR FROM recorded) AS year,
		EXTRACT(MONTH FROM recorded) AS month,
		ROUND(max(monthlyrainin)::numeric, 2) AS avg_rain,
		ROUND(AVG(tempf)::numeric, 2) AS avg_temp,
		ROUND(MAX(tempf)::numeric, 2) AS max_temp,
		ROUND(MIN(tempf)::numeric, 2) AS min_temp
	FROM records
	GROUP BY year, month
	ORDER BY year, month`

	rows, err := db.Query(query)
	if err != nil {
		logger.WithError(err).Error("almanacQueries query")
		return nil
	}
	defer rows.Close()

	climate := make([]ClimateRaw, 0)
	for rows.Next() {
		r := ClimateRaw{}
		if err := rows.Scan(&r.Year, &r.Month, &r.AvgRain, &r.AvgTemp, &r.MaxTemp, &r.MinTemp); err != nil {
			logger.WithError(err).Error("almanacQueries scan")
			continue
		}
		climate = append(climate, r)
	}
	if err := rows.Err(); err != nil {
		logger.WithError(err).Error("almanacQueries rows.Err")
	}
	return climate
}

func firstFreeze() []FirstFreeze {
	const query = `WITH freezes_by_year AS (
			SELECT EXTRACT(YEAR FROM recorded) AS year, recorded
			FROM records WHERE tempf <= 32
		),
		spring_freezes AS (
			SELECT year, MAX(recorded) AS last_spring_freeze
			FROM freezes_by_year
			WHERE EXTRACT(MONTH FROM recorded) BETWEEN 1 AND 6
			GROUP BY year
		),
		fall_freezes AS (
			SELECT year, MIN(recorded) AS first_fall_freeze
			FROM freezes_by_year
			WHERE EXTRACT(MONTH FROM recorded) BETWEEN 7 AND 12
			GROUP BY year
		)
		SELECT
			COALESCE(sf.year, ff.year) AS year,
			COALESCE(sf.last_spring_freeze, '2000-01-01') AS last_spring_freeze,
			COALESCE(ff.first_fall_freeze, '2000-01-01') AS first_fall_freeze
		FROM spring_freezes sf
		FULL OUTER JOIN fall_freezes ff ON sf.year = ff.year
		ORDER BY year`

	rows, err := db.Query(query)
	if err != nil {
		logger.WithError(err).Error("firstFreeze query")
		return nil
	}
	defer rows.Close()

	ff := make([]FirstFreeze, 0)
	for rows.Next() {
		r := FirstFreeze{}
		if err := rows.Scan(&r.Year, &r.Spring, &r.Fall); err != nil {
			logger.WithError(err).Error("firstFreeze scan")
			continue
		}
		ff = append(ff, r)
	}
	if err := rows.Err(); err != nil {
		logger.WithError(err).Error("firstFreeze rows.Err")
	}
	return ff
}

func getHistory(ctx context.Context) (HistoryResponse, error) {
	loc := config.Location
	now := time.Now().In(loc)

	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	// exclusive upper bound: the start of the current (in-progress) hour
	currentHourStart := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)

	_, offsetSeconds := now.Zone()

	resp := HistoryResponse{
		Date:     dayStart.Format("2006-01-02"),
		TZOffset: offsetSeconds / 3600,
		Hours:    []HistoryHour{},
	}

	// For each local-time hour bucket:
	//   - temp/humidity/windspeed come from the record nearest to HH:00
	//     (DISTINCT ON the hour, ordered by distance from the hour boundary).
	//   - precip is MAX(dailyrainin) over the hour.
	// The two are computed separately and joined on the hour bucket.
	const query = `
			WITH nearest AS (
				SELECT DISTINCT ON (date_trunc('hour', recorded AT TIME ZONE $1))
					date_trunc('hour', recorded AT TIME ZONE $1) AS hour_bucket,
					tempf        AS temp,
					humidity     AS humidity,
					windspeedmph AS windspeed
				FROM records
				WHERE recorded >= $2 AND recorded < $3
				ORDER BY
					date_trunc('hour', recorded AT TIME ZONE $1),
					ABS(EXTRACT(EPOCH FROM (
						(recorded AT TIME ZONE $1)
						- date_trunc('hour', recorded AT TIME ZONE $1)
					)))
			),
			rain AS (
				SELECT date_trunc('hour', recorded AT TIME ZONE $1) AS hour_bucket,
					   MAX(dailyrainin) AS precip
				FROM records
				WHERE recorded >= $2 AND recorded < $3
				GROUP BY date_trunc('hour', recorded AT TIME ZONE $1)
			)
			SELECT
				to_char(n.hour_bucket, 'HH24:00:00') AS hour_label,
				n.temp,
				n.humidity,
				n.windspeed,
				COALESCE(r.precip, 0) AS precip
			FROM nearest n
			LEFT JOIN rain r USING (hour_bucket)
			ORDER BY n.hour_bucket`

	rows, err := db.QueryContext(ctx, query, config.LocationStr, dayStart, currentHourStart)
	if err != nil {
		return resp, fmt.Errorf("history query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var h HistoryHour
		var temp, humidity, windspeed, precip sql.NullFloat64
		var label string
		if err := rows.Scan(&label, &temp, &humidity, &windspeed, &precip); err != nil {
			return resp, fmt.Errorf("history scan: %w", err)
		}
		h.Datetime = label
		h.Temp = round2(temp.Float64)
		h.Humidity = round2(humidity.Float64)
		h.Windspeed = round2(windspeed.Float64)
		h.Precip = round2(precip.Float64)
		h.Source = "obs"
		resp.Hours = append(resp.Hours, h)
	}
	if err := rows.Err(); err != nil {
		return resp, fmt.Errorf("history rows: %w", err)
	}

	// latest raw record (full precision timestamp, UTC)
	latest, err := getLatestPoint(ctx)
	if err != nil {
		// non-fatal: return hours without latest
		logger.Error("history latest: ", err)
	} else {
		resp.Latest = latest
	}

	return resp, nil
}

func getLatestPoint(ctx context.Context) (*LatestPoint, error) {
	const q = `
		SELECT recorded, tempf
		FROM records
		ORDER BY recorded DESC
		LIMIT 1`
	var lp LatestPoint
	err := db.QueryRowContext(ctx, q).Scan(&lp.Datetime, &lp.Temp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	lp.Datetime = lp.Datetime.UTC()
	lp.Temp = round2(lp.Temp)
	return &lp, nil
}

func round2(f float64) float64 {
	return math.Round(f*100) / 100
}
