import React from "react";
import { showT } from "../../utils/format";

export default function TwentyFourChart({ hours, nowHour, units = "F" }) {
  const W = 1000, H = 220;
  const PAD = { l: 40, r: 40, t: 18, b: 36 };
  const innerW = W - PAD.l - PAD.r;
  const innerH = H - PAD.t - PAD.b;

  // Convert temps once, up front — everything below works in display units.
  const temps = hours.map((h) => showT(h.temp, units));
  const tmin = Math.floor(Math.min(...temps) - 2);
  const tmax = Math.ceil(Math.max(...temps) + 2);
  const tspan = tmax - tmin || 1;

  const precip = hours.map((h) => h.precipprob || 0);

  const xAt = (i) => PAD.l + (i / (hours.length - 1)) * innerW;
  const yAt = (t) => PAD.t + innerH - ((t - tmin) / tspan) * innerH;

  const linePts = temps.map((t, i) => [xAt(i), yAt(t)]);
  const lineD = linePts.map((p, i) => `${i ? "L" : "M"}${p[0].toFixed(1)},${p[1].toFixed(1)}`).join(" ");
  const areaD = `${lineD} L${linePts[linePts.length - 1][0]},${PAD.t + innerH} L${linePts[0][0]},${PAD.t + innerH} Z`;

  const ticks = [];
  const stepCandidates = [5, 10, 15, 20];
  let step = 10;
  for (const c of stepCandidates) if (tspan / c <= 6) { step = c; break; }
  for (let v = Math.ceil(tmin / step) * step; v <= tmax; v += step) ticks.push(v);

  const xLabels = hours.map((h, i) => ({ i, label: h.datetime.slice(0, 2) })).filter((_, i) => i % 3 === 0);
  const nowIdx = hours.findIndex((h) => h.datetime.startsWith(String(nowHour).padStart(2, "0")));

  return (
      <svg className="chart-svg" viewBox={`0 0 ${W} ${H}`} preserveAspectRatio="none">
        {ticks.map((t) => (
            <g key={t}>
              <line x1={PAD.l} x2={W - PAD.r} y1={yAt(t)} y2={yAt(t)} stroke="var(--rule-2)" strokeDasharray="2 4" />
              <text x={PAD.l - 8} y={yAt(t) + 3} textAnchor="end" fill="var(--muted)" fontFamily="var(--font-mono)" fontSize="10">{t}°</text>
            </g>
        ))}
        {nowIdx >= 0 && (
            <rect x={PAD.l} y={PAD.t} width={xAt(nowIdx) - PAD.l} height={innerH} fill="var(--ink)" fillOpacity="0.04" />
        )}
        {precip.map((p, i) => {
          if (!p) return null;
          const bw = (innerW / hours.length) * 0.6;
          const x = xAt(i) - bw / 2;
          const bh = (p / 100) * innerH * 0.4;
          return <rect key={i} x={x} y={PAD.t + innerH - bh} width={bw} height={bh} fill="var(--cool)" fillOpacity="0.5" rx="1" />;
        })}
        <path d={areaD} fill="var(--accent)" fillOpacity="0.12" />
        <path d={lineD} fill="none" stroke="var(--accent)" strokeWidth="1.8" strokeLinejoin="round" />
        {linePts.map((p, i) => i % 3 === 0 && (
            <g key={i}>
              <circle cx={p[0]} cy={p[1]} r="2" fill="var(--accent)" />
              <text x={p[0]} y={p[1] - 8} textAnchor="middle" fill="var(--ink)" fontFamily="var(--font-mono)" fontSize="10" fontWeight="500">
                {Math.round(temps[i])}°
              </text>
            </g>
        ))}
        {nowIdx >= 0 && (
            <g>
              <line x1={xAt(nowIdx)} y1={PAD.t} x2={xAt(nowIdx)} y2={PAD.t + innerH} stroke="var(--accent)" strokeWidth="1.2" strokeDasharray="3 3" />
              <text x={xAt(nowIdx)} y={PAD.t - 6} textAnchor="middle" fontFamily="var(--font-mono)" fontSize="9" fill="var(--accent)" letterSpacing="0.1em">NOW</text>
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
