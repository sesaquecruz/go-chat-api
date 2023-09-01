package services

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/pkg/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type RedisContainer struct {
	*redis.RedisContainer
	Host   string
	Port   string
	logger *log.Logger
}

func NewRedisContainer(ctx context.Context) (*RedisContainer, error) {
	logger := log.NewLogger("RedisContainer")

	container, err := redis.RunContainer(ctx, testcontainers.WithImage("redis:7.2-alpine"))
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	host, err := container.Host(ctx)
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	port, err := container.MappedPort(ctx, "6379/tcp")
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	return &RedisContainer{
		RedisContainer: container,
		Host:           host,
		Port:           port.Port(),
		logger:         logger,
	}, nil
}
