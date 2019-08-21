docker stop weather-data
docker rm weather-data
docker run --name weather-data \
-e DB_HOST=docker.for.mac.localhost \
-e DB_USER=user \
-e DB_PASSWORD=password \
-e DB_DATABASE=ambient \
-e AMBIENT_WEATHER_API_KEY=apikey \
-e AMBIENT_WEATHER_APPLICATION_KEY=appkey \
docker.zoms.net/weather-data:latest