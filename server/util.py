from collections import defaultdict
from datetime import datetime, timedelta
from zoneinfo import ZoneInfo
from models import Weather, Beaufort, Astro, Trend
from astral.sun import sun, elevation
from astral import LocationInfo
from sqlalchemy import text
from db import db


def get_astro():
    city = LocationInfo("Colorado Springs", "Colorado", "America/Denver", 38.735, -104.628)
    today = sun(city.observer, date=datetime.now(tz=ZoneInfo("America/Denver")  ))
    ele = elevation(city.observer, datetime.now(tz=ZoneInfo("America/Denver")))
    tomorrow = sun(city.observer, date=datetime.now(tz=ZoneInfo("America/Denver")) + timedelta(days=1))

    astro = Astro()
    astro.elevation = round(ele, 2)
    astro.sunrise = convert_to_local(today["sunrise"])
    astro.sunset = convert_to_local(today["sunset"])
    astro.sunrise_tomorrow = convert_to_local(tomorrow["sunrise"])
    astro.sunset_tomorrow = convert_to_local(tomorrow["sunset"])
    astro.darkness = tomorrow["sunrise"] - today["sunset"]
    astro.daylight = today["sunset"] - today["sunrise"]


    astro.has_sunset = False
    if ele < 0:
        astro.has_sunset = True
    return astro

def time_format(recorded):
    return recorded.strftime("%H:%M")

def full_date(recorded):
    return recorded.strftime("%Y-%m-%d %H:%M")


def temp_label(units):
    if units == "metric":
        return "&deg;C"
    else:
        return "&deg;F"

def box_format(units):
    """Get box configuration for different measurement types.
    
    Args:
        units (str): Unit system to use (metric/imperial)
        
    Returns:
        dict: Box configuration for different measurement types
    """
    box = {
        "temperature": {
            "icon": "fa-temperature-three-quarters",
            "title": "Temperature",
            "unit": temp_label(units),
            "style": {}
        },
        "forecast": {
            "icon": "fa-cloud-sun-rain",
            "title": "Forecast",
            "unit": temp_label(units),
            "style": "width: 570px"
        },
        "alerts": {
            "icon": "fa-triangle-exclamation",
            "title": "Alerts",
            "unit": "",
            "style": {}
        },
        "wind": {
            "icon": "fa-wind",
            "title": "Wind",
            "unit": wind_label(units),
            "style": {}
        },
        "rain": {
            "icon": "fa-cloud-showers-heavy",
            "title": "Rain",
            "unit": rain_label(units),
            "style": {}
        },
        "lightning": {
            "icon": "fa-bolt-lightning",
            "title": "Lightning",
            "unit": "",
            "style": {}
        },
        "humidity": {
            "icon": "fa-droplet",
            "title": "Humidity",
            "unit": "%",
            "style": {}
        },
        "barometer": {
            "icon": "fa-temperature-high",
            "title": "Barometer",
            "unit": baro_label(units),
            "style": {}
        },
        "sun": {
            "icon": "fa-sun",
            "title": "Sun",
            "unit": "",
            "style": {}
        },
        "uv": {
            "icon": "fa-cloud-sun",
            "title": "UV | Solar",
            "unit": "",
            "style": {}
        },
        "aqi": {
            "icon": "fa-lungs",
            "title": "Air Quality Index",
            "unit": "",
            "style": {}
        },
        "tempin": {
            "icon": "fa-temperature-half",
            "title": "Living",
            "unit": "",
            "style": {}
        },
        "temp1": {
            "icon": "fa-temperature-half",
            "title": "Basement",
            "unit": "",
            "style": {}
        },
        "temp2": {
            "icon": "fa-temperature-half",
            "title": "Master Bedroom",
            "unit": "",
            "style": {}
        },
        "temp3": {
            "icon": "fa-temperature-half",
            "title": "Office",
            "unit": "",
            "style": {}
        },
        "temp4": {
            "icon": "fa-temperature-half",
            "title": "Garage",
            "unit": "",
            "style": {}
        }
    }
    return box

def get_conditions(forecast):
    return forecast['days'][0]['conditions']

def get_icon(forecast):
    return forecast['days'][0]['icon']

def temp_color(temp: float) -> str:
    if temp <= -5:
        return "tempcolorminus10"
    elif temp <= 5:
        return "tempcolorminus5"
    elif temp <= 14:
        return "tempcolorminus"
    elif temp <= 23:
        return "tempcolor0-5"
    elif temp <= 32:
        return "tempcolorzero"
    elif temp <= 41:
        return "tempcolor0-5"
    elif temp < 50:
        return "tempcolor6-10"
    elif temp < 59:
        return "tempcolor11-15"
    elif temp < 68:
        return "tempcolor16-20"
    elif temp < 77:
        return "tempcolor21-25"
    elif temp < 86:
        return "tempcolor26-30"
    elif temp < 95:
        return "tempcolor31-35"
    elif temp < 104:
        return "tempcolor36-40"
    elif temp < 113:
        return "tempcolor41-45"
    elif temp < 212:
        return "tempcolor50"

    return ""

def temp_display(temp: float, units: str) -> str:
    t = temp
    if units == "metric":
        t = ((temp - 32) * 5) / 9
    return f"{round(t):.0f}"

