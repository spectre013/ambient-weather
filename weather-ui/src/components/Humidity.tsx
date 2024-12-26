import { HumidityData } from "../models/current";
import BoxData from "./BoxData";
import * as weather from '../util/weather'

import "./Humidity.css";



export interface Props {
    humidity: HumidityData
    units: string
}

const Humidity = (props:Props) => {

    return (
        <BoxData icon="fa-droplet" title="Humidity" unit="%" style={{}}>
            <>
                <div className="humidity-container">
                    <div className="humidity-wrap">
                        <div className={`humidity-text ${weather.humidityClass(props.humidity.humdity)}`}>
                            { props.humidity.humdity }%
                        </div>
                        <div className="dewpoint">Dewpoint:&nbsp;
                            <span className={weather.dewPointClass(props.humidity.dewpoint)}>
                                { props.humidity.dewpoint.toFixed(0) }&deg;
                            </span>
                        </div>
                    </div>
                    <div className="humiditymax">Max: {weather.tempDisplay(props.humidity.minmax.max.day.value, props.units)}%</div>
                    <div className="humiditymin">Min: {weather.tempDisplay(props.humidity.minmax.min.day.value, props.units)}%</div>
                </div>
            </>
        </BoxData>
    )
}
export default Humidity