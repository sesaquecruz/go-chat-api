// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/event"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler/impl/room"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/router"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/internal/usecase/impl"
	"github.com/sesaquecruz/go-chat-api/pkg/health"
)

// Injectors from wire.go:

// Factories
func NewRouter(db *config.DatabaseConfig, broker *config.BrokerConfig, api *config.ApiConfig) *gin.Engine {
	sqlDB := database.PostgresConnection(db)
	connection := event.RabbitMqConnection(broker)
	healthCheck := health.NewHealthCheck(sqlDB, connection)
	roomPostgresRepository := database.NewRoomPostgresRepository(sqlDB)
	createRoomUseCase := impl.NewCreateRoomUseCase(roomPostgresRepository)
	searchRoomUseCase := impl.NewSearchRoomUseCase(roomPostgresRepository)
	findRoomUseCase := impl.NewFindRoomUseCase(roomPostgresRepository)
	updateRoomUseCase := impl.NewUpdateRoomUseCase(roomPostgresRepository)
	deleteRoomUseCase := impl.NewDeleteRoomUseCase(roomPostgresRepository)
	messagePostgresRepository := database.NewMessagePostgresRepository(sqlDB)
	messageEventRabbitMqGateway := event.NewMessageEventRabbitMqGateway(connection)
	sendMessageUseCase := impl.NewSendMessageUseCase(roomPostgresRepository, messagePostgresRepository, messageEventRabbitMqGateway)
	roomHandler := room.NewRoomHandler(createRoomUseCase, searchRoomUseCase, findRoomUseCase, updateRoomUseCase, deleteRoomUseCase, sendMessageUseCase)
	engine := router.ApiRouter(api, healthCheck, roomHandler)
	return engine
}

// wire.go:

// Repositories
var setRoomRepository = wire.NewSet(database.NewRoomPostgresRepository, wire.Bind(new(repository.RoomRepository), new(*database.RoomPostgresRepository)))

var setMessageRepository = wire.NewSet(database.NewMessagePostgresRepository, wire.Bind(new(repository.MessageRepository), new(*database.MessagePostgresRepository)))

// Gateways
var setMessageEventGateway = wire.NewSet(event.NewMessageEventRabbitMqGateway, wire.Bind(new(gateway.MessageEventGateway), new(*event.MessageEventRabbitMqGateway)))

// Use Cases
var setCreateRoomUseCase = wire.NewSet(impl.NewCreateRoomUseCase, wire.Bind(new(usecase.CreateRoomUseCase), new(*impl.CreateRoomUseCase)))

var setSearchRoomUseCase = wire.NewSet(impl.NewSearchRoomUseCase, wire.Bind(new(usecase.SearchRoomUseCase), new(*impl.SearchRoomUseCase)))

var setFindRoomUseCase = wire.NewSet(impl.NewFindRoomUseCase, wire.Bind(new(usecase.FindRoomUseCase), new(*impl.FindRoomUseCase)))

var setUpdateRoomUseCase = wire.NewSet(impl.NewUpdateRoomUseCase, wire.Bind(new(usecase.UpdateRoomUseCase), new(*impl.UpdateRoomUseCase)))

var setDeleteRoomUseCase = wire.NewSet(impl.NewDeleteRoomUseCase, wire.Bind(new(usecase.DeleteRoomUseCase), new(*impl.DeleteRoomUseCase)))

var setSendMessageUseCase = wire.NewSet(impl.NewSendMessageUseCase, wire.Bind(new(usecase.SendMessageUseCase), new(*impl.SendMessageUseCase)))

// Health
var setHealth = wire.NewSet(health.NewHealthCheck, wire.Bind(new(health.Health), new(*health.HealthCheck)))

// Handlers
var setRoomHandler = wire.NewSet(room.NewRoomHandler, wire.Bind(new(handler.RoomHandler), new(*room.RoomHandler)))
