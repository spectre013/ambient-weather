const path = require('path');
const http = require('http');
const express = require('express');
const AmbientWeatherApi = require('ambient-weather-api');
const axios = require('axios');
const lune = require('lune');
const mysql = require('mysql2/promise');
const moment = require('moment');

if (process.env.NODE_ENV !== 'production') {
    require('dotenv').config();
}

const database_config = {
  host     : process.env.DB_HOST,
  user     : process.env.DB_USER,
  password : process.env.DB_PASSWORD,
  database : process.env.DB_DATABASE,
  multipleStatements: true
};

let connection = {};
mysql.createConnection(database_config).then(conn => {
  connection = conn;
});

const port = process.env.PORT || 3000;

let current = {};

const app = express();
const server = http.createServer(app);

const io = require('socket.io')(server, {path: '/api/ws'});

app.get('/api/forecast', async (req, res)  => {
    let data = {};
    try {
        let url = `https://api.darksky.net/forecast/${process.env.DARKSKY}/${process.env.LAT},${process.env.LON}`;
        await axios
          .get(url)
          .then(response => (data = response.data));
          res.json(data);
    } catch (e) {
      console.log(e);
    }
});

app.get('/api/current', async (req, res)  => {
  let data = [];
  try {
    const query  = 'select * from `records` order by `date` desc limit 0,1';
    await connection.query(query)
        .then(([rows]) => {
            if(rows.length > 0) {
              data = rows[0];
            } else {
              data = {};
            }
        });
    let rain = {};
    const start = moment().format('YYYY-MM-DD HH:mm:ss');
    const end = moment().subtract(60, 'minutes').format('YYYY-MM-DD HH:mm:ss');
    const hourRain = 'select MAX(dailyrainin) as hourlyrain from records where `date` between ? and ?';
    console.log(start,end)
      await connection.query(hourRain,[end,start])
          .then(([rows]) => {
              if(rows.length > 0) {
                  rain = rows[0];
              } else {
                  rain = {hourlyrain:0};
              }
          });
      data.hourlyrain = rain.hourlyrain;
  } catch (e) {
    res.status(500).send({erroMessage:e.message});
    console.log(e);
  }

  res.json(data);
});

app.get('/api/alltime/:calc/:type', async (req, res)  => {
    let type = req.params.type.replace(/^[\u0080-\uffff]/g, "");
    let calc = req.params.calc.replace(/^[\u0080-\uffff]/g, "");
    let data = {};
    let order = "";
    if(calc === "max") {
        order = `desc`;
    }
    try {
        const query  =  `select ${calc}(${type}) as value, \`date\` from records group by \`date\`,${type} order by ${type} ${order} limit 0,1`;
        await connection.query(query)
            .then(([rows]) => {
                if(rows.length > 0) {
                    data = rows[0];
                } else {
                    data = {};
                }
            });
    } catch (e) {
        res.status(500).send({erroMessage:e.message});
        console.log(e);
    }

    res.json(data);
});

