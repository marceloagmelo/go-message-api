#!/usr/bin/env bash

source setenv.sh

# Verificar rede
echo "Verificando se existe a rede $DOCKER_NETWORK..."
docker network ls | grep $DOCKER_NETWORK
if [ "$?" != 0 ]; then
   echo "Rede $DOCKER_NETWORK não existe!"
   exit 0
fi

# Rabbitmq message api
echo "Subindo o go-message-api..."
docker run -d --name go-message-api --network $DOCKER_NETWORK \
-p 8181:8080 \
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
