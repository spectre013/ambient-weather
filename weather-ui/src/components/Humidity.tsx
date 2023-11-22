import { Current } from "../models/current";
import BoxData from "./BoxData";
import * as weather from '../util/weather'
import { MinMax } from "../models/minmax";
import "./Humidity.css";
import {useEffect, useState} from "react";


export interface Props {
    live: Current
    units: string
}

const Humidity = (props:Props) => {
    const [minmax, setMinMax] = useState<MinMax>({"avg":{"day":{"value":68.53,"date":"2021-02-26T16:34:36Z"},"month":{"value":68.14,"date":"2021-02-26T16:34:36Z"},"year":{"value":68.68,"date":"2021-02-26T16:34:36Z"},"yesterday":{"value":67.9,"date":"2021-02-26T16:34:36Z"}},"max":{"day":{"value":70.9,"date":"2023-11-11T23:23:24Z"},"month":{"value":72.7,"date":"2023-11-06T01:31:55Z"},"year":{"value":73.6,"date":"2023-07-26T00:45:54Z"},"yesterday":{"value":69.6,"date":"2023-11-11T04:41:15Z"}},"min":{"day":{"value":66.9,"date":"2023-11-11T16:00:15Z"},"month":{"value":62.6,"date":"2023-11-05T12:36:59Z"},"year":{"value":0,"date":"2023-09-04T02:26:40Z"},"yesterday":{"value":66.2,"date":"2023-11-10T16:23:38Z"}}});
    const minmaxURL = "/api/minmax/humidity";
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
        <BoxData icon={'&#xF07A;'} title="Humidity" unit="%" style={{}}>
            <>
                <div className="humidity-container">
                    <div className="humidity-wrap">
                        <div className={`humidity-text ${weather.humidityClass(props.live.dewpoint)}`}>
                            { props.live.humidity }%
                        </div>
                        <div className="dewpoint">Dewpoint:&nbsp;
                            <span className={weather.dewPointClass(props.live.dewpoint)}>
                                { props.live.dewpoint.toFixed(0) }&deg;
                            </span>
                        </div>
                    </div>
                    <div className="humiditymax">Max: {weather.tempDisplay(minmax.max.day.value, props.units)}%</div>
                    <div className="humiditymin">Min: {weather.tempDisplay(minmax.min.day.value, props.units)}%</div>
                </div>
            </>
        </BoxData>
    )
}
export default Humidity