package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/di"
	"github.com/sesaquecruz/go-chat-api/pkg"

	_ "github.com/lib/pq"
)

func main() {
	logger := pkg.NewLogger("main")

	config.Load()
	cfg := config.GetConfig()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error(err)
	}

	err = db.Ping()
	if err != nil {
		logger.Error(err)
	}

	apiRouter := di.NewApiRouter(&cfg.API, db)

	go func() {
		if err := apiRouter.Run(); err != nil && err != http.ErrServerClosed {
			logger.Error(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := apiRouter.Stop(context.Background()); err != nil {
		logger.Error(err)
	}
}
