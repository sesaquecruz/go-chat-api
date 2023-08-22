package main

import (
	"fmt"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/di"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

func main() {
	logger := log.NewLogger("main")

	config.Load()
	cfg := config.GetConfig()

	router := di.NewApiRouter(&cfg.Database, &cfg.Broker, &cfg.API)
	addr := fmt.Sprintf(":%s", cfg.API.Port)

	logger.Infof("running on %s\n", addr)
	logger.Fatal((router.Run(addr)))
}
