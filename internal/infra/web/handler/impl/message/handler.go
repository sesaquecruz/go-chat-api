package message

import (
	"github.com/sesaquecruz/go-chat-api/internal/chat"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type MessageHandler struct {
	createMessageUseCase usecase.CreateMessageUseCase
	findRoomUseCase      usecase.FindRoomUseCase
	chat                 chat.Chat
	logger               *log.Logger
}

func NewMessageHandler(
	createMessageUseCase usecase.CreateMessageUseCase,
	findRoomUseCase usecase.FindRoomUseCase,
	chat chat.Chat,
) *MessageHandler {
	return &MessageHandler{
		createMessageUseCase: createMessageUseCase,
		findRoomUseCase:      findRoomUseCase,
		chat:                 chat,
		logger:               log.NewLogger("MessageHandler"),
	}
}
