const mysql = require('mysql2/promise');
const moment = require('moment');
console.log("WILL NOT RUN AGAIN!! ")
process.exit()


if (process.env.NODE_ENV !== 'production') {
    require('dotenv').config();
}

const database_config = {
    host     : process.env.DB_HOST,
    user     : process.env.DB_USER,
    password : process.env.DB_PASSWORD,
    database : process.env.DB_DATABASE
};

let connection = {};
async function getConnection()
{
    await mysql.createConnection(database_config).then(conn => {
        connection = conn;
    });
}

getConnection().then(() => {
    fetch().then((data) => {
        update(data).then();
    });
});

async function fetch() {
    let data = [];
    try {
        const query  = 'select id, date from records order by id desc';
        await connection.query(query)
            .then(([rows]) => {
                for(let i=0; i < rows.length; i++) {
                    //console.log(rows[i].id, rows[i].date, moment(rows[i].date).format("YYYY-MM-DD HH:mm:ss"))
                    data.push({id:rows[i].id,date:moment(rows[i].date).format("YYYY-MM-DD HH:mm:ss")})
                }
            });
    } catch (e) {
        console.log('Error Updating Records', e);
    }
    return data
}

function update(data) {
    let q = 'update records set date = ? where id = ?';
    for(let i=0; i < data.length; i++) {
        connection.execute(q, [moment.utc(data[i].date).local().format("YYYY-MM-DD HH:mm:ss"), data[i].id])
            .then((r) => {
                console.log(r);
            })
            .catch(error => {
                console.log(error.message)
            });
    }
}


