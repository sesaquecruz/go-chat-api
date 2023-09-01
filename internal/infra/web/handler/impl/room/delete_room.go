package room

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"
)

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	jwtClaims, err := middleware.JwtClaims(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	input := &usecase.DeleteRoomUseCaseInput{
		Id:      c.Param("id"),
		AdminId: jwtClaims.Subject,
	}

	err = h.deleteRoomUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(validation.ValidationError); ok {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if _, ok := err.(validation.NotFoundError); ok {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		if _, ok := err.(validation.UnauthorizedError); ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
