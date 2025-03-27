export interface Current {
    id: number;
    mac: string;
    date: string;
    barometer: BaroData;
    humidity: HumidityData;
    temp: TempData;
    tempin: TemperatureData;
    temp1: TemperatureData;
    temp2: TemperatureData;
    temp3: TemperatureData;
    temp4: TemperatureData;
    rain: RainData;
    lightning: LightningData;
    aqi: AirQualityIndex;
    wind: WindData;
    uv: UVData;
    astro: AstroData;
    alert: Alert[];
}

export interface minmax {
    avg: minmaxStats;
    max: minmaxStats;
    min: minmaxStats;
}

export interface minmaxStats {
    day: StatDetails;
    month: StatDetails;
    year: StatDetails;
    yesterday: StatDetails;
}

export interface StatDetails {
    value: number;
    date: string;
}

export interface BaroData {
    baromabsin: number;
    baromrelin: number;
    minmax: minmax;
    trend: TrendData;
}

export interface HumidityData {
    humidity: number;
    dewpoint: number;
    minmax: minmax;
}

export interface TrendData {
    trend: string;
    by: number;
}

export interface TempData {
    temp: number;
    humidity: number;
    battout: number;
    feelslike: number;
    dewpoint: number;
    minmax: minmax;
}

export interface TemperatureData {
    temp: number;
    humidity: number;
    battout: number;
    minmax: minmax;
}

export interface RainData {
    daily: number;
    event: number;
    hourly: number;
    yearly: number;
    monthly: number;
    weekly: number;
    total: number;
    lastrain: string;
}

export interface LightningData {
    day: number;
    hour: number;
    distance: number;
    time: string;
    month: number;
    minmax: minmax;
}

export interface AirQualityIndex {
    pm25: number;
    pm2524h: number;
    minmax: minmax;
}

export interface WindData {
    winddir: number;
    windgustmph: number;
    windgustdir: number;
    windspeedmph: number;
    maxdailygust: number;
    windavg: number;
    minmax: minmax;
}

export interface UVData {
    uv: number;
    solarradiation: number;
    minmax: minmax;
}

export interface AstroData {
    sunrise: string;
    sunset: string;
    sunriseTomorrow: string;
    sunsetTomorrow: string;
    darkness: number;
    daylight: number;
    elevation: number;
    hasSunset: boolean;
}

export interface Alert {
    id: number;
    alertid: string;
    wxtype: string;
    areadesc: string;
    sent: string;
    effective: string;
    onset: string;
    expires: string;
    end: string;
    status: string;
    messagetype: string;
    category: string;
    severity: string;
    certainty: string;
    urgency: string;
    event: string;
    sender: string;
    senderName: string;
    headline: string;
    description: string;
    instruction: string;
    response: string;
}
