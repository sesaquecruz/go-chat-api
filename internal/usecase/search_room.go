package usecase

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
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

type SearchRoomUseCaseInterface interface {
	Execute(ctx context.Context, input *SearchRoomUseCaseInput) (*pagination.Page[*SearchRoomUseCaseOutput], error)
}

type SearchRoomUseCase struct {
	roomRepository repository.RoomRepositoryInterface
	logger         *log.Logger
}

func NewSearchRoomUseCase(roomRepository repository.RoomRepositoryInterface) *SearchRoomUseCase {
	return &SearchRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("SearchRoomUseCase"),
	}
}

func (u *SearchRoomUseCase) Execute(ctx context.Context, input *SearchRoomUseCaseInput) (*pagination.Page[*SearchRoomUseCaseOutput], error) {
	query, err := pagination.NewQuery(input.Page, input.Size, input.Sort, input.Search)
	if err != nil {
		return nil, err
	}

	rooms, err := u.roomRepository.Search(ctx, query)
	if err != nil {
		u.logger.Error(err)
		return nil, validation.NewInternalError(err)
	}

	output := pagination.MapPage[*entity.Room, *SearchRoomUseCaseOutput](rooms, func(r *entity.Room) *SearchRoomUseCaseOutput {
		return &SearchRoomUseCaseOutput{
			Id:        r.Id().Value(),
			AdminId:   r.AdminId().Value(),
			Name:      r.Name().Value(),
			Category:  r.Category().Value(),
			CreatedAt: r.CreatedAt().String(),
			UpdatedAt: r.UpdatedAt().String(),
		}
	})

	return output, nil
}