def minmax(stats):
    minmax = defaultdict(lambda: defaultdict(dict))  # Nested defaultdict for automatic dictionary creation

    for stat in stats:
        parts = stat.id.split("_")  # Assuming v.ID is a dictionary key
        minmax[parts[2]][parts[1]][parts[0].lower()] = {
            "recorded": stat.recorded,
            "value": stat.value
        }

    return minmax

def get_day(date_string: str) -> str:
    # Parse the date string
    date = datetime.strptime(date_string, "%Y-%m-%d")
    return date.strftime("%a").upper()

def query_list(query_result):
    return [row.to_dict() if hasattr(row, 'to_dict') else dict(row) for row in query_result]


def get_beaufort(windspeed: float) -> int:
    """Convert wind speed to Beaufort scale number."""
    speed = windspeed / 1.151 if windspeed > 0 else windspeed
    
    if speed < 1:
        return 0
    elif speed < 4:
        return 1
    elif speed < 7:
        return 2
    elif speed < 11:
        return 3
    elif speed < 17:
        return 4
    elif speed < 22:
        return 5
    elif speed < 28:
        return 6
    elif speed < 34:
        return 7
    elif speed < 41:
        return 8
    elif speed < 48:
        return 9
    elif speed < 56:
        return 10
    elif speed < 64:
        return 11
    else:
        return 12