app.get('/api/trend/:type', async (req, res)  => {

  const type = req.params.type;

  if(type === 'temp') {
    let avg=0;
    let current = 0;
    let temptrend = {trend:'',by:0};
    const start = moment().format('YYYY-MM-DD HH:mm:ss');
    const end = moment().subtract(30, 'minutes').format('YYYY-MM-DD HH:mm:ss');
    const curTemp = 'select AVG(tempf) as temp from records where `date` between ? and ?';
    await connection.query(curTemp, [end,start])
        .then((result) => {
          avg = result[0][0].temp;
        }).catch((err) => {
          console.log(err);
          res.status(500).send(err);
    });

    const query  = 'select tempf from records order by date desc limit 0,1';
    await connection.query(query, [])
          .then((result) => {
            current = result[0][0].tempf;
          }).catch((err) => {
          console.log(err);
          res.status(500).send(err);
      });

      if(current > avg)  {
        //trend up
        temptrend.trend = 'up';
        temptrend.by = (current-avg).toFixed(1);
      } else {
        //trend down
        temptrend.trend = 'down';
        temptrend.by = (avg-current).toFixed(1);
      }
      res.json(temptrend);

  } else if(type === 'barometer') {
    let barotrend = {trend:'',by:0};
    let avg = 0;
    let current = 0;
    const start = moment().format('YYYY-MM-DD HH:mm:ss');
    const end = moment().subtract(3, 'hours').format('YYYY-MM-DD HH:mm:ss');
    const barAvg = 'select AVG(baromrelin) as baro from records where `date` between ? and ?';
    await connection.query(barAvg, [end,start])
        .then((result) => {
          avg = result[0][0].tempf;``
        }).catch((err) => {
          console.log(err);
          res.status(500).send(err);
        });
    const query  = 'select * from records order by date desc limit 1,1';
    await connection.query(query, [])
        .then((result) => {
          current = result[0][0].tempf;
        }).catch((err) => {
          console.log(err);
          res.status(500).send(err);
        });

        if(current > avg)  {
          //trend up
          barotrend.trend = 'Steady';
          if((current - avg) > .5 ) {
            barotrend.trend = 'Rising';
          }
        } else {
          //trend down
          barotrend.trend = 'Steady';
          if((avg - current) > .5 ) {
            barotrend.trend = 'Falling';
          }
        }

        res.json(barotrend);
  }
});

app.get('/api/chart/:type/:period', async (req, res)  => {
    const type = req.params.type;
    const period = req.params.period;
    let start = moment().startOf(period).format('YYYY-MM-DD 00:00:00');
    let end = moment().endOf(period).format('YYYY-MM-DD 23:59:59');
    let dateformat = '%Y-%m-%d';
    if(period === 'day') {
        dateformat = '%H:%d:%s';
    }

    let data = `select DATE_FORMAT(r.date,'${dateformat}') AS mmdd, max(${type}) max, min(${type}) min from records r where date between ? AND ? group by mmdd order by mmdd`;
    let json = {data1: [],data2: []};
    //console.log(data,start,end);
    await connection.query(data, [start, end])
      .then(async (result) => {
        await asyncForEach(result[0], (record) => {
          json.data1.push({label: record.mmdd,y:parseFloat(record.max)});
          json.data2.push({label: record.mmdd,y:parseFloat(record.min)});
        });

      }).catch((err) => {
        console.log(err);
        res.status(500).send(err);
      });
    res.json(json);
});

app.get('/api/wind', async (req, res)  => {
  let data = {};
  let start = moment().format('YYYY-MM-DD 00:00:00');
  let end = moment().format('YYYY-MM-DD 23:59:59');
  let wind = 'select max(windspeedmph) as value, `date` from records where `date` between ? and ? group by `date`,windspeedmph order by windspeedmph desc limit 0,1';
  let gust = 'select max(windgustmph) as value, `date` from records where `date` between ? and ? group by `date`,windgustmph order by windgustmph desc limit 0,1';
  let avg = 'select AVG(windspeedmph) as wind from records where `date` between ? and ?';
  let dir = 'select AVG(winddir) as dir from records where `date` between ? and ?';
  await connection.query(wind + ';'+ gust + ';'+ avg + ';'+ dir, [start,end,start,end,start,end,start,end])
      .then((results) => {
        let w = results[0];
        data = {wind:w[0][0], gust:w[1][0],avg:w[2][0],dir:w[3][0]};
      }).catch((err) => {
        console.log(err);
        res.status(500).send(err);
      });
  res.json(data);
});


