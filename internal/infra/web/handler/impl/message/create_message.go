package message

import (
	"fmt"
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func (h *MessageHandler) CreateMessage(c *gin.Context) {
	jwtClaims, err := middleware.JwtClaims(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var requestBody dto.MessageRequestDto

	err = c.BindJSON(&requestBody)
	if err != nil {
		return
	}

	input := &usecase.CreateMessageUseCaseInput{
		RoomId:     requestBody.RoomId,
		SenderId:   jwtClaims.Subject,
		SenderName: jwtClaims.Nickname,
		Text:       requestBody.Text,
	}

	output, err := h.createMessageUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(validation.ValidationError); ok {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		if _, ok := err.(validation.NotFoundError); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	location := fmt.Sprintf("%s/%s", c.Request.URL, output.MessageId)

	c.Header("Location", location)
	c.Status(http.StatusCreated)
}
