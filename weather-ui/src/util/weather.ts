import moment from "moment/moment";
import {ChartData} from "../models/DataSeries.ts";
import {WindDialOptions} from "../models/WindDial.ts";
import WindDial from "./windDial.ts";

export function processData(data: ChartData) {
    const dates = data[0].values.map(v => new Date(v.date));
    return dates.map((date, index) => {
        const dataPoint = { date };
        data.forEach(series => {
            // @ts-expect-error - series.key exists
            dataPoint[series.key] = series.values[index].value;
        });
        return dataPoint;
    });
}


export function full(date: string) {
    return moment(date).format('YYYY-MM-DD HH:mm:ss');
}

export function timeFormat(time: string):string {
    return moment(time).format("HH:mm");
}

export function timeFormatAMPM(time: string):string {
    return moment(time).format("HH:mm:ss A");
}

export function freezeDate(date: string):string {
    const date1 = moment('2010-01-1');
    const date2 = moment(date);
    if(date2.isBefore(date1)) {
        return "N/A";
    }
    return date2.format("YYYY-MM-DD HH:mm:ss");
}

export function getOtherUnit(units: string): string {
    if (units === 'metric') {
        return 'Imperial';
    } else {
        return 'Metric';
    }
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

export function degToCompass(num: number) {
    const val = Math.floor(num / 22.5 + 0.5);
    const arr = [
        'North',
        'NNE',
        'NE',
        'ENE',
        'East',
        'ESE',
        'SE',
        'SSE',
        'South',
        'SSW',
        'SW',
        'WSW',
        'West',
        'WNW',
        'NW',
        'NNW',
    ];
    return arr[val % 16];
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

export function distanceLabel(units: string): string {
    if (units === 'metric') {
        return 'km';
    } else {
        return 'mi';
    }
}

export function tempLabelAlt(units: string): string {
    if (units === 'metric') {
        return 'F';
    } else {
        return 'C';
    }
}
export function rainDisplay(rn: number, units : string) {
    if (units === 'metric') {
        return (rn * 25.4).toFixed(0);
    }
    return rn.toFixed(2);
}
export function baroDisplay(baro: number, units: string): string {
    const b = baro;
    if (units === 'metric') {
        return (baro * 33.86).toFixed(2);
    }
    return b.toFixed(2);
}
export function windDisplay(wind: number, units: string): string {
    if (units === 'imperial') {
        return (wind / 2.237).toFixed(0);
    }
    return wind.toFixed(0);
}
export function tempDisplay(temp: number, units:string) {
    let t = temp;
    if (units === 'metric') {
        t = ((temp - 32) * 5) / 9;
    }
    return t.toFixed(0);
}

/**
 * Converts a value from miles to kilometers.
 * @param {number} value - The distance in miles.
 * @param {string} units - The unit system, either 'imperial' or 'metric'.
 * @returns {number} The distance in kilometers.
 */
export function distanceDisplay(value: number, units:string) : number {
    const CONVERSION_FACTOR = 1.60934;
    let update = value;
    if (units === 'metric') {
        update = value * CONVERSION_FACTOR;
    }
    return update
}

/**
 * Maps a given temperature in Fahrenheit to a hex color value.
 * @param {number} tempFahrenheit The temperature in Fahrenheit.
 * @returns {string} The hex color string.
 */
export function tempToHex(tempFahrenheit: number) : string {
    if (tempFahrenheit >= 90) {
        return '#c62828';
    } else if (tempFahrenheit >= 80) {
        return '#e53935';
    } else if (tempFahrenheit >= 70) {
        return '#ef5350';
    } else if (tempFahrenheit >= 60) {
        return '#ffa726';
    } else if (tempFahrenheit >= 50) {
        return '#ffee58';
    } else if (tempFahrenheit >= 40) {
        return '#66bb6a';
    } else if (tempFahrenheit >= 30) {
        return '#42a5f5';
    } else if (tempFahrenheit >= 20) {
        return '#3949ab';
    } else if (tempFahrenheit >= 10) {
        return '#283593';
    } else {
        return '#1a237e';
    }
}

/**
 * Maps a given temperature in Fahrenheit to a CSS class name.
 * @param {number} tempFahrenheit The temperature in Fahrenheit.
 * @returns {string} The CSS class name.
 */
export function tempColor(tempFahrenheit: number): string {
    if (tempFahrenheit >= 90) {
        return 'color-range-90-and-above';
    } else if (tempFahrenheit >= 80) {
        return 'color-range-80-to-89';
    } else if (tempFahrenheit >= 70) {
        return 'color-range-70-to-79';
    } else if (tempFahrenheit >= 60) {
        return 'color-range-60-to-69';
    } else if (tempFahrenheit >= 50) {
        return 'color-range-50-to-59';
    } else if (tempFahrenheit >= 40) {
        return 'color-range-40-to-49';
    } else if (tempFahrenheit >= 30) {
        return 'color-range-30-to-39';
    } else if (tempFahrenheit >= 20) {
        return 'color-range-20-to-29';
    } else if (tempFahrenheit >= 10) {
        return 'color-range-10-to-19';
    } else {
        return 'color-range-10-below';
    }
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

/**
 * Returns a color based on the Beaufort wind scale.
 * @param {number} windSpeed The wind speed in knots.
 * @returns {string} A hex color string.
 */
export function BeaufortHex(windSpeed: number) {
    if (windSpeed < 1) {
        return '#A2E3B1'; // light green, for calm winds
    } else if (windSpeed <= 3) {
        return '#8EDD9B'; // greenish-yellow
    } else if (windSpeed <= 6) {
        return '#7AD685'; // bright green
    } else if (windSpeed <= 10) {
        return '#68C36E'; // medium green
    } else if (windSpeed <= 16) {
        return '#5BA860'; // dark green
    } else if (windSpeed <= 21) {
        return '#4F9456'; // blue-green
    } else if (windSpeed <= 27) {
        return '#417B46'; // darker blue-green
    } else if (windSpeed <= 33) {
        return '#35633A'; // very dark green
    } else if (windSpeed <= 40) {
        return '#D3B11A'; // yellow
    } else if (windSpeed <= 47) {
        return '#E19C1D'; // orange
    } else if (windSpeed <= 55) {
        return '#EA792D'; // dark orange
    } else if (windSpeed <= 63) {
        return '#F45543'; // reddish-orange
    } else {
        return '#B02525'; // red, for hurricane force winds
    }
}

export function windClass(beaufortValue: number) : string {
    if (beaufortValue >= 12) {
        return 'beaufort6';
    } else if (beaufortValue >= 9) {
        return 'beaufort6';
    } else if (beaufortValue >= 8) {
        return 'beaufort6';
    } else if (beaufortValue >= 7) {
        return 'beaufort6';
    } else if (beaufortValue >= 6) {
        return 'beaufort4-5';
    } else if (beaufortValue >= 5) {
        return 'beaufort4-5';
    } else if (beaufortValue >= 4) {
        return 'beaufort3-4';
    } else if (beaufortValue >= 3) {
        return 'beaufort1-3';
    } else if (beaufortValue >= 2) {
        return 'beaufort1-3';
    } else if (beaufortValue >= 1) {
        return 'beaufort1-3';
    } else if (beaufortValue >= 0) {
        return 'beaufort1-3';
    }
    return 'beaufort1-3';
}

export const formatDay = (dateString: string):string => {
    const date = new Date(dateString);
    return date.getDate().toString();
};

export function createDial(speed: number, direction: number, gusts: number, color: string, radiusColor: string, tickColor: string) : string  {
    const windData: WindDialOptions = {
        speed: speed,
        direction: direction, // East
        gusts: gusts,
        color: color,
        radiusColor: radiusColor,
        tickColor: tickColor
    };

    const windDial = new WindDial(windData);
    return windDial.generateSvg();
}