#!/bin/bash

docker_compose="docker compose"
if ! [ $(command -v "$docker_compose") ]; then
    docker_compose="docker-compose"
fi

if [ "$1" == "rebuild" ]; then
  "$docker_compose" down -v && "$docker_compose" compose up -d --build --force-recreate
else
 "$docker_compose" up -d
fi
