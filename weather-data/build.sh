#!/usr/bin/env bash

docker build . --no-cache -t docker.zoms.net/weather-data:1.1
if [[ $1 == 'push' ]]
then
    docker push docker.zoms.net/weather-data:1.1
fi
