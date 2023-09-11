package room

import (
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

// FindRoom godoc
//
// @Summary		Find a room
// @Description	Find a chat room.
// @Tags		rooms
// @Accept		json
// @Produce		json
// @Param		id					path				string	true	"Room Id"
// @Success		200 {object}		dto.RoomResponse
// @Failure		400	{object}		dto.HttpError
// @Failure		401
// @Failure		404	{object}		dto.HttpError
// @Failure		500
// @Security	Bearer token
// @Router		/rooms/{id} 		[get]
func (h *RoomHandler) FindRoom(c *gin.Context) {
	input := &usecase.FindRoomUseCaseInput{
		Id: c.Param("id"),
	}

	output, err := h.findRoomUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(validation.ValidationError); ok {
			dto.AbortWithHttpError(c, http.StatusBadRequest, err)
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

	responseBody := &dto.RoomResponse{
		Id:       output.Id,
		Name:     output.Name,
		Category: output.Category,
	}

	c.JSON(http.StatusOK, responseBody)
}
