package usecase

import (
	"context"
)

type FindRoomUseCaseInput struct {
	Id string
}

type FindRoomUseCaseOutput struct {
	Id        string
	AdminId   string
	Name      string
	Category  string
	CreatedAt string
	UpdatedAt string
}

type FindRoomUseCase interface {
	Execute(ctx context.Context, input *FindRoomUseCaseInput) (*FindRoomUseCaseOutput, error)
}
