import os
import requests
from cache import cache


@cache.memoize(timeout=1800)
def get_forecast():
    """Get weather forecast with automatic 30-minute caching.
    
    Returns:
        dict: Forecast data
    """
    try:
        forecast_url = f'https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Colorado+Springs?unitGroup=us&iconSets=icon2&include=days&key={os.getenv("WEATHER_API")}&contentType=json'
        response = requests.get(forecast_url)
        if response.status_code == 200:
            return response.json()
        raise Exception(f"API returned status code {response.status_code}")
            
    except Exception as e:
        raise Exception(f"Error fetching forecast: {str(e)}")