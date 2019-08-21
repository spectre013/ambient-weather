const axios = require('axios');
const mysql = require('mysql');
const moment = require('moment');
require('dotenv').config()

const database_config = {
    host     : process.env.DB_HOST,
    user     : process.env.DB_USER,
    password : process.env.DB_PASSWORD,
    database : process.env.DB_DATABASE
};

const ambient = 'https://api.ambientweather.net/v1/devices';
const mac = '84:F3:EB:20:DA:9E';
const apiurl = `${ambient}/${mac}?apiKey=${process.env.AMBIENT_WEATHER_API_KEY}&applicationKey=${process.env.AMBIENT_WEATHER_APPLICATION_KEY}&endDate={date}&limit=288`;
const connection = mysql.createConnection(database_config);

let lastApiCall = moment();

connection.connect();
const q = 'select `date` from records order by `date` desc';
connection.query(q, [], function(err, results) {
    if (err) throw err;
    let prevDate = moment();
    results.forEach(async (data) => {
        await getHistoric(prevDate, data)
        prevDate = data.date;
    });

});


function getHistoric(date, data) {
    let lastDate = moment(data.date);
    console.log("Date check 1", moment.duration(date.diff(lastDate)).asMinutes(),date.format('YYYY-MM-DD HH:mm:ss'),lastDate.format('YYYY-MM-DD HH:mm:ss'))
    if(moment.duration(date.diff(lastDate)).asMinutes() > 5) {
        // make api call and get records
        let apiresult = fetchData(date);
        for (let i = 0; i <apiresult.length; i++) {
            let newRecords = apiresult[i];
            const nd = moment(newRecords.date);
            //check record date against last db date if > 5 min add
            console.log("Date check 2", moment.duration(moment(nd).diff(currentDate)).asMinutes())
            if(moment.duration(moment(nd).diff(currentDate)).asMinutes() > 5) {
                // Add record to date base
                console.log('adding record');
            } else {
                break;
            }
        }
    }
}

async function fetchData(date) {
    let result = [];
    try {
        if(moment.duration(lastApiCall.diff(moment())).asMinutes() < 1) {
            //sleep for 1 second to make sure we dont make more api calls then 1 per second
            setTimeout(function(){
            }, 1000);
        }
    // await axios
    //     .get(apiurl)
    //     .then(response => (result = response.data));
    //     lastApiCall = moment();
    } catch (e) {
        console.log(e);
    }
    return result;
}