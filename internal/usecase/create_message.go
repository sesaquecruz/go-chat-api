package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type CreateMessageUseCaseInput struct {
	RoomId     string
	SenderId   string
	SenderName string
	Text       string
}

type CreateMessageUseCaseOutput struct {
	MessageId string
}

type CreateMessageUseCaseInterface interface {
	Execute(ctx context.Context, input *CreateMessageUseCaseInput) (*CreateMessageUseCaseOutput, error)
}

type CreateMessageUseCase struct {
	roomRepository       repository.RoomRepositoryInterface
	messageRepository    repository.MessageRepositoryInterface
	messageSenderGateway gateway.MessageSenderGatewayInterface
	logger               *log.Logger
}

func NewCreateMessageUseCase(
	roomRepository repository.RoomRepositoryInterface,
	messageRepository repository.MessageRepositoryInterface,
	messageSenderGateway gateway.MessageSenderGatewayInterface,
) *CreateMessageUseCase {
	return &CreateMessageUseCase{
		roomRepository:       roomRepository,
		messageRepository:    messageRepository,
		messageSenderGateway: messageSenderGateway,
		logger:               log.NewLogger("CreateMessageUseCase"),
	}
}

func (u *CreateMessageUseCase) Execute(ctx context.Context, input *CreateMessageUseCaseInput) (*CreateMessageUseCaseOutput, error) {
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

	message, err := entity.NewMessage(roomId, senderId, senderName, text)
	if err != nil {
		return nil, err
	}

	_, err = u.roomRepository.FindById(ctx, roomId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, validation.ErrNotFoundRoom
		}

		u.logger.Error(err)
		return nil, validation.NewInternalError(err)
	}

	err = u.messageRepository.Save(ctx, message)
	if err != nil {
		u.logger.Error(err)
		return nil, validation.NewInternalError(err)
	}

	err = u.messageSenderGateway.Send(ctx, message)
	if err != nil {
		u.logger.Error(err)
		return nil, validation.NewInternalError(err)
	}

	output := &CreateMessageUseCaseOutput{
		MessageId: message.Id().Value(),
	}

	return output, nil
}
