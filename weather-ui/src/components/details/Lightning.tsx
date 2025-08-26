import './Lightning.css'
import {useNavigate} from "react-router-dom";
import {distanceLabel, formatDay, full} from "../../util/weather.ts";
import {ChartData} from "../../models/DataSeries.ts";
import {useEffect, useState} from "react";
import {LightningData} from "../../models/current.ts";
import {
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    ResponsiveContainer, BarChart, Bar, Legend
} from 'recharts';
import moment from "moment";
import CircleGauge from "../../util/Circlegauge.tsx";
import {CustomTooltip} from "../../util/utilities.tsx";

const Temperature = () => {
    const navigate = useNavigate();
    const [fLoaded, setFLoaded] = useState(false);
    const [units, setUnits] = useState<string>("imperial");
    const [lightning, setLightning] = useState<LightningData>({} as LightningData);
    const [chart, setChart] = useState<ChartData>({} as ChartData);
    const urls = [
        '/api/current',
        '/api/chart/lightning/3h',
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
                setLightning(data[0].lightning);
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

    function lightningClass(cnt: number) :string {
        if(cnt == 0 || cnt < 50) {
            return "green"
        } else if (cnt >= 50 && cnt < 250) {
            return "yellow"
        } else if (cnt >= 250 && cnt < 500) {
            return "orange"
        } else if (cnt >= 500) {
            return "red"
        }
        return "green"
    }

    function lightningToHex(cnt: number) :string {
        if(cnt == 0 || cnt < 50) {
            return "#0B6623"
        } else if (cnt >= 50 && cnt < 250) {
            return "#ff0"
        } else if (cnt >= 250 && cnt < 500) {
            return "#ff7e00"
        } else if (cnt >= 500) {
            return "#FF0000"
        }
        return "#0B6623"
    }

    function distanceClass(d :number) :string {
        if (d == 0) {
            return "green"
        } else if (d >= 1 && d < 5) {
            return "red"
        } else if (d >= 5 && d < 10) {
            return "orange"
        } else if (d >= 10 && d < 15) {
            return "yellow"
        } else if (d >= 15) {
            return "green"
        }
        return "green"
    }
    
    function month(date: Date) {
        return moment(date).format('MMM');
    }

    return (
        <>
            <div className="details-dashboard">
                <header className="details-header">
                    <h1><span className="material-symbols-sharp">bolt</span> Losron Ranch -  Lightning</h1>
                    <div className="hasclick" onClick={handleClick}><span className="material-symbols-sharp">home</span></div>
                </header>
                <div className="content">
                    <div className="details-content">
                        <div className="details">
                            <div className="detail-item">
                                <div className="label">Hour:</div>
                                <div className="value">
                                    <span className={lightningClass(lightning.hour)}>
                                         {lightning.hour.toLocaleString('en-US')}</span></div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Day:</div>
                                <div className="value">
                                    <span className={lightningClass(lightning.day)}>
                                        {lightning.day.toLocaleString('en-US')}</span></div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Yesterday:</div>
                                <div className="value">
                                    <span className={lightningClass(lightning.minmax.max.yesterday.value)}>
                                        {lightning.minmax.max.yesterday.value.toLocaleString('en-US')}</span></div>
                            </div>
                            <div className="detail-item">
                                    <div className="label">{month(new Date())}:</div>
                                <div className="value">
                                    <span className={lightningClass(lightning.month)}>
                                        { lightning.month.toLocaleString('en-US') }</span></div>
                            </div>
                            <div className="detail-item">
                                <div className="label">{new Date().getFullYear()}:</div>
                                <div className="value">
                                    <span className={lightningClass(lightning.minmax.max.year.value)}>
                                        { lightning.minmax.max.year.value.toLocaleString('en-US') }</span></div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Last Strike Distance:</div>
                                <div className="value">
                                    <span className={distanceClass(lightning.distance)}>
                                        {lightning.distance}</span> &nbsp;{distanceLabel(units)}</div>
                            </div>
                            <div className="detail-item">
                                <div className="label">Last Strike:</div>
                                <div className="value">
                                    <span className="">{full(lightning.time)}</span></div>
                            </div>
                        </div>
                        <div className="dial">
                            <CircleGauge size={275} value={lightning.day} color={lightningToHex(lightning.day)}/>
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
                                    fill="#ff7e00"
                                    name="Strikes"
                                />
                                <Legend />
                            </BarChart>
                        </ResponsiveContainer>
                    </div>
                </div>
            </div>
        </>
    )
}
export default Temperature