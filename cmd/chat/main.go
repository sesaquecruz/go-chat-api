package main

import (
	"fmt"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/di"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

func main() {
	logger := log.NewLogger("Main")

	cfg := config.Load()

	router := di.NewRouter(&cfg.Database, &cfg.Broker, &cfg.Api)
	addr := fmt.Sprintf(":%s", cfg.Api.Port)

	logger.Infof("running on %s\n", addr)
	router.Run(addr)
}
