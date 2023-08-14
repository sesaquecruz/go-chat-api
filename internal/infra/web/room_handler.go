package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/validation"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type RoomHandlerInterface interface {
	CreateRoom(c *gin.Context)
	FindRoom(c *gin.Context)
	UpdateRoom(c *gin.Context)
	DeleteRoom(c *gin.Context)
}

type RoomHandler struct {
	createRoomUseCase usecase.CreateRoomUseCaseInterface
	findRoomUseCase   usecase.FindRoomUseCaseInterface
	updateRoomUseCase usecase.UpdateRoomUseCaseInterface
	deleteRoomUseCase usecase.DeleteRoomUseCaseInterface
	logger            *log.Logger
}

func NewRoomHandler(
	createRoomUseCase usecase.CreateRoomUseCaseInterface,
	findRoomUseCase usecase.FindRoomUseCaseInterface,
	updateRoomUseCase usecase.UpdateRoomUseCaseInterface,
	deleteRoomUseCase usecase.DeleteRoomUseCaseInterface,
) *RoomHandler {
	return &RoomHandler{
		createRoomUseCase: createRoomUseCase,
		findRoomUseCase:   findRoomUseCase,
		updateRoomUseCase: updateRoomUseCase,
		deleteRoomUseCase: deleteRoomUseCase,
		logger:            log.NewLogger("RoomHandler"),
	}
}

func (h *RoomHandler) CreateRoom(c *gin.Context) {
	jwtClaims, err := middleware.JwtClaims(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var requestBody RoomRequestDto
	if err := c.BindJSON(&requestBody); err != nil {
		return
	}

	input := usecase.CreateRoomUseCaseInput{
		AdminId:  jwtClaims.Subject,
		Name:     requestBody.Name,
		Category: requestBody.Category,
	}

	output, err := h.createRoomUseCase.Execute(c.Request.Context(), &input)
	if err != nil {
		if _, ok := err.(*validation.InternalError); ok {
			h.logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.Header("Location", fmt.Sprintf("%s/%s", c.Request.URL, output.RoomId))
	c.Status(http.StatusCreated)
}

func (h *RoomHandler) FindRoom(c *gin.Context) {
	id := c.Param("id")

	input := usecase.FindRoomUseCaseInput{RoomId: id}
	output, err := h.findRoomUseCase.Execute(c.Request.Context(), &input)
	if err != nil {
		if _, ok := err.(*validation.InternalError); ok {
			h.logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	res := RoomResponseDto{
		Id:       output.Id,
		Name:     output.Name,
		Category: output.Category,
	}

	c.JSON(http.StatusOK, res)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	id := c.Param("id")

	jwtClaims, err := middleware.JwtClaims(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var requestBody RoomRequestDto
	if err := c.BindJSON(&requestBody); err != nil {
		return
	}

	input := usecase.UpdateRoomUseCaseInput{
		Id:       id,
		AdminId:  jwtClaims.Subject,
		Name:     requestBody.Name,
		Category: requestBody.Category,
	}

	if err = h.updateRoomUseCase.Execute(c.Request.Context(), &input); err != nil {
		if _, ok := err.(*validation.InternalError); ok {
			h.logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if errors.Is(err, validation.ErrInvalidRoomAdmin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, validation.ErrNotFoundRoom) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *RoomHandler) DeleteRoom(c *gin.Context) {
	id := c.Param("id")

	jwtClaims, err := middleware.JwtClaims(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	input := usecase.DeleteRoomUseCaseInput{
		Id:      id,
		AdminId: jwtClaims.Subject,
	}

	if err = h.deleteRoomUseCase.Execute(c.Request.Context(), &input); err != nil {
		if _, ok := err.(*validation.InternalError); ok {
			h.logger.Error(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if errors.Is(err, validation.ErrInvalidRoomAdmin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		if errors.Is(err, validation.ErrNotFoundRoom) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
