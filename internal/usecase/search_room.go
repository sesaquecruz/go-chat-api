package usecase

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/pagination"
)

type SearchRoomUseCaseInput struct {
	Page   string
	Size   string
	Sort   string
	Search string
}

type SearchRoomUseCaseOutput struct {
	Id        string
	AdminId   string
	Name      string
	Category  string
	CreatedAt string
	UpdatedAt string
}

type SearchRoomUseCase interface {
	Execute(ctx context.Context, input *SearchRoomUseCaseInput) (*pagination.Page[*SearchRoomUseCaseOutput], error)
}
