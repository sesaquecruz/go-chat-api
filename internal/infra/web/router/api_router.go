package router

import (
	"github.com/sesaquecruz/go-chat-api/config"
	_ "github.com/sesaquecruz/go-chat-api/docs"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-chat-api/pkg/health"
	"github.com/sesaquecruz/go-chat-api/pkg/middleware"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ApiRouter(
	cfg *config.ApiConfig,
	healthCheck health.Health,
	roomHandler handler.RoomHandler,
) *gin.Engine {
	gin.SetMode(cfg.Mode)

	r := gin.New()
	r.Use(middleware.CorsMiddleware(cfg.AllowOrigins))

	api := r.Group(cfg.Path)
	{

		api.GET("/healthz", gin.WrapH(healthCheck.Handler()))
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		api.Use(middleware.JwtMiddleware(cfg.JwtIssuer, []string{cfg.JwtAudience}))

		RoomRouter(api, roomHandler)
	}

	return r
}
