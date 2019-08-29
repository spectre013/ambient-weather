#!/usr/bin/env bash

docker build . --no-cache -t docker.zoms.net/weather-server:latest
if [[ $1 == 'push' ]]
then
    docker push docker.zoms.net/weather-server:latest
fi
