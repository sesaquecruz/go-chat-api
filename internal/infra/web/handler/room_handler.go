package handler

import (
	"fmt"
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/pkg"

	"github.com/gin-gonic/gin"
)

type RoomHandlerInterface interface {
	CreateRoom(c *gin.Context)
}

type RoomHandler struct {
	createRoomUseCase usecase.CreateRoomUseCaseInterface
	logger            *pkg.Logger
}

func NewRoomHandler(createRoomUseCase usecase.CreateRoomUseCaseInterface) *RoomHandler {
	return &RoomHandler{
		createRoomUseCase: createRoomUseCase,
		logger:            pkg.NewLogger("RoomHandler"),
	}
}

func (rh *RoomHandler) CreateRoom(c *gin.Context) {
	jwtClaims, err := pkg.JwtClaims(c)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var requestBody CreateRoomRequest
	if err := c.BindJSON(&requestBody); err != nil {
		return
	}

	input := usecase.CreateRoomUseCaseInput{
		AdminId:  jwtClaims.Subject,
		Name:     requestBody.Name,
		Category: requestBody.Category,
	}

	output, err := rh.createRoomUseCase.Execute(c.Request.Context(), &input)
	if err != nil {
		if _, ok := err.(*domain.ValidationError); ok {
			c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}

		rh.logger.Error(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Header("Location", fmt.Sprintf("/rooms/%s", output.RoomId))
	c.Status(http.StatusCreated)
}
