import React, { useEffect, useMemo, useState } from "react";

// Climate page — multi-year monthly aggregates for the station.
// Source: /data/climate.json (array of { Year, Data: { avgrain, avgtemp,
// maxtemp, mintemp } }). Each array is length-13 with index 0 unused; we
// slice [1..12] for Jan..Dec. A value of 0 in months that are clearly
// off-station (e.g. 2020 before the station was installed) is treated as
// missing.

const MONTHS = ["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"];

const METRICS = [
  { key: "avgtemp", label: "Avg temperature", unit: "°F", color: "var(--warm)", domain: "temp" },
  { key: "maxtemp", label: "Max temperature", unit: "°F", color: "var(--danger)", domain: "temp" },
  { key: "mintemp", label: "Min temperature", unit: "°F", color: "var(--cool)", domain: "temp" },
  { key: "avgrain", label: "Rainfall",        unit: "in",  color: "var(--cool)", domain: "rain" },
];

// A value is "present" if it isn't exactly 0, OR if at least one other metric
// for the same year/month is non-zero. That handles legitimate zero rain.
function isPresent(d, metric, monthIdx) {
  const v = d[metric][monthIdx];
  if (v !== 0) return true;
  return [d.avgtemp[monthIdx], d.maxtemp[monthIdx], d.mintemp[monthIdx]]
    .some((x) => x !== 0);
}

function buildSeries(climate, metric) {
  return climate.map((row) => {
    const months = [];
    for (let m = 1; m <= 12; m++) {
      const present = isPresent(row.Data, metric, m);
      months.push({ month: m - 1, value: present ? row.Data[metric][m] : null });
    }
    return { year: row.Year, months };
  });
}

function flatten(series) {
  const pts = [];
  series.forEach((row) => {
    row.months.forEach((m) => {
      if (m.value !== null) pts.push(m.value);
    });
  });
  return pts;
}

// ---------- Heatmap ----------

function Heatmap({ climate, metric }) {
  const series = useMemo(() => buildSeries(climate, metric.key), [climate, metric]);
  const all = useMemo(() => flatten(series), [series]);
  const min = Math.min(...all);
  const max = Math.max(...all);

  const ramp = (v) => {
    if (v === null) return "var(--rule-2)";
    const t = (v - min) / (max - min || 1);
    if (metric.domain === "rain") {
      // pale → cool blue
      return `oklch(${0.96 - t * 0.4} ${0.02 + t * 0.1} 235)`;
    }
    // cool → warm gradient through neutral
    const hue = 235 - t * 200; // 235 → 35
    const l = 0.92 - t * 0.42;
    const c = 0.04 + Math.abs(t - 0.5) * 0.18;
    return `oklch(${l} ${c} ${hue})`;
  };

  return (
    <div className="climate-heatmap">
      <div className="hm-grid">
        <div className="hm-corner" />
        {MONTHS.map((m) => <div key={m} className="hm-mlabel">{m}</div>)}

        {series.map((row) => (
          <React.Fragment key={row.year}>
            <div className="hm-ylabel">{row.year}</div>
            {row.months.map((cell) => (
              <div
                key={cell.month}
                className={`hm-cell ${cell.value === null ? "empty" : ""}`}
                style={{ background: ramp(cell.value) }}
                title={cell.value !== null
                  ? `${MONTHS[cell.month]} ${row.year} · ${cell.value.toFixed(2)} ${metric.unit}`
                  : `${MONTHS[cell.month]} ${row.year} · no data`}
              >
                <span className="hm-v">
                  {cell.value === null ? "—" : cell.value.toFixed(metric.domain === "rain" ? 2 : 0)}
                </span>
              </div>
            ))}
          </React.Fragment>
        ))}
      </div>

      <div className="hm-legend">
        <span className="lab">{min.toFixed(metric.domain === "rain" ? 2 : 0)} {metric.unit}</span>
        <div className="ramp" style={{
          background: metric.domain === "rain"
            ? "linear-gradient(to right, oklch(0.96 0.02 235), oklch(0.56 0.12 235))"
            : "linear-gradient(to right, oklch(0.92 0.04 235), oklch(0.85 0.04 200), oklch(0.78 0.06 110), oklch(0.70 0.13 60), oklch(0.50 0.20 35))"
        }} />
        <span className="lab">{max.toFixed(metric.domain === "rain" ? 2 : 0)} {metric.unit}</span>
      </div>
    </div>
  );
}

