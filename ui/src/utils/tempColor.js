// Temperature colours on an absolute °F scale, the way weather maps draw them:
// white in deep cold, through blues to freezing, greens through the mild
// range, then yellow, orange and red as it gets hot. Absolute rather than
// relative to the week, so 40° is the same blue whichever days it sits beside.
const STOPS = [
  [0, [244, 248, 255]],   // white
  [20, [169, 205, 245]],  // pale blue
  [32, [61, 127, 214]],   // blue, freezing
  [45, [47, 163, 138]],   // teal
  [58, [108, 192, 74]],   // green
  [70, [227, 210, 60]],   // yellow
  [80, [243, 154, 43]],   // orange
  [90, [229, 83, 42]],    // red-orange
  [100, [194, 31, 44]],   // red
  [110, [142, 19, 53]],   // deep red
];

const rgb = ([r, g, b]) => `rgb(${Math.round(r)}, ${Math.round(g)}, ${Math.round(b)})`;

// tempColor returns the colour for a temperature in °F, interpolating between
// stops and clamping beyond either end.
export function tempColor(f) {
  if (f <= STOPS[0][0]) return rgb(STOPS[0][1]);
  for (let i = 1; i < STOPS.length; i++) {
    const [t1, c1] = STOPS[i];
    if (f <= t1) {
      const [t0, c0] = STOPS[i - 1];
      const k = (f - t0) / (t1 - t0);
      return rgb(c0.map((v, j) => v + (c1[j] - v) * k));
    }
  }
  return rgb(STOPS[STOPS.length - 1][1]);
}

// tempRangeGradient paints one day's low-to-high span on a track shared by the
// whole list. lo and hi are the day's temperatures and min/max the list's, all
// in °F; the span is coloured by the absolute scale, with every scale stop that
// falls inside it placed where that temperature sits on the track, and the
// rest of the track left as `track`.
export function tempRangeGradient(lo, hi, min, max, track) {
  const span = max - min || 1;
  const at = (t) => `${(((t - min) / span) * 100).toFixed(2)}%`;
  const parts = [`${track} 0%`, `${track} ${at(lo)}`, `${tempColor(lo)} ${at(lo)}`];
  for (const [t] of STOPS) {
    if (t > lo && t < hi) parts.push(`${tempColor(t)} ${at(t)}`);
  }
  parts.push(`${tempColor(hi)} ${at(hi)}`, `${track} ${at(hi)}`, `${track} 100%`);
  return `linear-gradient(to right, ${parts.join(", ")})`;
}
