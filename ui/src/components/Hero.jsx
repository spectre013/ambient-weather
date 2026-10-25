import React from "react";
import WeatherIcon from "./WeatherIcon";
import SunArc from "./charts/SunArc";
import {showT, tUnit, fmtTime, statVal, fmtOr} from "../utils/format";

export default function Hero({ current, forecast, units }) {
  const today = forecast.days[0];
  // Prefer the observed station conditions (from /api/current websocket); fall
  // back to the forecast's today values when no observation exists yet.
  const conditions = current.conditions || today.conditions;
  const icon = current.icon || today.icon;
  const hi = today.tempmax;
  const lo = today.tempmin;
  const t = current.temp.temp;
  const feels = current.temp.feelslike;
  const dew = current.temp.dewpoint;
  const aHi = statVal(current.temp.minmax, "max", "day");
  const aLo = statVal(current.temp.minmax, "min", "day");

  const now = new Date(current.date).getTime() / 1000;
  const progress = Math.max(0, Math.min(1,
    (now - today.sunriseEpoch) / (today.sunsetEpoch - today.sunriseEpoch)));

  return (
    <div className="card hero-card">
      <div className="label">
        <span><span className="dot" />Now · Conditions</span>
        <span className="right">Updated continuously</span>
      </div>
      <div className="hero-current">
        <div className="tempblock">
          <div>
            <div className="stack">
              <span className="num">{fmtOr(t, (v) => Math.round(showT(v, units)))}</span>
              <span className="deg">{tUnit(units)}</span>
            </div>
            <div className="cond">{conditions}<br/>feels like {fmtOr(feels, (v) => Math.round(showT(v, units)))}{tUnit(units)}</div>
          </div>
          <div className="row">
            <div className="kv"><span className="k">High today</span><span className="v">{fmtOr(aHi, (v) => Math.round(showT(v, units)))}{tUnit(units)}</span></div>
            <div className="kv"><span className="k">Low today</span><span className="v">{fmtOr(aLo, (v) => Math.round(showT(v, units)))}{tUnit(units)}</span></div>
            <div className="kv"><span className="k">Dew point</span><span className="v">{fmtOr(current.temp.dewpoint, (v) => Math.round(showT(v, units)))}{tUnit(units)}</span></div>
            <div className="kv"><span className="k">Humidity</span><span className="v">{fmtOr(current.humidity.humidity)}%</span></div>
          </div>
        </div>
        <div className="iconbox">
          <WeatherIcon kind={icon} size={120} />
          <div className="hilo">
            <span><b>H</b> {Math.round(showT(hi, units))}°</span>
            <span><b>L</b> {Math.round(showT(lo, units))}°</span>
          </div>
        </div>
      </div>
      <div className="suntrack">
        <div className="ends">
          <span className="t">{fmtTime(current.astro.sunrise)}</span>
          <span className="lbl">Sunrise</span>
        </div>
        <SunArc progress={progress} width={400} height={70} />
        <div className="ends">
          <span className="t">{fmtTime(current.astro.sunset,-6)}</span>
          <span className="lbl">Sunset</span>
        </div>
      </div>
    </div>
  );
}
