import {TempData } from "../models/current";
import BoxData from "./BoxData";
import * as weather from '../util/weather'
import "./Temperature.css";
import {tempColor} from "../util/weather";

export interface Props {
    temp: TempData
    icon: string
    units: string
}

const Temperature = (props:Props) => {
    return (
        <BoxData icon="fa-temperature-three-quarters" title="Temperature" unit="&deg;F" style={{}}>
            <>
                <div className="temp-container">
                    <div className="icon"><img alt={props.icon} src={'/images/icons/' + props.icon + '.png'}/> </div>
                    <div className="temp">
                        <div className={`temp-text ${weather.tempColor(props.temp.temp)}`}>{ weather.tempDisplay(props.temp.temp, props.units)}&deg;</div>
                        <div className="feels">Feels: <span className={weather.tempColor(props.temp.feelslike)}>{ weather.tempDisplay( props.temp.feelslike , props.units)}&deg;</span> </div>
                    </div>
                    <div className="max">Max: <span className={tempColor(props.temp.minmax.max.day.value)}>{weather.tempDisplay(props.temp.minmax.max.day.value, props.units)}&deg;</span>
                    </div>
                    <div className="min">Min: <span className={tempColor(props.temp.minmax.min.day.value)}>{weather.tempDisplay(props.temp.minmax.min.day.value, props.units)}&deg;</span>
                    </div>
                </div>
            </>
        </BoxData>
    )
}
export default Temperature