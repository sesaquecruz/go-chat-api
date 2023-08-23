package event

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var brokerMessageGateway, _ = test.NewRabbitmqContainer(context.Background(), "../../../")

type MessageGatewayTestSuite struct {
	suite.Suite
	ctx            context.Context
	messageGateway *MessageGateway
}

func (s *MessageGatewayTestSuite) SetupSuite() {
	conn := BrokerConnection(&config.BrokerConfig{
		Host:     brokerMessageGateway.Host,
		Port:     brokerMessageGateway.Port,
		User:     brokerMessageGateway.User,
		Password: brokerMessageGateway.Password,
	})

	s.ctx = context.Background()
	s.messageGateway = NewMessageGateway(conn)
}

func (s *MessageGatewayTestSuite) TearDownSuite() {
	if err := brokerMessageGateway.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating broker container: %s", err)
	}
}

func TestMessageGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(MessageGatewayTestSuite))
}

func (s *MessageGatewayTestSuite) TestSend_ShouldSendAndReceiveAMessage() {
	t := s.T()

	roomId := valueobject.NewId()
	senderId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533g")
	senderName, _ := valueobject.NewUserNameWith("username")
	text, _ := valueobject.NewMessageTextWith("a message")
	message, _ := entity.NewMessage(
		roomId,
		senderId,
		senderName,
		text,
	)

	err := s.messageGateway.Send(s.ctx, message)
	assert.Nil(t, err)

	msgs, err := s.messageGateway.Receive()
	assert.Nil(t, err)

	var data MessageEvent

	select {
	case msg := <-msgs:
		err := json.Unmarshal(msg.Body, &data)
		if err != nil {
			t.Fatal(err)
		}
		err = msg.Ack(true)
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second * 5):
		t.FailNow()
	}

	assert.Equal(t, message.Id().Value(), data.Id)
	assert.Equal(t, message.RoomId().Value(), data.RoomId)
	assert.Equal(t, message.SenderId().Value(), data.SenderId)
	assert.Equal(t, message.SenderName().Value(), data.SenderName)
	assert.Equal(t, message.Text().Value(), data.Text)
	assert.Equal(t, message.CreatedAt().String(), data.CreatedAt)
}
