import React from "react";
import Card from "../Card";
import {timeAgo, fmtOr} from "../../utils/format.js";

export default function LightningCard({ lightning }) {
  const ago = fmtOr(lightning.time, timeAgo) // null until the first strike
  const rows = [
    ["Hour", lightning.hour, 5],
    ["Day", lightning.day, 5],
    ["Month", lightning.month, 1.25],
  ];
  return (
    <Card className="metric" label="Lightning" right={lightning.day === 0 ? "QUIET" : "ACTIVE"}>
      <div className="big">{lightning.day}<span className="unit"> strikes today</span></div>
      <div className="sub">Last Distance <b>{lightning.distance} mi</b> away</div>
        <div className="sub">Last Time {ago}</div>
      <div className="scale-track" style={{ marginTop: "auto", paddingTop: 6 }}>
        {rows.map(([lbl, v, scale]) => (
          <div key={lbl} className="row">
            <span className="lbl">{lbl}</span>
            <div className="bar"><i style={{ width: `${Math.min(100, v * scale)}%`, background: "var(--gold)" }} /></div>
            <span className="v">{v}</span>
          </div>
        ))}
      </div>
    </Card>
  );
}
