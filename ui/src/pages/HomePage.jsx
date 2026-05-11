import React from "react";
import { useOutletContext } from "react-router-dom";
import AlertBar from "../components/AlertBar";
import Hero from "../components/Hero";
import QuickStats from "../components/QuickStats";
import TwentyFourChart from "../components/charts/TwentyFourChart";
import WindCard from "../components/cards/WindCard";
import PressureCard from "../components/cards/PressureCard";
import HumidityCard from "../components/cards/HumidityCard";
import UvCard from "../components/cards/UvCard";
import AqiCard from "../components/cards/AqiCard";
import RainCard from "../components/cards/RainCard";
import LightningCard from "../components/cards/LightningCard";
import MoonCard from "../components/cards/MoonCard";
import SensorGrid from "../components/SensorGrid";
import HourlyStrip from "../components/HourlyStrip";
import DailyList from "../components/DailyList";

export default function HomePage() {
  const { current, forecast, units, tweaks } = useOutletContext();
  const today = forecast.days[0];
  const nowDate = new Date(current.date);
  const tzOffsetH = forecast.tzoffset;
  const todayHours = today.hours;

  return (
    <>
      {current.alert && current.alert[0] && <AlertBar alerts={current.alert} />}

      <div className="hero">
        <Hero current={current} forecast={forecast} units={units} />
        <QuickStats current={current} today={today} units={units} />
      </div>

      <div className="card chart-strip" style={{ marginBottom: "var(--gap)" }}>
        <div className="label">
          <span>Today · 24-hour record</span>
          <span className="right" style={{ display: "flex", gap: 14 }}>
            <span style={{ display: "flex", alignItems: "center", gap: 6 }}>
              <i style={{ display: "inline-block", width: 12, height: 2, background: "var(--accent)", borderRadius: 2 }} />
              <span>Temperature</span>
            </span>
            <span style={{ display: "flex", alignItems: "center", gap: 6 }}>
              <i style={{ display: "inline-block", width: 8, height: 8, background: "var(--cool)", opacity: 0.5, borderRadius: 1 }} />
              <span>Precip prob.</span>
            </span>
          </span>
        </div>
        <TwentyFourChart hours={todayHours} nowHour={nowDate.getUTCHours() + tzOffsetH} units={units} />
      </div>

      <div className="grid">
        <div className="section-rule">
          <span className="h">§ Live readings</span>
          <span className="line" />
          <span className="meta">8 instruments · indoor unit nominal</span>
        </div>

        <div className="col-7"><WindCard wind={current.wind} hours={todayHours} units={units} /></div>
        <div className="col-5"><PressureCard baro={current.barometer} /></div>

        <div className="col-3"><HumidityCard humidity={current.humidity} /></div>
        <div className="col-3"><UvCard uv={current.uv} /></div>
        <div className="col-3"><AqiCard aqi={current.aqi} /></div>
        <div className="col-3"><LightningCard lightning={current.lightning} /></div>

        <div className="col-7"><RainCard rain={current.rain} units={units} /></div>
        <div className="col-5"><MoonCard astro={current.astro} today={today} /></div>
      </div>
      <div className="section-rule" style={{ marginTop: 18 }}>
        <span className="h">§ Remote sensors</span>
        <span className="line" />
        <span className="meta">5 channels · all reporting</span>
      </div>
      <SensorGrid current={current} units={units} />

      <div className="section-rule" style={{ marginTop: 18 }}>
        <span className="h">§ Hourly · next 24 hours</span>
        <span className="line" />
        <span className="meta">temperature & precipitation by hour</span>
      </div>
      <div className="card hourly-card">
        <HourlyStrip days={forecast.days} nowDate={nowDate} tzOffsetH={tzOffsetH} units={units} />
      </div>

      <div className="section-rule" style={{ marginTop: 18 }}>
        <span className="h">§ 14-day outlook</span>
        <span className="line" />
        <span className="meta">Visual Crossing · NOAA NDFD ensemble</span>
      </div>
      <div className="card">
        <DailyList days={forecast.days} tzOffsetH={tzOffsetH} units={units} />
      </div>
    </>
  );
}