// ---------- Year-on-year line chart ----------

function YearLines({ climate, metric }) {
  const series = useMemo(() => buildSeries(climate, metric.key), [climate, metric]);
  const all = flatten(series);
  const min = Math.min(...all);
  const max = Math.max(...all);
  const pad = (max - min) * 0.08;
  const lo = min - pad;
  const hi = max + pad;

  const W = 1100, H = 280, ML = 44, MR = 16, MT = 14, MB = 26;
  const innerW = W - ML - MR;
  const innerH = H - MT - MB;
  const xOf = (mIdx) => ML + (innerW * mIdx) / 11;
  const yOf = (v) => MT + innerH - ((v - lo) / (hi - lo)) * innerH;

  // line color per year — interpolate cool → warm by year ordinal
  const colorOf = (i) => {
    const t = i / Math.max(series.length - 1, 1);
    return `oklch(${0.62 - t * 0.05} ${0.08 + t * 0.06} ${235 - t * 200})`;
  };

  // ticks
  const tickCount = 5;
  const ticks = Array.from({ length: tickCount }, (_, i) => lo + ((hi - lo) * i) / (tickCount - 1));

  return (
    <div className="climate-lines">
      <svg viewBox={`0 0 ${W} ${H}`} className="climate-svg" preserveAspectRatio="none">
        {/* y-grid */}
        {ticks.map((t, i) => (
          <g key={i}>
            <line x1={ML} x2={W - MR} y1={yOf(t)} y2={yOf(t)}
              stroke="var(--rule-2)" strokeWidth="1" />
            <text x={ML - 6} y={yOf(t) + 3} textAnchor="end"
              fontFamily="JetBrains Mono, monospace" fontSize="10" fill="var(--muted)">
              {t.toFixed(metric.domain === "rain" ? 1 : 0)}
            </text>
          </g>
        ))}
        {/* x-axis labels */}
        {MONTHS.map((m, i) => (
          <text key={m} x={xOf(i)} y={H - 6} textAnchor="middle"
            fontFamily="JetBrains Mono, monospace" fontSize="10" fill="var(--muted)">{m}</text>
        ))}

        {/* lines */}
        {series.map((row, i) => {
          const pts = row.months
            .filter((m) => m.value !== null)
            .map((m) => `${xOf(m.month)},${yOf(m.value)}`).join(" ");
          if (!pts) return null;
          const isLast = i === series.length - 1;
          return (
            <g key={row.year}>
              <polyline
                points={pts}
                fill="none"
                stroke={colorOf(i)}
                strokeWidth={isLast ? 2.2 : 1.4}
                strokeLinejoin="round"
                strokeLinecap="round"
                opacity={isLast ? 1 : 0.85}
              />
              {row.months.filter((m) => m.value !== null).map((m) => (
                <circle key={m.month} cx={xOf(m.month)} cy={yOf(m.value)}
                  r={isLast ? 3 : 2} fill={colorOf(i)} />
              ))}
            </g>
          );
        })}
      </svg>

      <div className="climate-legend">
        {series.map((row, i) => (
          <span key={row.year} className="leg">
            <i style={{ background: colorOf(i) }} />
            <b>{row.year}</b>
          </span>
        ))}
      </div>
    </div>
  );
}

// ---------- Top-line statistics ----------

