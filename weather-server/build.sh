#!/usr/bin/env bash

docker build . --no-cache -t docker.zoms.net/weather-server:1.0
if [[ $1 == 'push' ]]
then
    docker push docker.zoms.net/weather-server:1.0
fi
