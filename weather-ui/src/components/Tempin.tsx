import { Current } from "../models/current";
import BoxData from "./BoxData";
import * as weather from '../util/weather'
import { MinMax } from "../models/minmax";
import "./Tempin.css";
import {useEffect, useState} from "react";
import {tempColor} from "../util/weather";

export interface Props {
    live: Current
    sensor: string
    area: string
    units: string
}

const Tempin = (props:Props) => {
    const [minmax, setMinMax] = useState<MinMax>({"avg":{"day":{"value":68.53,"date":"2021-02-26T16:34:36Z"},"month":{"value":68.14,"date":"2021-02-26T16:34:36Z"},"year":{"value":68.68,"date":"2021-02-26T16:34:36Z"},"yesterday":{"value":67.9,"date":"2021-02-26T16:34:36Z"}},"max":{"day":{"value":70.9,"date":"2023-11-11T23:23:24Z"},"month":{"value":72.7,"date":"2023-11-06T01:31:55Z"},"year":{"value":73.6,"date":"2023-07-26T00:45:54Z"},"yesterday":{"value":69.6,"date":"2023-11-11T04:41:15Z"}},"min":{"day":{"value":66.9,"date":"2023-11-11T16:00:15Z"},"month":{"value":62.6,"date":"2023-11-05T12:36:59Z"},"year":{"value":0,"date":"2023-09-04T02:26:40Z"},"yesterday":{"value":66.2,"date":"2023-11-10T16:23:38Z"}}});
    const minmaxURL = "/api/minmax/temp" + props.sensor+"f";

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
    }, [minmaxURL]);
    function getSensor(sensor: string) {
        return props.live[sensor as keyof Current] as number;
    }
    return (
        <BoxData icon="fa-temperature-half" title={props.area} unit="&deg;F" style={{}}>
            <>
                <div className="tempin-container">
                    <div className="tempin">
                        <div className={`tempin-text ${weather.tempColor(getSensor("temp"+props.sensor+"f"))}`}>
                            { weather.tempDisplay(getSensor("temp"+props.sensor+"f"), props.units)}&deg;</div>
                        <div className="feels">Humidity:&nbsp;
                            <span className={weather.humidityClass(getSensor("humidity"+props.sensor))}>
                                { getSensor("humidity"+props.sensor) }%
                            </span>
                        </div>
                    </div>
                    <div className="maxin">Max: <span className={tempColor(minmax.max.day.value)}>{weather.tempDisplay(minmax.max.day.value, props.units)}&deg;</span></div>
                    <div className="minin">Min: <span className={tempColor(minmax.min.day.value)}>{weather.tempDisplay(minmax.min.day.value, props.units)}&deg;</span></div>
                </div>
            </>
        </BoxData>
    )
}
export default Tempin