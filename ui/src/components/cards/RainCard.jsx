import React from "react";
import Card from "../Card";
import { showP, pUnit, pDecimals, parseLocalTimestamp, STATION_TZ } from "../../utils/format";

export default function RainCard({ rain, units = "F" }) {
  // lastrain is a `timestamp without time zone` column (station time).
  // It is null until the first rain has been recorded.
  const last = rain.lastrain ? parseLocalTimestamp(rain.lastrain) : null;
  const sinceLast = last ? `${Math.floor((Date.now() - last.getTime()) / 86400000)}d ago` : "—";
  const yearMax = Math.max(rain.yearly, 1);
  const u = pUnit(units);
  const dp = pDecimals(units);
  const rows = [
    ["Hour", rain.hourly], ["Day", rain.daily], ["Week", rain.weekly],
    ["Month", rain.monthly], ["Year", rain.yearly],
  ];
  return (
    <Card className="metric" label="Precipitation" right={`last rain ${sinceLast}`}>
      <div className="big">{showP(rain.daily,units).toFixed(pDecimals(units))}<span className="unit"> {u} today</span></div>
      <div className="sub">Rate {showP(rain.rate ?? 0,units).toFixed(pDecimals(units))} {u}/hr · Last event {showP(rain.event,units).toFixed(pDecimals(units))} {u} · {last ? last.toLocaleDateString("en-US", { timeZone: STATION_TZ, month: "short", day: "numeric" }) : "—"}</div>
      <div className="rain-rows" style={{ marginTop: "auto", paddingTop: 6 }}>
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
