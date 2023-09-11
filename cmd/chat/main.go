package main

import (
	"fmt"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/di"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

//	@title			Chat API
//	@version		1.0.0
//	@description	A Rest API for Chat App.
//	@termsOfService	https://github.com/sesaquecruz/go-chat-api

//	@contact.name	API Support
//	@contact.url	https://github.com/sesaquecruz/go-chat-api

//	@license.name	MIT
//	@license.url	https://github.com/sesaquecruz/go-chat-api

//	@BasePath	/api/v1

//	@securityDefinitions.apikey	Bearer token
//	@in							header
//	@name						Authorization
//	@description				API authorization token

func main() {
	logger := log.NewLogger("Main")

	cfg := config.Load()

	router := di.NewRouter(&cfg.Database, &cfg.Broker, &cfg.Api)
	addr := fmt.Sprintf(":%s", cfg.Api.Port)

	logger.Infof("server started on %s\n", addr)
	router.Run(addr)
}
