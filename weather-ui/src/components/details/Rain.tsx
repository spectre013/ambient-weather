import './Rain.css'
import {rainLabel, rainDisplay, full, formatDay} from "../../util/weather.ts";
import {useEffect, useState} from "react";
import {RainData} from "../../models/current.ts";
import {ChartData} from "../../models/DataSeries.ts";
import RainGuage from "../../util/RainGauge.tsx";
import { CustomTooltip } from "../../util/utilities.tsx";
import * as weather from "../../util/weather.ts";
import moment from "moment/moment";
import {
    ResponsiveContainer,
    BarChart,
    CartesianGrid,
    XAxis,
    YAxis,
    Tooltip,
    Bar, Legend,
} from 'recharts';
import Header from "../Header.tsx";

const Rain = () => {
    const [fLoaded, setFLoaded] = useState(false);
    const [units, setUnits] = useState<string>("imperial");
    const [rain, setRain] = useState<RainData>({} as RainData);
    const [chart, setChart] = useState<ChartData>({} as ChartData);
    const urls = [
        '/api/current',
        '/api/chart/rain/3h',
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
                setRain(data[0].rain);
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

    function year(date: string) {
        return moment(date).format('YYYY');
    }
    function month(date: string) {
        return moment(date).format('MMM');
    }

    return (
        <>
            <div className="details-dashboard">
                <Header />
                <div className="content">
                    <div className="details-content">
                        <div className="details">
                            <div className="detail-item">
                                <div className="label">Current:</div>
                                <div className="value">
                                    <span className="rain-blue">{ rainDisplay(rain.daily, units) }</span>&nbsp;{rainLabel(units)}
                                </div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Rate:</div>
                                <div className="value">
                                    <span className="rain-blue">
                                        { weather.rainDisplay(rain.hourly,units) }</span>&nbsp;{rainLabel(units)}
                                </div>
                            </div>
                            <div className="detail-item">
                                <div className="label">{ year(rain.lastrain) }:</div>
                                <div className="value">
                                    <span className="rain-blue">{ rainDisplay(rain.yearly, units)}</span>&nbsp;{rainLabel(units)}
                                </div>
                            </div>
                            <div className="detail-item">
                                <div className="label">{ month(rain.lastrain) }:</div>
                                <div className="value">
                                    <span className="rain-blue">{rainDisplay(rain.monthly, units)}</span>&nbsp;{rainLabel(units)}
                                </div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Last 24h:</div>
                                <div className="value">
                                    <span className="rain-blue">{rainDisplay(rain.daily, units)}</span>&nbsp;{rainLabel(units)}
                                </div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Last Rain:</div>
                                <div className="value">
                                    <span className="">{ full(rain.lastrain) }</span>
                                </div>
                            </div>
                        </div>
                        <div className="rain-gauge">
                            <RainGuage rainAmount={rain.daily} size={300} gaugeColor="#555" waterColor="#3b9cac" tickColor="#fff"/>
                        </div>
                    </div>
                    <div className="bar-chart">
                        <h3>Last 30 days</h3>
                        <ResponsiveContainer width="100%" height="100%">
                            <BarChart
                                data={chart[0].values}
                                margin={{ top: 20, right: 30, left: 20, bottom: 5 }}
                            >
                                <CartesianGrid stroke="#ccc" />
                                <XAxis
                                    dataKey="date"
                                    tickFormatter={formatDay}
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
                                    content={CustomTooltip}
                                />
                                <Bar
                                    dataKey="value"
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
export default Rain