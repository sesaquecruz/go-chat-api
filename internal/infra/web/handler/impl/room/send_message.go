package room

import (
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

// SendMessage godoc
//
// @Summary		Send a message
// @Description	Send a message to the chat room.
// @Tags		rooms
// @Accept		json
// @Produce		json
// @Param		id					path			string				true	"Room Id"
// @Param		message				body			dto.MessageRequest	true	"Message"
// @Success		201
// @Failure		400
// @Failure		401
// @Failure		404	{object}		dto.HttpError
// @Failure		422	{object}		dto.HttpError
// @Failure		500
// @Security	Bearer token
// @Router		/rooms/{id}/send 	[post]
func (h *RoomHandler) SendMessage(c *gin.Context) {
	jwtClaims, err := middleware.JwtClaims(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var requestBody dto.MessageRequest

	err = c.BindJSON(&requestBody)
	if err != nil {
		return
	}

	input := &usecase.SendMessageUseCaseInput{
		RoomId:     c.Param("id"),
		SenderId:   jwtClaims.Subject,
		SenderName: jwtClaims.Nickname,
		Text:       requestBody.Text,
	}

	_, err = h.sendMessageUseCase.Execute(c.Request.Context(), input)
	if err != nil {
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
	}

	c.Status(http.StatusCreated)
}
