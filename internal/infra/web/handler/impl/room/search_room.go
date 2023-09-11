package room

import (
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"

	"github.com/gin-gonic/gin"
)

// SearchRoom godoc
//
// @Summary		Search rooms
// @Description	Search chat rooms.
// @Tags		rooms
// @Accept		json
// @Produce		json
// @Param		page				query				string	false	"Page"			default(0)
// @Param		size				query				string	false	"Size"			default(10)
// @Param		sort				query				string	false	"Sort"			default(asc)
// @Param		search				query				string	false	"Search Term"	default()
// @Success		200	{array}			dto.RoomPage
// @Failure		400 {object}		dto.HttpError
// @Failure		401
// @Failure		500
// @Security	Bearer token
// @Router		/rooms		 		[get]
func (h *RoomHandler) SearchRoom(c *gin.Context) {
	input := &usecase.SearchRoomUseCaseInput{
		Page:   c.Query("page"),
		Size:   c.Query("size"),
		Sort:   c.Query("sort"),
		Search: c.Query("search"),
	}

	output, err := h.searchRoomUseCase.Execute(c.Request.Context(), input)
	if err != nil {
		if _, ok := err.(validation.ValidationError); ok {
			dto.AbortWithHttpError(c, http.StatusBadRequest, err)
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	mapper := func(r *usecase.SearchRoomUseCaseOutput) *dto.RoomResponse {
		return &dto.RoomResponse{
			Id:       r.Id,
			Name:     r.Name,
			Category: r.Category,
		}
	}

	result := pagination.MapPage[*usecase.SearchRoomUseCaseOutput, *dto.RoomResponse](output, mapper)

	page := &dto.RoomPage{
		Page:  result.Page,
		Size:  result.Size,
		Total: result.Total,
		Rooms: result.Items,
	}

	c.JSON(http.StatusOK, page)
}
