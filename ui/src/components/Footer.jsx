import React from "react";

export default function Footer({ forecast }) {
  return (
    <div className="footer">
      <div className="col">
        <span className="h">Station</span>
        <span>Ambient WS-5000 · v4.3.7</span>
        <span></span>
        <span>Github - github.com/spectre013/ambient-weather</span>
      </div>
      <div className="col">
        <span className="h">Network</span>
        <span>NWS, METAR & nearby ASOS</span>
        <span>{forecast.stations}</span>
        <span>Forecast © Visual Crossing · NOAA NDFD</span>
      </div>
      <div className="col">
        <span className="h">Services</span>
        <span>UI · 4.0</span>   <span>Receiver · 1.0</span>
        <span>Server · 4.1.0</span> <span>Processor · 2.3</span>
      </div>
        <div>
            &copy; zoms.net - 2018 - {new Date().getFullYear()}
        </div>
    </div>
  );
}
