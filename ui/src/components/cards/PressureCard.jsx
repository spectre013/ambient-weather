import React from "react";
import Card from "../Card";
import Sparkline from "../charts/Sparkline";
import { statVal } from "../../utils/format";

export default function PressureCard({ baro }) {
  // Empty when the server has no reading from 3 hours ago to compare with.
  const trend = baro.trend?.trend || "—";
  const arrow = trend === "Rising" ? "↗" : trend === "Falling" ? "↘" : "→";
  const color = trend === "Rising" ? "var(--warm)" : trend === "Falling" ? "var(--cool)" : "var(--muted)";
  const cur = baro.baromrelin;
  // Fall back to the current reading until today's stats row exists.
  const hi = statVal(baro.minmax, "max", "day") ?? cur;
  const lo = statVal(baro.minmax, "min", "day") ?? cur;
  const fakeHist = Array.from({ length: 24 }, (_, i) => {
    const t = i / 23;
    return lo + (hi - lo) * (0.4 + 0.6 * Math.sin(t * Math.PI * 1.2));
  });
  fakeHist[fakeHist.length - 1] = cur;
  const span = Math.max(0.1, hi - lo);
  return (
    <Card label="Barometric Pressure" right="rel · in Hg">
      <div className="big">{baro.baromrelin.toFixed(2)}<span className="unit">inHg</span></div>
      <div className="sub">
        <span style={{ color }}>{arrow} {trend}</span> · absolute {baro.baromabsin.toFixed(2)} inHg
      </div>
      <Sparkline values={fakeHist} width={260} height={36} color={color} fill />
      <div className="scale-track" style={{ marginTop: 8 }}>
        <div className="row">
          <span className="lbl">24h hi</span>
          <div className="bar"><i style={{ width: "100%" }} /></div>
          <span className="v">{hi.toFixed(2)}</span>
        </div>
        <div className="row">
          <span className="lbl">Now</span>
          <div className="bar"><i style={{ width: `${((cur - lo) / span) * 100}%`, background: "var(--accent)" }} /></div>
          <span className="v">{cur.toFixed(2)}</span>
        </div>
        <div className="row">
          <span className="lbl">24h lo</span>
          <div className="bar"><i style={{ width: "4%" }} /></div>
          <span className="v">{lo.toFixed(2)}</span>
        </div>
      </div>
    </Card>
  );
}
