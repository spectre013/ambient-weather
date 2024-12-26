import { TemperatureData } from "../models/current";
import BoxData from "./BoxData";
import * as weather from '../util/weather'
import "./Tempin.css";
import {tempColor} from "../util/weather";

export interface Props {
    temp: TemperatureData
    area: string
    units: string
}

const Tempin = (props:Props) => {

    return (
        <BoxData icon="fa-temperature-half" title={props.area} unit="&deg;F" style={{}}>
            <>
                <div className="tempin-container">
                    <div className="tempin">
                        <div className={`tempin-text ${weather.tempColor(props.temp.temp)}`}>
                            { weather.tempDisplay(props.temp.temp, props.units)}&deg;</div>
                        <div className="feels">Humidity:&nbsp;
                            <span className={weather.humidityClass(props.temp.humidity)}>
                                { props.temp.humidity }%
                            </span>
                        </div>
                    </div>
                    <div className="maxin">Max: <span className={tempColor(props.temp.minmax.max.day.value)}>{weather.tempDisplay(props.temp.minmax.max.day.value, props.units)}&deg;</span></div>
                    <div className="minin">Min: <span className={tempColor(props.temp.minmax.min.day.value)}>{weather.tempDisplay(props.temp.minmax.min.day.value, props.units)}&deg;</span></div>
                </div>
            </>
        </BoxData>
    )
}
export default Tempin