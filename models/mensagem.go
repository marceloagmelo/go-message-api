package models

import (
	"errors"
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
func (m *Mensagem) Criar(mensagemModel db.Collection) (string, error) {
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

	mensagem = fmt.Sprintf("Mensagem %s enviada para o rabbitmq", strID)
	logger.Info.Println(mensagem)

	return strID, nil
}

//Atualizar uma mensagem no banco de dados
func (m *Mensagem) Atualizar(mensagemModel db.Collection) error {
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

	return nil
}

//TodasMensagens listar todos os mensagens
func TodasMensagens(usuarioModel db.Collection) ([]Mensagem, error) {

	var mensagens []Mensagem

	if err := usuarioModel.Find().All(&mensagens); err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao listar todos os mensagens", err)
		logger.Erro.Println(mensagem)
		return mensagens, err
	}

	return mensagens, nil
}

//Apagar um mensagem no banco de dados
func Apagar(mensagemModel db.Collection, id int) error {

	resultado := mensagemModel.Find("id=?", id)
	if count, err := resultado.Count(); count < 1 {
		mensagem := ""
		if err != nil {
			mensagem = fmt.Sprintf("%s: %s", "Erro ao recuperar mensagem", err)
		}
		if count > 0 {
		} else {
			mensagem = fmt.Sprintf("Mensagem [%v] não encontrada!", id)
			err = errors.New(mensagem)
		}

		if mensagem != "" {
			logger.Erro.Println(mensagem)
			return err
		}
	}
	if err := resultado.Delete(); err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao apagar mensagem", err)
		logger.Erro.Println(mensagem)
		return err
	}

	return nil
}

//ListarStatus listar mensagens por status
func ListarStatus(mensagemModel db.Collection, status int) ([]Mensagem, error) {

	var mensagens []Mensagem

	resultado := mensagemModel.Find("status", status)
	if count, err := resultado.Count(); count < 1 {
		mensagem := ""
		if err != nil {
			mensagem = fmt.Sprintf("%s: %s", "Erro ao listar status de mensagens", err)
		} else {
			mensagem = fmt.Sprintf("Mensagens com status [%v] não encontrados!", status)
			err = errors.New(mensagem)
		}

		if mensagem != "" {
			logger.Erro.Println(mensagem)
			return mensagens, err
		}
	}

	if err := resultado.All(&mensagens); err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao listar status de mensagens", err)
		logger.Erro.Println(mensagem)
	}

	return mensagens, nil
}

//UmaMensagem recuperar um mensagem no banco de dados
func UmaMensagem(mensagemModel db.Collection, id int) (Mensagem, error) {

	var mensagem Mensagem

	resultado := mensagemModel.Find("id=?", id)
	if count, err := resultado.Count(); count < 1 {
		msg := ""
		if err != nil {
			msg = fmt.Sprintf("%s: %s", "Erro ao recuperar mensagem", err)
		} else {
			msg = fmt.Sprintf("Mensagem [%v] não encontrada!", id)
			err = errors.New(msg)
		}

		if msg != "" {
			logger.Erro.Println(msg)
			return mensagem, err
		}
	}
	if err := resultado.One(&mensagem); err != nil {
		msg := ""
		if err != nil {
			msg = fmt.Sprintf("%s: %s", "Erro ao recuperar mensagem", err)
		} else {
			msg = fmt.Sprintf("Mensagem [%v] não encontrado!", id)
		}

		logger.Erro.Println(msg)
		return mensagem, err
	}

	return mensagem, nil
}
