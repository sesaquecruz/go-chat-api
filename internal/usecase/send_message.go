package usecase

import (
	"context"
)

type SendMessageUseCaseInput struct {
	RoomId     string
	SenderId   string
	SenderName string
	Text       string
}

type SendMessageUseCaseOutput struct {
	MessageId string
}

type SendMessageUseCase interface {
	Execute(ctx context.Context, input *SendMessageUseCaseInput) (*SendMessageUseCaseOutput, error)
}
