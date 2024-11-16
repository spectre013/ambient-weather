import moment from "moment";
import clock from  '../assets/clock.svg';
import info from '../assets/info.svg';
import {useEffect, useState} from "react";
import Box from "./Box";


const Time = () => {
    const [currentDate, setCurrentDate] = useState(moment());

    function showTime() {
        setCurrentDate( moment());
    }

    useEffect(() => {
        setInterval(() => {
            showTime();
        }, 1000);
    });

    function dateFormat() {
        return currentDate.format('ddd MMM Do YYYY');
    }

    function timeFormat() {
        return currentDate.format('HH:mm:ss');
    }

    return (
        <Box icon={info} maintitle="Station" subtitle="time"  max="" height={10} width={10} class="orange">
            <div className="timeContainer">
                <div className="calendar34">
                    <img src={clock} alt="Clock" height={24} width={24} />
                </div>
                <div className="theTime">
                    <div className="weatherclock34">{ dateFormat() }
                    <div className="orangeclock">{ timeFormat() }</div>
                    </div>
                </div>
            </div>
        </Box>
    )
}
export default Time
