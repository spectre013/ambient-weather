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

export const fmtTime = (iso, tzOffsetH   = -6) => {
  const localString = new Date(iso);
  return localString.toLocaleTimeString('en-US', {
    hour: 'numeric',
    minute: '2-digit',
    hour12: true
  });
};

const dirs = ["N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE",
              "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"];
export const dirToCompass = (d) => dirs[Math.round(d / 22.5) % 16];

const dayNames = ["SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"];
const monthNames = ["Jan", "Feb", "Mar", "Apr", "May", "Jun",
                    "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

export const dayName = (epoch, tzOffsetH) => {
  const d = new Date(epoch * 1000 + tzOffsetH * 3600 * 1000);
  return dayNames[d.getUTCDay()];
};
export const dayDate = (epoch, tzOffsetH) => {
  const d = new Date(epoch * 1000 + tzOffsetH * 3600 * 1000);
  return `${monthNames[d.getUTCMonth()]} ${d.getUTCDate()}`;
};

export const clamp = (v, a, b) => Math.max(a, Math.min(b, v));

export const slugifyEvent = (s) =>
  (s || "").toLowerCase().replace(/[^a-z0-9]+/g, "-").replace(/^-+|-+$/g, "");

export const  timeAgo = (dateString) => {
  const then = new Date(dateString);
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