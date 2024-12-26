import BoxData from "./BoxData";
import {Luna} from "../models/Luna.ts";
import "./Sun.css"
import moment, {Moment} from "moment";
import {timeFormat} from "../util/weather.ts";
import {AstroData} from "../models/current.ts";

export interface Props {
    astro: AstroData
    units: string
}

export interface times {
    [key: string]: number; // Allows any property with a string key
}


const Sun = (props:Props) => {
    function setDateTime(hours: string) {
        const h = hours.split(':');
        return moment().startOf('day').hour(parseInt(h[0])).minute(parseInt(h[1]));
    }
    function todayTomorrow(type: string): string {
        const eventTime = props.astro[type as keyof AstroData];
        const event = setDateTime(eventTime as string);
        if (moment() > event) {
            return 'Tomorrow';
        } else {
            return 'Today';
        }
    }

    function hasSunSetText(hasSunSet: boolean) {
        if (hasSunSet) {
            return 'Sunrise in';
        } else {
            return 'Sunset in';
        }
    }

    function riseSetClass(hasSunSet: boolean) {
        if(hasSunSet) {
            return "riseclr";
        } else {
            return "setclr"
        }
    }
    function sunBelow(hasSunSet: boolean) {
        if (hasSunSet) {
            return 'sunbelow';
        } else {
            return 'sunabove';
        }
    }

    function nanosecondsToTime(nanoseconds: number) : {hours: number, minutes: number} {
        const seconds = nanoseconds / 1e9;
        const hours = Math.floor(seconds / 3600);
        const minutes = Math.floor((seconds % 3600) / 60);
        return {"hours": hours, "minutes": minutes};
    }

    function getFullTime(nanoseconds: number) :string {
        const nano = nanosecondsToTime(nanoseconds);
        return `${nano.hours} hrs ${nano.minutes} min`;
    }

    function durationToHoursMinutes(durationSeconds: number) :times {
        const hours = Math.floor(durationSeconds / 3600);
        const minutes = Math.floor((durationSeconds % 3600) / 60);
        return { hours, minutes };
    }

    function getTime(part: string,hasSunSet: boolean, luna: Luna) : number {
        let t: Moment = moment();
        const now: Moment = moment().utc();
        const sunrise = moment(luna.sunrise).utc();
        const sunset = moment(luna.sunset).utc();
        const sunriseTomorrow = moment(luna.sunriseTomorrow).utc();

        if (hasSunSet) {
            if (now.day() == sunrise.day()) {
                t = sunrise;
            } else {
                t = sunriseTomorrow;
            }
        }

        if (!hasSunSet) {
            t = sunset;
        }
        const times :times = durationToHoursMinutes(moment.duration(t.diff(now)).asSeconds())
        return times[part]
    }

    return (
        <BoxData icon="fa-sun" title="Sun" unit="" style={{}}>
            <div className="sun-container">
                <div className="daylight">
                    <div><span className="riseclr">{getFullTime(props.astro.daylight)}</span></div>
                    <div>Total Daylight</div>
                </div>
                <div className="darkness">
                    <div>
                        <span className="setclr">{getFullTime(props.astro.darkness)}</span>
                    </div>
                    <div>Total Darkness</div>
                </div>
                <div className="remaining">
                    <div className="daylightvalue1">
                        <div>{hasSunSetText(props.astro.hasSunset)}</div>
                        <div><span className={riseSetClass(props.astro.hasSunset)}>{getTime("hours", props.astro.hasSunset, props.astro)}</span> hrs&nbsp;
                            <span className={riseSetClass(props.astro.hasSunset)}>{getTime("minutes", props.astro.hasSunset, props.astro)}</span> min
                        </div>
                    </div>
                </div>
                <div className="rise">
                    <div>Sun Rise</div>
                    <div>{todayTomorrow("sunrise")}</div>
                    <div className="riseclr">{timeFormat(props.astro.sunrise)}</div>
                </div>
                <div className="set">
                    <div>Sun Set</div>
                    <div>{todayTomorrow("sunset")}</div>
                    <div className="setclr">{timeFormat(props.astro.sunrise)}</div>
                </div>
                <div className="elevation">
                    <div>Elevation:</div>
                    <div className={sunBelow(props.astro.hasSunset)}>{props.astro.elevation.toFixed(2)}</div>
                </div>
            </div>
        </BoxData>
    )
}
export default Sun