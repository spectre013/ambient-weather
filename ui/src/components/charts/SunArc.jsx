import React from "react";

export default function SunArc({ progress = 0.5, width = 240, height = 64 }) {
  const cx = width / 2;
  const cy = height - 4;
  const r = Math.min(width / 2 - 8, height - 12);
  const angle = Math.PI - progress * Math.PI;
  const x = cx + Math.cos(angle) * r;
  const y = cy - Math.sin(angle) * r;
  const arcD = `M ${cx - r} ${cy} A ${r} ${r} 0 0 1 ${cx + r} ${cy}`;
  return (
    <svg width="100%" height={height} viewBox={`0 0 ${width} ${height}`}>
      <path d={arcD} fill="none" stroke="var(--rule)" strokeWidth="1" strokeDasharray="2 3" />
      <line x1={cx - r} y1={cy} x2={cx + r} y2={cy} stroke="var(--rule)" strokeWidth="1" />
      <circle cx={x} cy={y} r="6" fill="var(--accent)" />
      <circle cx={x} cy={y} r="11" fill="none" stroke="var(--accent)" strokeWidth="1" strokeOpacity="0.5" />
    </svg>
  );
}
