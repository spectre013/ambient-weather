docker stop weather-ui
docker rm weather-ui
docker run --name weather-ui -d -p 8080:80 \
docker.zoms.net/production/weather-ui