import {Current} from "../models/current.ts";
import {MinMax} from "../models/minmax.ts";
import BoxData from "./BoxData.tsx";
import {useEffect, useState} from "react";
import moment from "moment/moment";
import {Luna} from "../models/Luna.ts";
import "./Uv.css";


export interface Props {
    live: Current
    units: string
    luna: Luna
}

const Uv = (props:Props) => {
    const [minmax, setMinMax] = useState<MinMax>({"avg":{"day":{"value":68.53,"date":"2021-02-26T16:34:36Z"},"month":{"value":68.14,"date":"2021-02-26T16:34:36Z"},"year":{"value":68.68,"date":"2021-02-26T16:34:36Z"},"yesterday":{"value":67.9,"date":"2021-02-26T16:34:36Z"}},"max":{"day":{"value":70.9,"date":"2023-11-11T23:23:24Z"},"month":{"value":72.7,"date":"2023-11-06T01:31:55Z"},"year":{"value":73.6,"date":"2023-07-26T00:45:54Z"},"yesterday":{"value":69.6,"date":"2023-11-11T04:41:15Z"}},"min":{"day":{"value":66.9,"date":"2023-11-11T16:00:15Z"},"month":{"value":62.6,"date":"2023-11-05T12:36:59Z"},"year":{"value":0,"date":"2023-09-04T02:26:40Z"},"yesterday":{"value":66.2,"date":"2023-11-10T16:23:38Z"}}});
    const [hasSunset, setHasSunset] = useState(false);
    const minmaxURL = "/api/minmax/uv";

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
        function sunHasSet() {
            const ss = props.luna.sunset.split(':');
            const sunset = moment().startOf('day')
                .hour(parseInt(ss[0]))
                .minute(parseInt(ss[1]));
            const s = moment.duration(sunset.diff(moment())).minutes();

            const sr = props.luna.sunrise.split(':');
            const sunrise = moment().startOf('day').hour(parseInt(sr[0])).minute(parseInt(sr[1]));
            const r = moment.duration(sunrise.diff(moment())).minutes();
            setHasSunset(s <= 0 || r >= 0);
        }
        dataFetch().then(() => {});
        sunHasSet()
    }, [props]);
    function uvCaution(uv: number) {
        if (uv >= 10) {
            return 'Extreme';
        } else if (uv >= 8) {
            return 'Very High';
        } else if (uv >= 5) {
            return 'High';
        } else if (uv >= 3) {
            return 'Moderate';
        } else if (!hasSunset && uv >= 0) {
            return 'Low';
        } else if (hasSunset && uv <= 0) {
            return 'Below Horizon';
        }
        return '';
    }

    function uvToday() {
        if (props.live.uv >= 10) {
            return 'uvtoday11';
        } else if (props.live.uv >= 8) {
            return 'uvtoday9-10';
        } else if (props.live.uv >= 5) {
            return 'uvtoday6-8';
        } else if (props.live.uv >= 3) {
            return 'uvtoday4-5';
        } else if (props.live.uv >= 0) {
            return 'uvtoday1-3';
        }
        return '';
    }

    return (
        <BoxData icon="fa-cloud-sun" title="UV | Solar" unit="" style={{}}>
            <div className="uv-container">
                <div className="uvimax">Max: <span className={uvToday()}>{ minmax.max.day.value }</span> UVI</div>
                <div className="uvitext">
                    <span className="uvi-icon uvi-top">&#xF00D;</span> <span className="uvi-top">UVI</span> { uvCaution(props.live.uv) }
                </div>
                <div className="uvi">
                    <div><span className={`value-text ${uvToday()}`}>{ props.live.uv }</span> UVI</div>
                    <div>UV Index</div>
                </div>
                <div className="solar">
                    <div><span className={`solar-value ${uvToday()}`}>{ props.live.solarradiation.toFixed(0) }</span> W/m<sup>2</sup></div>
                    <div>Solar Radiation</div>
                </div>
            </div>
        </BoxData>
    )
}
export default Uv