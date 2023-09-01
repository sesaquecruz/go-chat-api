package message

import (
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
		RoomId:     c.Param("id"),
		SenderId:   jwtClaims.Subject,
		SenderName: jwtClaims.Nickname,
		Text:       requestBody.Text,
	}

	_, err = h.createMessageUseCase.Execute(c.Request.Context(), input)
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

	c.Status(http.StatusCreated)
}
