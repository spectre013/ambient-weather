import { Current } from "../models/current";
import BoxData from "./BoxData";
import moment from "moment";
import * as weather from '../util/weather'
import "./Rain.css"
export interface Props {
    live: Current
    units: string
}

const Rain = (props:Props) => {

    function year(date: string) {
        return moment(date).format('YYYY');
    }
    function month(date: string) {
        return moment(date).format('MMM');
    }
    function full(date: string) {
        return moment(date).format('YYYY-MM-DD HH:mm:ss');
    }

    return (
        <BoxData icon={'&#xF019;'} title="RAIN" unit={weather.rainLabel(props.units)} style={{}}>
            <div className="rain-container">
                <div className="lastrain"><span className="rain-blue">Last Rain:</span>&nbsp;{ full(props.live.lastrain) }</div>
                <div className="raintotal">
                    <span className="amount rain-blue">{ weather.rainDisplay(props.live.dailyrainin, props.units) }</span>&nbsp;
                    { weather.rainLabel(props.units) }
                </div>
                <div className="rate">
                    <span className="rain-blue"> Rate:</span>&nbsp;{ weather.rainDisplay(props.live.hourlyrainin, props.units) }&nbsp;
                    { weather.rainLabel(props.units) }
                </div>
                <div className="year">
                    { year(props.live.date) }&nbsp;
                    <span className="rain-blue">{weather.rainDisplay(props.live.yearlyrainin, props.units)}</span>&nbsp;
                    { weather.rainLabel(props.units) }
                </div>
                <div className="month">
                    { month(props.live.date) }:&nbsp;
                    <span className="rain-blue">{weather.rainDisplay(props.live.monthlyrainin, props.units)}</span>
                    &nbsp;{weather.rainLabel(props.units)}
                </div>
                <div className="hour">
                    Last Hour: <span className="rain-blue">{ weather.rainDisplay(props.live.hourlyrain, props.units)}</span>
                    &nbsp;{ weather.rainLabel(props.units) }
                </div>
                <div className="tfhour">
                    Last 24hr: <span className="rain-blue">{ weather.rainDisplay(props.live.dailyrainin, props.units)}</span>
                    &nbsp;{ weather.rainLabel(props.units)}
                </div>
            </div>
        </BoxData>
    )
}
export default Rain
