package main

import (
	"database/sql"
	"html/template"
	"time"
)

type Weather struct {
	DB *sql.DB
}

type BoxProps struct {
	Icon  string
	Title string
	Unit  string
	Style map[string]string
}

type TemplateData struct {
	Units     string
	Record    Record
	Minmax    map[string]map[string]map[string]StatValue
	Alerts    []Alert
	Forecast  ForecastImage
	Wind      map[string]StatValue
	Lightning LightningData
	Astro     Astro
	tTrend    Trend
	bTrend    Trend
}

// Stat Stat Table structure
type Stat struct {
	ID       string
	Recorded time.Time
	Value    float64
}

// StatValue date and value of minmax
type StatValue struct {
	Value    float64   `json:"value"`
	Recorded time.Time `json:"date"`
}

// Trend response
type Trend struct {
	Trend string  `json:"trend"`
	By    float64 `json:"by"`
}

// Location for Astro struct
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Astro Moon sun data
type Astro struct {
	Sunrise         time.Time
	Sunset          time.Time
	SunriseTomorrow time.Time
	SunsetTomorrow  time.Time
	Darkness        time.Duration
	Daylight        time.Duration
	Elevation       float64
	HasSunset       bool
}

// Tomorrow copy of moon sun data for getting sun rise for tomorrow
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

// Record data for main database  table
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
	Aqipm25           int       `json:"aqipm25" db:"aqi"`
	Aqipm2524h        int       `json:"aqipm2524h" db:"aqi24"`
}

type LightningData struct {
	LightningDay      int       `json:"lightningday"`
	LightningHour     int       `json:"lightninghour"`
	LightningDistance float64   `json:"lightningdistance"`
	LightningTime     time.Time `json:"lightningtime"`
	LightningMonth    int       `json:"lightningmonth"`
}

type Alert struct {
	ID          int       `json:"id"`
	Alertid     string    `json:"alertid"`
	Wxtype      string    `json:"wxtype"`
	Areadesc    string    `json:"areadesc"`
	Sent        time.Time `json:"sent"`
	Effective   time.Time `json:"effective"`
	Onset       time.Time `json:"onset"`
	Expires     time.Time `json:"expires"`
	Ends        time.Time `json:"end"`
	Status      string    `json:"status"`
	Messagetype string    `json:"messagetype"`
	Category    string    `json:"category"`
	Severity    string    `json:"severity"`
	Certainty   string    `json:"certainty"`
	Urgency     string    `json:"urgency"`
	Event       string    `json:"event"`
	Sender      string    `json:"sender"`
	SenderName  string    `json:"senderName"`
	Headline    string    `json:"headline"`
	Description string    `json:"description"`
	Instruction string    `json:"instruction"`
	Response    string    `json:"response"`
}

