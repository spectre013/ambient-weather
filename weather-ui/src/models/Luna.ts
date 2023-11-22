export interface Luna {
    location: Location
    date: string
    sunrise: string
    sunset: string
    solar_noon: string
    day_length: string
    sun_altitude: number
    sun_distance: number
    sun_azimuth: number
    moonrise: string
    moonset: string
    moon_altitude: number
    moon_distance: number
    moon_azimuth: number
    moon_parallactic_angle: number
    tomorrow: Tomorrow
    newmoon: string
    nextnewmoon: string
    fullmoon: string
    phase: string
    illuminated: number
    age: number
}

export interface Location {
    latitude: number
    longitude: number
}

export interface Tomorrow {
    location: Location2
    date: string
    sunrise: string
    sunset: string
    solar_noon: string
    day_length: string
    sun_altitude: number
    sun_distance: number
    sun_azimuth: number
    moonrise: string
    moonset: string
    moon_altitude: number
    moon_distance: number
    moon_azimuth: number
    moon_parallactic_angle: number
}

export interface Location2 {
    latitude: number
    longitude: number
}
