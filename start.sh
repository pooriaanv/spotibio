#!/bin/bash

if [ "$1" == "rebuild" ]; then
  if type docker-compose > /dev/null; then
    docker-compose down -v && docker-compose up -d --build --force-recreate
  else
    docker compose down -v && docker compose up -d --build --force-recreate
  fi
else
  if type docker-compose > /dev/null; then
    docker-compose up -d
  else
    docker compose up -d
  fi
fi
