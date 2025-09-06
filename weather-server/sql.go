package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

func getCurrent() Record {
	query := `select id,mac,recorded,baromabsin,baromrelin,battout,Batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,
				maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,battlightning,lightningday,lightninghour,lightningtime,lightningdistance,  aqipm25, aqipm2524h
				from records order by recorded desc limit 1`
	rec := getRecord(query)

	lightningMonth := lightningMonth()
	hourlyRain := getHourlyRain()
	rec.Hourlyrain = hourlyRain
	rec.LightningMonth = lightningMonth
	return rec
}

func getHourlyRain() float64 {
	start := time.Now()
	end := start.Add(-60 * time.Minute)
	var maxrain float64
	query := fmt.Sprintf("select dailyrainin from records where recorded BETWEEN '%s' AND '%s' order by dailyrainin desc limit 1", formatDate(end), formatDate(start))
	logger.Debug(query)
	crows := db.QueryRow(query)
	err := crows.Scan(&maxrain)
	if err != nil {
		if err == sql.ErrNoRows {
			//logger.Error("Zero Rows Found", query)
		} else {
			logger.Error("Scan:", err)
		}
	}
	return maxrain
}

func getRecord(sqlStatement string) Record {

	rows := db.QueryRow(sqlStatement)
	r := Record{}
	err := rows.Scan(&r.ID, &r.Mac, &r.Recorded, &r.Baromabsin, &r.Baromrelin, &r.Battout, &r.Batt1, &r.Batt2, &r.Batt3, &r.Batt4, &r.Batt5, &r.Batt6, &r.Batt7, &r.Batt8, &r.Batt9, &r.Batt10, &r.Co2, &r.Dailyrainin, &r.Dewpoint, &r.Eventrainin, &r.Feelslike, &r.Hourlyrainin, &r.Hourlyrain, &r.Humidity, &r.Humidity1, &r.Humidity2, &r.Humidity3, &r.Humidity4, &r.Humidity5, &r.Humidity6, &r.Humidity7, &r.Humidity8, &r.Humidity9, &r.Humidity10, &r.Humidityin, &r.Lastrain, &r.Maxdailygust, &r.Relay1, &r.Relay2, &r.Relay3, &r.Relay4, &r.Relay5, &r.Relay6, &r.Relay7, &r.Relay8, &r.Relay9, &r.Relay10, &r.Monthlyrainin, &r.Solarradiation, &r.Tempf, &r.Temp1f, &r.Temp2f, &r.Temp3f, &r.Temp4f, &r.Temp5f, &r.Temp6f, &r.Temp7f, &r.Temp8f, &r.Temp9f, &r.Temp10f, &r.Tempinf, &r.Totalrainin, &r.Uv, &r.Weeklyrainin, &r.Winddir, &r.Windgustmph, &r.Windgustdir, &r.Windspeedmph, &r.Yearlyrainin, &r.Battlightning, &r.Lightningday, &r.Lightninghour, &r.Lightningtime, &r.Lightningdistance, &r.Aqipm25, &r.Aqipm2524h)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.Error("Zero Rows Found ", sqlStatement)
		} else {
			logger.Error("Scan:", err)
		}
	}

	return r
}
func getStats() []Stat {
	stats := make([]Stat, 0)
	rows, err := db.Query("Select * from stats")
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		r := Stat{}
		err = rows.Scan(&r.ID, &r.Recorded, &r.Value)
		if err != nil {
			logger.Error("Scan:", err)
		}
		stats = append(stats, r)
	}
	err = rows.Err()
	if err != nil {
		logger.Error(err)
	}
	return stats
}

