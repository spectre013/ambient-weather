import {useContext} from "react";
import {TempData } from "../models/current";
import "./Temperature.css";
import {tempColor, tempDisplay, tempLabel} from "../util/weather";
import { useNavigate } from 'react-router-dom';
import {WeatherContext} from "../Context.ts";

export interface Props {
    temp: TempData
    icon: string
    conditions: string
    units: string
}

const Temperature = (props:Props) => {
    const navigate = useNavigate();
    const ctx = useContext(WeatherContext);

    const handleClick = () => {
        navigate('/details/temperature'); // Navigate to the /dashboard route
    };

    return (
        <>
        <section className="current-weather hasclick" onClick={handleClick}>
            <div className="weather-icon"><img alt={props.icon} src={'/images/icons/' + props.icon + '.png'}/></div>
            <div className="weather-info">
                <h1 className={`temp-text ${tempColor(props.temp.temp)}`}>{ tempDisplay(props.temp.temp, props.units)}<span className={'main-unit'}>°{ tempLabel(props.units) }</span></h1>
                <p>{ props.conditions }</p>
                <div className="location">{ctx.longname}</div>
            </div>
        </section>
        </>
    )
}
export default Temperature