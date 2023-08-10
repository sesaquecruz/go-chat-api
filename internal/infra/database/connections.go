package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/sesaquecruz/go-chat-api/config"

	_ "github.com/lib/pq"
)

func PostgresDb(cfg *config.DatabaseConfig) *sql.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
