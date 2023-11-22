import BoxData from "./BoxData";
import {Luna} from "../models/Luna.ts";
import "./Sun.css"
import moment from "moment";
export interface Props {
    luna: Luna
}

const Sun = (props:Props) => {
    let hasSunset = false;

    function timeToDecimal(t: string) {
        const arr: string[] = t.split(':');
        const dec: number = (parseInt(arr[1], 10) / 6) * 10;
        return parseFloat(parseInt(arr[0], 10) + '.' + (dec < 10 ? '0' : '') + dec);
    }

    function hoursTilSunSet() {
        let sunsetTime = moment(props.luna.date + ' ' + props.luna.sunset + ':00');
        const now = moment();
        const h = moment.duration(sunsetTime.diff(now)).hours();

        if (h < 1 && hasSunset) {
            sunsetTime = moment(props.luna.tomorrow.date + ' ' + props.luna.tomorrow.sunrise + ':00');
        }
        const duration = sunsetTime.diff(now);
        return moment.duration(duration).hours();
    }
    function minTilSunSet() {
        let sunsetTime = moment(props.luna.date + ' ' + props.luna.sunset + ':00');
        const now = moment();
        const t = moment.duration(sunsetTime.diff(now)).minutes();

        if (t <= 0 && hasSunset) {
            sunsetTime = moment(props.luna.tomorrow.date + ' ' + props.luna.tomorrow.sunrise + ':00');
        }
        const duration = sunsetTime.diff(now);

        return moment.duration(duration).minutes();
    }
    function setDateTime(hours: string) {
        const h = hours.split(':');
        return moment().startOf('day').hour(parseInt(h[0])).minute(parseInt(h[1]));
    }
    function todayTomorrow(type: string): string {
        const eventTime = props.luna[type as keyof Luna];
        const event = setDateTime(eventTime as string);
        if (moment() > event) {
            return 'Tomorrow';
        } else {
            return 'Today';
        }
    }
    function sunHasSet() {
        const ss = props.luna.sunset.split(':');
        const sunset = moment().startOf('day')
            .hour(parseInt(ss[0]))
            .minute(parseInt(ss[1]));
        const s = moment.duration(sunset.diff(moment())).minutes();

        const sr = props.luna.sunrise.split(':');
        const sunrise = moment().startOf('day').hour(parseInt(sr[0])).minute(parseInt(sr[1]));
        const r = moment.duration(sunrise.diff(moment())).minutes();
        hasSunset = s <= 0 || r >= 0;
    }

    function darkness() {
        const day = moment('1970-01-01 23:59:59');
        const hours = day.subtract(timeToDecimal(props.luna.day_length), 'hours');
        return hours.format('HH:mm');
    }
    function isSunSet() {
        sunHasSet();
        if (hasSunset) {
            return 'Time til Sunrise';
        } else {
            return 'Time til Sunset';
        }
    }

    function riseSetClass() {
        if(hasSunset) {
            return "riseclr";
        } else {
            return "setclr"
        }
    }
    function sunBelow() {
        if (hasSunset) {
            return 'sunbelow';
        } else {
            return 'sunabove';
        }
    }

    return (
        <BoxData icon={'&#xF00D;'} title="Sun" unit="" style={{}}>
            <div className="sun-container">
                <div className="daylight">
                    <div><span className="riseclr">{ props.luna.day_length }</span> hrs</div>
                    <div>Total Daylight</div>
                </div>
                <div className="darkness">
                    <div><span className="setclr">{ darkness() }</span> hrs</div>
                    <div>Total Darkness</div>
                </div>
                <div className="remaining">
                    <div className="daylightvalue1">
                        <div>{ isSunSet() }</div>
                        <div><span className={riseSetClass()}>{ hoursTilSunSet() }</span>&nbsp;hrs &nbsp;
                            <span className={riseSetClass()}>{ minTilSunSet() }</span> min</div>
                    </div>
                </div>
                <div className="rise">
                    <div>Sun Rise</div>
                    <div>{ todayTomorrow('sunrise')}</div>
                    <div className="riseclr">{ props.luna.sunrise}</div>
                </div>
                <div className="set">
                    <div>Sun Set</div>
                    <div>{ todayTomorrow('sunset') }</div>
                    <div className="setclr">{ props.luna.sunset }</div>
                </div>
                <div className="elevation">
                    <div>Elevation:</div>
                    <div className={sunBelow()}>{ props.luna.sun_altitude.toFixed(2) }</div>
                </div>
            </div>
        </BoxData>
    )
}
export default Sun