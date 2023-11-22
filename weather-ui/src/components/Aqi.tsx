import { Current } from "../models/current";
import BoxData from "./BoxData";
import { MinMax } from "../models/minmax";
import "./Aqi.css";
import {useEffect, useState} from "react";

export interface Props {
    live: Current
    units: string
}

const Aqi = (props:Props) => {
    const [minmax, setMinMax] = useState<MinMax>({"avg":{"day":{"value":68.53,"date":"2021-02-26T16:34:36Z"},"month":{"value":68.14,"date":"2021-02-26T16:34:36Z"},"year":{"value":68.68,"date":"2021-02-26T16:34:36Z"},"yesterday":{"value":67.9,"date":"2021-02-26T16:34:36Z"}},"max":{"day":{"value":70.9,"date":"2023-11-11T23:23:24Z"},"month":{"value":72.7,"date":"2023-11-06T01:31:55Z"},"year":{"value":73.6,"date":"2023-07-26T00:45:54Z"},"yesterday":{"value":69.6,"date":"2023-11-11T04:41:15Z"}},"min":{"day":{"value":66.9,"date":"2023-11-11T16:00:15Z"},"month":{"value":62.6,"date":"2023-11-05T12:36:59Z"},"year":{"value":0,"date":"2023-09-04T02:26:40Z"},"yesterday":{"value":66.2,"date":"2023-11-10T16:23:38Z"}}});
    const minmaxURL = "/api/minmax/aqipm25";


    const categories = [
        {max: 50, color: "green", name: "Good"},
        {max: 100, color: "yellow", name: "Moderate"},
        {max: 150, color: "orange", name: "Unhealthy for sensitive groups"},
        {max: 200, color: "red", name: "Unhealthy"},
        {max: 300, color: "purple", name: "Very unhealthy"},
        {max: 500, color: "maroon", name: "Hazardous"}
    ]

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

    function getDetails(aqi: number) {
        if(aqi < 0 || (aqi >= 0 && aqi <= 50)) {
            return categories[0];
        } else if(aqi > 50 || aqi <= 100) {
            return categories[1];
        } else if(aqi >100 || aqi <= 150) {
            return categories[2];
        } else if(aqi > 150 || aqi <= 200) {
            return categories[3];
        } else if(aqi > 200 || aqi <= 300) {
            return categories[4];
        } else if(aqi > 300) {
            return categories[5];
        }
        return categories[0];
    }

    function Linear(AQIhigh: number, AQIlow: number, Conchigh: number, Conclow: number, Concentration: number): number {
        const Conc= Concentration;
        const a = ((Conc-Conclow)/(Conchigh-Conclow))*(AQIhigh-AQIlow)+AQIlow;
        return  Math.round(a);
    }

    function Aqi(Concentration: number)
    {
        const  Conc= Concentration;
        let AQI = 0;
        const c: number=(Math.floor(10*Conc))/10;
        if (c >=0 && c <12.1) {
            AQI=Linear(50,0,12,0,c);
        } else if (c>=12.1 && c<35.5)
        {
            AQI=Linear(100,51,35.4,12.1,c);
        } else if (c>=35.5 && c<55.5) {
            AQI=Linear(150,101,55.4,35.5,c);
        }
        else if (c>=55.5 && c<150.5) {
            AQI=Linear(200,151,150.4,55.5,c);
        }
        else if (c>=150.5 && c<250.5) {
            AQI=Linear(300,201,250.4,150.5,c);
        }
        else if (c>=250.5 && c<350.5) {
            AQI=Linear(400,301,350.4,250.5,c);
        }
        else if (c>=350.5) {
            AQI=Linear(500,401,500.4,350.5,c);
        }
        return AQI;
    }


    return (
        <BoxData icon={'&#xF082;'} title="Air Quality Index" unit="&deg;F" navigate="aqi" style={{}}>
            <>
                <div className="aqi-container">
                    <div className="aqi-wrap">
                        <div className={`aqi-text ${getDetails(Aqi(props.live.aqipm2524h)).color}`}>
                            { Aqi(props.live.aqipm2524h) }
                        </div>
                        <div className="status">
                            <div>{getDetails(Aqi(props.live.aqipm2524h)).name}</div>
                            <div>{(props.live.aqipm25)} µg/m3</div>
                        </div>
                    </div>
                    <div className="aqimax">
                        Max: <span className={getDetails(Aqi(minmax.max.day.value)).color}>{ Aqi(minmax.max.day.value) }</span>
                    </div>
                    <div className="aqimin">
                        Min: <span className={getDetails(Aqi(minmax.min.day.value)).color}>{ Aqi(minmax.min.day.value) }</span>
                    </div>
                </div>
            </>
        </BoxData>
    )
}
export default Aqi