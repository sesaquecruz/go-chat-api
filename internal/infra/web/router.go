package web

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-chat-api/pkg"
)

type ApiRouterInterface interface {
	IsRunning() error
	Run() error
	Stop() error
}

type ApiRouter struct {
	path    string
	server  *http.Server
	running bool
	logger  *pkg.Logger
}

func NewApiRouter(cfg *config.APIConfig, roomHandler handler.RoomHandlerInterface) *ApiRouter {
	gin.SetMode(cfg.GinMode)
	router := gin.New()

	router.Use(pkg.CorsMiddleware(cfg.AllowOrigins))
	router.Use(pkg.JwtMiddleware(cfg.JwtIssuer, []string{cfg.JwtAudience}))

	apiPath := router.Group(cfg.Path)

	apiPath.POST("/rooms", roomHandler.CreateRoom)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router.Handler(),
	}

	return &ApiRouter{
		path:    cfg.Path,
		server:  &server,
		running: false,
		logger:  pkg.NewLogger("ApiRouter"),
	}
}

func (ar *ApiRouter) IsRunning() bool {
	return ar.running
}

func (ar *ApiRouter) Run() error {
	if !ar.running {
		ar.running = true
		ar.logger.Infof("using the path %s\n", ar.path)
		ar.logger.Infof("running server on %s\n", ar.server.Addr)
		return ar.server.ListenAndServe()
	}

	return errors.New("router already running")
}

func (ar *ApiRouter) Stop(ctx context.Context) error {
	if ar.running {
		ar.running = false
		ar.logger.Info("stopping server")
		return ar.server.Shutdown(ctx)
	}

	return errors.New("router is not running")
}
