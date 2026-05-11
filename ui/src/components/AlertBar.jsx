import React from "react";
import { slugifyEvent } from "../utils/format";

export default function AlertBar({ alerts }) {
  const a = alerts && alerts[0];
  if (!a) return null;
  const ends = new Date(a.end || a.expires).toLocaleString("en-US", {
    month: "short", day: "numeric", hour: "numeric", minute: "2-digit",
  });
  const cls = slugifyEvent(a.event);
  return (
    <div className={`alert nws-${cls}`}>
      <div className="sev">{a.severity}</div>
      <div>
        <div className="title">{a.event} — {a.senderName}</div>
        <div className="meta">In effect until {ends} · {a.urgency} · {a.certainty}</div>
      </div>
      <button className="cta" onClick={() => alert(a.description)}>Read advisory →</button>
    </div>
  );
}