// let warnings = `Tsunami Warning | 1 |   | Tomato | 253 99 71 | FD6347
// Tornado Warning | 2 |   | Red | 255 0 0 | FF0000
// Extreme Wind Warning | 3 |   | Darkorange | 255 140 0 | FF8C00
// Severe Thunderstorm Warning | 4 |   | Orange | 255 165 0 | FFA500
// Flash Flood Warning | 5 |   | Darkred | 139 0 0 | 8B0000
// Flash Flood Statement | 6 |   | Darkred | 139 0 0 | 8B0000
// Severe Weather Statement | 7 |   | Aqua | 0 255 255 | 00FFFF
// Shelter In Place Warning | 8 |   | Salmon | 250 128 114 | FA8072
// Evacuation - Immediate | 9 |   | Chartreuse | 127 255 0 | 7FFF00
// Civil Danger Warning | 10 |   | Lightpink | 255 182 193 | FFB6C1
// Nuclear Power Plant Warning | 11 |   | Indigo | 75 0 130 | 4B0082
// Radiological Hazard Warning | 12 |   | Indigo | 75 0 130 | 4B0082
// Hazardous Materials Warning | 13 |   | Indigo | 75 0 130 | 4B0082
// Fire Warning | 14 |   | Sienna | 160 82 45 | A0522D
// Civil Emergency Message | 15 |   | Lightpink | 255 182 193 | FFB6C1
// Law Enforcement Warning | 16 |   | Silver | 192 192 192 | C0C0C0
// Storm Surge Warning | 17 |   | SSWarning | 181 36 247 | B524F7
// Hurricane Force Wind Warning | 18 |   | Westernred | 205 92 92 | CD5C5C
// Hurricane Warning | 19 |   | Crimson | 220 20 60 | DC143C
// Typhoon Warning | 20 |   | Crimson | 220 20 60 | DC143C
// Special Marine Warning | 21 |   | Orange | 255 165 0 | FFA500
// Blizzard Warning | 22 |   | Orangered | 255 69 0 | FF4500
// Snow Squall Warning | 23 |   | Mediumvioletred | 199 21 133 | C71585
// Ice Storm Warning | 24 |   | Darkmagenta | 139 0 139 | 8B008B
// Winter Storm Warning | 25 |   | Hotpink | 255 105 180 | FF69B4
// High Wind Warning | 26 |   | Goldenrod | 218 165 32 | DAA520
// Tropical Storm Warning | 27 |   | Firebrick | 178 34 34 | B22222
// Storm Warning | 28 |   | Darkviolet | 148 0 211 | 9400D3
// Tsunami Advisory | 29 |   | Chocolate | 210 105 30 | D2691E
// Tsunami Watch | 30 |   | Fushsia | 255 0 255 | FF00FF
// Avalanche Warning | 31 |   | Dodgerblue | 30 144 255 | 1E90FF
// Earthquake Warning | 32 |   | Saddlebrown | 139 69 19 | 8B4513
// Volcano Warning | 33 |   | darkslategray | 47 79 79 | 2F4F4F
// Ashfall Warning | 34 |   | Darkgray | 169 169 169 | A9A9A9
// Coastal Flood Warning | 35 |   | Forestgreen | 34 139 34 | 228B22
// Lakeshore Flood Warning | 36 |   | Forestgreen | 34 139 34 | 228B22
// Flood Warning | 37 |   | Lime | 0 255 0 | 00FF00
// High Surf Warning | 38 |   | Forestgreen | 34 139 34 | 228B22
// Dust Storm Warning | 39 |   | Bisque | 255 228 196 | FFE4C4
// Blowing Dust Warning | 40 |   | Bisque | 255 228 196 | FFE4C4
// Lake Effect Snow Warning | 41 |   | Darkcyan | 0 139 139 | 008B8B
// Excessive Heat Warning | 42 |   | Mediumvioletred | 199 21 133 | C71585
// Tornado Watch | 43 |   | Yellow | 255 255 0 | FFFF00
// Severe Thunderstorm Watch | 44 |   | Palevioletred | 219 112 147 | DB7093
// Flash Flood Watch | 45 |   | Seagreen | 46 139 87 | 2E8B57
// Gale Warning | 46 |   | Plum | 221 160 221 | DDA0DD
// Flood Statement | 47 |   | Lime | 0 255 0 | 00FF00
// Wind Chill Warning | 48 |   | Lightsteelblue | 176 196 222 | B0C4DE
// Extreme Cold Warning | 49 |   | Blue | 0 0 255 | 0000FF
// Hard Freeze Warning | 50 |   | Darkviolet | 148 0 211 | 9400D3
// Freeze Warning | 51 |   | Darkslateblue | 72 61 139 | 483D8B
// Red Flag Warning | 52 |   | Deeppink | 255 20 147 | FF1493
// Storm Surge Watch | 53 |   | SSWatch | 219 127 247 | DB7FF7
// Hurricane Watch | 54 |   | Magenta | 255 0 255 | FF00FF
// Hurricane Force Wind Watch | 55 |   | Darkorchid | 153 50 204 | 9932CC
// Typhoon Watch | 56 |   | Magenta | 255 0 255 | FF00FF
// Tropical Storm Watch | 57 |   | Lightcoral | 240 128 128 | F08080
// Storm Watch | 58 |   | Moccasin | 255 228 181 | FFE4B5
// Hurricane Local Statement | 59 |   | Moccasin | 255 228 181 | FFE4B5
// Typhoon Local Statement | 60 |   | Moccasin | 255 228 181 | FFE4B5
// Tropical Storm Local Statement | 61 |   | Moccasin | 255 228 181 | FFE4B5
// Tropical Depression Local Statement | 62 |   | Moccasin | 255 228 181 | FFE4B5
// Avalanche Advisory | 63 |   | Peru | 205 133 63 | CD853F
// Freezing Rain Advisory | 64 |   | Orchid | 218 112 214 | DA70D6
// Winter Weather Advisory | 65 |   | Mediumslateblue | 123 104 238 | 7B68EE
// Lake Effect Snow Advisory | 66 |   | Mediumturquoise | 72 209 204 | 48D1CC
// Wind Chill Advisory | 67 |   | Paleturquoise | 175 238 238 | AFEEEE
// Heat Advisory | 68 |   | Coral | 255 127 80 | FF7F50
// Urban And Small Stream Flood Advisory | 69 |   | Springgreen | 0 255 127 | 00FF7F
// Small Stream Flood Advisory | 70 |   | Springgreen | 0 255 127 | 00FF7F
// Arroyo And Small Stream Flood Advisory | 71 |   | Springgreen | 0 255 127 | 00FF7F
// Flood Advisory | 72 |   | Springgreen | 0 255 127 | 00FF7F
// Hydrologic Advisory | 73 |   | Springgreen | 0 255 127 | 00FF7F
// Lakeshore Flood Advisory | 74 |   | Lawngreen | 124 252 0 | 7CFC00
// Coastal Flood Advisory | 75 |   | Lawngreen | 124 252 0 | 7CFC00
// High Surf Advisory | 76 |   | Mediumorchid | 186 85 211 | BA55D3
// Heavy Freezing Spray Warning | 77 |   | Deepskyblue | 0 191 255 | 00BFFF
// Dense Fog Advisory | 78 |   | Slategray | 112 128 144 | 708090
// Dense Smoke Advisory | 79 |   | Khaki | 240 230 140 | F0E68C
// Small Craft Advisory For Hazardous Seas | 80 |   | Thistle | 216 191 216 | D8BFD8
// Small Craft Advisory For Rough Bar | 81 |   | Thistle | 216 191 216 | D8BFD8
// Small Craft Advisory For Winds | 82 |   | Thistle | 216 191 216 | D8BFD8
// Small Craft Advisory | 83 |   | Thistle | 216 191 216 | D8BFD8
// Brisk Wind Advisory | 84 |   | Thistle | 216 191 216 | D8BFD8
// Hazardous Seas Warning | 85 |   | Thistle | 216 191 216 | D8BFD8
// Dust Advisory | 86 |   | Darkkhaki | 189 183 107 | BDB76B
// Blowing Dust Advisory | 87 |   | Darkkhaki | 189 183 107 | BDB76B
// Lake Wind Advisory | 88 |   | Tan | 210 180 140 | D2B48C
// Wind Advisory | 89 |   | Tan | 210 180 140 | D2B48C
// Frost Advisory | 90 |   | Cornflowerblue | 100 149 237 | 6495ED
// Ashfall Advisory | 91 |   | Dimgray | 105 105 105 | 696969
// Freezing Fog Advisory | 92 |   | Teal | 0 128 128 | 008080
// Freezing Spray Advisory | 93 |   | Deepskyblue | 0 191 255 | 00BFFF
// Low Water Advisory | 94 |   | Brown | 165 42 42 | A52A2A
// Local Area Emergency | 95 |   | Silver | 192 192 192 | C0C0C0
// Child Abduction Emergency | 96 |   | Gold | 255 215 0 | FFD700
// Avalanche Watch | 97 |   | Sandybrown | 244 164 96 | F4A460
// Blizzard Watch | 98 |   | Greenyellow | 173 255 47 | ADFF2F
// Rip Current Statement | 99 |   | Turquoise | 64 224 208 | 40E0D0
// Beach Hazards Statement | 100 |   | Turquoise | 64 224 208 | 40E0D0
// Gale Watch | 101 |   | Pink | 255 192 203 | FFC0CB
// Winter Storm Watch | 102 |   | Steelblue | 70 130 180 | 4682B4
// Hazardous Seas Watch | 103 |   | Darkslateblue | 72 61 139 | 483D8B
// Heavy Freezing Spray Watch | 104 |   | Rosybrown | 188 143 143 | BC8F8F
// Coastal Flood Watch | 105 |   | Mediumaquamarine | 102 205 170 | 66CDAA
// Lakeshore Flood Watch | 106 |   | Mediumaquamarine | 102 205 170 | 66CDAA
// Flood Watch | 107 |   | Seagreen | 46 139 87 | 2E8B57
// High Wind Watch | 108 |   | Darkgoldenrod | 184 134 11 | B8860B
// Excessive Heat Watch | 109 |   | Maroon | 128 0 0 | 800000
// Extreme Cold Watch | 110 |   | Blue | 0 0 255 | 0000FF
// Wind Chill Watch | 111 |   | Cadetblue | 95 158 160 | 5F9EA0
// Lake Effect Snow Watch | 112 |   | Lightskyblue | 135 206 250 | 87CEFA
// Hard Freeze Watch | 113 |   | Royalblue | 65 105 225 | 4169E1
// Freeze Watch | 114 |   | Cyan | 0 255 255 | 00FFFF
// Fire Weather Watch | 115 |   | Navajowhite | 255 222 173 | FFDEAD
// Extreme Fire Danger | 116 |   | Darksalmon | 233 150 122 | E9967A
// 911 Telephone Outage | 117 |   | Silver | 192 192 192 | C0C0C0
// Coastal Flood Statement | 118 |   | Olivedrab | 107 142 35 | 6B8E23
// Lakeshore Flood Statement | 119 |   | Olivedrab | 107 142 35 | 6B8E23
// Special Weather Statement | 120 |   | Moccasin | 255 228 181 | FFE4B5
// Marine Weather Statement | 121 |   | Peachpuff | 255 239 213 | FFDAB9
// Air Quality Alert | 122 |   | Gray | 128 128 128 | 808080
// Air Stagnation Advisory | 123 |   | Gray | 128 128 128 | 808080
// Hazardous Weather Outlook | 124 |   | Palegoldenrod | 238 232 170 | EEE8AA
// Hydrologic Outlook | 125 |   | Lightgreen | 144 238 144 | 90EE90
// Short Term Forecast | 126 |   | Palegreen | 152 251 152 | 98FB98
// Administrative Message | 127 |   | Silver | 192 192 192 | C0C0C0
// Test | 128 |   | Azure | 240 255 255 | F0FFFF`;
//
// let lines = warnings.split("\n");
// lines.forEach((line) => {
//     let alert = line.split(" | ");
//     if(alert[0].startsWith("911")) {
//         alert[0] = "Telephone Outage 911";
//     }
//     console.log(`.${alert[0].toLowerCase().replace(/\s+/g, "-")} {
//     color: #${alert[5]};
// }`)
// });