require('dotenv').config();
const moment = require('moment');
const mysql = require('mysql2/promise');

const database_config = {
    host     : process.env.DB_HOST,
    user     : process.env.DB_USER,
    password : process.env.DB_PASSWORD,
    database : process.env.DB_DATABASE,
    multipleStatements: true
};


let connection = {};
async function getConnection()
{
    await mysql.createConnection(database_config).then(conn => {
        connection = conn;
    });
}

async function updateStatistics() {
    await getConnection()
    let data = {};
    let queries = {};
    const types = ['MAX','MIN','AVG'];
    const periods = ['day','yesterday','month','year'];
    const fields = ['tempf','tempinf','temp1f','temp2f','baromrelin','uv','humidity','windspeedmph','windgustmph','dewpoint'];

    try {

        await asyncForEach(periods,async (period) => {
            await asyncForEach(types,async (type) => {
                await asyncForEach(fields,(field) => {
                    console.log('`'+period+'_'+type.toLowerCase()+'_'+field+'_value` double DEFAULT NULL,');
                    console.log('`'+period+'_'+type.toLowerCase()+'_'+field+'_date` datetime DEFAULT NULL,');
                });
            });
        });
        // periods.forEach((period) => {
        //     types.forEach((type) => {
        //         fields.forEach((field) => {
        //             let key = period+'_'+type.toLowerCase()+'_'+field;
        //
        //             let order = '';
        //             if(type === 'MAX') {
        //                 order = ' DESC';
        //             }
        //
        //             let q = `select ${type}(${field})  as value, \`date\` from records where \`date\` between ? and ? group by \`date\`,${field} order by ${field}${order} limit 0,1`;
        //             if(type === 'AVG') {
        //                 q = `select CAST(${type}(${field}) AS DECIMAL(10,2)) as value from records where \`date\` between ? and ? order by ${field} limit 0,1`;
        //             }
        //             queries[key] = {
        //                 query:q,
        //                 params: getTimeframe(period)
        //             };
        //         });
        //     });
        // });
        // //console.log(queries);
        // let name = '';
        // let qt = '';
        // await asyncForEach( Object.keys(queries),async (key) => {
        //     await asyncForEach(Object.keys(queries[key]),async (tf) => {
        //         name = key;
        //         qt = tf;
        //         let query = queries[key];
        //         await connection.query(query.query,query.params , (err) => {
        //             if (err) throw err;
        //         }).then(async (rows) => {
        //             if(rows.length > 0) {
        //                 //data[name][qt] = rows[0][0];
        //                 try {
        //
        //                     let value = key+'_value';
        //                     let date = ', '+key+'_date = ?';
        //                     let insertData = [];
        //                     insertData.push(rows[0][0].value)
        //                     if(key.includes('avg')) {
        //                         date = '';
        //                     } else {
        //                         insertData.push(rows[0][0].date)
        //                     }
        //
        //                     const update  = `UPDATE stats SET ${value} = ?${date} WHERE id = 1`;
        //                     console.log(update);
        //                     await connection.query(update, insertData)
        //                         .then(([rows]) => {
        //                             if(rows.length > 0) {
        //                                 data = rows[0];
        //                             } else {
        //                                 data = {};
        //                             }
        //                         });
        //                 } catch (e) {
        //                     console.log(e);
        //                 }
        //             }
        //         });
        //     });
        // });
    } catch (e) {
        console.log(e);
    }

}

updateStatistics()

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

async function asyncForEach(array, callback) {
    for (let index = 0; index < array.length; index++) {
        await callback(array[index], index, array);
    }
}



