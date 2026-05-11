import React from "react";
import Card from "../Card";

export default function HumidityCard({ humidity }) {
  const h = humidity.humidity;
  const dew = humidity.dewpoint;
  const comfort = h < 30 ? "Dry" : h > 60 ? "Humid" : "Comfortable";
  return (
    <Card label="Humidity" right={comfort.toUpperCase()}>
      <div className="big">{h}<span className="unit">%</span></div>
      <div className="sub">Dew point <b>{Math.round(dew)}°F</b></div>
      <div className="sub">24h range {humidity.minmax.min.day.value}–{humidity.minmax.max.day.value}%</div>
        <div className="meter" style={{ marginTop: 12, position: "relative" }}>
            <span style={{ width: `${h}%`, background: h < 30 ? "var(--warm)" : h > 60 ? "var(--cool)" : "var(--green)" }} />
            <i style={{ position: "absolute", left: "30%", top: -2, bottom: -2, width: 1, background: "var(--rule)" }} />
            <i style={{ position: "absolute", left: "60%", top: -2, bottom: -2, width: 1, background: "var(--rule)" }} />
        </div>
        <div style={{ position: "relative", height: 14, fontFamily: "var(--font-mono)", fontSize: 9, color: "var(--faint)", letterSpacing: "0.1em", marginTop: 4 }}>
            <span style={{ position: "absolute", left: 0, top: 0 }}>0%</span>
            <span style={{ position: "absolute", left: "30%", top: 0, transform: "translateX(-50%)" }}>30 DRY</span>
            <span style={{ position: "absolute", left: "60%", top: 0, transform: "translateX(-50%)" }}>60 HUMID</span>
            <span style={{ position: "absolute", right: 0, top: 0 }}>100%</span>
            <span style={{ position: "absolute", left: `${h}%`, top: `10px`, transform: "translateX(-50%)", color: "var(--ink)", fontWeight: 500 }}>{h}%</span>
        </div>
    </Card>
  );
}
