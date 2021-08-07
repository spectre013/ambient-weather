package main

import (
	"time"
)

//Stat Stat Table structure
type Stat struct {
	ID    string
	Recorded  time.Time
	Value float64
}

//Stats for building minmax values
type Stats struct {
	Day       StatValue `json:"day"`
	Yesterday StatValue `json:"yesterday"`
	Month     StatValue `json:"month"`
	Year      StatValue `json:"year"`
}

//StatValue date and value of minmax
type StatValue struct {
	Value float64   `json:"value"`
	Recorded  time.Time `json:"date"`
}

//Trend response
type Trend struct {
	Trend string  `json:"trend"`
	By    float64 `json:"by"`
}

//Location for Astro struct
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

//Astro Moon sun data
type Astro struct {
	Location             Location  `json:"location"`
	Date                 string    `json:"date"`
	Sunrise              string    `json:"sunrise"`
	Sunset               string    `json:"sunset"`
	SolarNoon            string    `json:"solar_noon"`
	DayLength            string    `json:"day_length"`
	SunAltitude          float64   `json:"sun_altitude"`
	SunDistance          float64   `json:"sun_distance"`
	SunAzimuth           float64   `json:"sun_azimuth"`
	Moonrise             string    `json:"moonrise"`
	Moonset              string    `json:"moonset"`
	MoonAltitude         float64   `json:"moon_altitude"`
	MoonDistance         float64   `json:"moon_distance"`
	MoonAzimuth          float64   `json:"moon_azimuth"`
	MoonParallacticAngle float64   `json:"moon_parallactic_angle"`
	Tomorrow             Tomorrow  `json:"tomorrow"`
	Newmoon              time.Time `json:"newmoon"`
	Nextnewmoon          time.Time `json:"nextnewmoon"`
	Fullmoon             time.Time `json:"fullmoon"`
	Phase                string    `json:"phase"`
	Illuminated          float64   `json:"illuminated"`
	Age                  float64   `json:"age"`
}
type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

//Tomorrow copy of moon sun data for getting sun rise for tomorrow
type Tomorrow struct {
	Location             Location `json:"location"`
	Date                 string   `json:"date"`
	Sunrise              string   `json:"sunrise"`
	Sunset               string   `json:"sunset"`
	SolarNoon            string   `json:"solar_noon"`
	DayLength            string   `json:"day_length"`
	SunAltitude          float64  `json:"sun_altitude"`
	SunDistance          float64  `json:"sun_distance"`
	SunAzimuth           float64  `json:"sun_azimuth"`
	Moonrise             string   `json:"moonrise"`
	Moonset              string   `json:"moonset"`
	MoonAltitude         float64  `json:"moon_altitude"`
	MoonDistance         float64  `json:"moon_distance"`
	MoonAzimuth          float64  `json:"moon_azimuth"`
	MoonParallacticAngle float64  `json:"moon_parallactic_angle"`
}

