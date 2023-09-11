package usecase

import (
	"context"
)

type UpdateRoomUseCaseInput struct {
	Id       string
	AdminId  string
	Name     string
	Category string
}

type UpdateRoomUseCase interface {
	Execute(ctx context.Context, input *UpdateRoomUseCaseInput) error
}