def beaufort_scale(wind_speed: float) -> Beaufort:
    """Get Beaufort scale information for given wind speed."""
    bft = get_beaufort(wind_speed)
    
    if bft >= 12:
        return Beaufort(
            svg='<svg id="weather34 bft12" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 29.96 C 5.55 36.68 11.21 43.31 16.80 50.00 C 18.26 49.99 19.73 49.99 21.19 49.99 C 18.93 47.26 16.67 44.53 14.40 41.79 C 15.94 40.54 17.47 39.27 19.00 38.00 C 22.34 42.00 25.66 46.01 29.01 50.00 C 42.72 49.98 56.43 50.03 70.14 49.98 C 71.17 47.82 72.07 45.50 73.83 43.81 C 77.91 39.62 84.85 39.15 89.85 41.94 C 93.15 43.97 95.29 47.56 96.00 51.33 L 96.00 54.56 C 95.35 58.38 93.17 62.01 89.84 64.06 C 85.44 66.52 79.67 66.42 75.46 63.60 C 72.81 61.81 71.37 58.87 70.15 56.02 C 46.76 55.98 23.38 56.01 0.00 56.00 L 0.00 29.96 Z" /></svg>',
            text="Hurricane",
            css="beaufort6"
        )
    elif bft >= 11:
        return Beaufort(
            svg='<svg id="weather34 bft11" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 29.97 C 5.57 36.68 11.20 43.33 16.81 50.00 C 34.60 49.99 52.38 50.02 70.16 49.99 C 71.98 43.63 78.44 39.00 85.10 40.36 C 90.77 40.90 95.07 45.87 96.00 51.29 L 96.00 54.67 C 95.15 59.33 91.95 63.89 87.21 65.17 C 82.45 66.67 76.62 65.56 73.32 61.64 C 71.87 60.01 71.03 57.98 70.16 56.01 C 46.77 55.99 23.39 56.01 0.00 56.00 L 0.00 29.97 Z" /></svg>',
            text="Violent Storm",
            css="beaufort6"
        )
    elif bft >= 10:
        return Beaufort(
            svg='<svg id="weather34 bft10" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 29.97 C 5.30 36.33 10.66 42.65 15.99 48.98 C 16.01 42.66 15.99 36.34 16.00 30.02 C 21.62 36.67 27.19 43.35 32.81 50.00 C 34.20 50.00 35.60 49.99 36.99 50.00 C 33.74 45.99 30.46 42.01 27.21 38.00 C 28.66 36.67 30.12 35.34 31.58 34.01 C 36.02 39.32 40.38 44.69 44.81 50.00 C 53.27 49.99 61.72 50.02 70.18 49.99 C 71.39 46.85 73.14 43.69 76.15 41.96 C 80.11 39.71 85.11 39.63 89.20 41.59 C 92.87 43.50 95.27 47.34 96.00 51.35 L 96.00 54.56 C 95.18 60.08 90.75 65.14 85.02 65.65 C 78.40 66.97 71.95 62.35 70.18 56.01 C 46.79 55.98 23.39 56.01 0.00 56.00 L 0.00 29.97 Z" /></svg>',
            text="Storm",
            css="beaufort6"
        )
    elif bft >= 9:
        return Beaufort(
            svg='<svg id="weather34 bft9" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 29.97 C 5.29 36.34 10.66 42.65 15.99 48.99 C 16.01 42.66 15.99 36.34 16.00 30.01 C 21.61 36.67 27.19 43.34 32.80 50.00 C 45.26 49.99 57.71 50.02 70.16 49.98 C 71.97 43.66 78.38 39.03 85.02 40.35 C 90.73 40.87 95.12 45.87 96.00 51.36 L 96.00 54.55 C 95.18 60.08 90.75 65.14 85.00 65.66 C 78.37 66.96 71.98 62.34 70.16 56.02 C 46.77 55.98 23.39 56.01 0.00 56.00 L 0.00 29.97 Z" /></svg>',
            text="Strong Gale",
            css="beaufort6"
        )
    elif bft >= 8:
        return Beaufort(
            svg='<svg id="weather34 bft8" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 4.64 30.07 C 10.05 36.70 15.41 43.37 20.82 50.01 C 22.21 50.00 23.60 50.00 25.00 49.99 C 20.66 44.66 16.33 39.33 12.00 34.00 C 13.54 32.67 15.07 31.34 16.60 30.00 C 22.01 36.67 27.40 43.35 32.81 50.01 C 34.21 50.00 35.60 49.99 37.00 49.99 C 32.66 44.66 28.33 39.33 24.00 34.00 C 25.54 32.67 27.07 31.34 28.60 30.00 C 34.01 36.67 39.40 43.35 44.82 50.01 C 46.21 50.00 47.60 50.00 49.00 49.99 C 44.66 44.66 40.33 39.33 36.00 34.00 C 37.54 32.67 39.07 31.34 40.60 30.00 C 46.01 36.67 51.40 43.35 56.81 50.01 C 58.34 50.00 59.86 50.00 61.39 49.99 C 58.60 46.59 55.80 43.20 53.00 39.80 C 54.54 38.53 56.07 37.27 57.61 36.01 C 61.73 40.79 65.44 45.94 69.89 50.42 C 71.21 47.70 72.41 44.73 74.89 42.83 C 79.11 39.58 85.30 39.36 89.89 41.99 C 93.19 43.96 95.20 47.55 96.00 51.23 L 96.00 54.77 C 95.21 58.43 93.21 62.00 89.94 63.98 C 85.52 66.55 79.63 66.43 75.40 63.55 C 72.77 61.77 71.38 58.81 70.11 56.01 C 52.74 55.99 35.38 56.01 18.01 56.00 C 11.92 48.95 6.57 41.23 0.00 34.64 L 0.00 33.40 C 1.68 32.49 3.18 31.30 4.64 30.07 Z" /></svg>',
            text="Gale",
            css="beaufort6"
        )
    elif bft >= 7:
        return Beaufort(
            svg='<svg id="weather34 bft7" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 34.01 C 1.53 32.68 3.03 31.30 4.61 30.03 C 10.06 36.65 15.35 43.40 20.85 49.98 C 22.23 50.02 23.60 50.02 24.98 49.97 C 20.67 44.64 16.31 39.36 12.04 34.00 C 13.53 32.64 15.05 31.31 16.61 30.03 C 22.05 36.65 27.35 43.39 32.84 49.98 C 34.22 50.02 35.60 50.02 36.98 49.98 C 32.69 44.64 28.30 39.37 24.04 34.00 C 25.53 32.64 27.05 31.31 28.61 30.03 C 34.05 36.65 39.36 43.39 44.83 49.98 C 46.35 50.02 47.86 50.02 49.38 49.99 C 46.62 46.57 43.78 43.22 41.03 39.80 C 42.53 38.52 44.05 37.24 45.61 36.03 C 49.51 40.65 53.29 45.38 57.22 49.98 C 61.55 50.03 65.88 50.00 70.21 49.99 C 71.17 47.29 72.62 44.67 74.86 42.84 C 78.91 39.72 84.66 39.43 89.20 41.60 C 92.85 43.49 95.26 47.32 96.00 51.30 L 96.00 54.66 C 95.11 60.04 90.82 65.13 85.16 65.58 C 78.59 67.06 71.90 62.43 70.21 56.01 C 52.82 55.97 35.43 56.04 18.04 55.98 C 11.96 48.71 6.04 41.31 0.00 34.01 L 0.00 34.01 Z" /></svg>',
            text="Near Gale",
            css="beaufort6"
        )
    elif bft >= 6:
        return Beaufort(
            svg='<svg id="weather34 bft6" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 4.55 30.01 C 10.03 36.62 15.37 43.35 20.81 50.00 C 22.20 50.00 23.60 50.00 24.99 49.99 C 20.67 44.65 16.33 39.34 12.01 34.00 C 13.53 32.66 15.07 31.33 16.60 30.00 C 22.02 36.67 27.39 43.38 32.84 50.02 C 34.22 50.01 35.60 49.99 36.98 49.98 C 32.67 44.64 28.31 39.34 24.01 33.99 C 25.54 32.66 27.07 31.33 28.60 30.01 C 34.01 36.67 39.39 43.35 44.81 50.00 C 53.26 49.99 61.71 50.01 70.15 49.99 C 71.04 48.00 71.89 45.95 73.36 44.31 C 76.67 40.43 82.45 39.34 87.19 40.83 C 91.91 42.08 95.07 46.60 96.00 51.22 L 96.00 54.75 C 95.20 58.73 92.83 62.57 89.13 64.44 C 84.81 66.48 79.42 66.27 75.43 63.58 C 72.80 61.79 71.34 58.86 70.15 56.01 C 52.77 55.99 35.39 56.01 18.01 56.00 C 11.92 48.94 6.51 41.22 0.00 34.56 L 0.00 33.45 C 1.83 32.80 3.11 31.23 4.55 30.01 Z" /></svg>',
            text="Strong Breeze",
            css="beaufort4-5"
        )
    elif bft >= 5:
        return Beaufort(
            svg='<svg id="weather34 bft5" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 4.55 30.01 C 10.04 36.62 15.37 43.37 20.82 50.01 C 22.21 50.00 23.60 49.99 25.00 49.99 C 20.67 44.66 16.33 39.33 12.00 34.00 C 13.53 32.67 15.07 31.34 16.60 30.01 C 22.01 36.67 27.39 43.35 32.82 50.01 C 45.26 49.98 57.71 50.02 70.15 49.99 C 71.41 46.91 73.07 43.77 76.03 42.02 C 79.40 40.12 83.56 39.63 87.24 40.85 C 91.95 42.11 95.08 46.63 96.00 51.23 L 96.00 55.03 C 95.11 58.56 93.16 61.97 90.02 63.95 C 85.60 66.53 79.71 66.45 75.44 63.58 C 72.80 61.79 71.34 58.86 70.15 56.01 C 52.77 55.99 35.39 56.00 18.02 56.00 C 11.93 48.90 6.44 41.24 0.00 34.48 L 0.00 33.53 C 1.72 32.64 3.15 31.32 4.55 30.01 Z" /></svg>',
            text="Fresh Breeze",
            css="beaufort4-5"
        )
    elif bft >= 4:
        return Beaufort(
            svg='<svg id="weather34 bft4" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 33.39 C 1.62 32.38 3.17 31.27 4.69 30.10 C 10.05 36.74 15.43 43.37 20.80 50.01 C 22.27 49.99 23.73 49.99 25.20 49.99 C 22.39 46.60 19.61 43.19 16.80 39.80 C 18.34 38.53 19.87 37.27 21.40 36.00 C 25.26 40.67 29.13 45.33 33.00 50.00 C 45.36 49.99 57.72 50.02 70.08 49.98 C 71.35 47.43 72.52 44.67 74.84 42.87 C 79.08 39.57 85.34 39.34 89.94 42.02 C 93.23 44.01 95.21 47.59 96.00 51.27 L 96.00 54.84 C 95.16 58.45 93.23 61.98 89.99 63.95 C 85.38 66.65 79.11 66.44 74.86 63.15 C 72.54 61.35 71.34 58.58 70.08 56.02 C 52.72 55.98 35.37 56.01 18.01 56.00 C 11.92 48.97 6.60 41.23 0.00 34.67 L 0.00 33.39 Z" /></svg>',
            text="Moderate Breeze",
            css="beaufort3-4"
        )
    elif bft >= 3:
        return Beaufort(
            svg='<svg id="weather34 bft3" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 33.44 C 1.67 32.50 3.17 31.28 4.64 30.06 C 10.04 36.70 15.41 43.36 20.80 50.00 C 37.24 49.99 53.68 50.02 70.12 49.98 C 71.39 47.19 72.76 44.24 75.38 42.46 C 79.66 39.55 85.61 39.46 90.05 42.08 C 93.25 44.09 95.22 47.60 96.00 51.23 L 96.00 54.90 C 95.16 58.48 93.20 61.96 90.01 63.95 C 85.59 66.53 79.71 66.44 75.44 63.58 C 72.79 61.80 71.39 58.83 70.12 56.02 C 52.75 55.98 35.38 56.01 18.01 56.00 C 11.92 48.94 6.53 41.24 0.00 34.58 L 0.00 33.44 Z" /></svg>',
            text="Gentle Breeze",
            css="beaufort1-3"
        )
    elif bft >= 2:
        return Beaufort(
            svg='<svg id="weather34 bft2" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 0.00 33.38 C 1.68 32.46 3.19 31.28 4.67 30.09 C 10.04 36.72 15.42 43.36 20.80 50.00 C 37.23 49.99 53.66 50.03 70.09 49.98 C 71.41 47.21 72.76 44.23 75.39 42.45 C 79.66 39.54 85.60 39.45 90.03 42.07 C 93.26 44.08 95.26 47.64 96.00 51.31 L 96.00 54.79 C 95.15 58.68 92.92 62.47 89.30 64.34 C 84.74 66.62 78.83 66.29 74.79 63.09 C 72.52 61.29 71.31 58.57 70.10 56.02 C 46.73 55.97 23.37 56.02 0.00 56.00 L 0.00 49.94 C 4.33 50.04 8.66 50.00 13.00 49.99 C 8.62 44.92 4.88 39.28 0.00 34.68 L 0.00 33.38 Z" /></svg>',
            text="Light Breeze",
            css="beaufort1-3"
        )
    elif bft >= 1:
        return Beaufort(
            svg='<svg id="weather34 bft1" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 73.92 43.89 C 77.12 40.10 82.80 39.45 87.34 40.81 C 91.48 42.01 93.99 45.85 96.00 49.39 L 96.00 56.58 C 94.00 60.14 91.49 63.99 87.34 65.19 C 82.80 66.55 77.13 65.90 73.92 62.11 C 72.32 60.28 71.03 58.19 69.69 56.16 C 46.47 55.76 23.23 56.12 0.00 56.00 L 0.00 50.00 C 23.23 49.88 46.47 50.24 69.69 49.84 C 71.03 47.81 72.31 45.73 73.92 43.89 Z" /></svg>',
            text="Light Air",
            css="beaufort1-3"
        )
    
    # Default return for calm conditions
    return Beaufort(
        svg='<svg id="weather34 bft0" width="12pt" height="12pt" viewBox="0 0 96 96"><path fill="currentcolor" stroke="currentcolor" strokeWidth="0.09375" opacity="1.00" d=" M 73.92 43.89 C 77.12 40.10 82.80 39.45 87.34 40.81 C 91.48 42.01 93.99 45.85 96.00 49.39 L 96.00 56.58 C 94.00 60.14 91.49 63.99 87.34 65.19 C 82.80 66.55 77.13 65.90 73.92 62.11 C 72.32 60.28 71.03 58.19 69.69 56.16 C 46.47 55.76 23.23 56.12 0.00 56.00 L 0.00 50.00 C 23.23 49.88 46.47 50.24 69.69 49.84 C 71.03 47.81 72.31 45.73 73.92 43.89 Z" /></svg>',
        text="Calm",
        css="beaufort1-3"
    )

