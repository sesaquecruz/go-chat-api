package database

import (
	"context"
	"testing"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test/services"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var postgresMessageRepository, _ = services.NewPostgresContainer(context.Background(), "file://../../../")

type MessagePostgresRepositoryTestSuite struct {
	suite.Suite
	ctx               context.Context
	roomRepository    repository.RoomRepository
	messageRepository repository.MessageRepository
}

func (s *MessagePostgresRepositoryTestSuite) SetupSuite() {
	postgresMessageRepository.Clear()

	db := PostgresConnection(&config.DatabaseConfig{
		Host:     postgresMessageRepository.Host,
		Port:     postgresMessageRepository.Port,
		User:     postgresMessageRepository.User,
		Password: postgresMessageRepository.Password,
		Name:     postgresMessageRepository.Name,
	})

	s.ctx = context.Background()
	s.roomRepository = NewRoomPostgresRepository(db)
	s.messageRepository = NewMessagePostgresRepository(db)
}

func (s *MessagePostgresRepositoryTestSuite) TearDownSuite() {
	if err := postgresMessageRepository.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating postgres container: %s", err)
	}
}

func TestMessagePostgresRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MessagePostgresRepositoryTestSuite))
}

func (s *MessagePostgresRepositoryTestSuite) TestShouldSaveAndFindAMessage() {
	defer postgresMessageRepository.Clear()
	t := s.T()

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533c")
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room := entity.NewRoom(adminId, name, category)

	err := s.roomRepository.Save(s.ctx, room)
	assert.Nil(t, err)

	senderId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	senderName, _ := valueobject.NewUserNameWith("An username")
	text, _ := valueobject.NewMessageTextWith("A text")
	message := entity.NewMessage(room.Id(), senderId, senderName, text)

	err = s.messageRepository.Save(s.ctx, message)
	assert.Nil(t, err)

	result, err := s.messageRepository.FindById(s.ctx, message.Id())
	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, message.Id().Value(), result.Id().Value())
	assert.Equal(t, message.RoomId().Value(), result.RoomId().Value())
	assert.Equal(t, message.SenderId().Value(), result.SenderId().Value())
	assert.Equal(t, message.SenderName().Value(), result.SenderName().Value())
	assert.Equal(t, message.Text().Value(), result.Text().Value())
	assert.Equal(t, message.CreatedAt().Value(), result.CreatedAt().Value())
}
