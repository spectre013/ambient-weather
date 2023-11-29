import {Current} from "../models/current";
import {Trend} from "../models/Trend";
import BoxData from "./BoxData";
import * as weather from '../util/weather'
import {MinMax} from "../models/minmax";
import "./Barometer.css";
import {useEffect, useState} from "react";
export interface Props {
    live: Current
    trend: Trend
    units: string
}
const Barometer = (props:Props) => {
    const [minmax, setMinMax] = useState<MinMax>({"avg":{"day":{"value":68.53,"date":"2021-02-26T16:34:36Z"},"month":{"value":68.14,"date":"2021-02-26T16:34:36Z"},"year":{"value":68.68,"date":"2021-02-26T16:34:36Z"},"yesterday":{"value":67.9,"date":"2021-02-26T16:34:36Z"}},"max":{"day":{"value":70.9,"date":"2023-11-11T23:23:24Z"},"month":{"value":72.7,"date":"2023-11-06T01:31:55Z"},"year":{"value":73.6,"date":"2023-07-26T00:45:54Z"},"yesterday":{"value":69.6,"date":"2023-11-11T04:41:15Z"}},"min":{"day":{"value":66.9,"date":"2023-11-11T16:00:15Z"},"month":{"value":62.6,"date":"2023-11-05T12:36:59Z"},"year":{"value":0,"date":"2023-09-04T02:26:40Z"},"yesterday":{"value":66.2,"date":"2023-11-10T16:23:38Z"}}});
    const minmaxURL = "/api/minmax/baromrelin";
    useEffect(() => {
        const dataFetch = async () => {
            const result = (
                await Promise.all([
                    fetch(minmaxURL),
                ])
            ).map((r) => r.json());

            const [ minmaxResult ] = await Promise.all(
                result
            );

            // when the data is ready, save it to state
            setMinMax(minmaxResult);
        }
        dataFetch();
    }, []);


    return (
        <BoxData icon={'&#xF079;'} title="Barometer" unit="inHG" style={{}}>
            <>
                <div className="barometer-container">
                    <div className="barometer-wrap">
                        <div className="barometer-text">
                            { weather.baroDisplay(props.live.baromrelin, props.units) }&nbsp;
                            <span className="units"> { weather.baroLabel(props.units) }</span>
                        </div>
                        <div>{ props.trend.trend}</div>
                    </div>
                    <div className="barometermax">Max: {weather.baroDisplay(minmax.max.day.value, props.units)}&nbsp;{ weather.baroLabel(props.units) }</div>
                    <div className="barometermin">Min: {weather.baroDisplay(minmax.min.day.value, props.units)}&nbsp;{ weather.baroLabel(props.units) }</div>
                </div>
            </>
        </BoxData>
    )
}
export default Barometer