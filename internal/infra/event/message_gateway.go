package event

import (
	"context"
	"encoding/json"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/pkg/log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageGateway struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	logger *log.Logger
}

func NewMessageGateway(conn *amqp.Connection) *MessageGateway {
	ch, _ := conn.Channel()

	return &MessageGateway{
		conn:   conn,
		ch:     ch,
		logger: log.NewLogger("MessageGateway"),
	}
}

func (g *MessageGateway) Send(ctx context.Context, message *entity.Message) error {
	event := MessageEvent{
		Id:         message.Id().Value(),
		RoomId:     message.RoomId().Value(),
		SenderId:   message.SenderId().Value(),
		SenderName: message.SenderName().Value(),
		Text:       message.Text().Value(),
		CreatedAt:  message.CreatedAt().String(),
	}

	payload, err := json.Marshal(event)
	if err != nil {
		g.logger.Error(err)
		return err
	}

	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        payload,
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

func (g *MessageGateway) Receive() (<-chan amqp.Delivery, error) {
	ch, err := g.ch.Consume(
		"messages.queue",
		"message-gateway",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		g.logger.Error(err)
		return nil, err
	}

	return ch, nil
}
