package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

func getCurrent() ([]byte, error) {
	query := `select id,mac,recorded,baromabsin,baromrelin,battout,Batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,
				hourlyrainin,hourlyrain,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,
				maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,
				tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,battlightning,lightningday,lightninghour,lightningtime,lightningdistance,  aqipm25, aqipm2524h
				from records order by recorded desc limit 1`
	rec := getRecord(query)

	hourlyRain := getHourlyRain()
	rec.Hourlyrain = hourlyRain
	b, err := json.Marshal(rec)
	if err != nil {
		log.Println(err)
		return []byte(""), err
	}
	return b, nil
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

func calculateStats() {
	type Query struct {
		Query  string
		Params []time.Time
	}

	types := []string{"MAX", "MIN", "AVG"}
	periods := []string{"day", "yesterday", "month", "year"}
	fields := []string{"tempf", "tempinf", "temp1f", "temp2f", "temp3f", "baromrelin", "uv", "humidity", "windspeedmph", "windgustmph", "dewpoint", "humidityin", "humidity1", "humidity2", "humidity3", "dailyrainin", "lightning", "aqipm25"}
	queries := make(map[string]Query)
	for _, p := range periods {
		for _, t := range types {
			for _, f := range fields {
				key := p + "_" + strings.ToLower(t) + "_" + f
				order := ""
				if t == "MAX" {
					order = " DESC"
				}
				if (f == "dailyrainin" && t == "MAX") || f != "dailyrainin" {
					if !strings.Contains(f, "lightning") {

						q := fmt.Sprintf("select '%s' as id, CAST(COALESCE(%s,0.0) AS decimal(10,2)) as value, recorded from records where recorded between ? and ? order by %s%s limit 1", key, f, f, order)
						if t == "AVG" {
							q = fmt.Sprintf("select '%s' as id,  CAST(COALESCE(%s(%s),0.0) AS decimal(10,2)) as value from records where recorded between ? and ? limit 1", key, t, f)
						}
						queries[key] = Query{
							Query:  q,
							Params: getTimeframe(p),
						}
					}
				}

				if strings.Contains(f, "lightning") && t == "MAX" {
					q := ""
					if p == "month" || p == "year" {
						q = fmt.Sprintf(`
						SELECT '%s' as id, SUM(A.value) as value
						FROM (SELECT TO_CHAR(recorded,'YYY-MM-DD') as ldate, CAST(COALESCE(MAX(lightningday),0.0) AS decimal(10,2)) as value FROM records where recorded between ? and ? GROUP BY ldate) A
						`, key)
					} else {
						q = fmt.Sprintf(`SELECT '%s' as id, CAST(COALESCE(lightningday,0.0) AS decimal(10,2)) as value, recorded FROM records where recorded between ? and ? order by value desc limit 1`, key)
					}
					queries[key] = Query{
						Query:  q,
						Params: getTimeframe(p),
					}
				}

			}
		}
	}
	start := time.Now()
	tx, err := db.Begin()
	if err != nil {
		logger.Error(err)
	}
	for k, v := range queries {
		v.Query = strings.Replace(v.Query, "?", "'%s'", -1)
		v.Query = fmt.Sprintf(v.Query, formatDate(v.Params[0]), formatDate(v.Params[1]))
		d := "recorded = src.recorded,"
		if strings.Contains(k, "avg") || (strings.Contains(v.Query, "MAX(lightningday)")) {
			d = ""
		}

		update := checkStat(k)
		if update {
			updateQuery := fmt.Sprintf(`
				UPDATE stats set
				%s
				value = src.value
				from (
					%s
					) AS src
				WHERE
					stats.id = '%s';
			`, d, v.Query, k)
			//logger.Debug(updateQuery)
			_, err := tx.Exec(updateQuery)
			if err != nil {
				logger.Debug(updateQuery)
				logger.Error(err)
				break
			}
		} else {
			insert := "insert into stats (id,value,recorded)"
			if strings.Contains(k, "avg") || (strings.Contains(v.Query, "MAX(lightningday)")) {
				insert = "insert into stats (id,value)"
			}
			updateQuery := fmt.Sprintf(`
				%s
				%s
			`, insert, v.Query)
			//logger.Debug(updateQuery)
			_, err := tx.Exec(updateQuery)
			if err != nil {
				logger.Debug(updateQuery)
				logger.Error(err)
				break
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		logger.Error(err)
	}

	elapsed := time.Since(start)
	log.Printf("Update took %s", elapsed)

}

func checkStat(id string) bool {
	rows := db.QueryRow(fmt.Sprintf("select id from stats where id = '%s'", id))
	var r string
	err := rows.Scan(&r)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		} else {
			logger.Error("Scan:", err)
		}
	}
	return true
}

func insertRecord(r Record) bool {
	query := fmt.Sprintf(`insert into records (id,mac,recorded,baromabsin,baromrelin,battout,batt1,Batt2,Batt3,Batt4,Batt5,Batt6,Batt7,Batt8,Batt9,Batt10,co2,dailyrainin,dewpoint,eventrainin,feelslike,hourlyrainin,humidity,humidity1,humidity2,humidity3,humidity4,humidity5,humidity6,humidity7,humidity8,humidity9,humidity10,humidityin,lastrain,maxdailygust,relay1,relay2,relay3,relay4,relay5,relay6,relay7,relay8,relay9,relay10,monthlyrainin,solarradiation,tempf,temp1f,temp2f,temp3f,temp4f,temp5f,temp6f,temp7f,temp8f,temp9f,temp10f,tempinf,totalrainin,uv,weeklyrainin,winddir,windgustmph,windgustdir,windspeedmph,yearlyrainin,hourlyrain,lightningday,lightninghour,lightningtime,lightningdistance,battlightning, aqipm25, aqipm2524h) values
			(DEFAULT,'%s','%s',%f,%f,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%f,%f,%f,%f,%f,%f,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,'%s',%f,%d,%d,%d,%d,%d,%d,%d,%d,%d,%d,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%f,%d,%f,%d,%f,%f,%f,%d,%d,'%s',%f,%d,%d,%d)				
	`, r.Mac, formatDate(r.Recorded), r.Baromabsin, r.Baromrelin, r.Battout, r.Batt1, r.Batt2, r.Batt3, r.Batt4, r.Batt5, r.Batt6, r.Batt7, r.Batt8, r.Batt9, r.Batt10, r.Co2, r.Dailyrainin, r.Dewpoint, r.Eventrainin, r.Feelslike, r.Hourlyrainin, r.Humidity, r.Humidity1, r.Humidity2, r.Humidity3, r.Humidity4, r.Humidity5, r.Humidity6, r.Humidity7, r.Humidity8, r.Humidity9, r.Humidity10, r.Humidityin, formatDate(r.Lastrain), r.Maxdailygust, r.Relay1, r.Relay2, r.Relay3, r.Relay4, r.Relay5, r.Relay6, r.Relay7, r.Relay8, r.Relay9, r.Relay10, r.Monthlyrainin, r.Solarradiation, r.Tempf, r.Temp1f, r.Temp2f, r.Temp3f, r.Temp4f, r.Temp5f, r.Temp6f, r.Temp7f, r.Temp8f, r.Temp9f, r.Temp10f, r.Tempinf, r.Totalrainin, r.Uv, r.Weeklyrainin, r.Winddir, r.Windgustmph, r.Windgustdir, r.Windspeedmph, r.Yearlyrainin, r.Hourlyrain, r.Lightningday, r.Lightninghour, formatDate(r.Lightningtime), r.Lightningdistance, r.Battlightning, r.Aqipm25, r.Aqipm2524h)
	logger.Debug(query)
	_, err := db.Exec(query)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
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
