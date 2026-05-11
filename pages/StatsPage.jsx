import React, { useMemo } from "react";
import { useOutletContext } from "react-router-dom";

// "Database" stats — synthetic but plausible. The real ambient-weather
// project persists each reading to Postgres; this page pretends to surface
// the operational stats: row counts per table, ingest rate, partition sizes,
// last ANALYZE, and a small write-volume sparkline.

const TABLES = [
  { name: "readings",       rows: 14_523_811, growth: "+938 / hr",  size: "5.18 GB",  retention: "indefinite", hot: true },
  { name: "readings_5m",    rows:  1_452_381, growth: "+12 / hr",   size: "612 MB",   retention: "indefinite" },
  { name: "readings_hourly", rows:    121_032, growth: "+1 / hr",    size:  "78 MB",   retention: "indefinite" },
  { name: "readings_daily", rows:      5_043, growth: "+1 / day",   size:   "4.2 MB", retention: "indefinite" },
  { name: "alerts",         rows:        268, growth: "—",          size:   "324 KB", retention: "5 yrs" },
  { name: "forecast_cache", rows:      6_842, growth: "+24 / day",  size:   "1.8 MB", retention: "30 d" },
  { name: "sensor_log",     rows:    245_612, growth: "+18 / hr",   size:    "62 MB", retention: "180 d" },
  { name: "calibration",    rows:         42, growth: "—",          size:    "12 KB", retention: "indefinite" },
];

// fake 30-day write-volume series, derived from a sine + noise
function writeSeries(seed = 7) {
  const out = [];
  let s = seed;
  for (let i = 0; i < 30; i++) {
    s = (s * 9301 + 49297) % 233280;
    const r = s / 233280;
    const base = 22000 + Math.sin(i * 0.45) * 2400;
    out.push(Math.round(base + (r - 0.5) * 3200));
  }
  return out;
}

function VolumeSpark({ series }) {
  const min = Math.min(...series), max = Math.max(...series);
  const W = 600, H = 80, PAD = 4;
  const pts = series.map((v, i) => {
    const x = PAD + (i * (W - PAD * 2)) / (series.length - 1);
    const y = PAD + (1 - (v - min) / (max - min || 1)) * (H - PAD * 2);
    return [x, y];
  });
  const path = pts.map(([x, y], i) => `${i === 0 ? "M" : "L"}${x.toFixed(1)},${y.toFixed(1)}`).join(" ");
  const area = `${path} L${pts[pts.length - 1][0]},${H - PAD} L${pts[0][0]},${H - PAD} Z`;
  return (
    <svg viewBox={`0 0 ${W} ${H}`} className="volume-spark" preserveAspectRatio="none">
      <path d={area} fill="var(--accent)" opacity="0.08" />
      <path d={path} fill="none" stroke="var(--accent)" strokeWidth="1.6" />
      {pts.map(([x, y], i) => (
        <circle key={i} cx={x} cy={y} r={i === pts.length - 1 ? 3 : 1.5}
          fill="var(--accent)" />
      ))}
    </svg>
  );
}

function fmt(n) { return n.toLocaleString("en-US"); }

