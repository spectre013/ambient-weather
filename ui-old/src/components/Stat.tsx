
import "./Stat.css"
import {
    baroDisplay,
    distanceDisplay,
    rainDisplay,
    tempColor,
    tempDisplay,
    windClass,
    windDisplay
} from "../util/weather";
import { useNavigate} from "react-router-dom";

export interface Props {
    icon: string;
    label: string;
    value: number
    valueType: string
    units: string
}

function setClassName(valueType: string, val: number): string {
    switch (valueType) {
        case "temperature":
            return tempColor(val);
        case "tempin":
            return tempColor(val);
        case "wind":
            return windClass(val);
        default:
            return ""
    }
}

function setValue(valueType: string, val: number, units:string): string {
    switch (valueType) {
        case "temperature":
            return tempDisplay(val,units);
        case "tempin":
            return tempDisplay(val,units);
        case "barometer":
            return baroDisplay(val,units);
        case "wind":
            return windDisplay(val, units);
        case "rain":
            return rainDisplay(val, units);
        case "lightning":
            return distanceDisplay(val, units).toFixed(0);
        case "forecast":
            return distanceDisplay(val, units).toFixed(0);
        default:
            return val.toString();
    }
}


const Stat = (props:Props) => {
    let hasClick = false;
    let statClass = "stat-card";
    const navigate = useNavigate();
    const clickTypes = ['temperature', 'wind','rain','lightning'];
    const units:string = localStorage.getItem("units") || "imperial";
    const handleClick = () => {
        navigate('/details/' + props.valueType); // Navigate to the details page for the specific stat
    }

    if (clickTypes.includes(props.valueType)) {
        statClass += " hasclick";
        hasClick = true;
    }

    return (
        <div className={statClass} onClick={hasClick ? handleClick: undefined}>
            <div className="stat-icon"><span className="material-symbols-sharp">{ props.icon}</span>
            </div>
            <div className="value"><span className={setClassName(props.valueType,props.value)}>
                { setValue(props.valueType, props.value,units) }</span>
                <span className="stat-value-unit">{ props.units }</span></div>
            <div className="label">{ props.label }</div>
        </div>
    )
}
export default Stat