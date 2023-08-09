package test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type PostgresContainer struct {
	*postgres.PostgresContainer
	DSN        string
	Migrations string
}

func NewPostgresContainer(ctx context.Context, migrationsPath string) (*PostgresContainer, error) {
	username := "postgres"
	password := "postgres"
	database := "test_db"

	container, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.WithDatabase("test_db"),
		testcontainers.WithWaitStrategy(
			wait.
				ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		log.Fatal(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port.Port(), database,
	)

	return &PostgresContainer{
		PostgresContainer: container,
		DSN:               dsn,
		Migrations:        migrationsPath,
	}, nil
}

func (pc *PostgresContainer) ClearDB() {
	m, err := migrate.New(pc.Migrations, pc.DSN)
	if err != nil {
		log.Fatal("migrate", pc.DSN, err)
	}
	defer m.Close()

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal("down", err)
	}

	if err := m.Up(); err != nil {
		log.Fatal("up", err)
	}
}
