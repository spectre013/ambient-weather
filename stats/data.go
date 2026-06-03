package main

import "time"

type Forecast struct {
	QueryCost       int     `json:"queryCost"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	ResolvedAddress string  `json:"resolvedAddress"`
	Address         string  `json:"address"`
	Timezone        string  `json:"timezone"`
	Tzoffset        float64 `json:"tzoffset"`
	Days            []Day   `json:"days"`
}

type Day struct {
	Datetime       string   `json:"datetime"`
	DatetimeEpoch  int64    `json:"datetimeEpoch"`
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
	SunriseEpoch   int64    `json:"sunriseEpoch"`
	Sunset         string   `json:"sunset"`
	SunsetEpoch    int64    `json:"sunsetEpoch"`
	Moonphase      float64  `json:"moonphase"`
	Conditions     string   `json:"conditions"`
	Description    string   `json:"description"`
	Icon           string   `json:"icon"`
	Stations       []string `json:"stations"`
	Source         string   `json:"source"`
	Hours          []Hour   `json:"hours"`
}

type Hour struct {
	Datetime       string   `json:"datetime"`
	DatetimeEpoch  int64    `json:"datetimeEpoch"`
	Temp           float64  `json:"temp"`
	Feelslike      float64  `json:"feelslike"`
	Humidity       float64  `json:"humidity"`
	Dew            float64  `json:"dew"`
	Precip         float64  `json:"precip"`
	Precipprob     float64  `json:"precipprob"`
	Snow           float64  `json:"snow"`
	Snowdepth      float64  `json:"snowdepth"`
	Preciptype     []string `json:"preciptype"`
	Windgust       float64  `json:"windgust"`
	Windspeed      float64  `json:"windspeed"`
	Winddir        float64  `json:"winddir"`
	Pressure       float64  `json:"pressure"`
	Visibility     float64  `json:"visibility"`
	Cloudcover     float64  `json:"cloudcover"`
	Solarradiation float64  `json:"solarradiation"`
	Solarenergy    float64  `json:"solarenergy"`
	Uvindex        float64  `json:"uvindex"`
	Severerisk     float64  `json:"severerisk"`
	Conditions     string   `json:"conditions"`
	Icon           string   `json:"icon"`
	Stations       []string `json:"stations"`
	Source         string   `json:"source"`
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

// ----------------------------------------------------------------------------
// NWS station observation (api.weather.gov/stations/{id}/observations/latest)
// ----------------------------------------------------------------------------

type NWSObservation struct {
	Properties NWSObservationProps `json:"properties"`
}

type NWSObservationProps struct {
	Timestamp        time.Time           `json:"timestamp"`
	TextDescription  string              `json:"textDescription"`
	Icon             string              `json:"icon"`
	PresentWeather   []NWSPresentWeather `json:"presentWeather"`
	CloudLayers      []NWSCloudLayer     `json:"cloudLayers"`
	Temperature      NWSQuantity         `json:"temperature"`
	RelativeHumidity NWSQuantity         `json:"relativeHumidity"`
}

// NWSQuantity is a measured value with units. Value is a pointer because the
// feed sends null when an instrument is down or the value is unavailable.
type NWSQuantity struct {
	UnitCode string   `json:"unitCode"`
	Value    *float64 `json:"value"`
}

type NWSPresentWeather struct {
	Intensity string `json:"intensity"`
	Modifier  string `json:"modifier"`
	Weather   string `json:"weather"`
	RawString string `json:"rawString"`
}

type NWSCloudLayer struct {
	Base   NWSQuantity `json:"base"`
	Amount string      `json:"amount"`
}

// conditionsDB mirrors the public.conditions table columns.
type conditionsDB struct {
	Station         string    `json:"station"`
	ObservedAt      time.Time `json:"observed_at"`
	Conditions      string    `json:"conditions"`
	Icon            string    `json:"icon"`
	TextDescription string    `json:"text_description"`
	PresentWeather  string    `json:"present_weather"`
	CloudLayers     string    `json:"cloud_layers"`
	RawIcon         string    `json:"raw_icon"`
	Temperature     *float64  `json:"temperature"`
	Humidity        *float64  `json:"humidity"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type forecastDB struct {
	Datetime       time.Time `json:"datetime"`
	DatetimeEpoch  int64     `json:"datetimeEpoch"`
	TempMax        float64   `json:"tempmax"`
	TempMin        float64   `json:"tempmin"`
	Temp           float64   `json:"temp"`
	FeelsLikeMax   float64   `json:"feelslikemax"`
	FeelsLikeMin   float64   `json:"feelslikemin"`
	FeelsLike      float64   `json:"feelslike"`
	Dew            float64   `json:"dew"`
	Humidity       float64   `json:"humidity"`
	Precip         float64   `json:"precip"`
	PrecipProb     float64   `json:"precipprob"`
	PrecipCover    float64   `json:"precipcover"`
	PrecipType     string    `json:"preciptype"`
	Snow           float64   `json:"snow"`
	SnowDepth      float64   `json:"snowdepth"`
	WindGust       float64   `json:"windgust"`
	WindSpeed      float64   `json:"windspeed"`
	WindDir        float64   `json:"winddir"`
	Pressure       float64   `json:"pressure"`
	CloudCover     float64   `json:"cloudcover"`
	Visibility     float64   `json:"visibility"`
	SolarRadiation float64   `json:"solarradiation"`
	SolarEnergy    float64   `json:"solarenergy"`
	UVIndex        float64   `json:"uvindex"`
	SevereRisk     float64   `json:"severerisk"`
	Sunrise        string    `json:"sunrise"`
	SunriseEpoch   int64     `json:"sunriseEpoch"`
	Sunset         string    `json:"sunset"`
	SunsetEpoch    int64     `json:"sunsetEpoch"`
	MoonPhase      float64   `json:"moonphase"`
	Conditions     string    `json:"conditions"`
	Description    string    `json:"description"`
	Icon           string    `json:"icon"`
	Stations       string    `json:"stations"`
	Source         string    `json:"source"`
	Hours          string    `json:"hours"`
	Summary        string    `json:"summary"`
}
