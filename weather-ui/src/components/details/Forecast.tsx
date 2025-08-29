import './Forecast.css'
import {useEffect, useState} from "react";
import {Day} from "../../models/Forecast.ts";
import {useParams} from "react-router-dom";
import Guage from "../../util/TemperatureGuage.tsx";
import {
    BeaufortHex,
    dewPointClass,
    distanceLabel,
    tempColor,
    tempLabel,
    tempToHex,
    createDial, windLabel, degToCompass, rainLabel
} from "../../util/weather.ts";
import moment from "moment";
import GaugeComponent from "react-gauge-component";
import {
    Bar,
    BarChart,
    CartesianGrid,
    Legend,
    Line,
    LineChart,
    ResponsiveContainer,
    Tooltip,
    XAxis,
    YAxis
} from "recharts";
import {CustomForecastTooltip } from "../../util/utilities.tsx";
import Header from "../Header.tsx";

const Forecast = () => {
    const { day } = useParams()
    const [fLoaded, setFLoaded] = useState(false);
    const [units, setUnits] = useState<string>(localStorage.getItem('units') || 'imperial');
    const [forecast, setForecast] = useState<Day>({} as Day);
    const urls = [
        '/api/forecast',
    ];

    useEffect(() => {
        setUnits(localStorage.getItem('units') || 'imperial');
        const fetchPromises = urls.map(url => fetch(url));
        Promise.all(fetchPromises)
            .then(responses => {
                // 'responses' will be an array of Response objects
                // Process each response to extract JSON data
                return Promise.all(responses.map(response => response.json()));
            })
            .then(data => {
                // 'data' will be an array containing the JSON data from each fetch
                console.log(day, data);
                const forecastData: Day = data[0].days[parseInt(day || '0')];
                setForecast(forecastData);
                setFLoaded(true);
            })
            .catch(error => {
                // Handle any errors that occurred during fetching or JSON parsing
                console.error('Error fetching data:', error);
            });
    },[]);

    if (!fLoaded) {
        return 'loading';
    }

    function sunLabel(value: number): string {
        if (value === 0) {
            return "Rise";
        } else {
            return "Set";
        }
    }

    return (
        <>
        <div className="forecast-details-dashboard">
            <Header />
            <div className="forecast-details-content">
                <div className="left">
                    <div>Forecast for {moment(forecast.datetime).format('dddd, MMMM Do YYYY')}</div>
                    <div className="forecast-info">
                        <div className="weather-icon">
                            <div>
                                <img alt="" src={'/images/icons/'+forecast.icon+'.png'} />
                            </div>
                            <div>
                                <h3>{forecast.conditions}</h3>
                                <div>{ forecast.description}</div>
                            </div>
                        </div>
                    </div>
                    <div className="details">
                        <div className="detail-item">
                            <div className="label">High: </div>
                            <div className="value">
                                <span className={tempColor(forecast.tempmax)}>&nbsp;{ forecast.tempmax }</span>&deg;{tempLabel(units)}</div>
                        </div>
                        <div className="detail-item">
                            <div className="label">Low:</div>
                            <div className="value">
                                <span className={tempColor(forecast.tempmin)}>&nbsp;{ forecast.tempmin}</span>&deg;{tempLabel(units)}</div>
                        </div>
                        <div className="detail-item">
                            <div className="label">Dewpoint: </div>
                            <div className="value">
                                <span className={dewPointClass(forecast.dew)}>&nbsp;{ forecast.dew }</span>&deg;{tempLabel(units)}</div>
                        </div>
                        <div className="detail-item">
                            <div className="label">Humidity</div>
                            <div className="value">
                                <span>{ forecast.humidity }</span>%</div>
                        </div>
                        <div className="detail-item">
                            <div className="label">Visibility</div>
                            <div className="value">
                                <span>{ forecast.visibility }</span>&nbsp;{ distanceLabel(units) }</div>
                        </div>
                        <div className="detail-item">
                            <div className="label">Wind Speed:</div>
                            <div className="value">
                                <span>{ forecast.windspeed }</span>&nbsp;{ windLabel(units)}</div>
                        </div>
                        <div className="detail-item">
                            <div className="label">Wind Direction</div>
                            <div className="value">
                                <span>{ forecast.winddir }</span>&nbsp;{ degToCompass(forecast.winddir)}</div>
                        </div>
                        <div className="detail-item">
                            <div className="label">Wind Gusts: </div>
                            <div className="value">
                                <span>{ forecast.windgust }</span>&nbsp;{ windLabel(units)}</div>
                        </div>
                        <div className="detail-item">
                            <div className="label">Precipitation: </div>
                            <div className="value">
                                <span>{ forecast.precip }</span>&nbsp;{ rainLabel(units)}</div>
                        </div>
                        <div className="detail-item">
                            <div className="label">Precipitation Prob: </div>
                            <div className="value">
                                <span>{ forecast.precipprob }</span>%</div>
                        </div>
                    </div>
                </div>
                <div className="right">
                    <div className="gauges-grid">
                        <div className="gauge-card dial">
                            <Guage startColor={tempToHex(forecast.tempmin)}
                                   endColor={tempToHex(forecast.tempmax)}
                                   currentColor={tempToHex(forecast.temp)}
                                   size={250}
                                   min={Number(forecast.tempmin.toFixed(0))}
                                   max={Number(forecast.tempmax.toFixed(0))}
                                   value={Number(forecast.temp.toFixed(0))} />
                            <div className="gauge-label">Temperature</div>
                        </div>
                        <div className="gauge-card dial">
                            <div className="dial"
                                 dangerouslySetInnerHTML={{__html: createDial(forecast.windspeed, forecast.winddir, forecast.windgust,
                                         BeaufortHex(forecast.windspeed))}}>
                            </div>
                            <div className="gauge-label">Wind information</div>
                        </div>
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
                                    value={forecast.uvindex}
                                    minValue={0}
                                    maxValue={12}
                                />
                                <div className="gauge-label">UV Index</div>
                            </div>
                        </div>
                        <div className="gauge-gauge">
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
                                            valueLabel: { formatTextValue: value => value + 'º' },
                                            tickLabels: {
                                                type: "outer",
                                                defaultTickValueConfig: {
                                                    formatTextValue: sunLabel,
                                                }
                                            }
                                        }}
                                        style={{ width: '100%', height: '100%' }}
                                        value={50}
                                        minValue={0}
                                        maxValue={100}
                                    />
                                </div>
                                <div className="gauge-label">Sunrise { forecast.sunrise } •
                                    Sunset { forecast.sunset}</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div className="charts">
                <div className="bar-chart">
                    <h3>Temperature for {moment(forecast.datetime).format('dddd, MMMM Do YYYY')}</h3>
                    <ResponsiveContainer width="100%" height="100%">
                        <LineChart
                            data={forecast.hours}
                            margin={{
                                top: 5,
                                right: 30,
                                left: 20,
                                bottom: 5,
                            }}
                        >
                            <CartesianGrid strokeDasharray="3 3" stroke="#4a5568" />
                            <XAxis
                                dataKey="datetimeEpoch"
                                stroke="#9ca3af"
                                tickFormatter={(date) => moment.unix(date).format('HH:mm')}
                                label={{ value: 'Time of Day', position: 'insideBottom', offset: -15, fill: '#e5e7eb' }}
                            />
                            <YAxis
                                stroke="#9ca3af"
                                label={{ value: 'Degrees (°'+tempLabel(units)+')', angle: -90, position: 'insideLeft', offset: 15, fill: '#e5e7eb' }}
                            />
                            <Tooltip
                                contentStyle={{ backgroundColor: '#2d3748', border: 'none' }}
                                labelFormatter={(label) => new Date(label).toLocaleTimeString()}
                            />
                            <Legend wrapperStyle={{ color: '#e5e7eb',position: 'relative' }} />
                            <Line
                                type="monotone"
                                dataKey="Dewpoint"
                                stroke="#60a5fa"
                                activeDot={{ r: 8 }}
                                dot={false}
                            />
                            <Line
                                type="monotone"
                                dataKey="temp"
                                stroke="#ef4444"
                                activeDot={{ r: 8 }}
                                dot={false}
                            />
                        </LineChart>
                    </ResponsiveContainer>
                </div>
                <div className="bar-chart">
                    <h3>Precip for {moment(forecast.datetime).format('dddd, MMMM Do YYYY')}</h3>
                    <ResponsiveContainer width="100%" height="100%">
                        <BarChart
                            data={forecast.hours}
                            margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
                        >
                            <CartesianGrid stroke="#ccc" />
                            <XAxis
                                dataKey="datetime"
                                stroke="#ffffff"
                                tick={{ fill: '#ffffff', fontSize: 12 }}
                                axisLine={false}
                            />
                            <YAxis
                                stroke="#ffffff"
                                tick={{ fill: '#ffffff', fontSize: 12 }}
                                axisLine={false}
                            />
                            <Tooltip
                                cursor={{ fill: '#4b5563', opacity: 0.5 }}
                                content={CustomForecastTooltip}
                            />
                            <Bar
                                dataKey="precip"
                                fill="#3b9cac"
                                name={rainLabel(units)}
                            />
                            <Legend />
                        </BarChart>
                    </ResponsiveContainer>
                </div>
            </div>
        </div>
            <div className='copyright'>&copy; 2018-{moment().format('YYYY')} : zoms.net</div>
        </>
    )
}


export default Forecast;