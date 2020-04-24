package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/marceloagmelo/go-message-api/lib"
	"github.com/marceloagmelo/go-message-api/logger"
	"github.com/marceloagmelo/go-message-api/models"
	"github.com/marceloagmelo/go-message-api/utils"
	"github.com/marceloagmelo/go-message-api/variaveis"
	"upper.io/db.v3"
)

type retorno struct {
	Status string `json:"mensagem"`
}

//Health testa conexão com o mysql e rabbitmq
func Health(db db.Database, w http.ResponseWriter, r *http.Request) {
	dataHoraFormatada := variaveis.DataHoraAtual.Format(variaveis.DataFormat)

	var mensagemModel = db.Collection("mensagem")
	_, err := models.TodasMensagens(mensagemModel)
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao conectar com o banco de dados", err)
		logger.Erro.Println(mensagem)
		respondError(w, http.StatusInternalServerError, mensagem)
		return
	}

	conn, err := lib.ConectarRabbitMQ()
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao conectar com o rabbitmq", err)
		respondError(w, http.StatusInternalServerError, mensagem)
		return
	}
	defer conn.Close()

	retorno := retorno{}
	retorno.Status = fmt.Sprintf("OK [%v] !", dataHoraFormatada)

	respondJSON(w, http.StatusOK, retorno)
}

//TodasMensagens listagem de todoos os mensagens
func TodasMensagens(db db.Database, w http.ResponseWriter, r *http.Request) {
	var mensagemModel = db.Collection("mensagem")

	mensagens, err := models.TodasMensagens(mensagemModel)
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao listar todas as mensagens", err)
		logger.Erro.Println(mensagem)
		respondError(w, http.StatusInternalServerError, mensagem)
		return
	}

	respondJSON(w, http.StatusOK, mensagens)
}

//Enviar mensagem
func Enviar(db db.Database, w http.ResponseWriter, r *http.Request) {
	var novaMensagem models.Mensagem

	if r.Method == "POST" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			mensagem := fmt.Sprintf("%s: %s", "Erro ao enviar a mensagem", err)
			logger.Erro.Println(mensagem)
		}

		json.Unmarshal(reqBody, &novaMensagem)
		novaMensagem.Status = 1

		if novaMensagem.Titulo != "" && novaMensagem.Texto != "" {
			var mensagemModel = db.Collection("mensagem")
			var interf models.Metodos

			interf = novaMensagem

			//err := interf.Criar(mensagemModel)
			strID, err := interf.Criar(mensagemModel)
			if err != nil {
				mensagem := fmt.Sprintf("%s: %s", "Erro ao enviar a mensagem", err)
				respondError(w, http.StatusInternalServerError, mensagem)
				return
			}

			id, err := strconv.Atoi(strID)
			if err != nil {
				if err != nil {
					mensagem := fmt.Sprintf("%s: %s", "Erro ao enviar a mensagem", err)
					logger.Erro.Println(mensagem)
					respondError(w, http.StatusInternalServerError, mensagem)
					return
				}
			}
			novaMensagem.ID = id
		} else {
			mensagem := fmt.Sprint("Titulo ou Texto obrigatórios!")
			logger.Erro.Println(mensagem)

			respondError(w, http.StatusLengthRequired, mensagem)
			return
		}

		respondJSON(w, http.StatusCreated, novaMensagem)
	}
}

//Atualizar atualizar mensagem
func Atualizar(db db.Database, w http.ResponseWriter, r *http.Request) {
	var novaMensagem models.Mensagem

	if r.Method == "PUT" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			mensagem := fmt.Sprintf("%s: %s", "Erro ao atualizara mensagem", err)
			logger.Erro.Println(mensagem)
		}

		json.Unmarshal(reqBody, &novaMensagem)

		if novaMensagem.ID > 0 && novaMensagem.Titulo != "" && novaMensagem.Texto != "" && utils.InBetween(novaMensagem.Status, 1, 2) {
			var mensagemModel = db.Collection("mensagem")
			var interf models.Metodos

			interf = novaMensagem

			err := interf.Atualizar(mensagemModel)
			if err != nil {
				mensagem := fmt.Sprintf("%s: %s", "Erro ao atualizar a mensagem", err)
				respondError(w, http.StatusInternalServerError, mensagem)
				return
			}
		} else {
			mensagem := fmt.Sprint("Campos obrigatórios!")

			if novaMensagem.ID <= 0 {
				mensagem = fmt.Sprint("ID da mensagem menor ou igual a zero!")
			} else if !utils.InBetween(novaMensagem.Status, 1, 2) {
				mensagem = fmt.Sprint("Status diferente de 1 e 2!")
			}
			logger.Erro.Println(mensagem)

			respondError(w, http.StatusLengthRequired, mensagem)
			return
		}

		respondJSON(w, http.StatusOK, novaMensagem)
	}
}

