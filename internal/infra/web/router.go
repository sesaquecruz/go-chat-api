package web

import (
	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Router(
	cfg *config.ApiConfig,
	roomHandler handler.RoomHandler,
	messageHandler handler.MessageHandler,
) *gin.Engine {
	gin.SetMode(cfg.Mode)

	r := gin.New()

	r.Use(middleware.CorsMiddleware(cfg.AllowOrigins))
	r.Use(middleware.JwtMiddleware(cfg.JwtIssuer, []string{cfg.JwtAudience}))

	path := r.Group(cfg.Path)

	path.POST("/rooms", roomHandler.CreateRoom)
	path.GET("/rooms", roomHandler.SearchRoom)
	path.GET("/rooms/:id", roomHandler.FindRoom)
	path.PUT("/rooms/:id", roomHandler.UpdateRoom)
	path.DELETE("/rooms/:id", roomHandler.DeleteRoom)

	path.POST("/rooms/:id/send", messageHandler.CreateMessage)

	return r
}
