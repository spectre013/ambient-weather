import './Temperature.css'
import {useNavigate} from "react-router-dom";
import Guage from "../../util/TemperatureGuage.tsx";
import {tempColor, tempLabel, tempToHex} from "../../util/weather.ts";
import {ChartData} from "../../models/DataSeries.ts";
import {useEffect, useState} from "react";
import { TempData } from "../../models/current.ts";
import {
    LineChart,
    Line,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
    ResponsiveContainer
} from 'recharts';
import {processData} from "../../util/weather.ts";

const Temperature = () => {
    const navigate = useNavigate();
    const [fLoaded, setFLoaded] = useState(false);
    const [units, setUnits] = useState<string>("imperial");
    const [temp, setTemp] = useState<TempData>({} as TempData);
    const [chart, setChart] = useState<ChartData>({} as ChartData);
    const urls = [
        '/api/current',
        '/api/chart/temperature/3h',
    ];

    useEffect(() => {
        setUnits('imperial'); // Set the units to imperial by default
        const fetchPromises = urls.map(url => fetch(url));
        Promise.all(fetchPromises)
            .then(responses => {
                // 'responses' will be an array of Response objects
                // Process each response to extract JSON data
                return Promise.all(responses.map(response => response.json()));
            })
            .then(data => {
                // 'data' will be an array containing the JSON data from each fetch
                setTemp(data[0].temp);
                setChart(data[1])
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

    const handleClick = () => {
        navigate('/'); // Navigate to the details page for the specific stat
    }

    const combinedChartData = processData(chart);

    return (
        <>
            <div className="details-dashboard">
                <header className="details-header">
                    <h1><span className="material-symbols-sharp">thermostat</span> Losron Ranch -  Temperature</h1>
                    <div className="hasclick" onClick={handleClick}><span className="material-symbols-sharp">home</span></div>
                </header>
                <div className="content">
                    <div className="details-content">
                        <div className="details">
                                <div className="detail-item">
                                    <div className="label">Current</div>
                                    <div className="value">
                                        <span className={tempColor(temp.temp)}>{temp.temp}</span>&deg;{tempLabel(units)}</div>
                                </div>
                                <div className="detail-item">
                                    <div className="label">Feels Like</div>
                                    <div className="value">
                                    <span className={tempColor(temp.feelslike)}>{ temp.feelslike }</span>&deg;{tempLabel(units)}</div>
                            </div>
                            <div className="detail-item">
                                <div className="label">High</div>
                                <div className="value">
                                    <span className={tempColor(temp.minmax.max.day.value)}>{temp.minmax.max.day.value }</span>&deg;{tempLabel(units)}</div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Low</div>
                                <div className="value">
                                    <span className={tempColor(temp.minmax.min.day.value)}>{ temp.minmax.min.day.value }</span>&deg;{tempLabel(units)}</div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Humidity</div>
                                <div className="value">
                                    <span>{ temp.humidity }</span>%</div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Dewpoint</div>
                                <div className="value">
                                    <span className={tempColor(temp.dewpoint)}>{ temp.dewpoint }</span>&deg;{tempLabel(units)}</div>
                            </div>
                        </div>
                        <div className="dial">
                            <Guage startColor={tempToHex(temp.minmax.min.day.value)}
                                   endColor={tempToHex(temp.minmax.max.day.value)}
                                   currentColor={tempToHex(temp.temp)}
                                   size={275}
                                   min={Number(temp.minmax.min.day.value.toFixed(0))}
                                   max={Number(temp.minmax.max.day.value.toFixed(0))}
                                   value={Number(temp.temp.toFixed(0))} />
                        </div>
                    </div>
                    <div className="chart">
                        <h3>Last 3 hours</h3>
                        <ResponsiveContainer width="100%" height="100%">
                            <LineChart
                                data={combinedChartData}
                                margin={{
                                    top: 5,
                                    right: 30,
                                    left: 20,
                                    bottom: 5,
                                }}
                            >
                                <CartesianGrid strokeDasharray="3 3" stroke="#4a5568" />
                                <XAxis
                                    dataKey="date"
                                    stroke="#9ca3af"
                                    tickFormatter={(date) => date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
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
                                    dataKey="Temperature"
                                    stroke="#ef4444"
                                    activeDot={{ r: 8 }}
                                    dot={false}
                                />
                            </LineChart>
                        </ResponsiveContainer>
                    </div>
                </div>
            </div>
        </>
    )
}
export default Temperature