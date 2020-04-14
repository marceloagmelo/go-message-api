package main

import (
	"github.com/marceloagmelo/go-message-api/logger"

	"github.com/marceloagmelo/go-message-api/app"
	"github.com/marceloagmelo/go-message-api/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	logger.Info.Println("Listen 8080...")
	app.Run(":8080")
}
