// Unit / time formatting helpers shared by all components.

export const cToF = (c) => (c * 9) / 5 + 32;
export const fToC = (f) => ((f - 32) * 5) / 9;
export const showT = (f, units) => (units === "C" ? fToC(f) : f);
export const tUnit = (units) => (units === "C" ? "°C" : "°F");
// Wind: metric system uses km/h (1 mph = 1.609344 km/h)
export const mphToKph = (m) => m * 1.609344;
export const showW = (mph, units) => (units === "C" ? mphToKph(mph) : mph);
export const wUnit = (units) => (units === "C" ? "km/h" : "mph");
// Precip: metric uses mm (1 inch = 25.4 mm)
export const inToMm = (i) => i * 25.4;
export const showP = (inches, units) => (units === "C" ? inToMm(inches) : inches);
export const pUnit = (units) => (units === "C" ? "mm" : "in");
export const pDecimals = (units) => (units === "C" ? 1 : 2);

// Distance: metric uses km (1 mile = 1.609344 km)
export const miToKm = (mi) => mi * 1.609344;
export const showD = (miles, units) => (units === "C" ? miToKm(miles) : miles);
export const dUnit = (units) => (units === "C" ? "km" : "mi");

// The station's time zone. Times are shown in station time whatever the
// viewer's own zone, and offsets come from it per instant so they follow DST.
export const STATION_TZ = "America/Denver";

const offsetFmt = new Intl.DateTimeFormat("en-US", {
  timeZone: STATION_TZ, hourCycle: "h23",
  year: "numeric", month: "numeric", day: "numeric",
  hour: "numeric", minute: "numeric", second: "numeric",
});

// tzOffsetAt returns the station's UTC offset in hours (e.g. -6 in summer,
// -7 in winter) at the instant ms (epoch milliseconds or a Date).
export const tzOffsetAt = (ms = Date.now()) => {
  const t = typeof ms === "number" ? ms : ms.getTime();
  const p = Object.fromEntries(offsetFmt.formatToParts(t).map((x) => [x.type, x.value]));
  const asUTC = Date.UTC(+p.year, +p.month - 1, +p.day, +p.hour, +p.minute, +p.second);
  return Math.round((asUTC - Math.floor(t / 1000) * 1000) / 60000) / 60;
};

// stationDate returns a Date whose UTC fields read as station wall-clock
// time at epoch seconds `epoch` (use getUTCHours/getUTCDay on it).
export const stationDate = (epoch) => {
  const ms = epoch * 1000;
  return new Date(ms + tzOffsetAt(ms) * 3600 * 1000);
};

// stationHour is the station-local hour (0-23) of a Date.
export const stationHour = (d) => stationDate(d.getTime() / 1000).getUTCHours();

export const fmtTime = (iso) =>
  new Date(iso).toLocaleTimeString("en-US", {
    timeZone: STATION_TZ,
    hour: "numeric",
    minute: "2-digit",
    hour12: true,
  });

// parseLocalTimestamp parses a value from a `timestamp without time zone`
// column. Those hold station wall-clock time but are serialised with a
// trailing "Z", so the "Z" is replaced with the station's offset at that
// time instead of letting the browser read them as UTC.
export const parseLocalTimestamp = (s) => {
  const str = String(s).replace(/Z$/, "");
  const naive = new Date(str + "Z"); // wall-clock fields, read as UTC
  if (isNaN(naive)) return naive;
  // Two passes so a time just after a DST change picks up the right offset.
  let off = tzOffsetAt(naive.getTime());
  off = tzOffsetAt(naive.getTime() - off * 3600 * 1000);
  return new Date(naive.getTime() - off * 3600 * 1000);
};

const dirs = ["N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
              "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"];
export const dirToCompass = (d) => dirs[Math.round(d / 22.5) % 16];

const dayNames = ["SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"];
const monthNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun",
                    "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

export const dayName = (epoch) => dayNames[stationDate(epoch).getUTCDay()];
export const dayDate = (epoch) => {
  const d = stationDate(epoch);
  return `${monthNames[d.getUTCMonth()]} ${d.getUTCDate()}`;
};

// statVal reads minmax[kind][period].value ("max"/"min", "day"/"month"/...).
// The server omits periods that have no stats row yet, so this returns null
// instead of throwing when any level is missing.
export const statVal = (minmax, kind, period) => minmax?.[kind]?.[period]?.value ?? null;

// fmtOr applies f to v, or returns a dash when v is missing.
export const fmtOr = (v, f = (x) => x) =>
  v == null || (typeof v === "number" && Number.isNaN(v)) ? "—" : f(v);

export const clamp = (v, a, b) => Math.max(a, Math.min(b, v));

export const slugifyEvent = (s) =>
  (s || "").toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");

export const  timeAgo = (dateString) => {
  // The backend stores lightningtime in a `timestamp without time zone`
  // column, so the value comes across as station wall-clock time mislabeled
  // with a trailing "Z".
  const then = parseLocalTimestamp(dateString);
  if (isNaN(then)) return "invalid date";
  let diff = Math.floor((Date.now() - then.getTime()) / 1000); // seconds
  if (diff < 0) return "in the future";
  if (diff < 60) return "just now";

  const days    = Math.floor(diff / 86400); diff %= 86400;
  const hours   = Math.floor(diff / 3600);  diff %= 3600;
  const minutes = Math.floor(diff / 60);

  const parts = [];
  if (days)    parts.push(`${days}d${days    === 1 ? "" : "s"}`);
  if (hours)   parts.push(`${hours}h${hours === 1 ? "" : "s"}`);
  if (minutes) parts.push(`${minutes}m${minutes === 1 ? "" : "s"}`);

  return parts.join(" ") + " ago";
}