// Map Visual Crossing icon kinds → bundled PNG URLs (served from /public/icons).
const KIND_TO_FILE = {
  "clear-day": "clear-day.png",
  "clear-night": "clear-night.png",
  "partly-cloudy-day": "partly-cloudy-day.png",
  "partly-cloudy-night": "partly-cloudy-night.png",
  "cloudy": "cloudy.png",
  "rain": "rain.png",
  "showers-day": "showers-day.png",
  "showers-night": "showers-night.png",
  "thunder-rain": "thunder-rain.png",
  "thunder-showers-day": "thunder-showers-day.png",
  "thunder-showers-night": "thunder-showers-night.png",
  "snow": "snow.png",
  "snow-showers-day": "snow-showers-day.png",
  "snow-showers-night": "snow-showers-night.png",
  "sleet": "sleet.png",
  "wind": "wind.png",
  "fog": "fog.png",
  "hail": "hail.png",
  "blizzard": "blizzard.png",
  "hot": "hot.png",
};

export const iconUrlFor = (kind) => {
  const file = KIND_TO_FILE[kind] || "partly-cloudy-day.png";
  return `${import.meta.env.BASE_URL}icons/${file}`;
};
