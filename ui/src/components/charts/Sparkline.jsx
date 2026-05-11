import React from "react";

export default function Sparkline({
  values, width = 120, height = 36, color = "var(--ink-2)", fill = false, padding = 2,
}) {
  if (!values || values.length === 0) return null;
  const min = Math.min(...values);
  const max = Math.max(...values);
  const span = max - min || 1;
  const w = width - padding * 2;
  const h = height - padding * 2;
  const pts = values.map((v, i) => {
    const x = padding + (i / (values.length - 1 || 1)) * w;
    const y = padding + h - ((v - min) / span) * h;
    return [x, y];
  });
  const d = pts.map((p, i) => `${i ? "L" : "M"}${p[0].toFixed(1)},${p[1].toFixed(1)}`).join(" ");
  const area = `${d} L${pts[pts.length - 1][0]},${height} L${pts[0][0]},${height} Z`;
  return (
    <svg className="spark" viewBox={`0 0 ${width} ${height}`} preserveAspectRatio="none">
      {fill && <path d={area} fill={color} fillOpacity="0.12" />}
      <path d={d} fill="none" stroke={color} strokeWidth="1.4" strokeLinejoin="round" strokeLinecap="round" />
    </svg>
  );
}