//Record data for main database  table
type Record struct {
	ID                int       `json:"id" db:"id"`
	Mac               string    `json:"mac" db:"mac"`
	Recorded          time.Time `json:"date" db:"recorded"`
	Baromabsin        float64   `json:"baromabsin" db:"baromabsin"`
	Baromrelin        float64   `json:"baromrelin" db:"baromrelin"`
	Battout           int       `json:"battout" db:"battout"`
	Batt1             int       `json:"batt1" db:"batt1"`
	Batt2             int       `json:"batt2" db:"batt2"`
	Batt3             int       `json:"batt3" db:"batt3"`
	Batt4             int       `json:"batt4" db:"batt4"`
	Batt5             int       `json:"batt5" db:"batt5"`
	Batt6             int       `json:"batt6" db:"batt6"`
	Batt7             int       `json:"batt7" db:"batt7"`
	Batt8             int       `json:"batt8" db:"batt8"`
	Batt9             int       `json:"batt9" db:"batt9"`
	Batt10            int       `json:"batt10" db:"batt10"`
	Battlightning     int       `json:"battlightning" db:"battlightning"`
	Co2               float64   `json:"co2" db:"co2"`
	Dailyrainin       float64   `json:"dailyrainin" db:"dailyrainin"`
	Dewpoint          float64   `json:"dewpoint" db:"dewpoint"`
	Eventrainin       float64   `json:"eventrain" db:"eventrainin"`
	Feelslike         float64   `json:"feelslike" db:"feelslike"`
	Hourlyrainin      float64   `json:"hourlyrainin" db:"hourlyrainin"`
	Hourlyrain        float64   `json:"hourlyrain" db:"hourlyrain"`
	Humidity          int       `json:"humidity" db:"humidity"`
	Humidity1         int       `json:"humidity1" db:"humidity1"`
	Humidity2         int       `json:"humidity2" db:"humidity2"`
	Humidity3         int       `json:"humidity3" db:"humidity3"`
	Humidity4         int       `json:"humidity4" db:"humidity4"`
	Humidity5         int       `json:"humidity5" db:"humidity5"`
	Humidity6         int       `json:"humidity6" db:"humidity6"`
	Humidity7         int       `json:"humidity7" db:"humidity7"`
	Humidity8         int       `json:"humidity8" db:"humidity8"`
	Humidity9         int       `json:"humidity9" db:"humidity9"`
	Humidity10        int       `json:"humidity10" db:"humidity10"`
	Humidityin        int       `json:"humidityin" db:"humidityin"`
	Lastrain          time.Time `json:"lastrain" db:"lastrain"`
	Maxdailygust      float64   `json:"maxdailygust" db:"maxdailygust"`
	Relay1            int       `json:"relay1" db:"relay1"`
	Relay2            int       `json:"relay2" db:"relay2"`
	Relay3            int       `json:"relay3" db:"relay3"`
	Relay4            int       `json:"relay4" db:"relay4"`
	Relay5            int       `json:"relay5" db:"relay5"`
	Relay6            int       `json:"relay6" db:"relay6"`
	Relay7            int       `json:"relay7" db:"relay7"`
	Relay8            int       `json:"relay8" db:"relay8"`
	Relay9            int       `json:"relay9" db:"relay9"`
	Relay10           int       `json:"relay10" db:"relay10"`
	Monthlyrainin     float64   `json:"monthlyrainin" db:"monthlyrainin"`
	Solarradiation    float64   `json:"solarradiation" db:"solarradiation"`
	Tempf             float64   `json:"tempf" db:"tempf"`
	Temp1f            float64   `json:"temp1f" db:"temp1f"`
	Temp2f            float64   `json:"temp2f" db:"temp2f"`
	Temp3f            float64   `json:"temp3f" db:"temp3f"`
	Temp4f            float64   `json:"temp4f" db:"temp4f"`
	Temp5f            float64   `json:"temp5f" db:"temp5f"`
	Temp6f            float64   `json:"temp6f" db:"temp6f"`
	Temp7f            float64   `json:"temp7f" db:"temp7f"`
	Temp8f            float64   `json:"temp8f" db:"temp8f"`
	Temp9f            float64   `json:"temp9f" db:"temp9f"`
	Temp10f           float64   `json:"temp10f" db:"temp10f"`
	Tempinf           float64   `json:"tempinf" db:"tempinf"`
	Totalrainin       float64   `json:"totalrainin" db:"totalrainin"`
	Uv                float64   `json:"uv" db:"uv"`
	Weeklyrainin      float64   `json:"weeklyrainin" db:"weeklyrainin"`
	Winddir           int       `json:"winddir" db:"winddir"`
	Windgustmph       float64   `json:"windgustmph" db:"windgustmph"`
	Windgustdir       int       `json:"windgustdir" db:"windgustdir"`
	Windspeedmph      float64   `json:"windspeedmph" db:"windspeedmph"`
	Yearlyrainin      float64   `json:"yearlyrainin" db:"yearlyrainin"`
	Lightningday      int       `json:"lightningday" db:"lightiningday"`
	Lightninghour     int       `json:"lightninghour" db:"lightininghour"`
	Lightningdistance float64   `json:"lightningdistance" db:"lightningdistance"`
	Lightningtime     time.Time `json:"lightningtime" db:"lightningtime"`
}


