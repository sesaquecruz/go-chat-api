package event

import (
	"fmt"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/pkg/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func BrokerConnection(cfg *config.BrokerConfig) *amqp.Connection {
	logger := log.NewLogger("BrokerConnection")

	url := fmt.Sprintf("amqp://%s:%s@%s:%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal(err)
		return nil
	}
	defer ch.Close()

	return conn
}
