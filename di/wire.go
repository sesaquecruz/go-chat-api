//go:build wireinject
// +build wireinject

package di

import (
	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/event"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	room_handler "github.com/sesaquecruz/go-chat-api/internal/infra/web/handler/impl/room"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/router"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	impl_usecase "github.com/sesaquecruz/go-chat-api/internal/usecase/impl"
	"github.com/sesaquecruz/go-chat-api/pkg/health"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Repositories
var setRoomRepository = wire.NewSet(
	database.NewRoomPostgresRepository,
	wire.Bind(new(repository.RoomRepository), new(*database.RoomPostgresRepository)),
)

var setMessageRepository = wire.NewSet(
	database.NewMessagePostgresRepository,
	wire.Bind(new(repository.MessageRepository), new(*database.MessagePostgresRepository)),
)

// Gateways
var setMessageEventGateway = wire.NewSet(
	event.NewMessageEventRabbitMqGateway,
	wire.Bind(new(gateway.MessageEventGateway), new(*event.MessageEventRabbitMqGateway)),
)

// Use Cases
var setCreateRoomUseCase = wire.NewSet(
	impl_usecase.NewCreateRoomUseCase,
	wire.Bind(new(usecase.CreateRoomUseCase), new(*impl_usecase.CreateRoomUseCase)),
)

var setSearchRoomUseCase = wire.NewSet(
	impl_usecase.NewSearchRoomUseCase,
	wire.Bind(new(usecase.SearchRoomUseCase), new(*impl_usecase.SearchRoomUseCase)),
)

var setFindRoomUseCase = wire.NewSet(
	impl_usecase.NewFindRoomUseCase,
	wire.Bind(new(usecase.FindRoomUseCase), new(*impl_usecase.FindRoomUseCase)),
)

var setUpdateRoomUseCase = wire.NewSet(
	impl_usecase.NewUpdateRoomUseCase,
	wire.Bind(new(usecase.UpdateRoomUseCase), new(*impl_usecase.UpdateRoomUseCase)),
)

var setDeleteRoomUseCase = wire.NewSet(
	impl_usecase.NewDeleteRoomUseCase,
	wire.Bind(new(usecase.DeleteRoomUseCase), new(*impl_usecase.DeleteRoomUseCase)),
)

var setSendMessageUseCase = wire.NewSet(
	impl_usecase.NewSendMessageUseCase,
	wire.Bind(new(usecase.SendMessageUseCase), new(*impl_usecase.SendMessageUseCase)),
)

// Health
var setHealth = wire.NewSet(
	health.NewHealthCheck,
	wire.Bind(new(health.Health), new(*health.HealthCheck)),
)

// Handlers
var setRoomHandler = wire.NewSet(
	room_handler.NewRoomHandler,
	wire.Bind(new(handler.RoomHandler), new(*room_handler.RoomHandler)),
)

// Factories
func NewRouter(
	db *config.DatabaseConfig,
	broker *config.BrokerConfig,
	api *config.ApiConfig,
) *gin.Engine {
	wire.Build(
		// Connections
		database.PostgresConnection,
		event.RabbitMqConnection,

		// Repositories
		setRoomRepository,
		setMessageRepository,

		// Gateways
		setMessageEventGateway,

		// Use Cases
		setCreateRoomUseCase,
		setSearchRoomUseCase,
		setFindRoomUseCase,
		setUpdateRoomUseCase,
		setDeleteRoomUseCase,
		setSendMessageUseCase,

		// Health
		setHealth,

		// Handlers
		setRoomHandler,

		// Router
		router.ApiRouter,
	)

	return &gin.Engine{}
}
