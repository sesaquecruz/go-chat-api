package usecase

import (
	"context"
)

type CreateRoomUseCaseInput struct {
	AdminId  string
	Name     string
	Category string
}

type CreateRoomUseCaseOutput struct {
	RoomId string
}

type CreateRoomUseCase interface {
	Execute(ctx context.Context, input *CreateRoomUseCaseInput) (*CreateRoomUseCaseOutput, error)
}
