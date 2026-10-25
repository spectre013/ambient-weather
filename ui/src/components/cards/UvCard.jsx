import React from "react";
import Card from "../Card";
import { clamp, statVal, fmtOr } from "../../utils/format";

const COLORS = ["#5db66d", "#c9d23a", "#f5b324", "#e87f2b", "#d34a4a", "#9c3aab"];

export default function UvCard({ uv }) {
  const idx = uv.uv; // null while the WS80 is not reporting
  const cat = idx == null ? "No data" : idx < 3 ? "Low" : idx < 6 ? "Moderate" : idx < 8 ? "High" : idx < 11 ? "Very high" : "Extreme";
  const colorIdx = clamp(Math.floor((idx ?? 0) / 2), 0, 5);
  return (
    <Card className="metric" label="UV index" right={cat.toUpperCase()}>
      <div className="big" style={{ color: COLORS[colorIdx] }}>{fmtOr(idx)}<span className="unit" style={{ color: "var(--muted)" }}>/11</span></div>
      <div className="sub">Solar radiation <b>{fmtOr(uv.solarradiation, Math.round)} W/m²</b> · day max {fmtOr(statVal(uv.minmax, "max", "day"))}</div>
      <div style={{ display: "grid", gridTemplateColumns: "repeat(11, 1fr)", gap: 2, marginTop: "auto", paddingTop: 6 }}>
        {Array.from({ length: 11 }).map((_, i) => (
          <div key={i} style={{
            height: 18,
            background: idx != null && i <= idx ? COLORS[Math.floor(i / 2)] : "var(--rule-2)",
            opacity: idx != null && i <= idx ? 1 : 0.6,
            borderRadius: 1
          }} />
        ))}
      </div>
      <div style={{ display: "flex", justifyContent: "space-between", fontFamily: "var(--font-mono)", fontSize: 9, color: "var(--faint)", letterSpacing: "0.1em", marginTop: 4 }}>
        <span>0</span><span>3</span><span>6</span><span>8</span><span>11+</span>
      </div>
    </Card>
  );
}
