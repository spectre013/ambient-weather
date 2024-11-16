import { useState, useEffect } from 'react'
import './App.css'
import {Current} from "./models/current";
import { MinMax } from "./models/minmax";
import { Wind as WindModel } from "./models/Wind";
import { ForecastModel } from "./models/Forecast.ts"
import Temperature  from "./components/Temperature.tsx";
import AlertInfo from "./components/AlertInfo.tsx";
import Forecast from "./components/Forecast.tsx";
import Wind from "./components/Wind.tsx"
import Tempin from "./components/Tempin.tsx";
import moment from "moment";
import Rain from "./components/Rain.tsx";
import Humidity from "./components/Humidity.tsx";
import Barometer from "./components/Barometer.tsx";
import {Trend} from "./models/Trend.ts";
import Sun from "./components/Sun.tsx";
import {Luna} from "./models/Luna.ts";
import Aqi from "./components/Aqi.tsx";
import Uv from "./components/Uv.tsx";
import {Alert} from "./models/Alert.ts";
import Lightning from "./components/Lightning.tsx";

function App() {
    const [minmax, setMinMax] = useState<MinMax>();
    const [alerts, setAlerts] = useState<Alert[]>([]);
    const [forecast, setForecast] = useState<ForecastModel>({} as ForecastModel);
    const [live, setLive] = useState<Current>({} as Current);
    const [units, setUnits] = useState<string>("imperial");
    const [wind, setWind] = useState<WindModel>({} as WindModel);
    const [barTrend, setBarTrend] = useState<Trend>({} as Trend);
    const [luna, setLuna] = useState<Luna>({} as Luna);

    const alertURL = "/api/alerts";
    const TempMinmaxURL = "/api/minmax/tempf";
    const forecastURL = "/api/forecast";
    const BarTrendURL = "/api/trend/baromrelin";
    const windURL = "/api/wind";
    const lunaURL = "/api/luna";


        useEffect(() => {
            setUnits("imperial")
            // waiting for allthethings in parallel
            let wsurl = 'wss://' + window.location.host;
            if (window.location.protocol === 'http:') {
                wsurl = 'ws://' + window.location.host;
            }
            wsurl += '/api/ws';
            let connection = new WebSocket(wsurl);

            connection.addEventListener('open', () => {
                console.log('Connection Open!');
            });
            connection.addEventListener('close', () => {
                console.log('Connection Close!');
                setTimeout(function () {
                    connection = new WebSocket(wsurl);
                }, 1000);
            });
            // Listen for messages
            connection.addEventListener('message', (event) => {
                setLive(JSON.parse(event.data));
            });

            connection.onerror = function (error) {
                console.log(`[error]`, error);
                connection.close();
            };


            const dataFetch = async () => {
                const result = (
                    await Promise.all([
                        fetch(alertURL),
                        fetch(TempMinmaxURL),
                        fetch(forecastURL),
                        fetch(windURL),
                        fetch(BarTrendURL),
                        fetch(lunaURL)
                    ])
                ).map((r) => r.json());

                const [alertResult, minmaxResult, forecastResult,windResult,
                        barTrendResult, lunaResults
                        ] = await Promise.all(
                    result
                );

                // when the data is ready, save it to state
                setAlerts(alertResult);
                setMinMax(minmaxResult);
                setForecast(forecastResult);
                setWind(windResult);
                setBarTrend(barTrendResult);
                setLuna(lunaResults)

            }
            dataFetch();
        }, []);
        console.log(live)
        if (!alerts || !minmax || !forecast || !wind || !Object.hasOwn(live, 'date')) return 'loading';

        return (
            <>
                <div className="header">
                    <div className="title"><i className="fa-solid fa-house"></i> Lorson Ranch, Colorado Springs, CO</div>
                    <div className="last-update">Last update:&nbsp;
                        <span className="update-time">{moment(live.date).format('HH:mm:ss')}</span>
                    </div>
                </div>
                <div className="container">
                    <div className="tempurature">
                        <Temperature live={live} icon={forecast.days[0].icon} units={units} avg={minmax as MinMax}/>
                    </div>
                    <div className="forecast">
                        <Forecast forecast={forecast} units={units}/>
                    </div>
                    <div className="alert">
                        <AlertInfo alerts={alerts}/>
                    </div>
                    <div className="wind">
                        <Wind live={live} wind={wind} units={units}/>
                    </div>
                    <div className="rain">
                        <Rain live={live} units={units}/>
                    </div>
                    <div className="lightning">
                        <Lightning live={live} units={units}/>
                    </div>
                    <div className="humidity">
                        <Humidity live={live} units={units}/>
                    </div>
                    <div className="baro">
                        <Barometer live={live} trend={barTrend} units={units}/>
                    </div>
                    <div className="sun">
                        <Sun luna={luna}/>
                    </div>
                    <div className="uv">
                        <Uv live={live} luna={luna} units={units}/>
                    </div>
                    <div className="aq">
                        <Aqi live={live} units={units}/>
                    </div>
                    <div className="living">
                        <Tempin live={live} sensor="in" area="Living" units={units}/>
                    </div>
                    <div className="master">
                        <Tempin live={live} sensor="2" area="Master" units={units}/>
                    </div>
                    <div className="office">
                        <Tempin live={live} sensor="3" area="Office" units={units}/>
                    </div>
                    <div className="basement">
                        <Tempin live={live} sensor="1" area="Basement" units={units}/>
                    </div>
                </div>

            </>
        )
}

export default App;
