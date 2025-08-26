export interface gradient {
    background: string
    width: string
}
/**
 * The main data structure for the weather API response.
 */
export interface ForecastModel {
    queryCost: number;
    latitude: number;
    longitude: number;
    resolvedAddress: string;
    address: string;
    timezone: string;
    tzoffset: number;
    days: Day[];
    alerts: Alert[];
    stations: Record<string, Station>;
}

/**
 * Interface representing a single day's weather data.
 */
export interface Day {
    datetime: string;
    datetimeEpoch: number;
    tempmax: number;
    tempmin: number;
    temp: number;
    feelslikemax: number;
    feelslikemin: number;
    feelslike: number;
    dew: number;
    humidity: number;
    precip: number;
    precipprob: number;
    precipcover: number;
    preciptype: string[];
    snow: number;
    snowdepth: number;
    windgust: number;
    windspeed: number;
    winddir: number;
    pressure: number;
    cloudcover: number;
    visibility: number;
    solarradiation: number;
    solarenergy: number;
    uvindex: number;
    severerisk: number;
    sunrise: string;
    sunriseEpoch: number;
    sunset: string;
    sunsetEpoch: number;
    moonphase: number;
    conditions: string;
    description: string;
    icon: string;
    stations: string[];
    source: string;
    hours: Hour[];
}

/**
 * Interface representing a single hour's weather data.
 */
export interface Hour {
    datetime: string;
    datetimeEpoch: number;
    temp: number;
    feelslike: number;
    humidity: number;
    dew: number;
    precip: number;
    precipprob: number;
    snow: number;
    snowdepth: number;
    preciptype: string[] | null;
    windgust: number;
    windspeed: number;
    winddir: number;
    pressure: number;
    visibility: number;
    cloudcover: number;
    solarradiation: number;
    solarenergy: number;
    uvindex: number;
    severerisk: number;
    conditions: string;
    icon: string;
    stations: string[] | null;
    source: string;
}

/**
 * Interface for a weather alert.
 */
export interface Alert {
    event: string;
    headline: string;
    ends: string;
    endsEpoch: number;
    onset: string;
    onsetEpoch: number;
    id: string;
    language: string;
    link: string;
    description: string;
}

/**
 * Interface for a weather station.
 */
export interface Station {
    distance: number;
    latitude: number;
    longitude: number;
    useCount: number;
    id: string;
    name: string;
    quality: number;
    contribution: number;
}
