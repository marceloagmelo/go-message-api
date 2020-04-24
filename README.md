# API de Mensageria usando Golang, RabbitMQ e MySQL

Este é um serviço de envio de mensagem para o **RabbitMQ** e gravação no **MySQL**. Este serviço possuem alguma funcionalidades.

- [Listar Mensagens](#listar-mensagens)
- [Enviar Mensagem](#enviar-mensagem)
- [Atualizar Mensagem](#atualizar-mensagem)
- [Reenviar Mensagem](#reenviar-mensagem)
- [Apagar Mensagem](#apagar-mensagem)
- [Listar por Status da Mensagem](#listar-por-status-da-mensagem)
- [Recuperar uma Mensagem](#recuperar-uma-mensagem)

----


# Instalação

```
go get -v github.com/marceloagmelo/go-message-api
```
```
cd go-message-api
```

## Build da Aplicação

```
./image-build.sh
```

## Iniciar as Aplicações de Dependências
```
./go-dependecy-start.sh
```

## Preparar o MySQL

```
docker  exec -it mysqldb bash -c "mysql -u root -p"
```
- Criar a tabela
	> use gomessagedb;
	
	> CREATE TABLE mensagem (
id INTEGER UNSIGNED NOT NULL AUTO_INCREMENT,
titulo VARCHAR(100), texto VARCHAR(255),
status INTEGER,
PRIMARY KEY (id)
);

## Iniciar a Aplicação
```
./start.sh
```
```
http://localhost:8181/go-message/api/v1/health
```

## Finalizar a Aplicação
```
./api-stop.sh
```

## Finalizar a Todas as Aplicações
```
./stop-all.sh
```

# Serviços
Lista dos serviços disponíveis:

### Listar Mensagens
[http://localhost:8181/go-message/api/v1/mensagens](http://localhost:8181/go-message/api/v1/mensagens)

### Enviar Mensagem
```
curl -v -d '{"id":0, "titulo":"titulo 02", "texto":"menagem 02", "status":1}' -H "Content-Type: application/json" -X POST http://localhost:8181/go-message/api/v1/mensagem/enviar
```

### Atualizar Mensagem
```
curl -v -d '{"id":1, "titulo":"titulo 02", "texto":"menagem 02", "status":1}' -H "Content-Type: application/json" -X PUT http://localhost:8181/go-message/api/v1/mensagem/atualizar
```

### Reenviar Mensagem
```
curl -v -d '{"id":1, "titulo":"titulo 02", "texto":"menagem 02", "status":1}' -H "Content-Type: application/json" -X PUT http://localhost:8181/go-message/api/v1/mensagem/reenviar
```

### Apagar Mensagem
```
curl -H "Content-Type: application/json" -X DELETE http://localhost:8181/go-message/api/v1/mensagem/apagar/1
```

### Listar por Status da Mensagem
[http://localhost:8181/go-message/api/v1/mensagem/status/2](http://localhost:8181/go-message/api/v1/mensagem/status/2)

### Recuperar uma Mensagem
[http://localhost:8181/go-message/api/v1/mensagem/1](http://localhost:8181/go-message/api/v1/mensagem/1)