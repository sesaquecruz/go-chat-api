package web

import (
	"fmt"
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/errors"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

type RoomHandlerInterface interface {
	CreateRoom(c *gin.Context)
	FindRoom(c *gin.Context)
	UpdateRoom(c *gin.Context)
}

type RoomHandler struct {
	createRoomUseCase usecase.CreateRoomUseCaseInterface
	findRoomUseCase   usecase.FindRoomUseCaseInterface
	updateRoomUseCase usecase.UpdateRoomUseCaseInterface
	logger            *log.Logger
}

func NewRoomHandler(
	createRoomUseCase usecase.CreateRoomUseCaseInterface,
	findRoomUseCase usecase.FindRoomUseCaseInterface,
	updateRoomUseCase usecase.UpdateRoomUseCaseInterface,
) *RoomHandler {
	return &RoomHandler{
		createRoomUseCase: createRoomUseCase,
		findRoomUseCase:   findRoomUseCase,
		updateRoomUseCase: updateRoomUseCase,
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
		if _, ok := err.(*errors.ValidationError); ok {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
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
		if _, ok := err.(*errors.ValidationError); ok {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		if _, ok := err.(*errors.GatewayError); ok && err.Error() == gateway.ErrNotFoundRoom {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
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

	err = h.updateRoomUseCase.Execute(c.Request.Context(), &input)
	if err != nil {
		if _, ok := err.(*errors.ValidationError); ok {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		if _, ok := err.(*errors.GatewayError); ok && err.Error() == gateway.ErrNotFoundRoom {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if _, ok := err.(*errors.AuthorizationError); ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}

		h.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}