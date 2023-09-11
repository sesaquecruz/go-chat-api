package impl

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type SearchRoomUseCase struct {
	roomRepository repository.RoomRepository
	logger         *log.Logger
}

func NewSearchRoomUseCase(roomRepository repository.RoomRepository) *SearchRoomUseCase {
	return &SearchRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("SearchRoomUseCase"),
	}
}

func (u *SearchRoomUseCase) Execute(
	ctx context.Context,
	input *usecase.SearchRoomUseCaseInput,
) (*pagination.Page[*usecase.SearchRoomUseCaseOutput], error) {

	query, err := pagination.NewQuery(input.Page, input.Size, input.Sort, input.Search)
	if err != nil {
		return nil, err
	}

	page, err := u.roomRepository.Search(ctx, query)
	if err != nil {
		u.logger.Error(err)
		return nil, err
	}

	mapper := func(r *entity.Room) *usecase.SearchRoomUseCaseOutput {
		return &usecase.SearchRoomUseCaseOutput{
			Id:        r.Id().Value(),
			AdminId:   r.AdminId().Value(),
			Name:      r.Name().Value(),
			Category:  r.Category().Value(),
			CreatedAt: r.CreatedAt().Value(),
			UpdatedAt: r.UpdatedAt().Value(),
		}
	}

	output := pagination.MapPage[*entity.Room, *usecase.SearchRoomUseCaseOutput](page, mapper)

	return output, nil
}
