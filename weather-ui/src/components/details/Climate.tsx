import {useEffect, useState} from "react";
import {ClimateData, ClimateRecord, FrostDateRecord, weatherDisplay} from "../../models/ClimateModel.ts";
import {freezeDate, rainDisplay, tempDisplay} from "../../util/weather.ts";
import "./Climate.css";
import Header from "../Header.tsx";
import Footer from "../Footer.tsx";


function Climate() {
    const [fLoaded, setFLoaded] = useState(false);
    const [units, setUnits] = useState<string>(localStorage.getItem('units') || 'imperial');
    const [climate, setClimate] = useState<ClimateRecord[]>([] as ClimateRecord[]);
    const [freeze, setFreeze] = useState<FrostDateRecord[]>([] as FrostDateRecord[]);
    const [selectedMetric, setSelectedMetric] = useState<MetricKey>('avgtemp');
    const urls = [
        '/api/climate',
        'api/firstfreeze'
    ];
    const calculateAnnualAverage = (monthlyData: number[],year:number): number => {
        let validMonths = monthlyData;
        const curYear = new Date().getFullYear()
        if(year === 2020) {
            validMonths = monthlyData.slice(7,12); // Only 6 months of data for 2020
        } else if(year === curYear) {
            validMonths = monthlyData.slice(1,new Date().getMonth()+1); // Months are 1-indexed, ignore index 0
        }
        if(selectedMetric.includes("temp") || selectedMetric.includes("wind")) {
            const sum = validMonths.reduce((acc, value) => acc + value, 0);
            return Math.round(sum / validMonths.length);
        } else {
            return validMonths.reduce((acc, value) => acc + value, 0);
        }
    };
    const monthHeaders = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    type MetricKey = keyof Omit<ClimateData, 'year'>;

    const metricLabels: Record<MetricKey, string> = {
        avgtemp: "Average Temperature (°F)",
        maxtemp: "Max Temperature (°F)",
        mintemp: "Min Temperature (°F)",
        avgrain: "Rainfall (in)",
    };

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
                setClimate(data[0] as ClimateRecord[]);
                setFreeze(data[1] as FrostDateRecord[]);
                setFLoaded(true);
            })
            .catch(error => {
                // Handle any errors that occurred during fetching or JSON parsing
                console.error('Error fetching data:', error);
            });
    },[]);

    if (!fLoaded) {
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


    function render(data:ClimateRecord)  {
        const values:number[] = data.Data[selectedMetric as keyof ClimateData] as number[];
        //@ts-expect-error - issue with JSX.Element type
        const elements: jsx.Element[] = [];
        let func: weatherDisplay = tempDisplay;
        switch (selectedMetric) {
            case "avgtemp":
                func = tempDisplay;
                break;
            case "maxtemp":
                func = tempDisplay;
                break;
            case "mintemp":
                func = tempDisplay;
                break;
            case "avgrain":
                func = rainDisplay;
                break;
        }

        for (let i = 1; i < values.length; i++) {
            elements.push(<td key={i}>{ func(values[i], units) }</td>);
        }

        return (
                <tr key={data.Year}>
                    <td>{data.Year}</td>
                    {elements}
                    <td>{ func(calculateAnnualAverage(values,data.Year),units) }</td>
                </tr>
        )
    }

    return (
        <>
        <div className="details-dashboard">
            <Header name="Historical Climate Data" icon="mode_heat_cool" />
            <div className="content">
                <div className="climate-details-content">
                    <div className="climate-table-container">
                        <div className="table-header">
                            <h2 className="table-title">Historical Climate Data</h2>
                            <div className="metric-selector-wrapper">
                                <select
                                    id="metric-selector"
                                    className="metric-selector"
                                    value={selectedMetric}
                                    onChange={(e) => setSelectedMetric(e.target.value as MetricKey)}
                                >
                                    {Object.entries(metricLabels).map(([key, label]) => (
                                        <option key={key} value={key}>{label}</option>
                                    ))}
                                </select>
                            </div>
                        </div>
                        <div className="table-wrapper">
                            <table>
                                <thead>
                                <tr>
                                    <th>Year</th>
                                    {monthHeaders.map(month => <th key={month}>{month}</th>)}
                                    <th>Annual</th>
                                </tr>
                                </thead>
                                <tbody>
                                    {climate.map((value: ClimateRecord) => render(value))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                    <div className="climate-table-container">
                        <h2 className="table-title">First and Last Freeze Dates</h2>
                        <table>
                            <thead>
                            <tr>
                                <th>Year</th>
                                <th>First Freeze</th>
                                <th>Last Freeze</th>
                            </tr>
                            </thead>
                            <tbody>
                            {freeze.map((value: FrostDateRecord) => (
                                <tr key={value.year}>
                                    <td>{value.year}</td>
                                    <td>{freezeDate(value.spring)}</td>
                                    <td>{freezeDate(value.fall)}</td>
                                </tr>
                            ))}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
        <Footer />
        </>
    );
}

export default Climate;