import "./Forecast.css";
import * as weather from '../util/weather'
import moment from "moment";
import BoxData from "./BoxData.tsx";
import * as CSS from 'csstype';
import {ForecastModel, Day } from "../models/Forecast.ts";
export interface Props {
    forecast: ForecastModel
    units: string
}

const Forecast = (props:Props) => {
    function getDayOfWeek(date: string) : string {
        return moment(date).format("ddd");
    }
    function render(day: Day)  {
        return (
            <div className="forecast-wrap" key={day.datetimeEpoch}>
                <span>{getDayOfWeek(day.datetime)}</span>
                <div className="forecast-container">
                    <div className="forecast-icon"><img src={'/images/icons/'+day.icon+'.png'} /></div>
                    <div className={`forecast-max ${weather.tempColor(day.tempmax)}`}>{ weather.tempDisplay(day.tempmax, props.units) }</div>
                    <div className={`forecast-min ${weather.tempColor(day.tempmin)}`}>{ weather.tempDisplay(day.tempmin, props.units)}</div>
                </div>
            </div>
        )
    }
    const days = props.forecast.days.slice(1, 8);
    const style: CSS.Properties = {
        width: '570px'
    };

    return (
        <BoxData icon={'&#xF002;'} title="FORECAST" unit="&deg;F" style={style}>
            {days.map((day: Day) => render(day))}
        </BoxData>
    )
}
export default Forecast