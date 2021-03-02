#!/usr/bin/env bash

docker build . --no-cache -t docker.zoms.net/production/weather-server:1.5
if [[ $1 == 'push' ]]
then
    docker push docker.zoms.net/production/weather-server:1.5
fi