def wind_label(units: str) -> str:
    """Get wind speed unit label."""
    return "MPH" if units == "imperial" else "M/S"

def wind_display(wind: float, units: str) -> str:
    """Format wind speed for display."""
    if units == "ms":
        return f"{round(wind/2.237)}"
    return f"{round(wind)}"

def deg_to_compass(num: float) -> str:
    """Convert degrees to compass direction."""
    val = int((num/22.5) + 0.5)
    arr = [
        "North", "NNE", "NE", "ENE",
        "East", "ESE", "SE", "SSE",
        "South", "SSW", "SW", "WSW",
        "West", "WNW", "NW", "NNW"
    ]
    return arr[val % 16]

def get_wind():
    start = datetime.now()
    end = datetime.now() + timedelta(hours=1)
    maxgust_query = f"select coalesce(windgustmph ,0) as value, recorded from records where recorded BETWEEN '{start}' AND '{end}' order by windgustmph desc limit 1"
    maxgust = db.session.execute(text(maxgust_query)).first()
    maxspeed_query = f"select coalesce(windspeedmph ,0) as value, recorded from records where recorded BETWEEN '{start}' AND '{end}' order by windspeedmph desc limit 1"
    maxspeed = db.session.execute(text(maxspeed_query)).first()
    avgspeed_query = f"select coalesce(AVG(windspeedmph) ,0) as value from records where recorded BETWEEN '{start}' AND '{end}'"
    avgspeed = db.session.execute(text(avgspeed_query)).first()
    avgdir_query = f"select coalesce(AVG(winddir) ,0) as value from records where recorded BETWEEN '{start}' AND '{end}'"
    avgdir = db.session.execute(text(avgdir_query)).first()

    return {"gust":maxgust, "wind":maxspeed,"dir":avgdir,"avg":avgspeed}

