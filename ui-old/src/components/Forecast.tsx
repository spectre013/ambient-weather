import "./Forecast.css";
import * as weather from '../util/weather'
import moment from "moment";
import {ForecastModel, Day, gradient} from "../models/Forecast.ts";
import {tempToHex} from "../util/weather";
import { useNavigate } from 'react-router-dom';

export interface Props {
    forecast: ForecastModel
    units: string
}


const Forecast = (props:Props) => {
    const navigate = useNavigate();

    function forecastDate(date: string) : string {
        return moment(date).format("MMM Do");
    }

    const forecastClick = (day: number) => {
        navigate('/details/forecast/'+day); // Navigate to the /dashboard route
    };
    function isAfterNoon():number {
        const now = new Date(); // Creates a new Date object with the current time.
        const hour = now.getHours(); // Extracts the hour (0-23) from the Date object.

        // The hour for noon is 12 (in 24-hour format).
        // Any hour greater than or equal to 12 is considered noon or after.
        if (hour >= 12) {
            return 1;
        } else {
            return 0;
        }
    }

    function render(day: Day, i: number)  {
        return (
            <div className="forecast-day" key={day.datetimeEpoch} onClick={()=>forecastClick(i+1)}>
                <div className="day-date">{ forecastDate(day.datetime) }</div>
                <div className="weather-icon"><img alt="" src={'/images/icons/'+day.icon+'.png'} /></div>
                <div className="day-info">
                    <div className="day-temp">
                        {weather.tempDisplay(day.tempmin, props.units)}°
                        <div className="temp-bar-container">
                            <div className="temp-bar" style={getTemperatureGradient(day.tempmin, day.tempmax)}></div>
                        </div>
                        {weather.tempDisplay(day.tempmax, props.units)}°
                    </div>
                </div>
            </div>
        )
    }
    const days = props.forecast.days.slice(isAfterNoon(), 10);
    return (
        <section className="forecast-section">
            <div className="forecast-header">
                <h2>10-Day Forecast</h2>
            </div>
            <div className="forecast-list">
            {days.map((day: Day, i:number) => render(day,i))}
            </div>
        </section>
    )

    function getTemperatureGradient(temp1: number, temp2: number): gradient {
        const color1 = tempToHex(temp1);
        const color2 = tempToHex(temp2);

        return {
            background: `linear-gradient(to right, ${color1}, ${color2})`,
            width: '100%',
        }
    }

}
export default Forecast