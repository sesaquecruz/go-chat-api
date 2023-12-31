package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sesaquecruz/go-chat-api/pkg/log"

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
	Host       string
	Port       string
	User       string
	Password   string
	Name       string
	logger     *log.Logger
}

func NewPostgresContainer(ctx context.Context, rootPath string) (*PostgresContainer, error) {
	logger := log.NewLogger("PostgresContainer")

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
		logger.Fatal(err)
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	port, err := container.MappedPort(ctx, "5432/tcp")
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port.Port(), database,
	)

	migrations := fmt.Sprintf("%smigrations", rootPath)

	return &PostgresContainer{
		PostgresContainer: container,
		DSN:               dsn,
		Migrations:        migrations,
		Host:              host,
		Port:              port.Port(),
		User:              "postgres",
		Password:          "postgres",
		Name:              "test_db",
		logger:            logger,
	}, nil
}

func (c *PostgresContainer) Clear() {
	m, err := migrate.New(c.Migrations, c.DSN)
	if err != nil {
		c.logger.Fatal(err)
		return
	}
	defer m.Close()

	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		c.logger.Fatal(err)
		return
	}

	if err := m.Up(); err != nil {
		c.logger.Fatal(err)
		return
	}
}
