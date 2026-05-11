import React, { useEffect } from "react";
import { Outlet } from "react-router-dom";
import Masthead from "./Masthead";
import Nav from "./Nav";
import Footer from "./Footer";
import { useLocalStorage } from "../hooks/useLocalStorage";
import { useWeatherData } from "../hooks/useWeatherData";

// Static config — not user-facing prefs, so they don't go in localStorage.
const DATA_SOURCE = "live"; // "fixture" | "live"
const WS_URL = "/api/ws";
const FORECAST_URL = "/api/forecast";

const FIXTURES = {
  current: `${import.meta.env.BASE_URL}data/current.json`,
  forecast: `${import.meta.env.BASE_URL}data/forecast.json`,
};

export default function Layout() {
  // User prefs — persisted to localStorage, namespaced under "ws.*"
  const [theme, setTheme]       = useLocalStorage("ws.theme", "light");        // "light" | "dark"
  const [units, setUnits]       = useLocalStorage("ws.units", "F");            // "F" | "C"
  const [density, setDensity]   = useLocalStorage("ws.density", "comfortable"); // "comfortable" | "compact"
  const [accent, setAccent]     = useLocalStorage("ws.accent", "#476aa8");
  const [showSensors, setShowSensors] = useLocalStorage("ws.showSensors", true);

  const { current, forecast, conn } = useWeatherData({
    source: DATA_SOURCE,
    wsUrl: WS_URL,
    forecastUrl: FORECAST_URL,
    fixturePaths: FIXTURES,
  });

  // Mirror prefs onto <html> so CSS can react without rerenders.
  useEffect(() => {
    document.documentElement.dataset.theme = theme;
    document.documentElement.dataset.density = density;
    document.documentElement.style.setProperty("--accent", accent);
  }, [theme, density, accent]);

  // Wait silently for data — no loading screen.
  if (!current || !forecast) return null;

  // Provide the same shape downstream pages used to consume from tweaks,
  // so HomePage / StatsPage etc don't have to change.
  const prefs = { theme, units, density, accent, showSensors };
  const setPref = {
    theme: setTheme, units: setUnits, density: setDensity,
    accent: setAccent, showSensors: setShowSensors,
  };

  return (
      <div className="app">
        <Masthead current={current} forecast={forecast} conn={conn} />
        <Nav
            units={units}
            onUnitsChange={setUnits}
            theme={theme}
            onThemeChange={setTheme}
        />
        <Outlet context={{ current, forecast, conn, prefs, setPref, units }} />
        <Footer forecast={forecast} />
      </div>
  );
}
