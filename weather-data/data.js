#!/usr/bin/env node

const mysql = require('mysql2/promise');
const moment = require('moment');
const AmbientWeatherApi = require("ambient-weather-api");

if (process.env.NODE_ENV !== 'production') {
    require('dotenv').config();
}

const database_config = {
  host     : process.env.DB_HOST,
  user     : process.env.DB_USER,
  password : process.env.DB_PASSWORD,
  database : process.env.DB_DATABASE
};

console.log(database_config);

let connection = {};
async function getConnection()
{
    await mysql.createConnection(database_config).then(conn => {
        connection = conn;
    });
}

getConnection();



const api = new AmbientWeatherApi({
  apiKey: process.env.AMBIENT_WEATHER_API_KEY,
  applicationKey: process.env.AMBIENT_WEATHER_APPLICATION_KEY
});

api.connect()
api.on('connect', () => console.log('Connected to Ambient Weather Realtime API!'));
api.on('subscribed', res => {
    console.log('Subscribed')
    var data = res.devices[0].lastData;
    update(data);
  });

api.on('data', data => {
    console.log(data.date + ' current outdoor temperature is: ' + data.tempf + '°F');
    update(data);
});

api.subscribe(process.env.AMBIENT_WEATHER_API_KEY);

async function update(data) {
    const randomNumber = Math.floor(Math.random() * 40) - 20;
    let insertData = {
        date: moment(data.dateutc).format("YYYY-MM-DD HH:mm:ss"),
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
        lastRain: moment(data.lastRain ).format("YYYY-MM-DD HH:mm:ss")};

        try {
            const query  = 'INSERT INTO records SET ?';
            await connection.query(query, insertData)
                .then(([rows]) => {
                    if(rows.length > 0) {
                        data = rows[0];
                    } else {
                        data = {};
                    }
                });
        } catch (e) {
            console.log('Error Updating Records', e);
        }
    updateStatistics().then(() => {
        console.log('Min Max Update complete');
    }) ;
}


async function updateStatistics() {
    let queries = {};
    const types = ['MAX','MIN','AVG'];
    const periods = ['day','yesterday','month','year'];
    const fields = ['tempf','tempinf','temp1f','temp2f','baromrelin','uv','humidity','windspeedmph','windgustmph','dewpoint'];

    try {
        periods.forEach((period) => {
            types.forEach((type) => {
                fields.forEach((field) => {
                    let key = period+'_'+type.toLowerCase()+'_'+field;

                    let order = '';
                    if(type === 'MAX') {
                        order = ' DESC';
                    }

                    let q = `select ${type}(${field})  as value, \`date\` from records where \`date\` between ? and ? group by \`date\`,${field} order by ${field}${order} limit 0,1`;
                    if(type === 'AVG') {
                        q = `select CAST(${type}(${field}) AS DECIMAL(10,2)) as value from records where \`date\` between ? and ? order by ${field} limit 0,1`;
                    }
                    queries[key] = {
                        query:q,
                        params: getTimeframe(period)
                    };
                });
            });
        });
        //console.log(queries);
        let name = '';
        let qt = '';
        await asyncForEach( Object.keys(queries),async (key) => {
            await asyncForEach(Object.keys(queries[key]),async (tf) => {
                name = key;
                qt = tf;
                let query = queries[key];
                await connection.query(query.query,query.params , (err) => {
                    if (err) throw err;
                }).then(async (rows) => {
                    if(rows.length > 0) {
                        try {

                            let value = key+'_value';
                            let date = ', '+key+'_date = ?';
                            let insertData = [];
                            let val = 0
                            let dt = moment().utc().utcOffset(6).format('YYYY-MM-DD HH:mm:ss')
                            if(typeof rows[0][0] !== 'undefined') {
                                val = rows[0][0].value;
                                dt = rows[0][0].date;
                            }
                            insertData.push(val);
                            if(key.includes('avg')) {
                                date = '';
                            } else {
                                insertData.push(dt);
                            }

                            const update  = `UPDATE stats SET ${value} = ?${date} WHERE id = 1`;
                            await connection.query(update, insertData)
                                .then(() => {

                                });
                        } catch (e) {
                            console.log('Error Updating stats', e);
                        }
                    }
                });
            });
        });
    } catch (e) {
        console.log('Error getting stats', e);
    }

}

function getTimeframe(timeframe) {
    let dates = [];
    if(timeframe === 'yesterday') {
        dates = [moment.utc().startOf('day').subtract(1,'days').utcOffset(6).format('YYYY-MM-DD HH:mm:ss'),
            moment.utc().endOf('day').subtract(1,'days').utcOffset(6).format('YYYY-MM-DD HH:mm:ss')];
    } else if(timeframe === 'day') {
        dates = [moment.utc().startOf('day').utcOffset(6).format('YYYY-MM-DD HH:mm:ss'),
            moment.utc().endOf('day').utcOffset(6).format('YYYY-MM-DD HH:mm:ss')];
    } else {
        dates = [moment.utc().startOf(timeframe).utcOffset(6).format('YYYY-MM-DD HH:mm:ss'),moment.utc().endOf(timeframe).utcOffset(6).format('YYYY-MM-DD HH:mm:ss')];
    }
    return dates;
}``

async function asyncForEach(array, callback) {
    for (let index = 0; index < array.length; index++) {
        await callback(array[index], index, array);
    }
}