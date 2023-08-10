package main

import (
	"fmt"
	"log"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/di"
)

func main() {
	config.Load()
	cfg := config.GetConfig()
	router := di.NewApiRouter(&cfg.Database, &cfg.API)
	addr := fmt.Sprintf(":%s", cfg.API.Port)

	log.Printf("running on %s\n", addr)
	log.Fatal((router.Run(addr)))
}
