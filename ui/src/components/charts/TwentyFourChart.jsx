import React from "react";
import { showT } from "../../utils/format";

/**
 * 24-hour temperature chart.
 *
 * Renders two temperature lines:
 *   - Forecast (dashed, muted) — full 24h, from the day's hourly forecast
 *   - Observed (solid, accent) — hours 00..now from the station's recorded history
 *
 * The current minute (from the live websocket `current` payload) is added as
 * the trailing point of the observed line so the chart stays current between
 * full-hour readings.
 */
export default function TwentyFourChart({ hours, history, current, nowHour, units = "F" }) {
  const W = 1000, H = 220;
  const PAD = { l: 40, r: 40, t: 18, b: 36 };
  const innerW = W - PAD.l - PAD.r;
  const innerH = H - PAD.t - PAD.b;
  // Forecast temps (display units) — drive the y-axis along with observed.
  const fTemps = hours.map((h) => showT(h.temp, units));

  // Observed: history.hours is keyed by "HH:00:00". Build a map for lookup.
  const obsByHour = new Map();
  if (history && Array.isArray(history.hours)) {
    for (const h of history.hours) {
      const hr = parseInt(h.datetime.slice(0, 2), 10);
      if (!Number.isNaN(hr)) obsByHour.set(hr, showT(h.temp, units));
    }
  }

  // Live minute reading — append as the trailing observed point at fractional
  // x-position `nowHour + minute/60`. Falls back to last completed hour if
  // current is missing.
  let liveX = null, liveY = null, liveTemp = null;
  if (current && current.temp && typeof current.temp.temp === "number") {
    liveTemp = showT(current.temp.temp, units);
    const d = new Date(current.date);
    // nowHour was computed as UTC hours + tzOffsetH; do the same for minutes.
    const localHour = nowHour;
    const minuteFrac = d.getUTCMinutes() / 60;
    liveX = localHour + minuteFrac;
  }

  // Y-scale: include forecast, observations, and the live point so nothing
  // clips when the live reading runs hotter than the forecast.
  const allTemps = [
    ...fTemps,
    ...Array.from(obsByHour.values()),
    ...(liveTemp !== null ? [liveTemp] : []),
  ];
  const tmin = Math.floor(Math.min(...allTemps) - 2);
  const tmax = Math.ceil(Math.max(...allTemps) + 2);
  const tspan = tmax - tmin || 1;

  const precip = hours.map((h) => h.precipprob || 0);

  const xAt = (i) => PAD.l + (i / (hours.length - 1)) * innerW;
  const yAt = (t) => PAD.t + innerH - ((t - tmin) / tspan) * innerH;

  // Forecast line covering all 24 hours.
  const fPts = fTemps.map((t, i) => [xAt(i), yAt(t)]);
  const fLineD = fPts.map((p, i) => `${i ? "L" : "M"}${p[0].toFixed(1)},${p[1].toFixed(1)}`).join(" ");

  // Observed line: every completed hour we have, in order, plus the live
  // reading on the trailing edge.
  const obsHoursSorted = [...obsByHour.keys()].sort((a, b) => a - b);
  const obsPts = obsHoursSorted.map((hr) => [xAt(hr), yAt(obsByHour.get(hr))]);
  if (liveX !== null) {
    liveY = yAt(liveTemp);
    // Only append if it advances past the last hourly point.
    const lastHr = obsHoursSorted.length ? obsHoursSorted[obsHoursSorted.length - 1] : -1;
    if (liveX > lastHr) obsPts.push([xAt(liveX), liveY]);
  }
  const obsLineD = obsPts.length
      ? obsPts.map((p, i) => `${i ? "L" : "M"}${p[0].toFixed(1)},${p[1].toFixed(1)}`).join(" ")
      : "";
  const obsAreaD = obsPts.length
      ? `${obsLineD} L${obsPts[obsPts.length - 1][0]},${PAD.t + innerH} L${obsPts[0][0]},${PAD.t + innerH} Z`
      : "";

  // Y-axis ticks.
  const ticks = [];
  const stepCandidates = [5, 10, 15, 20];
  let step = 10;
  for (const c of stepCandidates) if (tspan / c <= 6) { step = c; break; }
  for (let v = Math.ceil(tmin / step) * step; v <= tmax; v += step) ticks.push(v);

  const xLabels = hours.map((h, i) => ({ i, label: h.datetime.slice(0, 2) })).filter((_, i) => i % 3 === 0);
  const nowIdx = hours.findIndex((h) => h.datetime.startsWith(String(nowHour).padStart(2, "0")));
  const nowX = liveX !== null ? xAt(liveX) : (nowIdx >= 0 ? xAt(nowIdx) : null);

  return (
      <svg className="chart-svg" viewBox={`0 0 ${W} ${H}`} preserveAspectRatio="none">
        {ticks.map((t) => (
            <g key={t}>
              <line x1={PAD.l} x2={W - PAD.r} y1={yAt(t)} y2={yAt(t)} stroke="var(--rule-2)" strokeDasharray="2 4" />
              <text x={PAD.l - 8} y={yAt(t) + 3} textAnchor="end" fill="var(--muted)" fontFamily="var(--font-mono)" fontSize="10">{t}°</text>
            </g>
        ))}
        {nowX !== null && (
            <rect x={PAD.l} y={PAD.t} width={nowX - PAD.l} height={innerH} fill="var(--ink)" fillOpacity="0.04" />
        )}
        {precip.map((p, i) => {
          if (!p) return null;
          const bw = (innerW / hours.length) * 0.6;
          const x = xAt(i) - bw / 2;
          const bh = (p / 100) * innerH * 0.4;
          return <rect key={i} x={x} y={PAD.t + innerH - bh} width={bw} height={bh} fill="var(--cool)" fillOpacity="0.5" rx="1" />;
        })}

        {/* Forecast: dashed, muted — the reference trace */}
        <path d={fLineD} fill="none" stroke="var(--muted)" strokeWidth="1.2" strokeDasharray="3 3" strokeLinejoin="round" opacity="0.7" />
        {fPts.map((p, i) => i % 3 === 0 && (
            <text key={`fl-${i}`} x={p[0]} y={p[1] - 8} textAnchor="middle" fill="var(--muted)" fontFamily="var(--font-mono)" fontSize="9" opacity="0.7">
              {Math.round(fTemps[i])}°
            </text>
        ))}

        {/* Observed: solid accent line + soft fill */}
        {obsAreaD && <path d={obsAreaD} fill="var(--accent)" fillOpacity="0.12" />}
        {obsLineD && <path d={obsLineD} fill="none" stroke="var(--accent)" strokeWidth="2" strokeLinejoin="round" strokeLinecap="round" />}
        {obsHoursSorted.map((hr) => {
          const x = xAt(hr);
          const y = yAt(obsByHour.get(hr));
          return (
              <g key={`obs-${hr}`}>
                <circle cx={x} cy={y} r="2.2" fill="var(--accent)" />
                {hr % 3 === 0 && (
                    <text x={x} y={y - 8} textAnchor="middle" fill="var(--ink)" fontFamily="var(--font-mono)" fontSize="10" fontWeight="500">
                      {Math.round(obsByHour.get(hr))}°
                    </text>
                )}
              </g>
          );
        })}

        {/* NOW marker + live reading dot */}
        {nowX !== null && (
            <g>
              <line x1={nowX} y1={PAD.t} x2={nowX} y2={PAD.t + innerH} stroke="var(--accent)" strokeWidth="1.2" strokeDasharray="3 3" />
              <text x={nowX} y={PAD.t - 6} textAnchor="middle" fontFamily="var(--font-mono)" fontSize="9" fill="var(--accent)" letterSpacing="0.1em">NOW</text>
            </g>
        )}
        {liveX !== null && liveY !== null && (
            <g>
              <circle cx={xAt(liveX)} cy={liveY} r="6" fill="var(--accent)" fillOpacity="0.18" />
              <circle cx={xAt(liveX)} cy={liveY} r="3.2" fill="var(--accent)" stroke="var(--surface)" strokeWidth="1.2">
                <animate attributeName="r" values="3.2;4.2;3.2" dur="2.4s" repeatCount="indefinite" />
              </circle>
              <text x={xAt(liveX)} y={liveY - 10} textAnchor="middle" fill="var(--accent)" fontFamily="var(--font-mono)" fontSize="11" fontWeight="600">
                {liveTemp.toFixed(1)}°
              </text>
            </g>
        )}

        {xLabels.map(({ i, label }) => (
            <text key={i} x={xAt(i)} y={H - PAD.b + 16} textAnchor="middle" fill="var(--muted)" fontFamily="var(--font-mono)" fontSize="10">{label}:00</text>
        ))}
        <text x={PAD.l} y={H - 6} fill="var(--faint)" fontFamily="var(--font-mono)" fontSize="9" letterSpacing="0.1em">TEMPERATURE °{units}</text>
        <text x={W - PAD.r} y={H - 6} textAnchor="end" fill="var(--faint)" fontFamily="var(--font-mono)" fontSize="9" letterSpacing="0.1em">PRECIP PROBABILITY</text>
      </svg>
  );
}
