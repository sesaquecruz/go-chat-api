package event

import (
	"fmt"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/pkg/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func RabbitMqConnection(cfg *config.BrokerConfig) *amqp.Connection {
	logger := log.NewLogger("RabbitMqConnection")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	return conn
}
