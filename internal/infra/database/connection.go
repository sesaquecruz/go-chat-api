package database

import (
	"database/sql"
	"fmt"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/pkg/log"

	_ "github.com/lib/pq"
)

func DbConnection(cfg *config.DatabaseConfig) *sql.DB {
	logger := log.NewLogger("PostgresDb")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	return db
}
