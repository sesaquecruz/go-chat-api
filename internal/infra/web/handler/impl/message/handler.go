package message

import (
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type MessageHandler struct {
	createMessageUseCase usecase.CreateMessageUseCase
	findRoomUseCase      usecase.FindRoomUseCase
	logger               *log.Logger
}

func NewMessageHandler(
	createMessageUseCase usecase.CreateMessageUseCase,
	findRoomUseCase usecase.FindRoomUseCase,
) *MessageHandler {
	return &MessageHandler{
		createMessageUseCase: createMessageUseCase,
		findRoomUseCase:      findRoomUseCase,
		logger:               log.NewLogger("MessageHandler"),
	}
}
