package gateway

import amqp "github.com/rabbitmq/amqp091-go"

type MessageReceiverGatewayInterface interface {
	Receive() (<-chan amqp.Delivery, error)
}
