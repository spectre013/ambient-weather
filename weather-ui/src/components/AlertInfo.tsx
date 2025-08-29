import { Alert } from '../models/Alert'
import {useEffect, useState} from "react";
import './AlertInfo.css'

export interface Props {
    alerts: Alert[]
}
const AlertInfo = (props: Props) => {
    const [currentAlert, setCurrentAlert] = useState(0);
    const [alert, setAlert] = useState({} as Alert);

    function changeAlert(dir: string) {
        console.log(dir, currentAlert, props.alerts.length);
        if(dir === "down" && currentAlert < props.alerts.length) {
            setCurrentAlert(currentAlert + 1)
        }
        if(dir === "up" && currentAlert > 0) {
            setCurrentAlert(currentAlert - 1)
        }
        setCurrent()
        return 0
    }

    function setCurrent() {
        setAlert(props.alerts[currentAlert])
    }

    useEffect(() => {
        setAlert(props.alerts[currentAlert])
    }, [props, currentAlert]);
    function alertColor() {
        if(!Object.hasOwn(alert,"event")) {
            return
        }
        if (alert.event.startsWith('911')) {
            return 'Telephone Outage 911'.replace(/\s+/g, '-');
        }
        return alert.event.toLowerCase().replace(/\s+/g, '-');
    }

    if(props.alerts.length == 0) {
        return (
            <div className="no-alerts">No Alerts!</div>
        )
    } else {
        return (
                <div className="alert-container">
                    <div className="up">
                        <div className="uparrow"  onClick={() => changeAlert("up")}>
                            {currentAlert > 0 && <span className="material-symbols-sharp">arrow_drop_up</span>}
                        </div>
                    </div>
                    <div className={`event ` + alertColor()}>
                        <span className="material-symbols-sharp">warning</span>
                             &nbsp;{ alert.event }&nbsp;
                        <span className="material-symbols-sharp">warning</span>
                    </div>
                    <div className="headline">
                        {alert.headline}
                    </div>
                    <div className="down">
                        <div className="downarrow" onClick={() => changeAlert("down")}>
                            {currentAlert + 1 < props.alerts.length &&  <span className="material-symbols-sharp">arrow_drop_down</span> }
                        </div>
                    </div>
                </div>
        )
    }
}
export default AlertInfo