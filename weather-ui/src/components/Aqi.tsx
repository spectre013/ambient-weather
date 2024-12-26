import { AirQualityIndex } from "../models/current";
import BoxData from "./BoxData";
import "./Aqi.css";

export interface Props {
    aqi: AirQualityIndex
    units: string
}

const Aqi = (props:Props) => {

    const categories = [
        {max: 50, color: "green", name: "Good"},
        {max: 100, color: "yellow", name: "Moderate"},
        {max: 150, color: "orange", name: "Unhealthy for sensitive groups"},
        {max: 200, color: "red", name: "Unhealthy"},
        {max: 300, color: "purple", name: "Very unhealthy"},
        {max: 500, color: "maroon", name: "Hazardous"}
    ]

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
        <BoxData icon="fa-lungs" title="Air Quality Index" unit="&deg;F" style={{}}>
            <>
                <div className="aqi-container">
                    <div className="aqi-wrap">
                        <div className={`aqi-text ${getDetails(Aqi(props.aqi.pm2524h)).color}`}>
                            { Aqi(props.aqi.pm2524h) }
                        </div>
                        <div className="status">
                            <div>{getDetails(Aqi(props.aqi.pm2524h)).name}</div>
                            <div>{(props.aqi.pm25)} µg/m3</div>
                        </div>
                    </div>
                    <div className="aqimax">
                        Max: <span className={getDetails(Aqi(props.aqi.minmax.max.day.value)).color}>{ Aqi(props.aqi.minmax.max.day.value) }</span>
                    </div>
                    <div className="aqimin">
                        Min: <span className={getDetails(Aqi(props.aqi.minmax.min.day.value)).color}>{ Aqi(props.aqi.minmax.min.day.value) }</span>
                    </div>
                </div>
            </>
        </BoxData>
    )
}
export default Aqi