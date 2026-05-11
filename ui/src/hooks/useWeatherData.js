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
// with a connection state. Falls back to bundled fixture if the source is
// "fixture" or the network call fails.
export function useWeatherData({ source, wsUrl, forecastUrl, fixturePaths }) {
  const [current, setCurrent] = useState(null);
  const [forecast, setForecast] = useState(null);
  const [conn, setConn] = useState("idle");
  const wsRef = useRef(null);

  useEffect(() => {
    let cancelled = false;
    const loadOffline = () => {
      Promise.all([
        fetch(fixturePaths.current).then((r) => r.json()),
        fetch(fixturePaths.forecast).then((r) => r.json()),
      ]).then(([c, f]) => {
        if (cancelled) return;
        setCurrent(c);
        setForecast(normalizeForecast(f));
        setConn("offline");
      }).catch((e) => console.error("offline data load failed", e));
    };

    if (source === "fixture") { loadOffline(); return () => { cancelled = true; }; }

    fetch(forecastUrl)
      .then((r) => { if (!r.ok) throw new Error(`HTTP ${r.status}`); return r.json(); })
      .then((f) => { if (cancelled) return; setForecast(normalizeForecast(f)); })
      .catch((e) => {
        console.warn("forecast fetch failed, falling back to fixture:", e.message);
        loadOffline();
      });

    return () => { cancelled = true; };
  }, [source, forecastUrl, fixturePaths.current, fixturePaths.forecast]);

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

  return { current, forecast, conn };
}
