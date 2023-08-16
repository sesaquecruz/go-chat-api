package usecase

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/search"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type SearchRoomUseCaseInput struct {
	Page   string
	Size   string
	Sort   string
	Search string
}

type SearchRoomUseCaseInterface interface {
	Execute(ctx context.Context, input *SearchRoomUseCaseInput) (*search.Page[*entity.Room], error)
}

type SearchRoomUseCase struct {
	roomGateway gateway.RoomGatewayInterface
	logger      *log.Logger
}

func NewSearchRoomUseCase(roomGateway gateway.RoomGatewayInterface) *SearchRoomUseCase {
	return &SearchRoomUseCase{
		roomGateway: roomGateway,
		logger:      log.NewLogger("SearchRoomUseCase"),
	}
}

func (u *SearchRoomUseCase) Execute(ctx context.Context, input *SearchRoomUseCaseInput) (*search.Page[*entity.Room], error) {
	query, err := search.NewQuery(input.Page, input.Size, input.Sort, input.Search)
	if err != nil {
		return nil, err
	}

	page, err := u.roomGateway.Search(ctx, query)
	if err != nil {
		u.logger.Error(err)
		return nil, validation.NewInternalError(err)
	}

	return page, nil
}
