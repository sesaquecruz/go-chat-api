package room

import (
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// UpdateRoom godoc
//
// @Summary		Update a room
// @Description	Update a chat room if the user is room admin. The room categories are: [General, Tech, Game, Book, Movie, Music, Language, Science].
// @Tags		rooms
// @Accept		json
// @Produce		json
// @Param		id					path			string			true	"Room Id"
// @Param		room				body			dto.RoomRequest	true	"Room"
// @Success		204
// @Failure		400
// @Failure		401
// @Failure		401	{object}		dto.HttpError
// @Failure		404	{object}		dto.HttpError
// @Failure		422	{object}		dto.HttpError
// @Failure		500
// @Security	Bearer token
// @Router		/rooms/{id}			[put]
func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	jwtClaims, err := middleware.JwtClaims(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var requestBody dto.RoomRequest

	err = c.BindJSON(&requestBody)
	if err != nil {
		return
	}

	input := &usecase.UpdateRoomUseCaseInput{
		Id:       c.Param("id"),
		AdminId:  jwtClaims.Subject,
		Name:     requestBody.Name,
		Category: requestBody.Category,
	}

	err = h.updateRoomUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(validation.UnauthorizedError); ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if _, ok := err.(validation.NotFoundError); ok {
			dto.AbortWithHttpError(c, http.StatusNotFound, err)
			return
		}

		if _, ok := err.(validation.ValidationError); ok {
			dto.AbortWithHttpError(c, http.StatusUnprocessableEntity, err)
			return

		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return

	}

	c.Status(http.StatusNoContent)
}
