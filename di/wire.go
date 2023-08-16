//go:build wireinject
// +build wireinject

package di

import (
	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var setRoomGatewayInterface = wire.NewSet(
	database.NewRoomGateway,
	wire.Bind(new(gateway.RoomGatewayInterface), new(*database.RoomGateway)),
)

var setCreateRoomUseCaseInterface = wire.NewSet(
	usecase.NewCreateRoomUseCase,
	wire.Bind(new(usecase.CreateRoomUseCaseInterface), new(*usecase.CreateRoomUseCase)),
)

var setSearchRoomUseCaseInterface = wire.NewSet(
	usecase.NewSearchRoomUseCase,
	wire.Bind(new(usecase.SearchRoomUseCaseInterface), new(*usecase.SearchRoomUseCase)),
)

var setFindRoomUseCaseInterface = wire.NewSet(
	usecase.NewFindRoomUseCase,
	wire.Bind(new(usecase.FindRoomUseCaseInterface), new(*usecase.FindRoomUseCase)),
)

var setUpdateRoomUseCaseInterface = wire.NewSet(
	usecase.NewUpdateRoomUseCase,
	wire.Bind(new(usecase.UpdateRoomUseCaseInterface), new(*usecase.UpdateRoomUseCase)),
)

var setDeleteRoomUseCaseInterface = wire.NewSet(
	usecase.NewDeleteRoomUseCase,
	wire.Bind(new(usecase.DeleteRoomUseCaseInterface), new(*usecase.DeleteRoomUseCase)),
)

var setRoomHandlerInterface = wire.NewSet(
	web.NewRoomHandler,
	wire.Bind(new(web.RoomHandlerInterface), new(*web.RoomHandler)),
)

func NewApiRouter(db *config.DatabaseConfig, api *config.APIConfig) *gin.Engine {
	wire.Build(
		database.DbConnection,
		setRoomGatewayInterface,
		setCreateRoomUseCaseInterface,
		setSearchRoomUseCaseInterface,
		setFindRoomUseCaseInterface,
		setUpdateRoomUseCaseInterface,
		setDeleteRoomUseCaseInterface,
		setRoomHandlerInterface,
		web.ApiRouter,
	)
	return &gin.Engine{}
}
