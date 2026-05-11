import React from "react";
import Card from "../Card";
import { showP, pUnit, pDecimals } from "../../utils/format";

export default function RainCard({ rain, units = "F" }) {
  const last = new Date(rain.lastrain);
  const sinceLast = Math.floor((Date.now() - last.getTime()) / 86400000);
  const yearMax = Math.max(rain.yearly, 1);
  const u = pUnit(units);
  const dp = pDecimals(units);
  const rows = [
    ["Hour", rain.hourly], ["Day", rain.daily], ["Week", rain.weekly],
    ["Month", rain.monthly], ["Year", rain.yearly],
  ];
  return (
    <Card label="Precipitation" right={`last rain ${sinceLast}d ago`}>
      <div className="big">{showP(rain.daily,units).toFixed(pDecimals(units))}<span className="unit"> {u} today</span></div>
      <div className="sub">Last event {showP(rain.event,units).toFixed(pDecimals(units))} {u} · {last.toLocaleDateString("en-US", { month: "short", day: "numeric" })}</div>
      <div className="rain-rows" style={{ marginTop: 12 }}>
        {rows.map(([lbl, v]) => (
          <div key={lbl} className="r">
            <span className="lbl">{lbl}</span>
            <div className="bar"><i style={{ width: `${Math.min(100, (v / yearMax) * 100)}%` }} /></div>
            <span className="v">{showP(v, units).toFixed(pDecimals(units))} {u}</span>
          </div>
        ))}
      </div>
    </Card>
  );
}
