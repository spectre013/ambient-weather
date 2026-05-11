import React, { useState } from "react";
import WeatherIcon from "./WeatherIcon";
import {
  showT, showW, showP, wUnit, pUnit, pDecimals,
  dirToCompass, dayName, dayDate, fmtTime,
} from "../utils/format";

export default function DailyList({ days, tzOffsetH, units }) {
  const list = days.slice(0, 14);
  const allMin = Math.min(...list.map((d) => d.tempmin));
  const allMax = Math.max(...list.map((d) => d.tempmax));
  const span = allMax - allMin || 1;
  const [openIdx, setOpenIdx] = useState(null);

  return (
      <div className="daily">
        {list.map((d, i) => {
          const lo = (d.tempmin - allMin) / span;
          const hi = (d.tempmax - allMin) / span;
          const dn = dayName(d.datetimeEpoch, tzOffsetH);
          const label = i === 0 ? "Today" : i === 1 ? "Tomorrow"
              : dn.charAt(0) + dn.slice(1).toLowerCase();
          const isOpen = openIdx === i;
          const onToggle = () => setOpenIdx(isOpen ? null : i);
          const onKey = (e) => {
            if (e.key === "Enter" || e.key === " ") { e.preventDefault(); onToggle(); }
          };
          return (
              <div className={`day-wrap${isOpen ? " is-open" : ""}`} key={d.datetime}>
                <div
                    className="day"
                    role="button"
                    tabIndex={0}
                    aria-expanded={isOpen}
                    onClick={onToggle}
                    onKeyDown={onKey}
                >
                  <div className="d">
                    <span className="day-name">{label}</span>
                    <span className="day-date">{dayDate(d.datetimeEpoch, tzOffsetH)}</span>
                  </div>
                  <WeatherIcon kind={d.icon} size={26} />
                  <span className={`precip ${d.precipprob ? "" : "zero"}`}>
                {d.precipprob ? `${Math.round(d.precipprob)}%` : "—"}
                    {d.precip > 0 ? <span style={{ color: "var(--faint)", marginLeft: 4 }}>
                  {showP(d.precip, units).toFixed(pDecimals(units))}{pUnit(units) === "in" ? "″" : ` ${pUnit(units)}`}
                </span> : null}
              </span>
                  <div className="range">
                    <span className="min">{Math.round(showT(d.tempmin, units))}°</span>
                    <div className="bar" style={{
                      background: `linear-gradient(to right,
                    var(--rule-2) 0%, var(--rule-2) ${lo * 100}%,
                    var(--cool) ${lo * 100}%,
                    color-mix(in oklch, var(--cool) 50%, var(--warm)) ${((lo + hi) / 2) * 100}%,
                    var(--warm) ${hi * 100}%,
                    var(--rule-2) ${hi * 100}%, var(--rule-2) 100%)`
                    }} />
                    <span className="max">{Math.round(showT(d.tempmax, units))}°</span>
                  </div>
                  <span className="conds">{d.conditions}</span>
                  <span className="caret" aria-hidden="true">▾</span>
                </div>
                <div className="day-detail" aria-hidden={!isOpen}>
                  <div className="day-detail-inner">
                    <p className="dd-desc">{d.summary}</p>
                  </div>
                </div>
              </div>
          );
        })}
      </div>
  );
}

function moonLabel(p) {
  if (p === 0 || p === 1) return "New";
  if (p < 0.25) return "Waxing crescent";
  if (p === 0.25) return "First quarter";
  if (p < 0.5) return "Waxing gibbous";
  if (p === 0.5) return "Full";
  if (p < 0.75) return "Waning gibbous";
  if (p === 0.75) return "Last quarter";
  return "Waning crescent";
}
