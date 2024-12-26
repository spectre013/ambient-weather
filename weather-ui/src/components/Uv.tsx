import {AstroData, UVData} from "../models/current.ts";
import BoxData from "./BoxData.tsx";
import "./Uv.css";


export interface Props {
    uv: UVData
    units: string
    astro: AstroData
}

const Uv = (props:Props) => {

    function uvCaution(uv: number) {
        if (uv >= 10) {
            return 'Extreme';
        } else if (uv >= 8) {
            return 'Very High';
        } else if (uv >= 5) {
            return 'High';
        } else if (uv >= 3) {
            return 'Moderate';
        } else if (props.astro.hasSunset && uv >= 0) {
            return 'Low';
        } else if (props.astro.hasSunset && uv <= 0) {
            return 'Below Horizon';
        }
        return '';
    }

    function uvToday() {
        if (props.uv.uv >= 10) {
            return 'uvtoday11';
        } else if (props.uv.uv >= 8) {
            return 'uvtoday9-10';
        } else if (props.uv.uv >= 5) {
            return 'uvtoday6-8';
        } else if (props.uv.uv >= 3) {
            return 'uvtoday4-5';
        } else if (props.uv.uv >= 0) {
            return 'uvtoday1-3';
        }
        return '';
    }

    return (
        <BoxData icon="fa-cloud-sun" title="UV | Solar" unit="" style={{}}>
            <div className="uv-container">
                <div className="uvimax">Max: <span className={uvToday()}>{ props.uv.minmax.max.day.value }</span> UVI</div>
                <div className="uvitext">
                    <span className="uvi-icon uvi-top">&#xF00D;</span> <span className="uvi-top">UVI</span> { uvCaution(props.uv.uv) }
                </div>
                <div className="uvi">
                    <div><span className={`value-text ${uvToday()}`}>{ props.uv.uv }</span> UVI</div>
                    <div>UV Index</div>
                </div>
                <div className="solar">
                    <div><span className={`solar-value ${uvToday()}`}>{ props.uv.solarradiation.toFixed(0) }</span> W/m<sup>2</sup></div>
                    <div>Solar Radiation</div>
                </div>
            </div>
        </BoxData>
    )
}
export default Uv