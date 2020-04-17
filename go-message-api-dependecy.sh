#!/usr/bin/env bash

source setenv.sh

# Criar rede
echo "Criando a rede $DOCKER_NETWORK..."
docker network ls | grep $DOCKER_NETWORK
if [ "$?" != 0 ]; then
   docker network create $DOCKER_NETWORK
fi

# Mysql
echo "Subindo o mysql..."
docker run -d --name mysqldb --network $DOCKER_NETWORK  \
-p 3306:3306 \
-e MYSQL_USER=${MYSQL_USER} \
-e MYSQL_PASSWORD=${MYSQL_PASSWORD} \
-e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
-e MYSQL_DATABASE=${MYSQL_DATABASE} \
mysql:5.7

# RabbitMQ
echo "Subindo o rabbitmq..."
docker run -d --name rabbitmq --network $DOCKER_NETWORK  \
-p 5672:5672 -p 15672:15672 \
-e RABBITMQ_DEFAULT_USER=${RABBITMQ_DEFAULT_USER} \
-e RABBITMQ_DEFAULT_PASS=${RABBITMQ_DEFAULT_PASS} \
-e RABBITMQ_ERLANG_COOKIE=${RABBITMQ_ERLANG_COOKIE} \
-e RABBITMQ_DEFAULT_VHOST=${RABBITMQ_DEFAULT_VHOST} \
rabbitmq:3.6.16-management

# Listando os containers
docker ps
