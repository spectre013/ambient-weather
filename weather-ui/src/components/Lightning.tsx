import {Current} from "../models/current.ts";
import BoxData from "./BoxData.tsx";
import {useEffect, useState} from "react";
import {MinMax} from "../models/minmax.ts";



export interface Props {
    live: Current
    units: string
}

const Uv = (props:Props) => {
    const [minmax, setMinMax] = useState<MinMax>({"avg":{"day":{"value":68.53,"date":"2021-02-26T16:34:36Z"},"month":{"value":68.14,"date":"2021-02-26T16:34:36Z"},"year":{"value":68.68,"date":"2021-02-26T16:34:36Z"},"yesterday":{"value":67.9,"date":"2021-02-26T16:34:36Z"}},"max":{"day":{"value":70.9,"date":"2023-11-11T23:23:24Z"},"month":{"value":72.7,"date":"2023-11-06T01:31:55Z"},"year":{"value":73.6,"date":"2023-07-26T00:45:54Z"},"yesterday":{"value":69.6,"date":"2023-11-11T04:41:15Z"}},"min":{"day":{"value":66.9,"date":"2023-11-11T16:00:15Z"},"month":{"value":62.6,"date":"2023-11-05T12:36:59Z"},"year":{"value":0,"date":"2023-09-04T02:26:40Z"},"yesterday":{"value":66.2,"date":"2023-11-10T16:23:38Z"}}});
    const minmaxURL = "/api/minmax/lightning";

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
        <BoxData icon={'&#xF016;'} title="Lightning" unit="" style={{}}>
            <div className="lightning-container">
                <div className="laststrike"></div>
                <div className="hour">{ props.live.lightninghour }</div>
                <div className="day">{ props.live.lightningday }</div>
                <div className="yesterday">{ minmax.max.yesterday.value}</div>
                <div className="month"></div>
                <div className="year"></div>
            </div>
        </BoxData>
    )
}
export default Uv