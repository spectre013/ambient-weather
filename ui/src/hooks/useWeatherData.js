import { useEffect, useRef, useState } from "react";

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
      tzoffset: raw[0]?.tzoffset ?? -6,
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

    fetch(forecastUrl)
        .then((r) => { if (!r.ok) throw new Error(`HTTP ${r.status}`); return r.json(); })
        .then((f) => { if (cancelled) return; setForecast(normalizeForecast(f)); })
        .catch((e) => console.warn("forecast fetch failed:", e.message));

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

    return () => { cancelled = true; clearInterval(histTimer); };
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
  