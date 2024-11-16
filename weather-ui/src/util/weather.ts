import moment from "moment/moment";

export function full(date: string) {
    return moment(date).format('YYYY-MM-DD HH:mm:ss');
}

export function timeFormat(time: string):string {
    return moment(time).format("HH:mm");
}

export function rainLabel(units: string): string {
    if (units === 'metric') {
        return 'mm';
    } else {
        return 'in';
    }
}
export function baroLabel(units: string): string {
    if (units === 'metric') {
        return 'hPA';
    } else {
        return 'inHG';
    }
}
export function baroLabelAlt(units: string): string {
    if (units === 'metric') {
        return 'inHG';
    } else {
        return 'hPA';
    }
}
export function windLabel(units: string): string {
    if (units === 'imperial') {
        return 'MPH';
    } else {
        return 'M/S';
    }
}

export function windLabelALT(units: string): string {
    if (units === 'metric') {
        return 'MPH';
    } else {
        return 'M/S';
    }
}

export function tempLabel(units: string): string {
    if (units === 'metric') {
        return 'C';
    } else {
        return 'F';
    }
}
export function tempLabelAlt(units: string): string {
    if (units === 'metric') {
        return 'F';
    } else {
        return 'C';
    }
}
export function rainDisplay(rn: number, units = 'imperial') {
    if (units === 'metric') {
        return (rn * 25.4).toFixed(0);
    }
    return rn.toFixed(2);
}
export function baroDisplay(baro: number, units = 'imperial'): string {
    const b = baro;
    if (units === 'metric') {
        return (baro * 33.86).toFixed(2);
    }
    return b.toFixed(2);
}
export function windDisplay(wind: number, units = 'mph'): string {
    if (units === 'ms') {
        return (wind / 2.237).toFixed(0);
    }
    return wind.toFixed(0);
}
export function tempDisplay(temp: number, units = 'imperial') {
    let t = temp;
    if (units === 'metric') {
        t = ((temp - 32) * 5) / 9;
    }
    return t.toFixed(0);
}
export function tempColor(temp: number): string {
    if (temp <= 14) {
        return 'tempcolor10';
    } else if (temp <= 23) {
        return 'tempcolor0-5';
    } else if (temp <= 32) {
        return 'tempcolorzero';
    } else if (temp <= 41) {
        return 'tempcolor0-5';
    } else if (temp < 50) {
        return 'tempcolor6-10';
    } else if (temp < 59) {
        return 'tempcolor11-15';
    } else if (temp < 68) {
        return 'tempcolor16-20';
    } else if (temp < 77) {
        return 'tempcolor21-25';
    } else if (temp < 86) {
        return 'tempcolor26-30';
    } else if (temp < 95) {
        return 'tempcolor31-35';
    } else if (temp < 104) {
        return 'tempcolor36-40';
    } else if (temp < 113) {
        return 'tempcolor41-45';
    } else if (temp < 212) {
        return 'tempcolor50';
    }
    return ""
}
export function dewPointClass(dewpoint: number): string {
    if (dewpoint > 69.8) {
        return 'tempmodulehome25-30c';
    } else if (dewpoint >= 68) {
        return 'tempmodulehome20-25c';
    } else if (dewpoint >= 59) {
        return 'tempmodulehome15-20c';
    } else if (dewpoint >= 50) {
        return 'tempmodulehome10-15c';
    } else if (dewpoint > 41) {
        return 'tempmodulehome5-10c';
    } else if (dewpoint >= 32) {
        return 'tempmodulehome0-5c';
    } else if (dewpoint > 14) {
        return 'tempmodulehome-10-0c';
    } else if (dewpoint >= -50) {
        return 'tempmodulehome-50-10c';
    }
    return 'tempmodulehome0-5c';
}
export function humidityClass(humidity: number): string {
    if (humidity > 90) {
        return 'temphumcircle80-100';
    } else if (humidity > 70) {
        return 'temphumcircle60-80';
    } else if (humidity > 35) {
        return 'temphumcircle35-60';
    } else if (humidity > 25) {
        return 'temphumcircle25-35';
    } else if (humidity <= 25) {
        return 'temphumcircle0-25';
    }
    return '';
}
export function temperaturetoday(temp: number): string {
    if (temp >= 105.8) {
        return 'temperaturetoday41-45';
    } else if (temp >= 96.8) {
        return 'temperaturetoday36-40';
    } else if (temp >= 87.8) {
        return 'temperaturetoday31-35';
    } else if (temp >= 78.8) {
        return 'temperaturetoday26-30';
    } else if (temp >= 69.8) {
        return 'temperaturetoday21-25';
    } else if (temp >= 60.8) {
        return 'temperaturetoday16-20';
    } else if (temp >= 50) {
        return 'temperaturetoday11-15';
    } else if (temp > 42.8) {
        return 'temperaturetoday6-10';
    } else if (temp >= 32) {
        return 'temperaturetoday0-5';
    } else if (temp < 32) {
        return 'temperaturetodayminus';
    } else if (temp <= 23) {
        return 'temperaturetodayminus5';
    } else if (temp < -14) {
        return 'temperaturetodayminus10';
    }
    return '';
}
export function smallTempClass(temp: number): string {
    if (temp >= 104) {
        return 'tempmodulehome40-50c';
    } else if (temp >= 95) {
        return 'tempmodulehome35-40c';
    } else if (temp >= 86) {
        return 'tempmodulehome30-35c';
    } else if (temp >= 77) {
        return 'tempmodulehome25-30c';
    } else if (temp >= 68) {
        return 'tempmodulehome20-25c';
    } else if (temp >= 59) {
        return 'tempmodulehome15-20c';
    } else if (temp >= 50) {
        return 'tempmodulehome10-15c';
    } else if (temp > 41) {
        return 'tempmodulehome5-10c';
    } else if (temp >= 32) {
        return 'tempmodulehome0-5c';
    } else if (temp > 14) {
        return 'tempmodulehome-10-0c';
    } else if (temp > -50) {
        return 'tempmodulehome-50-10c';
    }
    return '';
}
export function tempClass(temp: number): string {
    if (temp < 14) {
        return 'outsideminus10';
    } else if (temp <= 23) {
        return 'outsideminus5';
    } else if (temp <= 32) {
        return 'outsidezero';
    } else if (temp < 41) {
        return 'outside0-5';
    } else if (temp < 50) {
        return 'outside6-10';
    } else if (temp < 59) {
        return 'outside11-15';
    } else if (temp < 68) {
        return 'outside16-20';
    } else if (temp < 77) {
        return 'outside21-25';
    } else if (temp < 86) {
        return 'outside26-30';
    } else if (temp < 95) {
        return 'outside31-35';
    } else if (temp < 104) {
        return ' outside36-40';
    } else if (temp < 113) {
        return ' outside41-45';
    } else if (temp < 150) {
        return ' outside50';
    }
    return '';
}
