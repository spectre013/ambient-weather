docker stop weather-proxy
docker rm weather-proxy
docker run --name weather-proxy -d  \
  -p 80:80 \
  -v /Users/brian.paulson/opt/darksky/nginx.conf:/etc/nginx/nginx.conf \
  nginx:1.15
