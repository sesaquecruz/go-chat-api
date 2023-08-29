package event

import (
	"context"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/event"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var rabbitmqMessageEventGateway, _ = services.NewRabbitmqContainer(context.Background(), "../../../")

type MessageEventRabbitMqGatewayTestSuite struct {
	suite.Suite
	ctx                 context.Context
	messageEventGateway gateway.MessageEventGateway
}

func (s *MessageEventRabbitMqGatewayTestSuite) SetupSuite() {
	conn := RabbitMqConnection(&config.BrokerConfig{
		Host:     rabbitmqMessageEventGateway.Host,
		Port:     rabbitmqMessageEventGateway.Port,
		User:     rabbitmqMessageEventGateway.User,
		Password: rabbitmqMessageEventGateway.Password,
	})

	s.ctx = context.Background()
	s.messageEventGateway = NewMessageEventRabbitMqGateway(conn)
}

func (s *MessageEventRabbitMqGatewayTestSuite) TearDownSuite() {
	if err := rabbitmqMessageEventGateway.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating rabbitmq container: %s", err)
	}
}

func TestMessageEventRabbitMqGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(MessageEventRabbitMqGatewayTestSuite))
}

func (s *MessageEventRabbitMqGatewayTestSuite) TestShouldSendAndReceiveAMessage() {
	t := s.T()

	roomId := valueobject.NewId()
	senderId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533g")
	senderName, _ := valueobject.NewUserNameWith("An username")
	text, _ := valueobject.NewMessageTextWith("A text")
	message := entity.NewMessage(roomId, senderId, senderName, text)
	messageEvent := event.NewMessageEvent(message)

	err := s.messageEventGateway.Send(s.ctx, messageEvent)
	assert.Nil(t, err)

	msgs := make(chan *event.MessageEvent)
	defer close(msgs)

	go func() {
		err = s.messageEventGateway.Receive(s.ctx, msgs)
		if err != nil {
			t.Error(err)
		}
	}()

	select {
	case msg := <-msgs:
		assert.Equal(t, message.Id().Value(), msg.Id)
		assert.Equal(t, message.RoomId().Value(), msg.RoomId)
		assert.Equal(t, message.SenderId().Value(), msg.SenderId)
		assert.Equal(t, message.SenderName().Value(), msg.SenderName)
		assert.Equal(t, message.Text().Value(), msg.Text)
		assert.Equal(t, message.CreatedAt().Value(), msg.CreatedAt)
	case <-time.After(10 * time.Second):
		t.Fail()
	}
}
