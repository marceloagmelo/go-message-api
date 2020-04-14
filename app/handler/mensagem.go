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
	var mensagens []models.Mensagem
	if err := mensagemModel.Find().All(&mensagens); err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao conectar com o banco de dados", err)
		logger.Erro.Println(mensagem)
		respondError(w, http.StatusInternalServerError, mensagem)
		return
	}

	conn, err := lib.ConectarRabbitMQ()
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao conectar com o rabbitmq", err)
		logger.Erro.Println(mensagem)
		respondError(w, http.StatusInternalServerError, mensagem)
		return
	}
	defer conn.Close()

	retorno := retorno{}
	retorno.Status = fmt.Sprintf("OK [%v] !", dataHoraFormatada)

	//logger.Info.Println("Serviço OK!")

	respondJSON(w, http.StatusOK, retorno)
}

//TodasMensagens listagem de todas as mensagens
func TodasMensagens(db db.Database, w http.ResponseWriter, r *http.Request) {
	var mensagens []models.Mensagem
	var mensagemModel = db.Collection("mensagem")

	if err := mensagemModel.Find().All(&mensagens); err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao listar todas as mensagens", err)
		logger.Erro.Println(mensagem)
		respondError(w, http.StatusInternalServerError, mensagem)
		return
	}

	respondJSON(w, http.StatusOK, mensagens)
}

//EnviarMensagem enviar mensagem
func EnviarMensagem(db db.Database, w http.ResponseWriter, r *http.Request) {
	var novaMensagem models.Mensagem

	if r.Method == "POST" {
		reqBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			mensagem := fmt.Sprintf("%s: %s", "Erro ao enviar a mensagem", err)
			logger.Erro.Println(mensagem)
		}

		json.Unmarshal(reqBody, &novaMensagem)

		if novaMensagem.Titulo != "" && novaMensagem.Texto != "" {
			var mensagemModel = db.Collection("mensagem")
			var interf models.Metodos

			interf = novaMensagem

			//err := interf.Criar(mensagemModel)
			strID, err := interf.Criar(mensagemModel)
			if err != nil {
				mensagem := fmt.Sprintf("%s: %s", "Erro ao enviar a mensagem", err)
				logger.Erro.Println(mensagem)
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
			mensagem := fmt.Sprint("Titulo ou Descrição obrigatórios!")
			logger.Erro.Println(mensagem)

			respondError(w, http.StatusLengthRequired, mensagem)
			return
		}

		respondJSON(w, http.StatusCreated, novaMensagem)
	}
}

//AtualizarMensagem atualizar mensagem
func AtualizarMensagem(db db.Database, w http.ResponseWriter, r *http.Request) {
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
				logger.Erro.Println(mensagem)
				respondError(w, http.StatusInternalServerError, mensagem)
				return
			}
		} else {
			mensagem := fmt.Sprint("AtualizarMensagem(): Campos obrigatórios!")

			if novaMensagem.ID <= 0 {
				mensagem = fmt.Sprint("AtualizarMensagem(): ID da mensagem menor ou igual a zero!")
			} else if !utils.InBetween(novaMensagem.Status, 1, 2) {
				mensagem = fmt.Sprint("AtualizarMensagem(): Status diferente de 1 e 2!")
			}
			logger.Erro.Println(mensagem)

			respondError(w, http.StatusLengthRequired, mensagem)
			return
		}

		respondJSON(w, http.StatusOK, novaMensagem)
	}
}

//ApagarMensagem apagar uma mensagem
func ApagarMensagem(db db.Database, w http.ResponseWriter, r *http.Request) {
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

			resultado := mensagemModel.Find("id=?", id)
			if count, err := resultado.Count(); count < 1 {
				mensagem := ""
				if err != nil {
					mensagem = fmt.Sprintf("%s: %s", "Erro ao encontrar mensagem", err)
				} else {
					mensagem = fmt.Sprintf("Mensagem [%v] não encontrada!", id)
				}

				logger.Erro.Println(mensagem)

				respondError(w, http.StatusNotFound, mensagem)
				return
			}

			if err := resultado.Delete(); err != nil {
				mensagem := fmt.Sprintf("%s: %s", "Erro ao apagar mensagem", err)

				logger.Erro.Println(mensagem)

				respondError(w, http.StatusNotFound, mensagem)
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

//StatusMensagens lista de mensagens por status
func StatusMensagens(db db.Database, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro status inválido", err)
		logger.Erro.Println(mensagem)

		respondError(w, http.StatusBadRequest, mensagem)
		return
	}

	var mensagens []models.Mensagem
	var mensagemModel = db.Collection("mensagem")

	resultado := mensagemModel.Find("status", id)
	if count, err := resultado.Count(); count < 1 {
		mensagem := ""
		if err != nil {
			mensagem = fmt.Sprintf("%s: %s", "Erro ao encontrar mensagens", err)
		} else {
			mensagem = fmt.Sprintf("Mensagens com status [%v] não encontradas!", id)
		}

		logger.Erro.Println(mensagem)

		respondError(w, http.StatusNotFound, mensagem)
		return
	}

	if err := resultado.All(&mensagens); err != nil {
		mensagem := fmt.Sprintf("%s: %s", "Erro ao listar todas as mensagens", err)
		logger.Erro.Println(mensagem)
		respondError(w, http.StatusInternalServerError, mensagem)
		return
	}

	respondJSON(w, http.StatusOK, mensagens)
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

	var mensagem models.Mensagem
	var mensagemModel = db.Collection("mensagem")

	resultado := mensagemModel.Find("id=?", id)
	if count, err := resultado.Count(); count < 1 {
		mensagem := ""
		if err != nil {
			mensagem = fmt.Sprintf("%s: %s", "Erro ao encontrar mensagem", err)
		} else {
			mensagem = fmt.Sprintf("Mensagem [%v] não encontrada!", id)
		}

		logger.Erro.Println(mensagem)

		respondError(w, http.StatusNotFound, mensagem)
		return
	}

	if err := resultado.One(&mensagem); err != nil {
		mensagem := ""
		if err != nil {
			mensagem = fmt.Sprintf("%s: %s", "Erro ao encontrar mensagem", err)
		} else {
			mensagem = fmt.Sprintf("Mensagem [%v] não encontrada!", id)
		}

		logger.Erro.Println(mensagem)

		respondError(w, http.StatusNotFound, mensagem)
		return
	}

	respondJSON(w, http.StatusOK, mensagem)
}
