package main

import (
	"database/sql"
	"log"
)

func getRecord(db *sql.DB, sqlStatement string) Record {

	rows := db.QueryRow(sqlStatement)
	r := Record{}
	err := rows.Scan(&r.ID, &r.Mac, &r.Recorded, &r.Baromabsin, &r.Baromrelin, &r.Battout, &r.Batt1, &r.Batt2, &r.Batt3, &r.Batt4, &r.Batt5, &r.Batt6, &r.Batt7, &r.Batt8, &r.Batt9, &r.Batt10, &r.Co2, &r.Dailyrainin, &r.Dewpoint, &r.Eventrainin, &r.Feelslike, &r.Hourlyrainin, &r.Hourlyrain, &r.Humidity, &r.Humidity1, &r.Humidity2, &r.Humidity3, &r.Humidity4, &r.Humidity5, &r.Humidity6, &r.Humidity7, &r.Humidity8, &r.Humidity9, &r.Humidity10, &r.Humidityin, &r.Lastrain, &r.Maxdailygust, &r.Relay1, &r.Relay2, &r.Relay3, &r.Relay4, &r.Relay5, &r.Relay6, &r.Relay7, &r.Relay8, &r.Relay9, &r.Relay10, &r.Monthlyrainin, &r.Solarradiation, &r.Tempf, &r.Temp1f, &r.Temp2f, &r.Temp3f, &r.Temp4f, &r.Temp5f, &r.Temp6f, &r.Temp7f, &r.Temp8f, &r.Temp9f, &r.Temp10f, &r.Tempinf, &r.Totalrainin, &r.Uv, &r.Weeklyrainin, &r.Winddir, &r.Windgustmph, &r.Windgustdir, &r.Windspeedmph, &r.Yearlyrainin, &r.Battlightning, &r.Lightningday, &r.Lightninghour, &r.Lightningtime, &r.Lightningdistance, &r.Aqipm25, &r.Aqipm2524h)
	sqlError(err, "Error Getting Record: ")

	return r
}
func getStats(db *sql.DB) []Stat {
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

func chartQueries(t string) string {
	q := map[string]string{}

	q["1h"] = `SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 2) * interval '10 minutes') AS ts,
					ROUND(max(tempf)::numeric,2) as max, ROUND(min(tempf)::numeric,2) as min
				FROM records
					WHERE recorded >= NOW() - interval '1 hour'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`
	q["6h"] = `SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 15) * interval '15 minutes') AS ts,
					ROUND(max(tempf)::numeric,2) as max, ROUND(min(tempf)::numeric,2) as min
				FROM records
					WHERE recorded >= NOW() - interval '6 hour'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`

	q["12h"] = `SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 30) * interval '30 minute') AS ts,
					ROUND(max(tempf)::numeric,2) as max, ROUND(min(tempf)::numeric,2) as min
				FROM records
					WHERE recorded >= NOW() - interval '12 hours'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`

	q["1d"] = `SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 30) * interval '1 hour') AS ts,
					ROUND(max(tempf)::numeric,2) as max, ROUND(min(tempf)::numeric,2) as min
				FROM records
					WHERE recorded >= NOW() - interval '1 day'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`
	q["1m"] = `SELECT date_trunc('day', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 60) * interval '1 day') AS ts,
					ROUND(max(tempf)::numeric,2) as max, ROUND(min(tempf)::numeric,2) as min
				FROM records
					WHERE recorded >= NOW() - interval '1 month'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`
	q["1y"] = `SELECT date_trunc('month', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 60) * interval '1 day') AS ts,
					ROUND(max(tempf)::numeric,2) as max, ROUND(min(tempf)::numeric,2) as min
				FROM records
					WHERE recorded >= NOW() - interval '1 year'
					AND recorded <= NOW()
				GROUP BY ts
				order by ts asc`
	q["at"] = `SELECT date_trunc('month', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 60) * interval '1 day') AS ts,
					ROUND(max(tempf)::numeric,2) as max, ROUND(min(tempf)::numeric,2) as min
				FROM records
				GROUP BY ts
				order by ts asc`

	return q[t]
}
