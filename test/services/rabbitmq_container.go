package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sesaquecruz/go-chat-api/pkg/log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RabbitmqContainer struct {
	testcontainers.Container
	Host     string
	Port     string
	User     string
	Password string
	logger   *log.Logger
}

func NewRabbitmqContainer(ctx context.Context, rootPath string) (*RabbitmqContainer, error) {
	logger := log.NewLogger("RabbitmqContainer")

	basePath, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	backPath := strings.Split(rootPath, "../")

	for i := 0; i < len(backPath)-1; i++ {
		basePath = filepath.Dir(basePath)
	}

	rabbitmq := testcontainers.ContainerRequest{
		Image:        "rabbitmq:3.12.3-management",
		ExposedPorts: []string{"5672/tcp"},
		Mounts: testcontainers.ContainerMounts{
			testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{
					HostPath: fmt.Sprintf("%s/rabbitmq.config", basePath),
				},
				Target: "/etc/rabbitmq/rabbitmq.config",
			},
			testcontainers.ContainerMount{
				Source: testcontainers.GenericBindMountSource{
					HostPath: fmt.Sprintf("%s/rabbitmq.json", basePath),
				},
				Target: "/etc/rabbitmq/definitions.json",
			},
		},
		WaitingFor: wait.ForLog("Server startup complete"),
	}

	container, err := testcontainers.GenericContainer(ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: rabbitmq,
			Started:          true,
		},
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

	port, err := container.MappedPort(ctx, "5672/tcp")
	if err != nil {
		logger.Fatal(err)
		return nil, err
	}

	return &RabbitmqContainer{
		Container: container,
		Host:      host,
		Port:      port.Port(),
		User:      "guest",
		Password:  "guest",
		logger:    logger,
	}, nil
}
