#!/usr/bin/env bash

# Mysql
echo "Finalizando o mysql..."
docker rm -f mysqldb

# RabbitMQ
echo "Finalizando o rabbitmq..."
docker rm -f rabbitmq

# Teste de conexao
echo "Finalizando o go-message-api..."
docker rm -f go-message-api

# Remover rede
echo "Removendo a rede message-net..."
docker network rm message-net
