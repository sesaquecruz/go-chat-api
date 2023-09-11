package room

import (
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// DeleteRoom godoc
//
// @Summary		Delete a room
// @Description	Delete a chat room if the user is room admin.
// @Tags		rooms
// @Accept		json
// @Produce		json
// @Param		id					path			string	true	"Room Id"
// @Success		204
// @Failure		400
// @Failure		400 {object}		dto.HttpError
// @Failure		401
// @Failure		401 {object}		dto.HttpError
// @Failure		404 {object}		dto.HttpError
// @Failure		500
// @Security	Bearer token
// @Router		/rooms/{id} 		[delete]
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
			dto.AbortWithHttpError(c, http.StatusBadRequest, err)
			return
		}

		if _, ok := err.(validation.UnauthorizedError); ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if _, ok := err.(validation.NotFoundError); ok {
			dto.AbortWithHttpError(c, http.StatusNotFound, err)
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}
