package logger

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/marceloagmelo/go-message-api/variaveis"
)

type writer struct {
	io.Writer
	timeFormat string
}

func (w writer) Write(b []byte) (n int, err error) {
	return w.Writer.Write(append([]byte(time.Now().Format(variaveis.DataFormat)), b...))
}

//Info log info
var Info = log.New(&writer{os.Stdout, "2006/01/02 15:04:05 "}, " [info] ", 0)

//Erro log erro
var Erro = log.New(&writer{os.Stdout, "2006/01/02 15:04:05 "}, " [error] ", log.Llongfile)
