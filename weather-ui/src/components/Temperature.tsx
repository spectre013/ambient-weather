import { Current } from "../models/current";
import BoxData from "./BoxData";
import * as weather from '../util/weather'
import {MinMax} from "../models/minmax";
import "./Temperature.css";

export interface Props {
    live: Current
    icon: string
    units: string
    avg: MinMax
}

const Temperature = (props:Props) => {

    return (
        <BoxData icon={'&#xF053;'} title="TEMP" unit="&deg;F" style={{}} navigate="temperature">
            <>
                <div className="temp-container">
                    <div className="icon"><img alt={props.icon} src={'/images/icons/' + props.icon + '.png'}/> </div>
                    <div className="temp">
                        <div className={`temp-text ${weather.tempColor(props.live.tempf)}`}>{ weather.tempDisplay(props.live.tempf, props.units)}&deg;</div>
                        <div className="feels">Feels: <span className={weather.tempColor(props.live.feelslike)}>{ weather.tempDisplay( props.live.feelslike , props.units)}&deg;</span> </div>
                    </div>
                    <div className="max">Max: {weather.tempDisplay(props.avg.max.day.value, props.units)}&deg;</div>
                    <div className="min">Min: {weather.tempDisplay(props.avg.min.day.value, props.units)}&deg;</div>
                </div>
            </>
        </BoxData>
    )
}
export default Temperature