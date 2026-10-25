import React from "react";
import { showT, tUnit, statVal, fmtOr } from "../utils/format";

const SENSOR_DEFS = [
  { name: "Living Room", key: "tempin", hint: "Main Level" },
  { name: "Basement", key: "temp1", hint: "Downstairs" },
  { name: "Hannah", key: "temp2", hint: "Front Bedroom" },
  { name: "Bedroom", key: "temp3", hint: "Master Bedroom" },
  { name: "Garage", key: "temp4", hint: "Attached Garage" },
];

export default function SensorGrid({ current, units }) {
  return (
    <div className="sensors">
      {SENSOR_DEFS.map((s) => {
        const data = current[s.key];
        if (!data) return null;
        return (
          <div className="sensor" key={s.name}>
            <div className="name">
              <span>{s.name} · <span style={{ color: "var(--faint)" }}>{s.hint}</span></span>
              <span className={`battery ${data.battout ? "low" : ""}`}>BATT {data.battout ? "LOW" : "OK"}</span>
            </div>
            <div className="v">
              {fmtOr(data.temp, (v) => Math.round(showT(v, units)))}
              <span style={{ fontFamily: "var(--font-mono)", fontSize: 14, color: "var(--muted)", marginLeft: 4 }}>{tUnit(units)}</span>
            </div>
            <div className="h">RH {fmtOr(data.humidity)}% · day hi {fmtOr(statVal(data.minmax, "max", "day"), (v) => Math.round(showT(v, units)))}°</div>
          </div>
        );
      })}
    </div>
  );
}
