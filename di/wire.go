//go:build wireinject
// +build wireinject

package di

import (
	"database/sql"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"

	"github.com/google/wire"
)

var setRoomGatewayInterface = wire.NewSet(
	database.NewRoomPostgresGateway,
	wire.Bind(new(gateway.RoomGatewayInterface), new(*database.RoomPostgresGateway)),
)

var setCreateRoomUseCaseInterface = wire.NewSet(
	usecase.NewCreateRoomUseCase,
	wire.Bind(new(usecase.CreateRoomUseCaseInterface), new(*usecase.CreateRoomUseCase)),
)

var setRoomHandlerInterface = wire.NewSet(
	handler.NewRoomHandler,
	wire.Bind(new(handler.RoomHandlerInterface), new(*handler.RoomHandler)),
)

func NewApiRouter(cfg *config.APIConfig, db *sql.DB) *web.ApiRouter {
	wire.Build(
		setRoomGatewayInterface,
		setCreateRoomUseCaseInterface,
		setRoomHandlerInterface,
		web.NewApiRouter,
	)
	return &web.ApiRouter{}
}
