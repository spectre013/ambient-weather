import React from "react";
import { iconUrlFor } from "../utils/icons";

export default function WeatherIcon({ kind, size = 28, className = "" }) {
  return (
    <img
      className={`weather-icon ${className}`}
      src={iconUrlFor(kind)}
      alt={kind || "weather"}
      width={size}
      height={size}
      style={{ width: size, height: size, objectFit: "contain", display: "block" }}
    />
  );
}
