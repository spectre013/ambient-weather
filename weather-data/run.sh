docker stop weather-data
docker rm weather-data
docker run --name weather-data -d \
-e DB_HOST=$(cat .env | grep DB_HOST= | cut -d '=' -f2) \
-e DB_USER=$(cat .env | grep DB_USER= | cut -d '=' -f2) \
-e DB_PASSWORD=$(cat .env | grep DB_PASSWORD= | cut -d '=' -f2) \
-e DB_DATABASE=$(cat .env | grep DB_DATABASE= | cut -d '=' -f2) \
-e AMBIENT_WEATHER_API_KEY=$(cat .env | grep AMBIENT_WEATHER_API_KEY= | cut -d '=' -f2) \
-e AMBIENT_WEATHER_APPLICATION_KEY=$(cat .env | grep AMBIENT_WEATHER_APPLICATION_KEY= | cut -d '=' -f2) \
docker.zoms.net/weather-data:latest