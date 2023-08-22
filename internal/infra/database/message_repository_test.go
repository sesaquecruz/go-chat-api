package database

import (
	"context"
	"testing"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var dbMessageRepository, _ = test.NewPostgresContainer(context.Background(), "file://../../../")

type MessageRepositoryTestSuite struct {
	suite.Suite
	ctx               context.Context
	roomRepository    repository.RoomRepositoryInterface
	messageRepository repository.MessageRepositoryInterface
}

func (s *MessageRepositoryTestSuite) SetupSuite() {
	dbMessageRepository.Clear()

	db := DbConnection(&config.DatabaseConfig{
		Host:     dbMessageRepository.Host,
		Port:     dbMessageRepository.Port,
		User:     dbMessageRepository.User,
		Password: dbMessageRepository.Password,
		Name:     dbMessageRepository.Name,
	})

	s.ctx = context.Background()
	s.roomRepository = NewRoomRepository(db)
	s.messageRepository = NewMessageRepository(db)
}

func (s *MessageRepositoryTestSuite) TearDownSuite() {
	if err := dbMessageRepository.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating dbMessageRepository container: %s", err)
	}
}

func TestMessageRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(MessageRepositoryTestSuite))
}

func (s *MessageRepositoryTestSuite) TestSave_ShouldSaveAndFindAMessage() {
	defer dbMessageRepository.Clear()
	t := s.T()

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Resident Evil")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room, err := entity.NewRoom(adminId, name, category)
	assert.Nil(t, err)

	err = s.roomRepository.Save(s.ctx, room)
	assert.Nil(t, err)

	senderId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	senderName, _ := valueobject.NewUserNameWith("username")
	text, _ := valueobject.NewMessageTextWith("A simple message")
	message, err := entity.NewMessage(room.Id(), senderId, senderName, text)
	assert.Nil(t, err)

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
