import { LightningData} from "../models/current.ts";
import BoxData from "./BoxData.tsx";
import {full} from "../util/weather.ts";
import moment from "moment/moment";
import "./Lightning.css";

export interface Props {
    lightning: LightningData
    date: string
    units: string
}

const Lightning = (props:Props) => {

    function lightningClass(cnt: number) :string {
        if(cnt == 0 || cnt < 50) {
            return "green"
        } else if (cnt >= 50 && cnt < 250) {
            return "yellow"
        } else if (cnt >= 250 && cnt < 500) {
            return "orange"
        } else if (cnt >= 500) {
            return "red"
        }
        return "green"
    }

    function distanceClass(d :number) :string {
        if (d == 0) {
            return "green"
        } else if (d >= 1 && d < 5) {
            return "red"
        } else if (d >= 5 && d < 10) {
            return "orange"
        } else if (d >= 10 && d < 15) {
            return "yellow"
        } else if (d >= 15) {
            return "green"
        }
        return "green"
    }

    function year(date: string) {
        return moment(date).format('YYYY');
    }
    function month(date: string) {
        return moment(date).format('MMM');
    }


    return (
        <BoxData icon="fa-bolt-lightning" title="Lightning" unit="" style={{}}>
            <div className="lightning-container">
                <div className="laststrike"><span
                    className="lightning-yellow">Last Strike:</span> {full(props.lightning.time)}
                </div>
                <div className="lhour">
                    <div className="lightning-yellow">Hour:</div>
                    <div
                        className={`lightning-value ${lightningClass(props.lightning.hour)}`}>
                        {props.lightning.hour}
                    </div>
                </div>
                <div className="lday">
                    <div className="lightning-yellow">Day:</div>
                    <div
                        className={`lightning-value ${lightningClass(props.lightning.day)}`}>
                        { props.lightning.day }
                    </div>
                </div>
                <div className="lyesterday">
                    <div className="lightning-yellow">Yesterday:</div>
                    <div
                        className={`lightning-value ${lightningClass(props.lightning.minmax.max.yesterday.value)}`}>
                        { props.lightning.minmax.max.yesterday.value }
                    </div>
                </div>
                <div className="lmonth">
                    <div className="lightning-yellow">{month(props.date)}:</div>
                    <div
                        className={`lightning-value ${lightningClass(props.lightning.month)}`}>
                        { props.lightning.month }
                    </div>
                </div>
                <div className="lastd">
                    <div className="lightning-yellow">Distance</div>
                    <div
                        className={`lightning-value ${distanceClass(props.lightning.distance)}`}>
                        { props.lightning.distance }
                    </div>
                </div>
                <div className="lyear">
                    <div className="lightning-yellow">{year(props.date)}:</div>
                    <div
                        className={`lightning-value ${lightningClass(props.lightning.minmax.max.year.value)}`}>
                        { props.lightning.minmax.max.year.value }
                    </div>
                </div>
            </div>
        </BoxData>
    )
}
export default Lightning