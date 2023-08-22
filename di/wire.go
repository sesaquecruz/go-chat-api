//go:build wireinject
// +build wireinject

package di

import (
	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/event"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Repositories
var setRoomRepositoryInterface = wire.NewSet(
	database.NewRoomRepository,
	wire.Bind(new(repository.RoomRepositoryInterface), new(*database.RoomRepository)),
)

var setMessageRepositoryInterface = wire.NewSet(
	database.NewMessageRepository,
	wire.Bind(new(repository.MessageRepositoryInterface), new(*database.MessageRepository)),
)

// Gateways
var setMessageGatewayInterface = wire.NewSet(
	event.NewMessageGateway,
	wire.Bind(new(gateway.MessageGatewayInterface), new(*event.MessageGateway)),
)

// Use Cases
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

var setCreateMessageUseCaseInterface = wire.NewSet(
	usecase.NewCreateMessageUseCase,
	wire.Bind(new(usecase.CreateMessageUseCaseInterface), new(*usecase.CreateMessageUseCase)),
)

// Handlers
var setRoomHandlerInterface = wire.NewSet(
	web.NewRoomHandler,
	wire.Bind(new(web.RoomHandlerInterface), new(*web.RoomHandler)),
)

var setMessageHandlerInterface = wire.NewSet(
	web.NewMessageHandler,
	wire.Bind(new(web.MessageHandlerInterface), new(*web.MessageHandler)),
)

// Factories
func NewApiRouter(db *config.DatabaseConfig, broker *config.BrokerConfig, api *config.APIConfig) *gin.Engine {
	wire.Build(
		// Connections
		database.DbConnection,
		event.BrokerConnection,

		// Repositories
		setRoomRepositoryInterface,
		setMessageRepositoryInterface,

		// Gateways
		setMessageGatewayInterface,

		// Use Cases
		setCreateRoomUseCaseInterface,
		setSearchRoomUseCaseInterface,
		setFindRoomUseCaseInterface,
		setUpdateRoomUseCaseInterface,
		setDeleteRoomUseCaseInterface,
		setCreateMessageUseCaseInterface,

		// Handlers
		setRoomHandlerInterface,
		setMessageHandlerInterface,

		// Targets
		web.ApiRouter,
	)
	return &gin.Engine{}
}
