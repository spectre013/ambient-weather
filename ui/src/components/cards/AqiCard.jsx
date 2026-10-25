import React from "react";
import Card from "../Card";
import { statVal, fmtOr } from "../../utils/format";

export default function AqiCard({ aqi }) {
  const pm = aqi.pm25; // null while the WH0290 is not reporting
  const cat = pm == null ? "No data" : pm < 12 ? "Good" : pm < 35 ? "Moderate" : pm < 55 ? "Unhealthy SG" : pm < 150 ? "Unhealthy" : "Hazardous";
  const colors = ["var(--green)", "var(--gold)", "var(--warm)", "var(--danger)"];
  const colorIdx = pm == null ? 0 : pm < 12 ? 0 : pm < 35 ? 1 : pm < 55 ? 2 : 3;
  const bands = [
    ["Good", 0, 12, "var(--green)"],
    ["Mod.", 12, 35, "var(--gold)"],
    ["USG", 35, 55, "var(--warm)"],
    ["Unh.", 55, 150, "var(--danger)"],
  ];
  return (
    <Card className="metric" label="Air quality" right={cat.toUpperCase()}>
      <div className="big" style={{ color: pm == null ? "var(--muted)" : colors[colorIdx] }}>{fmtOr(pm, (v) => v.toFixed(1))}<span className="unit" style={{ color: "var(--muted)" }}>µg/m³ PM2.5</span></div>
      <div className="sub">24h avg <b>{fmtOr(aqi.pm2524h, (v) => v.toFixed(1))}</b> · month max {fmtOr(statVal(aqi.minmax, "max", "month"))}</div>
      <div className="scale-track" style={{ marginTop: "auto", paddingTop: 6 }}>
        {bands.map(([lbl, lo, hi, c]) => {
          const active = pm != null && pm >= lo && pm < hi;
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