function HeadlineStats({ climate }) {
  // hottest reading, coldest reading, wettest month, total years on record
  let hottest = { v: -Infinity }, coldest = { v: Infinity }, wettest = { v: -Infinity };
  let temps = [], rains = [];
  climate.forEach((row) => {
    for (let m = 1; m <= 12; m++) {
      const present = isPresent(row.Data, "avgtemp", m);
      if (!present) continue;
      const ma = row.Data.maxtemp[m], mi = row.Data.mintemp[m], rn = row.Data.avgrain[m], av = row.Data.avgtemp[m];
      if (ma > hottest.v) hottest = { v: ma, year: row.Year, month: m - 1 };
      if (mi < coldest.v) coldest = { v: mi, year: row.Year, month: m - 1 };
      if (rn > wettest.v) wettest = { v: rn, year: row.Year, month: m - 1 };
      temps.push(av); rains.push(rn);
    }
  });
  const meanT = temps.reduce((a, b) => a + b, 0) / temps.length;
  const totalRain = rains.reduce((a, b) => a + b, 0);

  const cells = [
    { lbl: "Years on record", v: climate.length, u: "" },
    { lbl: "Mean reading", v: meanT.toFixed(1), u: "°F" },
    { lbl: "Hottest", v: hottest.v.toFixed(1), u: "°F",
      sub: `${MONTHS[hottest.month]} ${hottest.year}` },
    { lbl: "Coldest", v: coldest.v.toFixed(1), u: "°F",
      sub: `${MONTHS[coldest.month]} ${coldest.year}` },
    { lbl: "Wettest month", v: wettest.v.toFixed(2), u: "in",
      sub: `${MONTHS[wettest.month]} ${wettest.year}` },
    { lbl: "Cumulative rain", v: totalRain.toFixed(1), u: "in" },
  ];

  return (
    <div className="climate-headline">
      {cells.map((c) => (
        <div key={c.lbl} className="ch-cell">
          <span className="ch-lbl">{c.lbl}</span>
          <span className="ch-val">{c.v}<small>{c.u}</small></span>
          {c.sub && <span className="ch-sub">{c.sub}</span>}
        </div>
      ))}
    </div>
  );
}

// ---------- Page ----------

export default function ClimatePage() {
  const [climate, setClimate] = useState(null);
  const [metric, setMetric] = useState(METRICS[0]);

  useEffect(() => {
    fetch(`${import.meta.env.BASE_URL}data/climate.json`)
      .then((r) => r.json())
      .then((d) => setClimate(d.sort((a, b) => a.Year - b.Year)))
      .catch((e) => console.error("climate fetch failed", e));
  }, []);

  if (!climate) {
    return (
      <div className="page" style={{ minHeight: 200 }}>
        <div className="loading" style={{ minHeight: 200 }}>Loading climate record…</div>
      </div>
    );
  }

  return (
    <article className="page climate-page">
      <header className="page-head">
        <div className="page-meta">§ Vol 4 · Climate record</div>
        <h1 className="page-title">Seven years on the south slope.</h1>
        <p className="page-dek">
          Monthly aggregates from the station's installation in late July 2020
          through the present. Empty cells were before the station went up or
          have not yet been logged. Where readings exist they were derived
          from sixteen-second telemetry, averaged into months at midnight.
        </p>
      </header>

      <HeadlineStats climate={climate} />

      <div className="page-rule" />

      <section className="climate-section">
        <div className="climate-section-head">
          <h2 className="page-h2">§ Monthly grid</h2>
          <div className="climate-tabs">
            {METRICS.map((m) => (
              <button
                key={m.key}
                className={`climate-tab ${metric.key === m.key ? "active" : ""}`}
                onClick={() => setMetric(m)}
              >
                {m.label}
                <span className="climate-tab-unit">{m.unit}</span>
              </button>
            ))}
          </div>
        </div>
        <p className="climate-caption">
          Each cell is one month at the station. Values are {metric.label.toLowerCase()},
          measured in {metric.unit}.
        </p>
        <Heatmap climate={climate} metric={metric} />
      </section>

      <div className="page-rule" />

      <section className="climate-section">
        <div className="climate-section-head">
          <h2 className="page-h2">§ Year over year</h2>
          <span className="climate-caption-inline">Each line is a calendar year.</span>
        </div>
        <YearLines climate={climate} metric={metric} />
      </section>
    </article>
  );
}
