package room

import (
	usecase "github.com/sesaquecruz/go-chat-api/internal/usecase"

	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type RoomHandler struct {
	createRoomUseCase  usecase.CreateRoomUseCase
	searchRoomUseCase  usecase.SearchRoomUseCase
	findRoomUseCase    usecase.FindRoomUseCase
	updateRoomUseCase  usecase.UpdateRoomUseCase
	deleteRoomUseCase  usecase.DeleteRoomUseCase
	sendMessageUseCase usecase.SendMessageUseCase
	logger             *log.Logger
}

func NewRoomHandler(
	createRoomUseCase usecase.CreateRoomUseCase,
	searchRoomUseCase usecase.SearchRoomUseCase,
	findRoomUseCase usecase.FindRoomUseCase,
	updateRoomUseCase usecase.UpdateRoomUseCase,
	deleteRoomUseCase usecase.DeleteRoomUseCase,
	sendMessageUseCase usecase.SendMessageUseCase,
) *RoomHandler {
	return &RoomHandler{
		createRoomUseCase:  createRoomUseCase,
		searchRoomUseCase:  searchRoomUseCase,
		findRoomUseCase:    findRoomUseCase,
		updateRoomUseCase:  updateRoomUseCase,
		deleteRoomUseCase:  deleteRoomUseCase,
		sendMessageUseCase: sendMessageUseCase,
		logger:             log.NewLogger("RoomHandler"),
	}
}
