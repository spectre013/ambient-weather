import { useState, useEffect } from 'react'
import './App.css'
import {Current} from "./models/current";
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
import Sun from "./components/Sun.tsx";
import Aqi from "./components/Aqi.tsx";
import Uv from "./components/Uv.tsx";
import Lightning from "./components/Lightning.tsx";

function App() {
    const [cLoaded, setCLoaded] = useState(false);
    const [fLoaded, setFLoaded] = useState(false);
    const [forecast, setForecast] = useState<ForecastModel>({} as ForecastModel);
    const [conditions, setConditions] = useState<Current>({} as Current);
    const [units, setUnits] = useState<string>("imperial");

    const forecastURL = "/api/forecast";


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
                setConditions(JSON.parse(event.data));
                setCLoaded(true);
            });

            connection.onerror = function (error) {
                console.log(`[error]`, error);
                connection.close();
            };


            const dataFetch = async () => {
                const result = (
                    await Promise.all([
                        fetch(forecastURL),
                    ])
                ).map((r) => r.json());

                const [forecastResult] = await Promise.all(
                    result
                );
                // when the data is ready, save it to state
                setForecast(forecastResult);
                setFLoaded(true);
            }


            dataFetch().then(() => {
            });
        }, []);

        if (!cLoaded || !fLoaded) {
            return 'loading';
        }

        return (
            <>
                <div className="header">
                    <div className="title"><i className="fa-solid fa-house"></i> Lorson Ranch, Colorado Springs, CO</div>
                    <div className="last-update">Last update:&nbsp;
                        <span className="update-time">{moment(conditions.date).format('HH:mm:ss')}</span>
                    </div>
                </div>
                <div className="container">
                    <div className="tempurature">
                        <Temperature temp={conditions.temp} icon={forecast.days[0].icon} units={units} />
                    </div>
                    <div className="forecast">
                        <Forecast forecast={forecast} units={units}/>
                    </div>
                    <div className="alert">
                        <AlertInfo alerts={conditions.alert}/>
                    </div>
                    <div className="wind">
                        <Wind wind={conditions.wind} units={units}/>
                    </div>
                    <div className="rain">
                        <Rain rain={conditions.rain} units={units}/>
                    </div>
                    <div className="lightning">
                        <Lightning lightning={conditions.lightning} date={conditions.date} units={units}/>
                    </div>
                    <div className="humidity">
                        <Humidity humidity={conditions.humidity} units={units}/>
                    </div>
                    <div className="baro">
                        <Barometer baro={conditions.barometer} units={units}/>
                    </div>
                    <div className="sun">
                        <Sun astro={conditions.astro} units={units}/>
                    </div>
                    <div className="uv">
                        <Uv uv={conditions.uv} astro={conditions.astro} units={units}/>
                    </div>
                    <div className="aq">
                        <Aqi aqi={conditions.aqi} units={units}/>
                    </div>
                    <div className="living">
                        <Tempin temp={conditions.tempin} area="Living" units={units}/>
                    </div>
                    <div className="master">
                        <Tempin temp={conditions.temp1} area="Master" units={units}/>
                    </div>
                    <div className="office">
                        <Tempin temp={conditions.temp2} area="Office" units={units}/>
                    </div>
                    <div className="basement">
                        <Tempin temp={conditions.temp3} area="Basement" units={units}/>
                    </div>
                    <div className="garage">
                        <Tempin temp={conditions.temp4} area="Garage" units={units}/>
                    </div>
                </div>

            </>
        )
}

export default App;
