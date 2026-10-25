import React from "react";
import Card from "../Card";
import { statVal, fmtOr } from "../../utils/format";

export default function HumidityCard({ humidity }) {
  const h = humidity.humidity;
  const dew = humidity.dewpoint;
  const comfort = h == null ? "No data" : h < 30 ? "Dry" : h > 60 ? "Humid" : "Comfortable";
  const hp = h ?? 0; // meter position while the sensor is not reporting
  return (
    <Card className="metric" label="Humidity" right={comfort.toUpperCase()}>
      <div className="big">{fmtOr(h)}<span className="unit">%</span></div>
      <div className="sub">Dew point <b>{fmtOr(dew, Math.round)}°F</b></div>
      <div className="sub">24h range {fmtOr(statVal(humidity.minmax, "min", "day"))}–{fmtOr(statVal(humidity.minmax, "max", "day"))}%</div>
        <div className="meter" style={{ marginTop: "auto", position: "relative" }}>
            <span style={{ width: `${hp}%`, background: h < 30 ? "var(--warm)" : h > 60 ? "var(--cool)" : "var(--green)" }} />
            <i style={{ position: "absolute", left: "30%", top: -2, bottom: -2, width: 1, background: "var(--rule)" }} />
            <i style={{ position: "absolute", left: "60%", top: -2, bottom: -2, width: 1, background: "var(--rule)" }} />
        </div>
        <div style={{ position: "relative", height: 14, fontFamily: "var(--font-mono)", fontSize: 9, color: "var(--faint)", letterSpacing: "0.1em", marginTop: 4 }}>
            <span style={{ position: "absolute", left: 0, top: 0 }}>0%</span>
            <span style={{ position: "absolute", left: "30%", top: 0, transform: "translateX(-50%)" }}>30 DRY</span>
            <span style={{ position: "absolute", left: "60%", top: 0, transform: "translateX(-50%)" }}>60 HUMID</span>
            <span style={{ position: "absolute", right: 0, top: 0 }}>100%</span>
        </div>
    </Card>
  );
}
