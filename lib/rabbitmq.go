package lib

import (
	"fmt"
	"os"

	"github.com/marceloagmelo/go-message-api/logger"
	"github.com/streadway/amqp"
)

const (
	fila string = "go-rabbitmq"
)

//ConectarRabbitMQ no rabbitmq
func ConectarRabbitMQ() (*amqp.Connection, error) {
	// Conectar com o rabbitmq
	var connectionString = fmt.Sprintf("amqp://%s:%s@%s:%s%s", os.Getenv("RABBITMQ_USER"), os.Getenv("RABBITMQ_PASS"), os.Getenv("RABBITMQ_HOSTNAME"), os.Getenv("RABBITMQ_PORT"), os.Getenv("RABBITMQ_VHOST"))
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Conectando com o rabbitmq", err)
		logger.Erro.Println(mensagem)

		return nil, err
	}

	return conn, nil
}

//EnviarMensagemRabbitMQ no rabbitmq
func EnviarMensagemRabbitMQ(conn *amqp.Connection, conteudoEnviar []byte) error {

	// Abrir o canal
	ch, err := conn.Channel()
	defer ch.Close()
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Abrindo canal no rabbitmq", err)
		logger.Erro.Println(mensagem)
		return err
	}

	// Declarara fila
	q, err := ch.QueueDeclare(
		fila,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		mensagemErro := fmt.Sprintf("%s: %s", "Declarando fila", err)
		logger.Erro.Println(mensagemErro)
		return err
	}

	//body := bytes.NewBuffer(conteudoEnviar)
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        conteudoEnviar,
		})
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Publicando mensagem", err)
		logger.Erro.Println(mensagem)
		return err
	}

	return nil
}
