package web

import (
	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRouter(
	cfg *config.APIConfig,
	roomHandler RoomHandlerInterface,
	messageHandler MessageHandlerInterface,
) *gin.Engine {
	gin.SetMode(cfg.Mode)

	r := gin.New()

	r.Use(middleware.CorsMiddleware(cfg.AllowOrigins))
	r.Use(middleware.JwtMiddleware(cfg.JwtIssuer, []string{cfg.JwtAudience}))

	api := r.Group(cfg.Path)

	api.POST("/rooms", roomHandler.CreateRoom)
	api.GET("/rooms", roomHandler.SearchRoom)
	api.GET("/rooms/:id", roomHandler.FindRoom)
	api.PUT("/rooms/:id", roomHandler.UpdateRoom)
	api.DELETE("/rooms/:id", roomHandler.DeleteRoom)

	api.POST("/messages", messageHandler.CreateMessage)

	return r
}
