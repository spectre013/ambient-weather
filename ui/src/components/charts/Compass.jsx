import React from "react";

export default function Compass({ dir, speed, size = 140 }) {
  const r = size / 2 - 6;
  const cx = size / 2, cy = size / 2;
  const cardinals = ["N", "E", "S", "W"];
  const angle = (deg) => (deg - 90) * (Math.PI / 180);

  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
      <circle cx={cx} cy={cy} r={r} fill="none" stroke="var(--rule)" strokeWidth="1" />
      <circle cx={cx} cy={cy} r={r - 8} fill="none" stroke="var(--rule-2)" strokeWidth="1" strokeDasharray="2 4" />
      {Array.from({ length: 12 }).map((_, i) => {
        const a = angle(i * 30);
        const r1 = r;
        const r2 = i % 3 === 0 ? r - 6 : r - 3;
        return (
          <line key={i}
            x1={cx + Math.cos(a) * r1} y1={cy + Math.sin(a) * r1}
            x2={cx + Math.cos(a) * r2} y2={cy + Math.sin(a) * r2}
            stroke="var(--muted)" strokeWidth={i % 3 === 0 ? 1.2 : 0.8} />
        );
      })}
      {cardinals.map((c, i) => {
        const a = angle(i * 90);
        const tr = r - 16;
        return (
          <text key={c}
            x={cx + Math.cos(a) * tr} y={cy + Math.sin(a) * tr + 3.5}
            textAnchor="middle" fontFamily="var(--font-mono)" fontSize="10"
            fill={c === "N" ? "var(--accent)" : "var(--muted)"} fontWeight="500" letterSpacing="0.05em">
            {c}
          </text>
        );
      })}
      <g transform={`rotate(${dir} ${cx} ${cy})`}>
        <line x1={cx} y1={cy} x2={cx} y2={cy - (r - 14)} stroke="var(--accent)" strokeWidth="2" strokeLinecap="round" />
        <polygon points={`${cx - 3},${cy - r + 14} ${cx + 3},${cy - r + 14} ${cx},${cy - r + 6}`} fill="var(--accent)" />
        <line x1={cx} y1={cy} x2={cx} y2={cy + (r - 24)} stroke="var(--ink-2)" strokeWidth="1.4" />
      </g>
      <circle cx={cx} cy={cy} r="3" fill="var(--ink)" />
    </svg>
  );
}
