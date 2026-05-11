import React from "react";
import Card from "../Card";
import MoonPhase from "../charts/MoonPhase";

export default function MoonCard({ astro, today }) {
  // Guard: bail until we have astro — the websocket payload may arrive
  // a tick after the initial render.
  if (!astro) return null;

  const illum = Math.round(astro.moonIlluminance ?? 0);
  const phase = astro.moonPhase ?? 0;
  const name  = astro.moonPhaseName ?? "—";

  return (
      <Card label="Moon" right={`${illum}% ILLUM.`}>
        <div style={{ display: "flex", alignItems: "center", gap: 14, marginBottom: 10 }}>
          <MoonPhase phase={phase} size={44} />
          <div style={{ minWidth: 0, flex: 1 }}>
            <div style={{ fontFamily: "var(--font-display)", fontSize: 22, lineHeight: 1.05, letterSpacing: "-0.01em" }}>
              {name}
            </div>
            <div className="sub" style={{ marginTop: 4 }}>
              Phase {phase.toFixed(2)} · {illum}% lit
            </div>
          </div>
        </div>
        <div className="scale-track">
          <div className="row" style={{ gridTemplateColumns: "80px 1fr", gap: 10 }}>
            <span className="lbl">Visibility</span>
            <span style={{ fontFamily: "var(--font-mono)", fontSize: 11, fontVariantNumeric: "tabular-nums" }}>{today.visibility} mi</span>
          </div>
          <div className="row" style={{ gridTemplateColumns: "80px 1fr", gap: 10 }}>
            <span className="lbl">Cloud cover</span>
            <span style={{ fontFamily: "var(--font-mono)", fontSize: 11, fontVariantNumeric: "tabular-nums" }}>{today.cloudcover.toFixed(0)}%</span>
          </div>
        </div>
      </Card>
  );
}
