import React from "react";
import WeatherIcon from "./WeatherIcon";
import { showT } from "../utils/format";

export default function HourlyStrip({ days, nowDate, tzOffsetH, units }) {
  const allHours = [];
  for (const d of days.slice(0, 2)) {
    for (const h of d.hours) allHours.push(h);
  }
  const currentEpoch = nowDate.getTime() / 1000;
  let startIdx = allHours.findIndex((h) => h.datetimeEpoch >= currentEpoch - 1800);
  if (startIdx < 0) startIdx = 0;
  const slice = allHours.slice(startIdx, startIdx + 24);

  return (
    <div className="hourly">
      {slice.map((h, i) => {
        const d = new Date(h.datetimeEpoch * 1000 + tzOffsetH * 3600 * 1000);
        const hh = d.getUTCHours();
        const ampm = hh >= 12 ? "PM" : "AM";
        const h12 = hh % 12 || 12;
        return (
          <div className={`hour${i === 0 ? " now" : ""}`} key={i}>
            <span className="h">{i === 0 ? "Now" : `${h12}${ampm}`}</span>
            <WeatherIcon kind={h.icon} size={22} />
            <span className="t">{Math.round(showT(h.temp, units))}°</span>
            <span className="p">{h.precipprob > 0 ? `${Math.round(h.precipprob)}%` : ""}</span>
          </div>
        );
      })}
    </div>
  );
}
