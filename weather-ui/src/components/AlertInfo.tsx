import { Alert } from '../models/Alert'
import {useEffect, useState} from "react";
import './AlertInfo.css'
import BoxData from "./BoxData";

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
            <BoxData icon={'&#xF0B1;'} title="ALERTS" unit="">
                <div className="no-alerts">No Alerts!</div>
            </BoxData>
        )
    } else {
        return (
            <BoxData icon={'&#xF0B1;'} title="ALERTS" unit="">
                <div className="alert-box">
                    <div className="alert-container">
                        <div className="up">
                            <div className="uparrow"  onClick={() => changeAlert("up")}>
                                {currentAlert > 0 && <svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="16" height="16">
                                    <path d="m 13,6 -5,5 -5,-5 z" fill="#797979" />
                                </svg>}
                            </div>
                        </div>
                        <div className={`event ` + alertColor()}>
                            <svg className="alertsvg" viewBox="0 0 32 32"  fill="none" stroke="currentcolor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2">
                                <path d="M16 3 L30 29 2 29 Z M16 11 L16 19 M16 23 L16 25"></path>
                            </svg>
                            &nbsp;
                            { alert.event }
                            &nbsp;
                            <svg className="alertsvg" viewBox="0 0 32 32"  fill="none" stroke="currentcolor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2">
                                <path d="M16 3 L30 29 2 29 Z M16 11 L16 19 M16 23 L16 25"></path>
                            </svg>
                        </div>
                        <div className="headline">
                            <div>{alert.headline}</div>

                        </div>
                        <div className="down">
                            <div className="downarrow" onClick={() => changeAlert("down")}>
                                {currentAlert + 1 < props.alerts.length &&  <svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="16" height="16">
                                    <path d="m 13,6 -5,5 -5,-5 z" fill="#797979" />
                                </svg> }
                            </div>
                        </div>
                    </div>
                </div>
            </BoxData>
        )
    }
}
export default AlertInfo