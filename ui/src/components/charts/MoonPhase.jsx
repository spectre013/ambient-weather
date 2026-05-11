import React from "react";

export default function MoonPhase({ phase = 0, size = 36 }) {
  const cx = size / 2, cy = size / 2;
  const r = size / 2 - 2;
  const clipId = `moon-clip-${size}`;
  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
      <circle cx={cx} cy={cy} r={r} fill="var(--rule-2)" />
      <defs>
        <clipPath id={clipId}>
          <circle cx={cx} cy={cy} r={r} />
        </clipPath>
      </defs>
      {phase < 0.5 ? (
        <g clipPath={`url(#${clipId})`}>
          <rect x={cx} y={0} width={cx} height={size} fill="var(--ink)" />
          <ellipse cx={cx} cy={cy} rx={Math.abs((0.5 - phase) * 2 * r)} ry={r} fill="var(--rule-2)" />
        </g>
      ) : (
        <g clipPath={`url(#${clipId})`}>
          <rect x={0} y={0} width={cx} height={size} fill="var(--ink)" />
          <ellipse cx={cx} cy={cy} rx={Math.abs((phase - 0.5) * 2 * r)} ry={r} fill="var(--rule-2)" />
        </g>
      )}
      <circle cx={cx} cy={cy} r={r} fill="none" stroke="var(--ink)" strokeWidth="1" />
    </svg>
  );
}