def get_lightning_month():
    start = datetime.now()
    end = start - timedelta(days=30)
    month_sql = f'''SELECT coalesce(SUM(A.value),0) as value
			FROM (SELECT TO_CHAR(recorded,'YYY-MM-DD') as ldate, 
				  MAX(lightningday) as value 
				  FROM records 
				  where recorded between '{start}' and '{end}' 
			GROUP BY ldate) A'''
    return db.session.execute(text(month_sql)).first()

def wind_run(wind: float) -> float:
    """Calculate wind run (wind speed * hours elapsed in day).
    
    Args:
        wind (float): Wind speed
        
    Returns:
        float: Wind run value
    """
    return wind * datetime.now().hour

def rain_label(units: str) -> str:
    """Get rain measurement unit label.
    
    Args:
        units (str): Unit system (metric/imperial)
        
    Returns:
        str: Rain unit label
    """
    return "mm" if units == "metric" else "in"

def baro_label(units: str) -> str:
    """Get barometric pressure unit label.
    
    Args:
        units (str): Unit system (metric/imperial)
        
    Returns:
        str: Pressure unit label
    """
    return "hPa" if units == "metric" else "inHg"

def get_field_value(instance, field_name):
    return getattr(instance, field_name, None)

def baro_display(baro: float, units: str) -> str:
    """Format barometric pressure for display.
    
    Args:
        baro (float): Barometric pressure value
        units (str): Unit system (metric/imperial)
        
    Returns:
        str: Formatted pressure value
    """
    if units == "metric":
        return f"{baro * 33.86:.2f}"
    return f"{baro:.2f}"

def dew_point_class(dewpoint: float) -> str:
    """Get CSS class for dew point temperature range.
    
    Args:
        dewpoint (float): Dew point temperature in Fahrenheit
        
    Returns:
        str: CSS class name for the temperature range
    """
    if dewpoint > 69.8:
        return "tempmodulehome25-30c"
    elif dewpoint >= 68:
        return "tempmodulehome20-25c"
    elif dewpoint >= 59:
        return "tempmodulehome15-20c"
    elif dewpoint >= 50:
        return "tempmodulehome10-15c"
    elif dewpoint > 41:
        return "tempmodulehome5-10c"
    elif dewpoint >= 32:
        return "tempmodulehome0-5c"
    elif dewpoint > 14:
        return "tempmodulehome-10-0c"
    elif dewpoint >= -50:
        return "tempmodulehome-50-10c"
    return "tempmodulehome0-5c"

