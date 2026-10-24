import { useEffect, useRef, useState } from "react";
import { tzOffsetAt } from "../utils/format";

// Visual Crossing's API serves day records as a flat array, with `hours`
// stringified. The bundled fixture is already in the proper shape. This
// helper makes both work.
export const normalizeForecast = (raw) => {
  if (Array.isArray(raw)) {
    const days = raw.map((d) => ({
      ...d,
      hours: typeof d.hours === "string" ? JSON.parse(d.hours) : (d.hours || []),
    }));
    return {
      days,
      // The forecast rows carry no tzoffset; use the station offset now.
      tzoffset: raw[0]?.tzoffset ?? tzOffsetAt(),
      latitude: raw[0]?.latitude ?? 38.8674,
      longitude: raw[0]?.longitude ?? -104.7605,
      resolvedAddress: raw[0]?.resolvedAddress ?? "Colorado Springs",
      stations: raw[0]?.stations ?? [],
    };
  }
  return raw;
};

// Loads `current` (websocket) and `forecast` (HTTP) and exposes them along
// with a connection state. In "fixture" mode it serves the bundled JSON;
// in "live" mode it fetches the configured endpoints.
export function useWeatherData({ source, wsUrl, forecastUrl, historyUrl, fixturePaths }) {
  const [current, setCurrent] = useState(null);
  const [forecast, setForecast] = useState(null);
  const [history, setHistory] = useState(null);
  const [conn, setConn] = useState("idle");
  const wsRef = useRef(null);

  useEffect(() => {
    let cancelled = false;

    if (source === "fixture") {
      Promise.all([
        fetch(fixturePaths.current).then((r) => r.json()),
        fetch(fixturePaths.forecast).then((r) => r.json()),
        fetch(fixturePaths.history).then((r) => r.json()),
      ]).then(([c, f, h]) => {
        if (cancelled) return;
        setCurrent(c);
        setForecast(normalizeForecast(f));
        setHistory(h);
      }).catch((e) => console.error("fixture load failed", e));
      return () => { cancelled = true; };
    }

    // Forecast: retry failures with backoff (the page waits for it), and
    // refresh periodically and when the tab becomes visible again so a page
    // left open past midnight moves on to the new "today".
    let forecastRetry = null;
    let forecastDelay = 5000;
    const fetchForecast = () => {
      clearTimeout(forecastRetry);
      fetch(forecastUrl)
          .then((r) => { if (!r.ok) throw new Error(`HTTP ${r.status}`); return r.json(); })
          .then((f) => {
            if (cancelled) return;
            // An empty forecast table comes back as null or [].
            if (!Array.isArray(f) || f.length === 0) throw new Error("empty forecast");
            forecastDelay = 5000;
            setForecast(normalizeForecast(f));
          })
          .catch((e) => {
            if (cancelled) return;
            console.warn("forecast fetch failed:", e.message);
            forecastRetry = setTimeout(fetchForecast, forecastDelay);
            forecastDelay = Math.min(forecastDelay * 2, 300000);
          });
    };
    fetchForecast();
    const forecastTimer = setInterval(fetchForecast, 30 * 60 * 1000);
    const onVisible = () => { if (document.visibilityState === "visible") fetchForecast(); };
    document.addEventListener("visibilitychange", onVisible);

    // History: hourly observations recorded so far today. Poll on a slow
    // heartbeat so a new completed hour appears within a minute even when
    // the websocket only carries minute-level `current`.
    const fetchHistory = () =>
        fetch(historyUrl)
            .then((r) => { if (!r.ok) throw new Error(`HTTP ${r.status}`); return r.json(); })
            .then((h) => { if (!cancelled) setHistory(h); })
            .catch((e) => console.warn("history fetch failed:", e.message));
    fetchHistory();
    const histTimer = setInterval(fetchHistory, 60000);

    return () => {
      cancelled = true;
      clearInterval(histTimer);
      clearInterval(forecastTimer);
      clearTimeout(forecastRetry);
      document.removeEventListener("visibilitychange", onVisible);
    };
  }, [source, forecastUrl, historyUrl, fixturePaths.current, fixturePaths.forecast, fixturePaths.history]);

  useEffect(() => {
    if (source === "fixture") return;
    let cancelled = false;
    let retry = null;
    const connect = () => {
      if (cancelled) return;
      setConn("connecting");
      let socket;
      try { socket = new WebSocket(wsUrl); } catch { setConn("error"); return; }
      wsRef.current = socket;
      socket.onopen = () => setConn("open");
      socket.onmessage = (evt) => {
        try { setCurrent(JSON.parse(evt.data)); }
        catch (e) { console.warn("bad ws payload", e); }
      };
      socket.onerror = () => setConn("error");
      socket.onclose = () => {
        setConn("closed");
        if (cancelled) return;
        retry = setTimeout(connect, 5000);
      };
    };
    connect();
    return () => {
      cancelled = true;
      if (retry) clearTimeout(retry);
      if (wsRef.current) try { wsRef.current.close(); } catch {}
    };
  }, [source, wsUrl]);

  return { current, forecast, history, conn };
}
  