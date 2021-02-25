const moment = require('moment');
const AmbientWeatherApi = require("ambient-weather-api");
const request = require('request');

if (process.env.NODE_ENV !== 'production') {
    require('dotenv').config();
}

const api = new AmbientWeatherApi({
    apiKey: process.env.AMBIENT_WEATHER_API_KEY,
    applicationKey: process.env.AMBIENT_WEATHER_APPLICATION_KEY
});

api.connect()
api.on('connect', () => console.log('Connected to Ambient Weather Realtime API!'));
api.on('subscribed', res => {
    console.log('Subscribed')
    let data = res.devices[0].lastData;
    update(data)

});

api.on('data', data => {
    update(data)
});

api.subscribe(process.env.AMBIENT_WEATHER_API_KEY);


function update(data) {
    console.log(data)
    const randomNumber = Math.floor(Math.random() * 40) - 20;
    let d = {
        date: moment(data.dateutc).local().format("YYYY-MM-DDTHH:mm:ssZ"),
        mac: '',
        winddir: data.winddir ,
        windgustdir: data.winddir - randomNumber,
        windspeedmph: data.windspeedmph ,
        windgustmph: data.windgustmph ,
        maxdailygust: data.maxdailygust ,
        tempf: data.tempf ,
        hourlyrainin: data.hourlyrainin ,
        eventrainin: data.eventrainin ,
        dailyrainin: data.dailyrainin ,
        weeklyrainin: data.weeklyrainin ,
        monthlyrainin: data.monthlyrainin ,
        yearlyrainin: data.yearlyrainin ,
        baromrelin: data.baromrelin ,
        baromabsin: data.baromabsin ,
        humidity: data.humidity ,
        tempinf: data.tempinf ,
        humidityin: data.humidityin ,
        uv: data.uv ,
        temp1f: data.temp1f ,
        temp2f: data.temp2f ,
        humidity1: data.humidity1 ,
        humidity2: data.humidity2 ,
        solarradiation: data.solarradiation ,
        feelsLike: data.feelsLike ,
        dewPoint: data.dewPoint ,
        lightninghour: data.lightning_hour,
        lightningday: data.lightning_day,
        lightningdistance: data.lightning_distance,
        lightningtime: moment(data.lightning_time).local().format("YYYY-MM-DDTHH:mm:ssZ"),
        lastRain: moment(data.lastRain).local().format("YYYY-MM-DDTHH:mm:ssZ"),
        batt1: data.batt1,
        batt2: data.batt2,
        battout: data.battout,
        battin: data.battin,
        battlightning: data.batt_lightning
    };
    let options = {
        uri: 'http://localhost:3000/api/apiin',
        method: 'POST',
        json: d
    };

    request(options, function (error, response, body) {
        if (!error && response.statusCode == 200) {
            //console.log(body.id) // Print the shortened url.
        }
    });
    console.log(d.date, d.tempf+'°F', d.humidity+'%');
}