def humidity_class(humidity: float) -> str:
    """Get CSS class for humidity range.
    
    Args:
        humidity (float): Humidity percentage
        
    Returns:
        str: CSS class name for the humidity range
    """
    if humidity > 90:
        return "temphumcircle80-100"
    elif humidity > 70:
        return "temphumcircle60-80"
    elif humidity > 35:
        return "temphumcircle35-60"
    elif humidity > 25:
        return "temphumcircle25-35"
    elif humidity <= 25:
        return "temphumcircle0-25"
    return ""

def rain_display(rain: float, units: str) -> str:
    """Format rain measurement for display.
    
    Args:
        rain (float): Rain measurement value
        units (str): Unit system (metric/imperial)
        
    Returns:
        str: Formatted rain value with 2 decimal places
    """
    if units == "metric":
        return f"{rain * 25.4:.2f}"
    return f"{rain:.2f}"

def year(date: datetime) -> str:
    """Get year from datetime as string.
    
    Args:
        date (datetime): Date to extract year from
        
    Returns:
        str: Year as string
    """
    return str(date.year)

def month(date: datetime) -> str:
    """Get three-letter month name from datetime.
    
    Args:
        date (datetime): Date to extract month from
        
    Returns:
        str: Three-letter month name (e.g., 'Jan')
    """
    return date.strftime("%b")

def lightning_class(cnt: int) -> str:
    """Get CSS color class based on lightning strike count.
    
    Args:
        cnt (int): Lightning strike count
        
    Returns:
        str: CSS color class ('green', 'yellow', 'orange', or 'red')
    """
    if cnt == 0 or cnt < 50:
        return "green"
    elif 50 <= cnt < 250:
        return "yellow"
    elif 250 <= cnt < 500:
        return "orange"
    elif cnt >= 500:
        return "red"
    return "green"

def distance_class(d: int) -> str:
    """Get CSS color class based on lightning distance.
    
    Args:
        d (int): Distance in miles/kilometers
        
    Returns:
        str: CSS color class ('green', 'yellow', 'orange', or 'red')
    """
    if d == 0:
        return "green"
    elif 1 <= d < 5:
        return "red"
    elif 5 <= d < 10:
        return "orange"
    elif 10 <= d < 15:
        return "yellow"
    elif d >= 15:
        return "green"
    return "green"


def trend(t: str) -> Trend:
    """Calculate trend for temperature or barometric pressure.
    
    Args:
        t (str): Field name to calculate trend for (e.g. 'tempf' or 'baromrelin')
        
    Returns:
        Trend: Trend object containing direction and amount of change
    """
    sel = f"AVG({t})"
    
    start = datetime.now()
    end = start - timedelta(minutes=30)
    
    # Get average value for last 30 minutes
    avg_query = f"SELECT {sel} FROM records WHERE recorded BETWEEN '{end}' AND '{start}'"
    avg_result = db.session.execute(text(avg_query)).first()
    avg = avg_result[0] if avg_result else 0.0
    
    # Get most recent record
    current_query = "SELECT id, baromrelin, tempf FROM records ORDER BY recorded DESC LIMIT 1"
    current = db.session.execute(text(current_query)).first()
    
    trend = Trend()
    
    if 'temp' in t:
        if current.tempf > avg:
            # Trend up
            trend.trend = "up"
            trend.by = round(current.tempf - avg, 2)
        else:
            # Trend down
            trend.trend = "down"
            trend.by = round(avg - current.tempf, 2)
    else:
        if current.baromrelin > avg:
            # Trend up
            trend.trend = "Steady"
            if (current.baromrelin - avg) > 0.5:
                trend.trend = "Rising"
        else:
            # Trend down
            trend.trend = "Steady"
            if (avg - current.baromrelin) > 0.5:
                trend.trend = "Falling"
    
    return trend

def parse_duration(d: str) -> dict:
    """Parse duration string into hours and minutes.
    
    Args:
        d (str): Duration string (e.g. '5h30m' or '45m')
        
    Returns:
        dict: Dictionary with 'hour' and 'min' keys
    """
    hours = "0"
    minutes = "0"
    
    if "h" in d:
        h = d.split("h")
        m = h[1].split("m")
        hours = h[0]
        minutes = m[0]
    else:
        m = d.split("m")
        hours = "0"
        minutes = m[0]
    
    return {
        "hour": hours,
        "min": minutes
    }

