package room

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sesaquecruz/go-chat-api/internal/domain/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
)

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
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	mapper := func(r *usecase.SearchRoomUseCaseOutput) *dto.RoomResponseDto {
		return &dto.RoomResponseDto{
			Id:       r.Id,
			Name:     r.Name,
			Category: r.Category,
		}
	}

	res := pagination.MapPage[*usecase.SearchRoomUseCaseOutput, *dto.RoomResponseDto](output, mapper)

	c.JSON(http.StatusOK, res)
}
