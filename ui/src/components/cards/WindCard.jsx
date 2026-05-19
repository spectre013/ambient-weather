import React from "react";
import Card from "../Card";
import Compass from "../charts/Compass";
import Sparkline from "../charts/Sparkline";
import {dirToCompass, fmtTime, showW, wUnit} from "../../utils/format";

export default function WindCard({ wind, hours, units = "F" }) {
    const sparkVals = hours.map((h) => showW(h.windspeed, units));
    const u = wUnit(units);
    return (
        <Card label="Wind" right={`gust ${showW(wind.gustminmax.max.day.value, units).toFixed(0)} ${u} today`}>
            <div className="compass-wrap">
                <Compass dir={wind.winddir} speed={wind.windspeedmph} />
                <div className="compass-readout">
                    <div className="head">{Math.round(showW(wind.windspeedmph, units))}<small>{u}</small></div>
                    <div className="row"><span>Direction</span><span>{dirToCompass(wind.winddir)} · {Math.round(wind.winddir)}°</span></div>
                    <div className="row"><span>Gust</span><span>{showW(wind.windgustmph, units).toFixed(1)} {u}</span></div>
                    <div className="row"><span>10-min avg</span><span>{showW(wind.windavg, units).toFixed(1)} {u}</span></div>
                    <div className="row"><span>Day max</span><span>[{fmtTime(wind.gustminmax.max.day.date)}] {showW(wind.gustminmax.max.day.value, units).toFixed(1)} {u}</span></div>
                </div>
            </div>
            <div style={{ marginTop: 12, borderTop: "1px solid var(--rule-2)", paddingTop: 10 }}>
                <div style={{ display: "flex", justifyContent: "space-between", marginBottom: 4 }}>
                    <span style={{ fontFamily: "var(--font-mono)", fontSize: 10, color: "var(--muted)", letterSpacing: "0.08em", textTransform: "uppercase" }}>Speed · 24h forecast</span>
                </div>
                <Sparkline values={sparkVals} width={300} height={28} color="var(--cool)" fill />
            </div>
        </Card>
    );
}
