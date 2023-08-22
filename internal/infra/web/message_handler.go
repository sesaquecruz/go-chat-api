package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type MessageHandlerInterface interface {
	CreateMessage(c *gin.Context)
}

type MessageHandler struct {
	createMessageUseCase usecase.CreateMessageUseCaseInterface
	logger               *log.Logger
}

func NewMessageHandler(
	createMessageUseCase usecase.CreateMessageUseCaseInterface,
) *MessageHandler {
	return &MessageHandler{
		createMessageUseCase: createMessageUseCase,
		logger:               log.NewLogger("MessageHandler"),
	}
}

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	jwtClaims, err := middleware.JwtClaims(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	userId := jwtClaims.Subject
	userNickname := jwtClaims.Nickname

	var requestBody MessageRequestDto
	if err := c.BindJSON(&requestBody); err != nil {
		return
	}

	input := usecase.CreateMessageUseCaseInput{
		RoomId:     requestBody.RoomId,
		SenderId:   userId,
		SenderName: userNickname,
		Text:       requestBody.Text,
	}

	output, err := h.createMessageUseCase.Execute(c.Request.Context(), &input)
	if err != nil {
		if _, ok := err.(*validation.InternalError); ok {
			h.logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if errors.Is(err, validation.ErrNotFoundRoom) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	location := fmt.Sprintf("%s/%s", c.Request.URL, output.MessageId)

	c.Header("Location", location)
	c.Status(http.StatusCreated)
}
