#!/usr/bin/env bash

docker build . --no-cache -t docker.zoms.net/production/weather-server:1.4
if [[ $1 == 'push' ]]
then
    docker push docker.zoms.net/production/weather-server:1.4
fi
