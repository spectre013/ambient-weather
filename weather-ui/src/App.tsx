import { useState, useEffect     } from 'react'
import './App.css'
import {Current} from "./models/current";
import { ForecastModel } from "./models/Forecast.ts"
import Temperature from "./components/Temperature.tsx";
import Forecast from "./components/Forecast.tsx";
import Stat from "./components/Stat.tsx";
import {
    baroLabel,
    distanceLabel, getOtherUnit,
    rainLabel,
    tempLabel, timeFormatAMPM,
    windLabel
} from "./util/weather.ts";
import GaugeComponent from 'react-gauge-component';
import moment from "moment/moment";
import AlertInfo from "./components/AlertInfo.tsx";
localStorage.setItem('units', 'imperial');
import {useNavigate} from "react-router-dom";

function App() {
    const [cLoaded, setCLoaded] = useState(false);
    const [fLoaded, setFLoaded] = useState(false);
    const [forecast, setForecast] = useState<ForecastModel>({} as ForecastModel);
    const [conditions, setConditions] = useState<Current>({} as Current);
    const [units, setUnits] = useState<string>("imperial");
    const [connected, setConnected] = useState(false);
    const navigate = useNavigate()
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
                setConnected(true);
            });
            connection.addEventListener('close', () => {
                console.log('Connection Close!');
                setConnected(false);
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
                setConnected(false);
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

        const updateUnit = () => {
            const newUnit = getOtherUnit(units).toLowerCase();
            setUnits(newUnit);
            localStorage.setItem('units', newUnit);
        };

        const about = () => {
            navigate('/about');
        };
        const lightdark = () => {
            const theme = localStorage.getItem('theme') || 'dark';
            const body: HTMLElement = document.body;
            if(theme === 'dark') {
                localStorage.setItem('theme', 'light');
                body.dataset['theme'] = 'light';
            } else {
                localStorage.setItem('theme', 'dark');
                body.dataset['theme'] = 'dark';
            }
        };

        function sunriseLabel(value: number): string {
            if(value === 0) {
                return "Rise";
            } else {
                return "Set";
            }
        }

        function isConnected(status: boolean): string {
            if(status) {
                return "material-symbols-sharp connected"
            } else {
                return "material-symbols-sharp disconnected"
            }
        }



        if (!cLoaded || !fLoaded) {
            return (
                <div className="loading-body">
                    <div className="loading-container">
                        <div>Lorson Ranch, Colorado Springs - Weather</div>
                        <div className="spinner"></div>
                        <p className="loading-text">Loading...</p>
                    </div>
                </div>
            )
        }


    return (
        <>
            <div className="dashboard">
                <div className="content">
                    <main className="main-content">
                        <div className="left-panel">
                        <Temperature temp={conditions.temp} icon={forecast.days[0].icon} conditions={forecast.days[0].conditions} units={units}/>
                                <div className="info-container">
                                    <div className="info-card" onClick={updateUnit}>
                                        Use {getOtherUnit(units)}
                                    </div>
                                    <div className="info-card" onClick={lightdark}>
                                        <span className="material-symbols-sharp">contrast</span>
                                    </div>
                                    <div className="info-card">
                                       <span className={isConnected(connected)}>wifi</span>&nbsp;Last Update: {timeFormatAMPM(conditions.date)}
                                    </div>
                                    <div className="info-card" onClick={about}>
                                        About
                                    </div>
                                </div>
                                <section className="alerts-section">
                                    <div className="alerts-info">
                                    <AlertInfo alerts={conditions.alert} />
                                    </div>
                                </section>
                                <Forecast forecast={forecast} units={units} />
                            </div>
                            <div className="right-panel">
                                <div className="stats-grid">
                                    <Stat icon="thermostat" valueType="temperature" label="Feels Like" value={conditions.temp.feelslike} units={"°"+tempLabel(units)} />
                                    <Stat icon="air" valueType="wind" label="Wind Speed" value={conditions.wind.windspeedmph} units={windLabel(units)} />
                                    <Stat icon="humidity_mid" valueType="humidity" label="Humidity" value={conditions.humidity.humidity} units="%" />
                                    <Stat icon="water_drop" valueType="barometer" label="Barometer" value={conditions.barometer.baromrelin} units={baroLabel(units)} />
                                    <Stat icon="thermostat" valueType="temperature" label="Dewpoint" value={conditions.humidity.dewpoint} units={"°"+tempLabel(units)} />
                                    <Stat icon="water_drop" valueType="rain" label="Precipitation" value={conditions.rain.daily} units={rainLabel(units)} />
                                    <Stat icon="bolt" valueType="lightning" label="Lightning Today" value={conditions.lightning.day} units="" />
                                    <Stat icon="bolt" valueType="lightning" label="Lightning Distance" value={conditions.lightning.distance} units={distanceLabel(units)} />
                                    <Stat icon="foggy" valueType="forecast" label="Visibility" value={forecast.days[0].visibility} units={distanceLabel(units)} />
                                    {/*<Stat icon="aq_indoor" valueType="" label="Air Quality" value={conditions.aqi.pm25} units="pm2.5" />*/}
                                </div>

                                <div className="gauges-grid">
                                    <div className="gauge-card uv-gauge">
                                        <div className="gauge-container">
                                            <GaugeComponent
                                                arc={{
                                                    subArcs: [
                                                        {
                                                            limit: 2,
                                                            color: '#5BE12C',
                                                            showTick: true
                                                        },
                                                        {
                                                            limit: 5,
                                                            color: '#F5CD19',
                                                            showTick: true
                                                        },
                                                        {
                                                            limit: 7,
                                                            color: '#F58B19',
                                                            showTick: true
                                                        },
                                                        {
                                                            limit: 10,
                                                            color: '#EA4228',
                                                            showTick: true
                                                        },
                                                        {
                                                            limit: 12,
                                                            color: 'violet',
                                                            showTick: true
                                                        },
                                                    ]
                                                }}
                                                labels={{
                                                    valueLabel: {
                                                        matchColorWithArc: true
                                                    }
                                                }}
                                                value={conditions.uv.uv}
                                                minValue={0}
                                                maxValue={12}
                                            />
                                        <div className="gauge-label">UV Index</div>
                                        </div>
                                    </div>
                                    <div className="gauge-card sunrise-gauge">
                                        <div className="gauge-container">
                                            <GaugeComponent
                                                type="semicircle"
                                                arc={{

                                                    gradient: true,
                                                    padding: 0.02,
                                                    width: 0.1,
                                                    subArcs:
                                                        [
                                                            { limit: 0, color: '#1E98D1'},
                                                            { limit: 25, color: '#CBE5F3'},
                                                            { limit: 50, color: '#F0E71A'},
                                                            { limit: 75, color: '#F5560C'},
                                                            { limit: 100, color: '#C57F51'},
                                                        ]
                                                }}
                                                pointer={{type: "blob",
                                                    animationDelay: 0,
                                                    color: '#ffd319',
                                                    baseColor: "#ffd319",
                                                    strokeWidth: 0,
                                                    }}
                                                labels={{
                                                    valueLabel: { formatTextValue: value => value + 'º',
                                                        matchColorWithArc: true},
                                                    tickLabels: {
                                                        type: "outer",
                                                        defaultTickValueConfig: {
                                                            formatTextValue: sunriseLabel,
                                                        }
                                                    }
                                                }}
                                                style={{ width: '100%', height: '100%' }}
                                                value={conditions.astro.elevation > 100 ? 100 : conditions.astro.elevation }
                                                minValue={0}
                                                maxValue={100}
                                            />
                                        </div>
                                        <div className="gauge-label">Sunrise { moment(conditions.astro.sunrise).format('LTS') } •
                                            Sunset { moment(conditions.astro.sunset).format('LTS')}</div>
                                    </div>
                                </div>

                                <div className="stats-grid">
                                    <Stat icon="thermostat" valueType="tempin" label="Living room" value={conditions.tempin.temp} units={"°"+tempLabel(units)} />
                                    <Stat icon="thermostat" valueType="tempin" label="Master Bedroom" value={conditions.temp2.temp} units={"°"+tempLabel(units)} />
                                    <Stat icon="thermostat" valueType="tempin" label="Hannah's room" value={conditions.temp3.temp} units={"°"+tempLabel(units)} />
                                    <Stat icon="thermostat" valueType="tempin" label="Basement" value={conditions.temp1.temp} units={"°"+tempLabel(units)} />
                                    <Stat icon="thermostat" valueType="tempin" label="Garage" value={conditions.temp4.temp} units={"°"+tempLabel(units)} />
                                </div>
                            </div>
                        </main>
                    </div>
                </div>
                <div className='copyright'>&copy; 2018-{moment().format('YYYY')} : zoms.net</div>
            </>
        )
}

export default App;
