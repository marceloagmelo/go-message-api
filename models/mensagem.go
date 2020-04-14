package models

import (
	"fmt"

	"github.com/marceloagmelo/go-message-api/lib"
	"github.com/marceloagmelo/go-message-api/logger"
	"upper.io/db.v3"
)

//Mensagem estrutura de mensagem
type Mensagem struct {
	ID     int    `db:"id" json:"id"`
	Titulo string `db:"titulo" json:"titulo"`
	Texto  string `db:"texto" json:"texto"`
	Status int    `db:"status" json:"status"`
}

// Metodos interface
type Metodos interface {
	Criar(mensagemModel db.Collection) (string, error)
	Atualizar(mensagemModel db.Collection) error
}

//Criar uma mensagem no banco de dados
func (m Mensagem) Criar(mensagemModel db.Collection) (string, error) {
	novoID, err := mensagemModel.Insert(m)
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Gravando a mensagem no banco de dados", err)
		logger.Erro.Println(mensagem)
		return "", err
	}
	strID := fmt.Sprintf("%v", novoID)
	mensagem := fmt.Sprintf("Mensagem %s adicionada no banco de dados", strID)
	logger.Info.Println(mensagem)

	conn, err := lib.ConectarRabbitMQ()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	err = lib.EnviarMensagemRabbitMQ(conn, strID)
	if err != nil {
		return strID, err
	}

	return strID, nil
}

//Atualizar uma mensagem no banco de dados
func (m Mensagem) Atualizar(mensagemModel db.Collection) error {
	var novaMensagem = Mensagem{
		ID:     m.ID,
		Titulo: m.Titulo,
		Texto:  m.Texto,
		Status: m.Status,
	}

	resultado := mensagemModel.Find("id", m.ID)
	if count, err := resultado.Count(); count < 1 {
		mensagem := ""
		if err != nil {
			mensagem = fmt.Sprintf("%s: %s", "Recuperando mensagem no banco de dados", err)
		} else {
			mensagem = fmt.Sprintf("Mensagem [%v] não encontrada!", m.ID)
		}

		logger.Erro.Println(mensagem)

		return err
	}

	if err := resultado.Update(&novaMensagem); err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Gravando a mensagem no banco de dados", err)
		logger.Erro.Println(mensagem)
		return err
	}

	strID := fmt.Sprintf("%v", m.ID)
	mensagem := fmt.Sprintf("Mensagem %s atualizada no banco de dados", strID)
	logger.Info.Println(mensagem)

	conn, err := lib.ConectarRabbitMQ()
	if err != nil {
		return err
	}
	defer conn.Close()

	err = lib.EnviarMensagemRabbitMQ(conn, strID)
	if err != nil {
		return err
	}

	return nil
}
