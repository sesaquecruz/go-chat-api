package room

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
)

func (h *RoomHandler) FindRoom(c *gin.Context) {
	input := &usecase.FindRoomUseCaseInput{
		Id: c.Param("id"),
	}

	output, err := h.findRoomUseCase.Execute(c.Request.Context(), input)
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

	responseBody := &dto.RoomResponseDto{
		Id:       output.Id,
		Name:     output.Name,
		Category: output.Category,
	}

	c.JSON(http.StatusOK, responseBody)
}
