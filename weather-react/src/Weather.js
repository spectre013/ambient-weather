import React, { useState, useEffect } from 'react';
import Time from "./elements/Time";
import SBox from "./containers/SBOX";
import RainSmall from "./elements/RainSmall";

const Weather = () => {
    const [conditions, setConditions] = useState([]);


    useEffect(() => {

        const ws = new WebSocket('ws://localhost:3000/api/ws');

        ws.onmessage = (event) => {
            const response = JSON.parse(event.data);
            setConditions(response);
        };
        ws.onclose = () => {
            ws.close();
        };

        return () => {
            ws.close();
        };
    }, []);

    return (
        <div>
            <div className="weather2-container">
                <SBox title={<> Weather Station <span className="orange"> Time </span></>}>
                    <Time />
                </SBox>
                <SBox title={<> <span className="oblue">Min</span> | <span className="ored">Max</span>  Temperatures</>}>

                </SBox>
                <SBox title={<> Rainfall<span className="oblue"> Data </span></>}>
                    <RainSmall conditions={conditions}  />
                </SBox>
                <SBox title={<> Rainfall<span className="blue"> Data </span></>}>

                </SBox>
            </div>
            {conditions.tempf}
        </div>
    );
};

export default Weather;
