package gateway

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageGatewayInterface interface {
	Send(ctx context.Context, message *entity.Message) error
	Receive() (<-chan amqp.Delivery, error)
}
