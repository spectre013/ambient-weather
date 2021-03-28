#!/usr/bin/env bash


docker build . --no-cache -t docker.zoms.net/weather-ui:1.5.1
if [[ $1 == 'push' ]]
then
    docker push docker.zoms.net/weather-ui:1.5.1
fi