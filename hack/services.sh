#!/bin/bash

docker run --rm -d \
  --name=sqlkit-mysql \
  -e MYSQL_USER=travis \
  -e MYSQL_PASSWORD="" \
  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
  -e MYSQL_DATABASE=sqlkit \
  -p 3306:3306 \
  mysql

docker run --rm -d \
  --name=sqlkit-postgres \
  -e POSTGRES_PASSWORD="" \
  -e POSTGRES_USER=travis \
  -e POSTGRES_DB=sqlkit \
  -p 5432:5432 \
  postgres
