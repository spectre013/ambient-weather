import {Current} from "../models/current.ts";
import BoxData from "./BoxData.tsx";
import {useEffect, useState} from "react";
import {MinMax} from "../models/minmax.ts";
import {full} from "../util/weather.ts";
import moment from "moment/moment";
import "./Lightning.css";

export interface Props {
    live: Current
    units: string
}

const Lightning = (props:Props) => {
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
        dataFetch().then(() => {});
    }, []);

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

    function year(date: string) {
        return moment(date).format('YYYY');
    }
    function month(date: string) {
        return moment(date).format('MMM');
    }


    return (
        <BoxData icon="fa-bolt-lightning" title="Lightning" unit="" style={{}}>
            <div className="lightning-container">
                <div className="laststrike"><span
                    className="lightning-yellow">Last Strike:</span> {full(props.live.lightningtime)}
                </div>
                <div className="lhour">
                    <div className="lightning-yellow">Hour:</div>
                    <div
                        className={`lightning-value ${lightningClass(props.live.lightninghour)}`}>
                        {props.live.lightninghour}
                    </div>
                </div>
                <div className="lday">
                    <div className="lightning-yellow">Day:</div>
                    <div
                        className={`lightning-value ${lightningClass(props.live.lightningday)}`}>
                        { props.live.lightningday }
                    </div>
                </div>
                <div className="lyesterday">
                    <div className="lightning-yellow">Yesterday:</div>
                    <div
                        className={`lightning-value ${lightningClass(minmax.max.yesterday.value)}`}>
                        { minmax.max.yesterday.value }
                    </div>
                </div>
                <div className="lmonth">
                    <div className="lightning-yellow">{month(props.live.date)}:</div>
                    <div
                        className={`lightning-value ${lightningClass(props.live.lightningmonth)}`}>
                        { props.live.lightningmonth }
                    </div>
                </div>
                <div className="lastd">
                    <div className="lightning-yellow">Distance</div>
                    <div
                        className={`lightning-value ${distanceClass(props.live.lightningdistance)}`}>
                        { props.live.lightningdistance }
                    </div>
                </div>
                <div className="lyear">
                    <div className="lightning-yellow">{year(props.live.date)}:</div>
                    <div
                        className={`lightning-value ${lightningClass(minmax.max.year.value)}`}>
                        { minmax.max.year.value }
                    </div>
                </div>
            </div>
        </BoxData>
    )
}
export default Lightning