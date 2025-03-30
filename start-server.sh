"start-server.sh" 6L, 221B                                                                                                                                4,69          All
docker stop weather-server
docker rm weather-server
docker run --name weather-server -d -p 6000:3000 \
 --env-file ./server/.env docker.zoms.net/production/weather-server:3.0
