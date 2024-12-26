import { RainData } from "../models/current";
import BoxData from "./BoxData";
import moment from "moment";
import * as weather from '../util/weather'
import "./Rain.css"
import {full} from "../util/weather";
export interface Props {
    rain: RainData
    units: string
}

const Rain = (props:Props) => {

    function year(date: string) {
        return moment(date).format('YYYY');
    }
    function month(date: string) {
        return moment(date).format('MMM');
    }


    return (
        <BoxData icon="fa-cloud-showers-heavy" title="RAIN" unit={weather.rainLabel(props.units)} style={{}}>
            <div className="rain-container">
                <div className="lastrain"><span className="rain-blue">Last Rain:</span>&nbsp;{ full(props.rain.lastrain) }</div>
                <div className="raintotal">
                    <span className="amount rain-blue">{ weather.rainDisplay(props.rain.daily, props.units) }</span>&nbsp;
                    { weather.rainLabel(props.units) }
                </div>
                <div className="rate">
                    <span className="rain-blue"> Rate:</span>&nbsp;{ weather.rainDisplay(props.rain.hourly, props.units) }&nbsp;
                    { weather.rainLabel(props.units) }
                </div>
                <div className="year">
                    { year(props.rain.lastrain) }&nbsp;
                    <span className="rain-blue">{weather.rainDisplay(props.rain.yearly, props.units)}</span>&nbsp;
                    { weather.rainLabel(props.units) }
                </div>
                <div className="month">
                    { month(props.rain.lastrain) }:&nbsp;
                    <span className="rain-blue">{weather.rainDisplay(props.rain.monthly, props.units)}</span>
                    &nbsp;{weather.rainLabel(props.units)}
                </div>
                <div className="hour">
                    Last Hour: <span className="rain-blue">{ weather.rainDisplay(props.rain.hourly, props.units)}</span>
                    &nbsp;{ weather.rainLabel(props.units) }
                </div>
                <div className="tfhour">
                    Last 24hr: <span className="rain-blue">{ weather.rainDisplay(props.rain.daily, props.units)}</span>
                    &nbsp;{ weather.rainLabel(props.units)}
                </div>
            </div>
        </BoxData>
    )
}
export default Rain
