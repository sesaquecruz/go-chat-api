package usecase

import (
	"context"
)

type DeleteRoomUseCaseInput struct {
	Id      string
	AdminId string
}

type DeleteRoomUseCase interface {
	Execute(ctx context.Context, input *DeleteRoomUseCaseInput) error
}
