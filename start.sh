#!/bin/bash

if [ "$1" == "rebuild" ]; then
  docker-compose down -v && docker compose up -d --build --force-recreate
else
  docker-compose up -d
fi
