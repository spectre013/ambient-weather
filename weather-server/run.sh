host=$(cat .env | grep DB_HOST= | cut -d '=' -f2)

docker stop weather-server
docker rm weather-server
docker run --name weather-server -d -p 3000:3000 \
-e PORT=3000 \
-e DB_HOST=$(cat .env | grep DB_HOST= | cut -d '=' -f2) \
-e DB_USER=$(cat .env | grep DB_USER= | cut -d '=' -f2) \
-e DB_PASSWORD=$(cat .env | grep DB_PASSWORD= | cut -d '=' -f2) \
-e DB_DATABASE=$(cat .env | grep DB_DATABASE= | cut -d '=' -f2) \
-e AMBIENT_WEATHER_API_KEY=$(cat .env | grep AMBIENT_WEATHER_API_KEY= | cut -d '=' -f2) \
-e AMBIENT_WEATHER_APPLICATION_KEY=$(cat .env | grep AMBIENT_WEATHER_APPLICATION_KEY= | cut -d '=' -f2) \
-e DARKSKY=$(cat .env | grep DARKSKY= | cut -d '=' -f2) \
-e IPGEO=$(cat .env | grep IPGEO= | cut -d '=' -f2) \
-e LAT=38.725798 \
-e LON=-104.66783 \
-e LOGLEVEL=Debug \
docker.zoms.net/weather-server:latest