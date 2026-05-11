import React from "react";
import { fmtTime } from "../utils/format";

const connMap = {
  open: { label: "Live", color: "var(--green)" },
  connecting: { label: "Connecting", color: "var(--warn)" },
  closed: { label: "Reconnecting", color: "var(--warn)" },
  error: { label: "Offline", color: "var(--danger)" },
  offline: { label: "Fixture", color: "var(--muted)" },
  idle: { label: "Idle", color: "var(--muted)" },
};

export default function Masthead({ current, forecast, conn }) {
  const date = new Date(current.date);
  const ds = date.toLocaleDateString("en-US", {
    weekday: "long", month: "long", day: "numeric", year: "numeric",
  });
  const ts = fmtTime(current.date, forecast.tzoffset);
  const c = connMap[conn] || connMap.idle;
  const pulsing = conn === "connecting" || conn === "closed";
  return (
    <div className="masthead">
      <div className="left">
        <div className="vol">Weather Station · Vol 4 · {date.getUTCFullYear()}</div>
        <div className="title">Lorson Ranch</div>
      </div>
      <div className="center">
        <b>{forecast.resolvedAddress}</b><br />
        {forecast.latitude.toFixed(3)}°N · {Math.abs(forecast.longitude).toFixed(3)}°W · 5,730 FT
      </div>
      <div className="right">
        <div className="kv"><span className="k">As of</span><span className="v">{ts}</span></div>
        <div className="kv"><span className="k">Date</span><span className="v">{ds}</span></div>
        <div className="kv">
          <span className="k">Status</span>
          <span className="v" style={{ display: "flex", alignItems: "center", gap: 6 }}>
            <span className="dot live" style={{
              width: 6, height: 6, borderRadius: 999, background: c.color,
              display: "inline-block",
              animation: pulsing ? "pulse 1.4s ease-in-out infinite" : "none"
            }} />
            {c.label}
          </span>
        </div>
      </div>
    </div>
  );
}
