import React from "react";
import Card from "../Card";

export default function AqiCard({ aqi }) {
  const pm = aqi.pm25;
  const cat = pm < 12 ? "Good" : pm < 35 ? "Moderate" : pm < 55 ? "Unhealthy SG" : pm < 150 ? "Unhealthy" : "Hazardous";
  const colors = ["var(--green)", "var(--gold)", "var(--warm)", "var(--danger)"];
  const colorIdx = pm < 12 ? 0 : pm < 35 ? 1 : pm < 55 ? 2 : 3;
  const bands = [
    ["Good", 0, 12, "var(--green)"],
    ["Mod.", 12, 35, "var(--gold)"],
    ["USG", 35, 55, "var(--warm)"],
    ["Unh.", 55, 150, "var(--danger)"],
  ];
  return (
    <Card label="Air quality" right={cat.toUpperCase()}>
      <div className="big" style={{ color: colors[colorIdx] }}>{pm.toFixed(1)}<span className="unit" style={{ color: "var(--muted)" }}>µg/m³ PM2.5</span></div>
      <div className="sub">24h avg <b>{aqi.pm2524h.toFixed(1)}</b> · month max {aqi.minmax.max.month.value}</div>
      <div className="scale-track" style={{ marginTop: 12 }}>
        {bands.map(([lbl, lo, hi, c]) => {
          const active = pm >= lo && pm < hi;
          return (
            <div key={lbl} className="row">
              <span className="lbl">{lbl}</span>
              <div className="bar" style={{ background: active ? c : "var(--rule-2)", opacity: active ? 1 : 0.4 }} />
              <span className="v" style={{ fontSize: 9, color: "var(--faint)" }}>{lo}–{hi}</span>
            </div>
          );
        })}
      </div>
    </Card>
  );
}
