package room

import (
	"fmt"
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// CreateRoom godoc
//
// @Summary		Create a room
// @Description	Create a new chat room. The room categories are: [General, Tech, Game, Book, Movie, Music, Language, Science].
// @Tags		rooms
// @Accept		json
// @Produce		json
// @Param		room				body			dto.RoomRequest		true	"Room"
// @Success		201	{string} 		string			"Location"
// @Failure		400
// @Failure		401
// @Failure		422	{object}		dto.HttpError
// @Failure		500
// @Security	Bearer token
// @Router		/rooms 				[post]
func (h *RoomHandler) CreateRoom(c *gin.Context) {
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

	input := &usecase.CreateRoomUseCaseInput{
		AdminId:  jwtClaims.Subject,
		Name:     requestBody.Name,
		Category: requestBody.Category,
	}

	output, err := h.createRoomUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(validation.ValidationError); ok {
			dto.AbortWithHttpError(c, http.StatusUnprocessableEntity, err)
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	location := fmt.Sprintf("%s/%s", c.Request.URL, output.RoomId)

	c.Header("Location", location)
	c.Status(http.StatusCreated)
}
