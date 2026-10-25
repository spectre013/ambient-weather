import React from "react";
import { showT, tUnit, dirToCompass, showW, wUnit, showP, pUnit, pDecimals, fmtOr } from "../utils/format";

export default function QuickStats({ current, today, units }) {
  // Readings are null while their sensor is not reporting.
  const w = current.wind;
  const uv = current.uv.uv;
  const items = [
    { lbl: "Wind", v: fmtOr(w.windspeedmph, (v) => Math.round(showW(v, units))), unit: wUnit(units),
      sub: `${fmtOr(w.winddir, dirToCompass)} · gust ${fmtOr(w.windgustmph, (v) => Math.round(showW(v, units)))} ${wUnit(units)}` },
    { lbl: "Pressure", v: fmtOr(current.barometer.baromrelin, (v) => v.toFixed(2)), unit: "inHg",
      sub: current.barometer.trend?.trend || "" },
    { lbl: "Humidity", v: fmtOr(current.humidity.humidity), unit: "%",
      sub: `dew ${fmtOr(current.humidity.dewpoint, (v) => Math.round(showT(v, units)))}${tUnit(units)}` },
    { lbl: "UV", v: fmtOr(uv), unit: "",
      sub: uv == null ? "" : uv < 3 ? "Low" : uv < 6 ? "Moderate" : uv < 8 ? "High" : "Very high" },
    { lbl: "Precip", v: showP(current.rain.daily,units).toFixed(pDecimals(units)), unit: pUnit(units),
      sub: `${today.precipprob || 0}% chance today` },
    { lbl: "Air", v: fmtOr(current.aqi.pm25, (v) => v.toFixed(2)), unit: "µg",
      sub: current.aqi.pm25 == null ? "" : current.aqi.pm25 < 12 ? "Good" : current.aqi.pm25 < 35 ? "Moderate" : "Unhealthy" },
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
