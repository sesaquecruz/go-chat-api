package event

import (
	"context"
	"encoding/json"

	"github.com/sesaquecruz/go-chat-api/internal/domain/event"
	"github.com/sesaquecruz/go-chat-api/pkg/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageEventRabbitMqGateway struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger *log.Logger
}

func NewMessageEventRabbitMqGateway(conn *amqp.Connection) *MessageEventRabbitMqGateway {
	ch, _ := conn.Channel()

	return &MessageEventRabbitMqGateway{
		conn:   conn,
		ch:     ch,
		logger: log.NewLogger("MessageRabbitMqGateway"),
	}
}

func (g *MessageEventRabbitMqGateway) Send(ctx context.Context, messageEvent *event.MessageEvent) error {
	body, err := json.Marshal(messageEvent)
	if err != nil {
		g.logger.Error(err)
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	}

	err = g.ch.PublishWithContext(
		ctx,
		"messages",
		"",
		false,
		false,
		msg,
	)
	if err != nil {
		g.logger.Error(err)
		return err
	}

	return nil
}

func (g *MessageEventRabbitMqGateway) Receive(ctx context.Context, messageEvents chan<- *event.MessageEvent) error {
	ch, err := g.conn.Channel()
	if err != nil {
		g.logger.Error(err)
		return err
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"messages.queue",
		"message-rabbitmq-gateway",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		g.logger.Error(err)
		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		case msg := <-msgs:
			messageEvent := &event.MessageEvent{}

			err = json.Unmarshal(msg.Body, messageEvent)
			if err != nil {
				g.logger.Error(err)
			} else {
				messageEvents <- messageEvent
			}

			msg.Ack(true)
		}
	}
}