type ForecastImage struct {
	Querycost       int     `json:"queryCost"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Resolvedaddress string  `json:"resolvedAddress"`
	Address         string  `json:"address"`
	Timezone        string  `json:"timezone"`
	Tzoffset        float64 `json:"tzoffset"`
	Days            []Days  `json:"days"`
	Stations        struct {
		Kfcs struct {
			Distance     float64 `json:"distance"`
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			Usecount     int     `json:"useCount"`
			ID           string  `json:"id"`
			Name         string  `json:"name"`
			Quality      int     `json:"quality"`
			Contribution float64 `json:"contribution"`
		} `json:"KFCS"`
		Kcwn struct {
			Distance     float64 `json:"distance"`
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			Usecount     int     `json:"useCount"`
			ID           string  `json:"id"`
			Name         string  `json:"name"`
			Quality      int     `json:"quality"`
			Contribution float64 `json:"contribution"`
		} `json:"KCWN"`
		Kaff struct {
			Distance     float64 `json:"distance"`
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			Usecount     int     `json:"useCount"`
			ID           string  `json:"id"`
			Name         string  `json:"name"`
			Quality      int     `json:"quality"`
			Contribution float64 `json:"contribution"`
		} `json:"KAFF"`
		C8796 struct {
			Distance     float64 `json:"distance"`
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			Usecount     int     `json:"useCount"`
			ID           string  `json:"id"`
			Name         string  `json:"name"`
			Quality      int     `json:"quality"`
			Contribution float64 `json:"contribution"`
		} `json:"C8796"`
		Kcos struct {
			Distance     float64 `json:"distance"`
			Latitude     float64 `json:"latitude"`
			Longitude    float64 `json:"longitude"`
			Usecount     int     `json:"useCount"`
			ID           string  `json:"id"`
			Name         string  `json:"name"`
			Quality      int     `json:"quality"`
			Contribution float64 `json:"contribution"`
		} `json:"KCOS"`
	} `json:"stations"`
}

type Days struct {
	Datetime       string   `json:"datetime"`
	Datetimeepoch  int      `json:"datetimeEpoch"`
	Tempmax        float64  `json:"tempmax"`
	Tempmin        float64  `json:"tempmin"`
	Temp           float64  `json:"temp"`
	Feelslikemax   float64  `json:"feelslikemax"`
	Feelslikemin   float64  `json:"feelslikemin"`
	Feelslike      float64  `json:"feelslike"`
	Dew            float64  `json:"dew"`
	Humidity       float64  `json:"humidity"`
	Precip         float64  `json:"precip"`
	Precipprob     float64  `json:"precipprob"`
	Precipcover    float64  `json:"precipcover"`
	Preciptype     []string `json:"preciptype"`
	Snow           float64  `json:"snow"`
	Snowdepth      float64  `json:"snowdepth"`
	Windgust       float64  `json:"windgust"`
	Windspeed      float64  `json:"windspeed"`
	Winddir        float64  `json:"winddir"`
	Pressure       float64  `json:"pressure"`
	Cloudcover     float64  `json:"cloudcover"`
	Visibility     float64  `json:"visibility"`
	Solarradiation float64  `json:"solarradiation"`
	Solarenergy    float64  `json:"solarenergy"`
	Uvindex        float64  `json:"uvindex"`
	Severerisk     float64  `json:"severerisk"`
	Sunrise        string   `json:"sunrise"`
	Sunriseepoch   int      `json:"sunriseEpoch"`
	Sunset         string   `json:"sunset"`
	Sunsetepoch    int      `json:"sunsetEpoch"`
	Moonphase      float64  `json:"moonphase"`
	Conditions     string   `json:"conditions"`
	Description    string   `json:"description"`
	Icon           string   `json:"icon"`
	Stations       []string `json:"stations"`
	Source         string   `json:"source"`
}

type RecordApp struct {
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
	Sunrise           string    `json:"sunrise"`
	Sunset            string    `json:"sunset"`
	Conditions        string    `json:"conditions"`
	Precipprob        float64   `json:"precipprob"`
	Visibility        float64   `json:"visibility"`
	MaxTemp           float64   `json:"maxTemp"`
	MinTemp           float64   `json:"minTemp"`
	AvgWind           float64   `json:"windavg"`
}

type Alerts struct {
	Type     string `json:"type"`
	Features []struct {
		ID         string   `json:"id"`
		Type       string   `json:"type"`
		Geometry   Geometry `json:"geometry"`
		Properties Property `json:"properties"`
	} `json:"features"`
	Title   string    `json:"title"`
	Updated time.Time `json:"updated"`
}
type Geometry struct {
	Type        string        `json:"type"`
	Coordinates []interface{} `json:"coordinates,omitempty"`
	Geometries  []*Geometry   `json:"geometries,omitempty"`
}
type Property struct {
	IDURI    string `json:"@id"`
	Type     string `json:"@type"`
	ID       string `json:"id"`
	AreaDesc string `json:"areaDesc"`
	Geocode  struct {
		UGC  []string `json:"UGC"`
		SAME []string `json:"SAME"`
	} `json:"geocode"`
	AffectedZones []string      `json:"affectedZones"`
	References    []interface{} `json:"references"`
	Sent          time.Time     `json:"sent"`
	Effective     time.Time     `json:"effective"`
	Onset         time.Time     `json:"onset"`
	Expires       time.Time     `json:"expires"`
	Ends          time.Time     `json:"ends"`
	Status        string        `json:"status"`
	MessageType   string        `json:"messageType"`
	Category      string        `json:"category"`
	Severity      string        `json:"severity"`
	Certainty     string        `json:"certainty"`
	Urgency       string        `json:"urgency"`
	Event         string        `json:"event"`
	Sender        string        `json:"sender"`
	SenderName    string        `json:"senderName"`
	Headline      string        `json:"headline"`
	Description   string        `json:"description"`
	Instruction   string        `json:"instruction"`
	Response      string        `json:"response"`
	Parameters    struct {
		NWSheadline  []string `json:"NWSheadline"`
		PIL          []string `json:"PIL"`
		BLOCKCHANNEL []string `json:"BLOCKCHANNEL"`
	} `json:"parameters"`
}

type Beaufort struct {
	SVG   template.HTML
	Text  string
	Class string
}

type AqiCategory struct {
	Max   int64
	Color string
	Name  string
}

type AlmanacInfo struct {
	TempMax     AlmanacData
	TempMin     AlmanacData
	FeelMax     AlmanacData
	FeelMin     AlmanacData
	HumidityMax AlmanacData
	HumidityMin AlmanacData
	BaroMax     AlmanacData
	BaroMin     AlmanacData
	DewpointMax AlmanacData
	DewpointMin AlmanacData
	WindMax     AlmanacData
	WindMin     AlmanacData
	WindGustMax AlmanacData
	WindGustMin AlmanacData
}
type AlmanacData struct {
	Label    string
	Value    float64
	Recorded time.Time
}

type ChartData struct {
	Series     []Series     `json:"series"`
	Legend     Legend       `json:"legend"`
	Colors     []string     `json:"colors"`
	Chart      ChartDef     `json:"chart"`
	Responsive []Responsive `json:"responsive"`
	Stroke     Stroke       `json:"stroke"`
	Markers    Markers      `json:"markers"`
	Labels     Labels       `json:"labels"`
	Grid       Grid         `json:"grid"`
	DataLabels DataLabels   `json:"dataLabels"`
	Tooltip    Tooltip      `json:"tooltip"`
	Xaxis      Xaxis        `json:"xaxis"`
	Yaxis      Yaxis        `json:"yaxis"`
}
type Series struct {
	Name string    `json:"name"`
	Data []float64 `json:"data"`
}
type Legend struct {
	Show            bool         `json:"show"`
	Position        string       `json:"position"`
	HorizontalAlign string       `json:"horizontalAlign"`
	Labels          legendLabels `json:"labels"`
}
type legendLabels struct {
	Colors          []string `json:"colors"`
	UseSeriesColors bool     `json:"useSeriesColors"`
}
type DropShadow struct {
	Enabled bool    `json:"enabled"`
	Color   string  `json:"color"`
	Top     int     `json:"top"`
	Blur    int     `json:"blur"`
	Left    int     `json:"left"`
	Opacity float64 `json:"opacity"`
}
type Toolbar struct {
	Show bool `json:"show"`
}
type ChartDef struct {
	FontFamily string     `json:"fontFamily"`
	Height     int        `json:"height"`
	Width      int        `json:"width"`
	Type       string     `json:"type"`
	DropShadow DropShadow `json:"dropShadow"`
	Toolbar    Toolbar    `json:"toolbar"`
}

type Options struct {
	Chart ChartOptions `json:"chart"`
}

type ChartOptions struct {
	Height int `json:"height"`
}
type Responsive struct {
	Breakpoint int     `json:"breakpoint"`
	Options    Options `json:"options"`
}
type Stroke struct {
	Width []int  `json:"width"`
	Curve string `json:"curve"`
}
type Hover struct {
	Size       int `json:"size"`
	SizeOffset int `json:"sizeOffset"`
}
type Markers struct {
	Size            int      `json:"size"`
	Colors          string   `json:"colors"`
	StrokeColors    []string `json:"strokeColors"`
	StrokeWidth     int      `json:"strokeWidth"`
	StrokeOpacity   float64  `json:"strokeOpacity"`
	StrokeDashArray int      `json:"strokeDashArray"`
	FillOpacity     int      `json:"fillOpacity"`
	Discrete        []any    `json:"discrete"`
	Hover           Hover    `json:"hover"`
}
type Labels struct {
	Show     bool   `json:"show"`
	Position string `json:"position"`
}
type Lines struct {
	Show bool `json:"show"`
}
type Gridaxis struct {
	Lines Lines `json:"lines"`
}

type Grid struct {
	Xaxis  Gridaxis      `json:"xaxis"`
	Yaxis  Gridaxis      `json:"yaxis"`
	Row    GridRowColumn `json:"row"`
	Column GridRowColumn `json:"column"`
}
type GridRowColumn struct {
	Colors  string
	Opacity float64
}

type DataLabels struct {
	Enabled bool `json:"enabled"`
}
type AxisBorder struct {
	Show bool `json:"show"`
}
type AxisTicks struct {
	Show bool `json:"show"`
}
type Xaxis struct {
	Type       string     `json:"type"`
	Categories []string   `json:"categories"`
	AxisBorder AxisBorder `json:"axisBorder"`
	AxisTicks  AxisTicks  `json:"axisTicks"`
	Labels     AxisLabel  `json:"labels"`
}
type Style struct {
	FontSize string `json:"fontSize"`
}
type Title struct {
	Style Style `json:"style"`
}
type Yaxis struct {
	Title  Title     `json:"title"`
	Min    int       `json:"min"`
	Max    int       `json:"max"`
	Labels AxisLabel `json:"labels"`
}

type AxisLabel struct {
	Show  bool       `json:"show"`
	Align string     `json:"align"`
	Style LabelStyle `json:"style"`
}

type Tooltip struct {
	Theme string
}
type LabelStyle struct {
	Colors   string `json:"colors"`
	FontSize string `json:"fontSize"`
}

type Pill struct {
	Interval string
	Name     string
}
