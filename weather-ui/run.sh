docker stop weather-ui
docker rm weather-ui
docker run --name weather-ui -d -p 8080:80 \
-e DB_HOST=docker.for.mac.localhost \
-e DB_USER=user \
-e DB_PASSWORD=password \
-e DB_DATABASE=ambient \
-e AMBIENT_WEATHER_API_KEY=apikey \
-e AMBIENT_WEATHER_APPLICATION_KEY=appkey \
docker.zoms.net/weather-ui:latest