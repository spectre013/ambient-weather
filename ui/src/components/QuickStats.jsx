import React from "react";
import { showT, tUnit, dirToCompass, showW, wUnit, showP, pUnit, pDecimals } from "../utils/format";

export default function QuickStats({ current, today, units }) {
  const w = current.wind;
  const items = [
    { lbl: "Wind", v: Math.round(showW(w.windspeedmph, units)), unit: wUnit(units),
      sub: `${dirToCompass(w.winddir)} · gust ${Math.round(showW(w.windgustmph, units))} ${wUnit(units)}` },
    { lbl: "Pressure", v: current.barometer.baromrelin.toFixed(2), unit: "inHg",
      sub: current.barometer.trend.trend },
    { lbl: "Humidity", v: current.humidity.humidity, unit: "%",
      sub: `dew ${Math.round(showT(current.humidity.dewpoint, units))}${tUnit(units)}` },
    { lbl: "UV", v: current.uv.uv, unit: "",
      sub: current.uv.uv < 3 ? "Low" : current.uv.uv < 6 ? "Moderate" : current.uv.uv < 8 ? "High" : "Very high" },
    { lbl: "Precip", v: showP(current.rain.daily,units).toFixed(pDecimals(units)), unit: pUnit(units),
      sub: `${today.precipprob || 0}% chance today` },
    { lbl: "Air", v: current.aqi.pm25.toFixed(2), unit: "µg",
      sub: current.aqi.pm25 < 12 ? "Good" : current.aqi.pm25 < 35 ? "Moderate" : "Unhealthy" },
  ];
  return (
    <div className="card quick-card">
      <div className="label"><span>At a glance</span><span className="right">live</span></div>
      <div className="quick-grid">
        {items.map((it) => (
          <div className="quick-cell" key={it.lbl}>
            <div className="ql">{it.lbl}</div>
            <div className="qv">
              <span className="qn">{it.v}</span>
              {it.unit && <span className="qu">{it.unit}</span>}
            </div>
            <div className="qs">{it.sub}</div>
          </div>
        ))}
      </div>
    </div>
  );
}
