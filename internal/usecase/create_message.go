package usecase

import (
	"context"
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

type CreateMessageUseCase interface {
	Execute(ctx context.Context, input *CreateMessageUseCaseInput) (*CreateMessageUseCaseOutput, error)
}