//Reenviar atualizar mensagem
func Reenviar(db db.Database, w http.ResponseWriter, r *http.Request) {
	var novaMensagem models.Mensagem

	if r.Method == "PUT" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			mensagem := fmt.Sprintf("%s: %s", "Erro ao atualizara mensagem", err)
			logger.Erro.Println(mensagem)
		}

		json.Unmarshal(reqBody, &novaMensagem)
		novaMensagem.Status = 1

		if novaMensagem.ID > 0 && novaMensagem.Titulo != "" && novaMensagem.Texto != "" && utils.InBetween(novaMensagem.Status, 1, 2) {
			var mensagemModel = db.Collection("mensagem")
			var interf models.Metodos

			interf = novaMensagem

			err := interf.Atualizar(mensagemModel)
			if err != nil {
				mensagem := fmt.Sprintf("%s: %s", "Erro ao atualizar a mensagem", err)
				respondError(w, http.StatusInternalServerError, mensagem)
				return
			}

			conn, err := lib.ConectarRabbitMQ()
			if err != nil {
				return
			}
			defer conn.Close()

			strID := fmt.Sprintf("%v", novaMensagem.ID)
			err = lib.EnviarMensagemRabbitMQ(conn, strID)
			if err != nil {
				return
			}

			mensagem := fmt.Sprintf("Mensagem %s enviada para o rabbitmq", strID)
			logger.Info.Println(mensagem)
		} else {
			mensagem := fmt.Sprint("Campos obrigatórios!")

			if novaMensagem.ID <= 0 {
				mensagem = fmt.Sprint("ID da mensagem menor ou igual a zero!")
			} else if !utils.InBetween(novaMensagem.Status, 1, 2) {
				mensagem = fmt.Sprint("Status diferente de 1 e 2!")
			}
			logger.Erro.Println(mensagem)

			respondError(w, http.StatusLengthRequired, mensagem)
			return
		}

		respondJSON(w, http.StatusOK, novaMensagem)
	}
}

//Apagar apagar uma mensagem
func Apagar(db db.Database, w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			mensagem := fmt.Sprintf("%s: %s", "Erro ID inválido", err)
			logger.Erro.Println(mensagem)

			respondError(w, http.StatusBadRequest, mensagem)
			return
		}

		if id > 0 {
			var mensagemModel = db.Collection("mensagem")

			err := models.Apagar(mensagemModel, id)
			if err != nil {
				mensagem := fmt.Sprintf("%s: %s", "Erro ao apagar o mensagem", err)
				respondError(w, http.StatusInternalServerError, mensagem)
				return
			}
		} else {
			mensagem := fmt.Sprint("ID da mensagem menor ou igual a zero!")
			logger.Erro.Println(mensagem)

			respondError(w, http.StatusLengthRequired, mensagem)
			return

		}
		retorno := retorno{}
		retorno.Status = fmt.Sprintf("Mensagem [%v] apagada com sucesso!", id)

		logger.Info.Println(retorno.Status)

		respondJSON(w, http.StatusOK, retorno)
	}
}

//ListarStatus lista de mensagens por status
func ListarStatus(db db.Database, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status, err := strconv.Atoi(vars["status"])
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro status inválido", err)
		logger.Erro.Println(mensagem)

		respondError(w, http.StatusBadRequest, mensagem)
		return
	}

	if status > 0 {
		if !utils.InBetween(status, 1, 2) {
			mensagem := fmt.Sprint("Status diferente de 1 e 2!")
			respondError(w, http.StatusInternalServerError, mensagem)
			return
		}

		var usuarioModel = db.Collection("mensagem")

		mensagens, err := models.ListarStatus(usuarioModel, status)
		if err != nil {
			mensagem := fmt.Sprintf("%s: %s", "Erro ao listar status de mensagens", err)
			respondError(w, http.StatusInternalServerError, mensagem)
			return
		}
		respondJSON(w, http.StatusOK, mensagens)
	} else {
		mensagem := fmt.Sprint("Status do mensagem menor ou igual a zero!")
		logger.Erro.Println(mensagem)

		respondError(w, http.StatusLengthRequired, mensagem)
		return

	}
}

//UmaMensagem recuperar mensagem
func UmaMensagem(db db.Database, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ID inválido", err)
		logger.Erro.Println(mensagem)

		respondError(w, http.StatusBadRequest, mensagem)
		return
	}

	if id > 0 {
		var mensagemModel = db.Collection("mensagem")

		mensagem, err := models.UmaMensagem(mensagemModel, id)
		if err != nil {
			msg := fmt.Sprintf("%s: %s", "Erro ao recuperar mensagem", err)
			respondError(w, http.StatusInternalServerError, msg)
			return
		}
		msg := fmt.Sprintf("Mensagem [%v] recuperado no banco de dados", id)
		logger.Info.Println(msg)

		respondJSON(w, http.StatusOK, mensagem)
	} else {
		msg := fmt.Sprint("ID do mensagem menor ou igual a zero!")
		logger.Erro.Println(msg)

		respondError(w, http.StatusLengthRequired, msg)
		return
	}
}