export default function StatsPage() {
  const { forecast } = useOutletContext();
  const series = useMemo(() => writeSeries(), []);

  const totalRows = TABLES.reduce((a, t) => a + t.rows, 0);
  const totalSize = "5.94 GB";
  const ingestRate = "972 / hr";
  const dailyWrites = series[series.length - 1];

  return (
    <article className="page stats-page">
      <header className="page-head">
        <div className="page-meta">§ Vol 4 · Database</div>
        <h1 className="page-title">A small Postgres on a small Pi.</h1>
        <p className="page-dek">
          Operational telemetry from the station's persistence layer. Counts
          are exact at the most recent <code>VACUUM ANALYZE</code>; sizes are
          from <code>pg_total_relation_size</code>. Live-readings table is
          partitioned by month and hot.
        </p>
      </header>

      <div className="stats-headline">
        <div className="sh-cell">
          <span className="sh-lbl">Total rows</span>
          <span className="sh-val">{fmt(totalRows)}</span>
        </div>
        <div className="sh-cell">
          <span className="sh-lbl">On disk</span>
          <span className="sh-val">{totalSize}</span>
        </div>
        <div className="sh-cell">
          <span className="sh-lbl">Ingest rate</span>
          <span className="sh-val">{ingestRate}</span>
        </div>
        <div className="sh-cell">
          <span className="sh-lbl">Today's writes</span>
          <span className="sh-val">{fmt(dailyWrites)}</span>
        </div>
        <div className="sh-cell">
          <span className="sh-lbl">Replicas</span>
          <span className="sh-val">1<small className="sh-sub"> · async</small></span>
        </div>
        <div className="sh-cell">
          <span className="sh-lbl">Last analyze</span>
          <span className="sh-val mono">04:12 UTC</span>
        </div>
      </div>

      <div className="page-rule" />

      <section className="stats-section">
        <div className="stats-section-head">
          <h2 className="page-h2">§ Tables</h2>
          <span className="climate-caption-inline">Row counts, partition sizes, retention.</span>
        </div>
        <div className="stats-table-wrap">
          <table className="stats-table">
            <thead>
              <tr>
                <th>Relation</th>
                <th className="num">Rows</th>
                <th className="num">Growth</th>
                <th className="num">Size</th>
                <th>Retention</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {TABLES.map((t) => (
                <tr key={t.name}>
                  <td className="rel"><code>{t.name}</code></td>
                  <td className="num">{fmt(t.rows)}</td>
                  <td className="num grow">{t.growth}</td>
                  <td className="num">{t.size}</td>
                  <td className="ret">{t.retention}</td>
                  <td>{t.hot && <span className="pill hot">hot</span>}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </section>

      <div className="page-rule" />

      <section className="stats-section">
        <div className="stats-section-head">
          <h2 className="page-h2">§ Write volume · last 30 days</h2>
          <span className="climate-caption-inline">Rows inserted into <code>readings</code>, daily totals.</span>
        </div>
        <div className="stats-volume-card">
          <VolumeSpark series={series} />
          <div className="volume-meta">
            <div><span className="vk">High</span><span className="vv">{fmt(Math.max(...series))}</span></div>
            <div><span className="vk">Low</span><span className="vv">{fmt(Math.min(...series))}</span></div>
            <div><span className="vk">Avg</span><span className="vv">{fmt(Math.round(series.reduce((a,b)=>a+b,0)/series.length))}</span></div>
          </div>
        </div>
      </section>

      <div className="page-rule" />

      <section className="stats-section">
        <div className="stats-section-head">
          <h2 className="page-h2">§ Process</h2>
          <span className="climate-caption-inline">Background jobs and their cadence.</span>
        </div>
        <ul className="proc-list">
          <li>
            <span className="proc-name">station_ingest</span>
            <span className="proc-meta">WebSocket → readings · running</span>
            <span className="proc-status ok">ok</span>
          </li>
          <li>
            <span className="proc-name">rollup_5m</span>
            <span className="proc-meta">cron · every 5 minutes · last ran 00:35 ago</span>
            <span className="proc-status ok">ok</span>
          </li>
          <li>
            <span className="proc-name">rollup_hourly</span>
            <span className="proc-meta">cron · top of hour · last ran 12 m ago</span>
            <span className="proc-status ok">ok</span>
          </li>
          <li>
            <span className="proc-name">rollup_daily</span>
            <span className="proc-meta">cron · 00:05 local · last ran 7 h ago</span>
            <span className="proc-status ok">ok</span>
          </li>
          <li>
            <span className="proc-name">forecast_refresh</span>
            <span className="proc-meta">Visual Crossing · every 60 m · last ran 22 m ago</span>
            <span className="proc-status ok">ok</span>
          </li>
          <li>
            <span className="proc-name">vacuum_full_weekly</span>
            <span className="proc-meta">cron · Sun 03:00 · 4 d until next</span>
            <span className="proc-status warn">scheduled</span>
          </li>
        </ul>
      </section>
    </article>
  );
}
