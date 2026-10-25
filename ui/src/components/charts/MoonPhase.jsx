import React from "react";

// phase: 0 = new, 0.25 = first quarter, 0.5 = full, 0.75 = last quarter.
// The lit area is the semicircle on the lit side (right while waxing, left
// while waning) closed by the terminator, a half-ellipse with x-radius
// r·|cos 2πphase| that bulges toward the lit side for a crescent and away
// from it for a gibbous moon.
export default function MoonPhase({ phase = 0, size = 36 }) {
  const cx = size / 2, cy = size / 2;
  const r = size / 2 - 2;
  const p = ((phase % 1) + 1) % 1;
  const c = Math.cos(2 * Math.PI * p);
  const rx = Math.abs(c) * r;
  const waxing = p < 0.5;
  const crescent = c > 0;
  const limbSweep = waxing ? 1 : 0;
  const termSweep = waxing === crescent ? 0 : 1;
  const d = [
    `M ${cx} ${cy - r}`,
    `A ${r} ${r} 0 0 ${limbSweep} ${cx} ${cy + r}`,
    `A ${rx} ${r} 0 0 ${termSweep} ${cx} ${cy - r}`,
    "Z",
  ].join(" ");
  return (
    <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
      <circle cx={cx} cy={cy} r={r} fill="var(--ink)" />
      <path d={d} fill="var(--rule-2)" />
      <circle cx={cx} cy={cy} r={r} fill="none" stroke="var(--ink)" strokeWidth="1" />
    </svg>
  );
}
