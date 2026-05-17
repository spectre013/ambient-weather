import React, { useEffect, useState } from "react";
import { slugifyEvent } from "../utils/format";

// Rotates through every alert in the array. Auto-advances every 8 s,
// pauses on hover, and gives the user prev/next + a position pill.
export default function AlertBar({ alerts, interval = 8000 }) {
    const list = (alerts || []).filter(Boolean);
    const [idx, setIdx] = useState(0);
    const [paused, setPaused] = useState(false);

    // Clamp index if the alerts array changes underneath us.
    useEffect(() => {
        if (idx >= list.length) setIdx(0);
    }, [list.length, idx]);

    // Auto-advance.
    useEffect(() => {
        if (list.length < 2 || paused) return;
        const t = setInterval(() => setIdx((i) => (i + 1) % list.length), interval);
        return () => clearInterval(t);
    }, [list.length, paused, interval]);

    if (!list.length) return null;
    const a = list[idx];
    const ends = new Date(a.end || a.expires).toLocaleString("en-US", {
        month: "short", day: "numeric", hour: "numeric", minute: "2-digit",
    });
    const cls = slugifyEvent(a.event);
    const multi = list.length > 1;
    const go = (n) => setIdx((i) => (i + n + list.length) % list.length);

    return (
        <div
            className={`alert nws-${cls}`}
            onMouseEnter={() => setPaused(true)}
            onMouseLeave={() => setPaused(false)}
        >
            <div className="sev">{a.severity}</div>
            <div className="alert-body">
                <div className="title">{a.event} — {a.senderName}</div>
                <div className="meta">In effect until {ends} · {a.urgency} · {a.certainty}</div>
                {a.summary && <p className="alert-summary">{a.summary}</p>}
            </div>
            <div className="alert-actions">
                {multi && (
                    <div className="alert-pager" role="group" aria-label="Alerts">
                        <button className="pg" onClick={() => go(-1)} aria-label="Previous alert">‹</button>
                        <span className="pg-count">{idx + 1}<span className="sep">/</span>{list.length}</span>
                        <button className="pg" onClick={() => go(1)} aria-label="Next alert">›</button>
                    </div>
                )}
                <button className="cta" onClick={() => alert(a.description)}>Read advisory →</button>
            </div>
        </div>
    );
}
