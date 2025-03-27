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

	q := map[string]string{}
	q["1h"] = `SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst'))) * interval '1 minutes') AS ts,
					%s as value
				FROM records
					WHERE recorded >= NOW() - interval '1 hour'
					AND recorded <= NOW()
				order by ts asc`
	q["6h"] = `SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 15) * interval '15 minutes') AS ts,
					%s as value
				FROM records
					WHERE recorded >= NOW() - interval '6 hour'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`

	q["12h"] = `SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 30) * interval '30 minute') AS ts,
					%s as value
				FROM records
					WHERE recorded >= NOW() - interval '12 hours'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`

	q["1d"] = `SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 30) * interval '1 hour') AS ts,
					%s as value
				FROM records
					WHERE recorded >= NOW() - interval '1 day'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`
	q["1m"] = `SELECT date_trunc('day', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 60) * interval '1 day') AS ts,
					%s as value
				FROM records
					WHERE recorded >= NOW() - interval '1 month'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`
	q["1y"] = `SELECT date_trunc('month', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 60) * interval '1 day') AS ts,
					%s as value
				FROM records
					WHERE recorded >= NOW() - interval '1 year'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`
	q["at"] = `SELECT date_trunc('month', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 60) * interval '1 day') AS ts,
					%s as value
				FROM records
				GROUP BY ts
				order by ts asc`

	return fmt.Sprintf(q[t], sensor)
}
