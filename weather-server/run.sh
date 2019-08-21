docker stop weather-server
docker rm weather-server
docker run --name weather-server -d -p 3000:3000 \
-e PORT=3000 \
-e DB_HOST=docker.for.mac.localhost \
-e DB_USER=user \
-e DB_PASSWORD=password \
-e DB_DATABASE=ambient \
-e AMBIENT_WEATHER_API_KEY=apikey \
-e AMBIENT_WEATHER_APPLICATION_KEY=appkey \
-e DARKSKY=apikey \
-e IPGEO=apikey \
-e LAT=38.725798 \
-e LON=38.725798 \
docker.zoms.net/weather-server:latest