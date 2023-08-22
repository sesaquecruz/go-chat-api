package usecase

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type FindRoomUseCaseInput struct {
	RoomId string
}

type FindRoomUseCaseOutput struct {
	Id        string
	AdminId   string
	Name      string
	Category  string
	CreatedAt string
	UpdatedAt string
}

type FindRoomUseCaseInterface interface {
	Execute(ctx context.Context, input *FindRoomUseCaseInput) (*FindRoomUseCaseOutput, error)
}

type FindRoomUseCase struct {
	roomRepository repository.RoomRepositoryInterface
	logger         *log.Logger
}

func NewFindRoomUseCase(roomRepository repository.RoomRepositoryInterface) *FindRoomUseCase {
	return &FindRoomUseCase{
		roomRepository: roomRepository,
		logger:         log.NewLogger("FindRoomUseCase"),
	}
}

func (u *FindRoomUseCase) Execute(ctx context.Context, input *FindRoomUseCaseInput) (*FindRoomUseCaseOutput, error) {
	id, err := valueobject.NewIdWith(input.RoomId)
	if err != nil {
		return nil, err
	}

	room, err := u.roomRepository.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, validation.ErrNotFoundRoom
		}

		u.logger.Error(err)
		return nil, validation.NewInternalError(err)
	}

	ouput := &FindRoomUseCaseOutput{
		Id:        room.Id().Value(),
		AdminId:   room.AdminId().Value(),
		Name:      room.Name().Value(),
		Category:  room.Category().Value(),
		CreatedAt: room.CreatedAt().String(),
		UpdatedAt: room.UpdatedAt().String(),
	}

	return ouput, nil
}
