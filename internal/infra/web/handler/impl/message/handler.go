package message

import (
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

type MessageHandler struct {
	createMessageUseCase usecase.CreateMessageUseCase
	logger               *log.Logger
}

func NewMessageHandler(
	createMessageUseCase usecase.CreateMessageUseCase,
) *MessageHandler {
	return &MessageHandler{
		createMessageUseCase: createMessageUseCase,
		logger:               log.NewLogger("MessageHandler"),
	}
}