func chartQueries(t string, sensor string) string {

	sqlString := `SELECT recorded AS ts, %s as value
	FROM records
	WHERE recorded >= NOW() - interval '%s'
	AND recorded <= NOW()
	order by ts asc`
	date := time.Now().Format("2006-01-02 15:04:05")
	rainsSQL := fmt.Sprintf(`SELECT EXTRACT(DAY FROM recorded) AS ts, max(%s) as value
				FROM records
				WHERE recorded >= DATE_TRUNC('month', '%s'::DATE)
				  AND recorded < (DATE_TRUNC('month', '%s'::DATE) + INTERVAL '1 month')
				  group by ts
				  order by ts;`, sensor, date, date)

	if sensor == "dailyrainin" || sensor == "lightningday" {
		return rainsSQL
	}

	timeFrame := "1 hour"
	switch t {
	case "1h":
		timeFrame = "1 hour"
		break
	case "3h":
		timeFrame = "3 hours"
		break
	case "6h":
		timeFrame = "6 hours"
		break
	case "12h":
		timeFrame = "12 hours"
		break
	case "24h":
		timeFrame = "24 hours"
		break
	case "1m":
		timeFrame = "1 month"
	}

	return fmt.Sprintf(sqlString, sensor, timeFrame)
}

func almanacQueries(data string) []ClimateRaw {
	avgs := `SELECT
		EXTRACT(YEAR FROM recorded) AS year,
		EXTRACT(MONTH FROM recorded) AS month,
		ROUND(max(monthlyrainin)::numeric, 2) AS avg_rain,
		ROUND(AVG(tempf)::numeric, 2) AS avg_temp,
		ROUND(MAX(tempf)::numeric, 2) AS max_temp,
		ROUND(MIN(tempf)::numeric, 2) AS min_temp
	FROM
		records
	GROUP BY
		year,
		month
	ORDER BY
		year,
		month`

	query := avgs
	rows, err := db.Query(query)
	logger.Println(err, query)
	if err != nil {
		log.Println(err)
	}
	climate := make([]ClimateRaw, 0)
	for rows.Next() {
		r := ClimateRaw{}
		err = rows.Scan(&r.Year, &r.Month, &r.AvgRain, &r.AvgTemp, &r.MaxTemp, &r.MinTemp)
		if err != nil {
			logger.Error("Climate Scan:", err)
		}
		climate = append(climate, r)
	}
	err = rows.Err()
	if err != nil {
		logger.Error(err)
	}

	return climate
}

func firstFreeze() []FirstFreeze {
	query := `WITH freezes_by_year AS (
			SELECT
				EXTRACT(YEAR FROM recorded) AS year,
				recorded
			FROM
				records
			WHERE
				tempf <= 32
		),
		spring_freezes AS (
			SELECT
				year,
				MAX(recorded) AS last_spring_freeze
			FROM
				freezes_by_year
			WHERE
				EXTRACT(MONTH FROM recorded) BETWEEN 1 AND 6 -- January to June
			GROUP BY
				year
		),
		fall_freezes AS (
			SELECT
				year,
				MIN(recorded) AS first_fall_freeze
			FROM
				freezes_by_year
			WHERE
				EXTRACT(MONTH FROM recorded) BETWEEN 7 AND 12 -- July to December
			GROUP BY
				year
		)
		SELECT
			COALESCE(sf.year,ff.year) AS year,
			COALESCE(sf.last_spring_freeze,'2000-01-01') AS last_spring_freeze,
			COALESCE(ff.first_fall_freeze,'2000-01-01') AS first_fall_freeze
		FROM
			spring_freezes sf
		FULL OUTER JOIN
			fall_freezes ff ON sf.year = ff.year
		ORDER BY
			year;`

	ff := make([]FirstFreeze, 0)

	rows, err := db.Query(query)
	logger.Println(err, query)
	if err != nil {
		log.Println(err)
	}
	for rows.Next() {
		r := FirstFreeze{}
		err = rows.Scan(&r.Year, &r.Spring, &r.Fall)
		if err != nil {
			logger.Error("Climate Scan:", err)
		}
		ff = append(ff, r)
	}
	return ff
}
