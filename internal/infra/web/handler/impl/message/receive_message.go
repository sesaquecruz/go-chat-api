package message

import (
	"io"
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

func (h *MessageHandler) ReceiveMessage(c *gin.Context) {
	roomId := c.Param("id")

	input := &usecase.FindRoomUseCaseInput{
		Id: roomId,
	}

	_, err := h.findRoomUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(validation.ValidationError); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, ok := err.(validation.NotFoundError); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	msgs, err := h.chat.Subscribe(c.Request.Context(), roomId)
	if err != nil {
		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	c.Stream(func(w io.Writer) bool {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return false
			}

			c.SSEvent("message", string(msg))
			return true

		default:
			return true
		}
	})

	c.Status(http.StatusOK)
}
