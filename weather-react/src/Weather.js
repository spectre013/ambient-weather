import React, { useState, useEffect } from 'react';
import Time from "./Time";

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
                <Time />
            </div>
        </div>
    );
};

export default Weather;