def light_dark(t: timedelta, time: str) -> str:
    """Format timedelta as hours and minutes string.
    
    Args:
        t (timedelta): Time duration
        
    Returns:
        str: Formatted string like '5 hrs 30 min'
    """
    total_seconds = int(t.total_seconds())
    results = {"hours":total_seconds // 3600,"min":(total_seconds % 3600) // 60}
    return results[time]

def is_sun_set(luna: 'Astro') -> str:
    """Get appropriate message based on sun state.
    
    Args:
        luna (Astro): Astronomical data object
        
    Returns:
        str: Message indicating next sun event
    """
    msg = "Time til Sunrise"
    if not luna.has_sunset:
        msg = "Time til Sunset"
    
    return msg

def rise_set_class(luna: 'Astro') -> str:
    """Get CSS class for sunrise/sunset.
    
    Args:
        luna (Astro): Astronomical data object
        
    Returns:
        str: CSS class name
    """
    return "riseclr" if luna.has_sunset else "setclr"

def sun_below(luna: 'Astro') -> str:
    """Get sun position class.
    
    Args:
        luna (Astro): Astronomical data object
        
    Returns:
        str: CSS class for sun position
    """
    return "sunbelow" if luna.has_sunset else "sunabove"

def sun_times(luna: 'Astro') -> dict:
    """Calculate time until next sun event.
    
    Args:
        luna (Astro): Astronomical data object
        
    Returns:
        dict: Dictionary with hours and minutes until next event
    """

    now = datetime.now()
    times = time_difference(now, luna.sunset)
    if luna.has_sunset:
        times = time_difference(now, luna.sunrise)
    
    return times

def today_tomorrow(t: str, luna: 'Astro') -> str:
    """Determine if sun event is today or tomorrow.
    
    Args:
        t (str): Event type ('sunrise' or 'sunset')
        luna (Astro): Astronomical data object
        
    Returns:
        str: 'Today' or 'Tomorrow'
    """
    event = {
        "sunrise": luna.sunrise,
        "sunset": luna.sunset
    }
    return "Tomorrow" if datetime.now() > event[t] else "Today"

def uv_today(data: 'Weather') -> str:
    """Get CSS class for UV index.
    
    Args:
        data (Record): Weather record with UV data
        
    Returns:
        str: CSS class name for UV level
    """
    if data.uv >= 10:
        return "uvtoday11"
    elif data.uv >= 8:
        return "uvtoday9-10"
    elif data.uv >= 5:
        return "uvtoday6-8"
    elif data.uv >= 3:
        return "uvtoday4-5"
    elif data.uv >= 0:
        return "uvtoday1-3"
    return ""

def uv_caution(data: 'Weather',sunset: bool) -> str:
    """Get UV caution level description.
    
    Args:
        data (TemplateData): Template data containing Record and Astro objects
        
    Returns:
        str: UV caution level description
    """
    uv = data.uv
    if uv >= 10:
        return "Extreme"
    elif uv >= 8:
        return "Very High"
    elif uv >= 5:
        return "High"
    elif uv >= 3:
        return "Moderate"
    elif not sunset and uv >= 0:
        return "Low"
    elif sunset and uv <= 0:
        return "Below Horizon"
    return ""

def time_difference(date1: datetime, date2: datetime) -> dict:
    """Calculate time difference between two dates in hours and minutes.
    
    Args:
        date1 (datetime): First date
        date2 (datetime): Second date
        
    Returns:
        dict: Dictionary with 'hour' and 'min' keys containing the time difference
    """
    delta = abs(date2 - date1)
    total_seconds = int(delta.total_seconds())
    hours = total_seconds // 3600
    minutes = (total_seconds % 3600) // 60
    
    return {
        "hour": str(hours),
        "min": str(minutes)
    }

def convert_to_local(date_time):
    """Converts a datetime object to local time and removes timezone information.

    Args:
        date_time: A datetime object, potentially with timezone information.

    Returns:
        A datetime object representing local time without timezone information.
    """
    local_time = date_time.astimezone(None)  # Convert to local timezone
    local_time_naive = local_time.replace(tzinfo=None)  # Remove timezone info
    return local_time_naive

def pm25_to_aqi(pm25: float) -> int:
    pm25 = float(pm25)
    """Convert PM2.5 concentration to AQI value.
    
    Args:
        pm25 (float): PM2.5 concentration in μg/m3
        
    Returns:
        int: AQI value
    """
    # AQI breakpoints for PM2.5
    breakpoints = [
        (0.0, 12.0, 0, 50),      # Good
        (12.1, 35.4, 51, 100),   # Moderate
        (35.5, 55.4, 101, 150),  # Unhealthy for Sensitive Groups
        (55.5, 150.4, 151, 200), # Unhealthy
        (150.5, 250.4, 201, 300), # Very Unhealthy
        (250.5, 500.4, 301, 500), # Hazardous
    ]
    
    # Handle values above scale
    if pm25 > 500.4:
        return 500
    
    # Handle values below scale
    if pm25 < 0:
        return 0
        
    # Find the appropriate breakpoint
    for low_pm25, high_pm25, low_aqi, high_aqi in breakpoints:
        if low_pm25 <= pm25 <= high_pm25:
            # Linear interpolation formula
            aqi = (((high_aqi - low_aqi) / (high_pm25 - low_pm25)) 
                  * (pm25 - low_pm25) + low_aqi)
            return round(aqi)
            
    return 0

def get_aqi_category(aqi: int) -> dict:
    """Get AQI category information.
    
    Args:
        aqi (int): AQI value
        
    Returns:
        dict: Dictionary containing category name and color
    """
    if aqi <= 50:
        return {
            "category": "Good",
            "color": "green",
            "description": "Air quality is satisfactory"
        }
    elif aqi <= 100:
        return {
            "category": "Moderate",
            "color": "yellow",
            "description": "Air quality is acceptable"
        }
    elif aqi <= 150:
        return {
            "category": "Unhealthy for Sensitive Groups",
            "color": "orange",
            "description": "Members of sensitive groups may experience health effects"
        }
    elif aqi <= 200:
        return {
            "category": "Unhealthy",
            "color": "red",
            "description": "Everyone may begin to experience health effects"
        }
    elif aqi <= 300:
        return {
            "category": "Very Unhealthy",
            "color": "purple",
            "description": "Health warnings of emergency conditions"
        }
    else:
        return {
            "category": "Hazardous",
            "color": "maroon",
            "description": "Health alert: everyone may experience serious health effects"
        }

def pill_selected(pill: str, timeframe: str) -> str:
    return "pill pill-selected" if pill == timeframe else "pill"

def chart_queries(t: str, sensor: str) -> str:
    query = chart_sql(t, sensor)
    
    result = db.session.execute(text(query))
    data = [dict(row) for row in result.mappings().all()]
    rowx = []
    rowy = []
    for row in data:
        # print(row['x'], row['y'])
        rowx.append(row['x'].astimezone(ZoneInfo("America/Denver")).strftime('%Y-%m-%d %H:%M'))
        rowy.append(float(row['y']))

    #data = [list(row) for row in result]
    return {"x":rowx,"y":rowy}

def chart_format(t: str, sensor: str) -> str:
    
    charts = {
        "temperature": {"sensors":[
            {"sensor":"ROUND(max(tempf)::numeric,2) as y", "color":"#EE4B2B","title":"Max Temperature"},
            {"sensor":"ROUND(max(dewpoint)::numeric,2) as y", "color":"yellow","title":"Dewpoint"},
            {"sensor":"ROUND(min(tempf)::numeric,2) as y", "color":"blue","title":"Min Temperature"}]},
        "humidity": {"sensors":[
            {"sensor":"ROUND(avg(humidity)::numeric,2) as y", "color":"green","title":"Humidity"}]},
        "windspeed": {"sensors":[
            {"sensor":"ROUND(max(windgustmph)::numeric,2) as y", "color":"red","title":"Wind Gust"},
            {"sensor":"ROUND(max(windspeedmph)::numeric,2) as y", "color":"orange","title":"Wind Speed"}]},
        "barometer": {"sensors":[
            {"sensor":"ROUND(avg(baromrelin)::numeric,2) as y", "color":"purple","title":"Barometric Pressure"}]},
        "lightning": {"sensors":[
            {"sensor":"ROUND(avg(lightninghour)::numeric,2) as y", "color":"yellow","title":"Lightning"}]},
    }
    retval = {"xaxis":[],"yaxis":[],"titles":[],"chartTitle":sensor.capitalize()}
    for sensor in charts[sensor]["sensors"]:
       retval["titles"].append(sensor["title"])
       values = chart_queries(t, sensor["sensor"])
       retval["xaxis"] = values["x"]
       retval['yaxis'].append({
        "data": values["y"],
        "name": sensor["title"],
        "color": sensor['color'],
        "type": "line",
       })
    return retval

def chart_sql(t: str, sensor: str) -> str:
    """Generate SQL query for chart data based on time interval and sensor.
    
    Args:
        t (str): Time interval ('1h', '6h', '12h', '1d', '1m', '1y', 'at')
        sensor (str): Sensor name to query
        
    Returns:
        str: Formatted SQL query
    """
    queries = {
        "1h": """SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst'))) * interval '10 minutes') AS x,
                    %s
                FROM records
                    WHERE recorded >= NOW() - interval '1 hour'
                    AND recorded <= NOW()
                GROUP BY x
                order by x asc""",
                
        "6h": """SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 2) * interval '15 minutes') AS x,
                    %s
                FROM records
                    WHERE recorded >= NOW() - interval '6 hour'
                    AND recorded <= NOW()
                GROUP BY x
                order by x asc""",
                
        "12h": """SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 2) * interval '15 minutes') AS x,
                    %s
                FROM records
                    WHERE recorded >= NOW() - interval '6 hour'
                    AND recorded <= NOW()
                GROUP BY x
                order by x asc""",

        "1d": """SELECT date_trunc('hour', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 20) * interval '15 minute') AS x,
                    %s
                FROM records
                    WHERE recorded >= NOW() - interval '1 day'
                    AND recorded <= NOW()
                GROUP BY x
                order by x asc""",
                
        "1m": """SELECT date_trunc('day', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst'))) * interval '24 hour') AS x,
                    %s
                FROM records
                    WHERE recorded >= NOW() - interval '1 month'
                    AND recorded <= NOW()
                GROUP BY x
                order by x asc""",
                
        "1y": """SELECT date_trunc('month', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 60) * interval '1 day') AS x,
                    %s
                FROM records
                    WHERE recorded >= NOW() - interval '1 year'
                    AND recorded <= NOW()
                GROUP BY x
                order by x asc""",
                
        "at": """SELECT date_trunc('month', (recorded at time zone 'mst')) + (floor(date_part('minute', (recorded at time zone 'mst')) / 60) * interval '1 day') as x,
                    %s
                FROM records
                group by x
                order by x asc"""
    }
    
    return queries[t] % sensor