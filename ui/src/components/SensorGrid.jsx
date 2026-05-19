import React from "react";
import { showT, tUnit } from "../utils/format";

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
        console.log(s.key)
        console.log(current[s.key].minmax.max.day.value);
        return (
          <div className="sensor" key={s.name}>
            <div className="name">
              <span>{s.name} · <span style={{ color: "var(--faint)" }}>{s.hint}</span></span>
              <span className={`battery ${data.battout ? "low" : ""}`}>BATT {data.battout ? "LOW" : "OK"}</span>
            </div>
            <div className="v">
              {Math.round(showT(data.temp, units))}
              <span style={{ fontFamily: "var(--font-mono)", fontSize: 14, color: "var(--muted)", marginLeft: 4 }}>{tUnit(units)}</span>
            </div>
            <div className="h">RH {data.humidity}% · day hi {Math.round(showT(data.minmax.max.day.value, units))}°</div>
          </div>
        );
      })}
    </div>
  );
}