app.get('/api/minmax/:field', async (req, res)  => {
  const field = req.params.field;
  let data = {};
  try {
      await connection.query("select * from stats where id = ?",[1])
          .then((results) => {
              let r = results[0][0];
              let keys = Object.keys(r)
              for (let index = 0; index < keys.length; index++) {
                let key = keys[index];
                if(key.includes(field)) {
                    //month_max_humidity_value
                    let parts = key.split("_");
                    let period = parts[0];
                    let type = parts[1];
                    let value = parts[3]
                    if(!data.hasOwnProperty(type)) {
                        data[type] = {};
                    }
                    if(!data[type].hasOwnProperty(period)) {
                        data[type][period] = {};
                    }
                    if(!data[type][period].hasOwnProperty(value)) {
                        data[type][period][value] = {};
                    }
                    if(value === 'value' && !r[key]) {
                        r[key] = 0.0;
                    }
                    data[type][period][value] = r[key];

                }
              }

          }).catch((err) => {
              console.log(err);
              res.status(500).send(err);
          });
  } catch (e) {
    res.status(500).send({erroMessage:e.message});
    console.log(e);
  }

  res.json(data);
});


async function asyncForEach(array, callback) {
  for (let index = 0; index < array.length; index++) {
    await callback(array[index], index, array);
  }
}

function getTimeframe(timeframe) {
    let dates = [];
    if(timeframe === 'yesterday') {
        dates = [moment().startOf('day').subtract(1,'days').format('YYYY-MM-DD HH:mm:ss'),
                 moment().endOf('day').subtract(1,'days').format('YYYY-MM-DD HH:mm:ss')];
    } else {
        dates = [moment().startOf(timeframe).format('YYYY-MM-DD HH:mm:ss'),moment().endOf(timeframe).format('YYYY-MM-DD HH:mm:ss')];
    }
  return dates;
}


app.get('/api/luna', async (req, res)  => {
  let lunardata = {tomorrow:{}};
  //var phases = lune.phase_hunt()
  const fullmoon = lune.phase_range(moment().toDate(),moment().add(30,'days').toDate(),lune.PHASE_FULL);
  const newmoon = lune.phase_range(moment().toDate(),moment().add(30,'days').toDate(),lune.PHASE_NEW);
  const currentPhase = lune.phase()
  let names= ['New Moon','Waxing Crescent','First Quarter','Waxing Gibbous','Full Moon','Waning Gibbous',
  'Third Quarter','Waning Crescent','New Moon'];
  try {
  await axios
    .get(`https://api.ipgeolocation.io/astronomy?apiKey=${process.env.IPGEO}&lat=${process.env.LAT}&long=${process.env.LON}`)
    .then(response => (data = response.data));
    lunardata = data;

    let tomorrow= moment().add(2,'days').format('YYYY-MM-DD');
    let uri = `https://api.ipgeolocation.io/astronomy?apiKey=${process.env.IPGEO}&lat=${process.env.LAT}&long=${process.env.LON}&date=${tomorrow}`;
    await axios
    .get(uri)
    .then(response => (lunardata.tomorrow = response.data));

  } catch (e) {
      //this will eventually be handled by your error handling middleware
      console.log(e);
  }
  lunardata.newmoon = (newmoon.length > 0 ) ? newmoon[0] : null;
  lunardata.fullmoon = (fullmoon.length > 0 ) ? fullmoon[0] : null;
  lunardata.phase = names[ Math.floor((currentPhase.phase + 0.0625)* 8)];
  lunardata.illuminated = Math.round(currentPhase.illuminated * 100);
  lunardata.age = Math.round(currentPhase.age);
  res.json(lunardata);
});

const api = new AmbientWeatherApi({
  apiKey: process.env.AMBIENT_WEATHER_API_KEY,
  applicationKey: process.env.AMBIENT_WEATHER_APPLICATION_KEY
});

api.connect();
api.on('connect', () => console.log('Connected to Ambient Weather Realtime API!'));
api.on('subscribed', data => {
    current = data.devices[0].lastData;
  });



io.on('connection', (socket) => {
    console.log('New User Connected');
    socket.emit('data', current);

    api.on('data', data => {
        current = data;
        socket.emit('data',data);
    });


    socket.on('disconnect', () => {
        console.log('User disconnected')
    });

});

api.subscribe(process.env.AMBIENT_WEATHER_API_KEY);
server.listen(port, () => {
    console.log(`Servers is up on port ${port}`)
});