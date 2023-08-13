package web

import (
	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func ApiRouter(
	cfg *config.APIConfig,
	roomHandler RoomHandlerInterface,
) *gin.Engine {
	gin.SetMode(cfg.Mode)

	r := gin.New()

	r.Use(middleware.CorsMiddleware(cfg.AllowOrigins))
	r.Use(middleware.JwtMiddleware(cfg.JwtIssuer, []string{cfg.JwtAudience}))

	api := r.Group(cfg.Path)

	api.POST("/rooms", roomHandler.CreateRoom)
	api.GET("/rooms/:id", roomHandler.FindRoom)
	api.PUT("/rooms/:id", roomHandler.UpdateRoom)

	return r
}
