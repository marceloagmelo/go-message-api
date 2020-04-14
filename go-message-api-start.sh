#!/usr/bin/env bash

source setenv.sh

# Criar rede
echo "Criando a rede message-net..."
docker network create message-net 

# Mysql
echo "Subindo o mysql..."
docker run -d --name mysqldb --network message-net  \
-p 3306:3306 \
-e MYSQL_USER=${MYSQL_USER} \
-e MYSQL_PASSWORD=${MYSQL_PASSWORD} \
-e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
-e MYSQL_DATABASE=${MYSQL_DATABASE} \
mysql:5.7

# RabbitMQ
echo "Subindo o rabbitmq..."
docker run -d --name rabbitmq --network message-net  \
-p 5672:5672 -p 15672:15672 \
-e RABBITMQ_DEFAULT_USER=${RABBITMQ_DEFAULT_USER} \
-e RABBITMQ_DEFAULT_PASS=${RABBITMQ_DEFAULT_PASS} \
-e RABBITMQ_ERLANG_COOKIE=${RABBITMQ_ERLANG_COOKIE} \
-e RABBITMQ_DEFAULT_VHOST=${RABBITMQ_DEFAULT_VHOST} \
rabbitmq:3.6.16-management

# Rabbitmq message api
echo "Subindo o go-message-api..."
docker run -d --name go-message-api --network message-net  \
-p 8080:8080 \
-e MYSQL_USER=${MYSQL_USER} \
-e MYSQL_PASSWORD=${MYSQL_PASSWORD} \
-e MYSQL_HOSTNAME=${MYSQL_HOSTNAME} \
-e MYSQL_DATABASE=${MYSQL_DATABASE} \
-e MYSQL_PORT=${MYSQL_PORT} \
-e RABBITMQ_USER=${RABBITMQ_USER} \
-e RABBITMQ_PASS=${RABBITMQ_PASS} \
-e RABBITMQ_HOSTNAME=${RABBITMQ_HOSTNAME} \
-e RABBITMQ_PORT=${RABBITMQ_PORT} \
-e RABBITMQ_VHOST=${RABBITMQ_VHOST} \
-e TZ=America/Sao_Paulo \
marceloagmelo/go-message-api

# Listando os containers
docker ps
