package impl

import (
	"context"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/event"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type SendMessageUseCase struct {
	roomRepository      repository.RoomRepository
	messageRepository   repository.MessageRepository
	messageEventGateway gateway.MessageEventGateway
	logger              *log.Logger
}

func NewSendMessageUseCase(
	roomRepository repository.RoomRepository,
	messageRepository repository.MessageRepository,
	messageEventGateway gateway.MessageEventGateway,
) *SendMessageUseCase {
	return &SendMessageUseCase{
		roomRepository:      roomRepository,
		messageRepository:   messageRepository,
		messageEventGateway: messageEventGateway,
		logger:              log.NewLogger("SendMessageUseCase"),
	}
}

func (u *SendMessageUseCase) Execute(
	ctx context.Context,
	input *usecase.SendMessageUseCaseInput,
) (*usecase.SendMessageUseCaseOutput, error) {

	roomId, err := valueobject.NewIdWith(input.RoomId)
	if err != nil {
		return nil, err
	}

	senderId, err := valueobject.NewUserIdWith(input.SenderId)
	if err != nil {
		return nil, err
	}

	senderName, err := valueobject.NewUserNameWith(input.SenderName)
	if err != nil {
		return nil, err
	}

	text, err := valueobject.NewMessageTextWith(input.Text)
	if err != nil {
		return nil, err
	}

	room, err := u.roomRepository.FindById(ctx, roomId)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFoundRoom) {
			u.logger.Error(err)
		}

		return nil, err
	}

	if room.IsDeleted() {
		return nil, repository.ErrNotFoundRoom
	}

	message := entity.NewMessage(roomId, senderId, senderName, text)
	messageEvent := event.NewMessageEvent(message)

	err = u.messageRepository.Save(ctx, message)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	err = u.messageEventGateway.Send(ctx, messageEvent)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	output := &usecase.SendMessageUseCaseOutput{
		MessageId: message.Id().Value(),
	}

	return output, nil
}
