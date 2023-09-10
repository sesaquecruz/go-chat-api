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
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler/impl/message"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler/impl/room"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/internal/usecase/impl"

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
	impl.NewCreateRoomUseCase,
	wire.Bind(new(usecase.CreateRoomUseCase), new(*impl.CreateRoomUseCase)),
)

var setSearchRoomUseCase = wire.NewSet(
	impl.NewSearchRoomUseCase,
	wire.Bind(new(usecase.SearchRoomUseCase), new(*impl.SearchRoomUseCase)),
)

var setFindRoomUseCase = wire.NewSet(
	impl.NewFindRoomUseCase,
	wire.Bind(new(usecase.FindRoomUseCase), new(*impl.FindRoomUseCase)),
)

var setUpdateRoomUseCase = wire.NewSet(
	impl.NewUpdateRoomUseCase,
	wire.Bind(new(usecase.UpdateRoomUseCase), new(*impl.UpdateRoomUseCase)),
)

var setDeleteRoomUseCase = wire.NewSet(
	impl.NewDeleteRoomUseCase,
	wire.Bind(new(usecase.DeleteRoomUseCase), new(*impl.DeleteRoomUseCase)),
)

var setCreateMessageUseCase = wire.NewSet(
	impl.NewCreateMessageUseCase,
	wire.Bind(new(usecase.CreateMessageUseCase), new(*impl.CreateMessageUseCase)),
)

// Handlers
var setRoomHandler = wire.NewSet(
	room.NewRoomHandler,
	wire.Bind(new(handler.RoomHandler), new(*room.RoomHandler)),
)

var setMessageHandler = wire.NewSet(
	message.NewMessageHandler,
	wire.Bind(new(handler.MessageHandler), new(*message.MessageHandler)),
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
		setCreateMessageUseCase,

		// Handlers
		setRoomHandler,
		setMessageHandler,

		// Router
		web.Router,
	)

	return &gin.Engine{}
}
