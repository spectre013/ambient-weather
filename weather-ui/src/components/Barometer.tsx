import { BaroData } from "../models/current";
import BoxData from "./BoxData";
import * as weather from '../util/weather'
import "./Barometer.css";
export interface Props {
    baro: BaroData
    units: string
}
const Barometer = (props:Props) => {

    return (
        <BoxData icon="fa-temperature-high" title="Barometer" unit="inHG" style={{}}>
            <>
                <div className="barometer-container">
                    <div className="barometer-wrap">
                        <div className="barometer-text">
                            { weather.baroDisplay(props.baro.baromrelin, props.units) }&nbsp;
                            <span className="units"> { weather.baroLabel(props.units) }</span>
                        </div>
                        <div>{ props.baro.trend.trend}</div>
                    </div>
                    <div className="barometermax">Max: {weather.baroDisplay(props.baro.minmax.max.day.value, props.units)}&nbsp;{ weather.baroLabel(props.units) }</div>
                    <div className="barometermin">Min: {weather.baroDisplay(props.baro.minmax.min.day.value, props.units)}&nbsp;{ weather.baroLabel(props.units) }</div>
                </div>
            </>
        </BoxData>
    )
}
export default Barometer
