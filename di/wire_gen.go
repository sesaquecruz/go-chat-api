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
	"github.com/sesaquecruz/go-chat-api/internal/infra/web"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
)

// Injectors from wire.go:

// Factories
func NewApiRouter(db *config.DatabaseConfig, broker *config.BrokerConfig, api *config.APIConfig) *gin.Engine {
	sqlDB := database.DbConnection(db)
	roomRepository := database.NewRoomRepository(sqlDB)
	createRoomUseCase := usecase.NewCreateRoomUseCase(roomRepository)
	searchRoomUseCase := usecase.NewSearchRoomUseCase(roomRepository)
	findRoomUseCase := usecase.NewFindRoomUseCase(roomRepository)
	updateRoomUseCase := usecase.NewUpdateRoomUseCase(roomRepository)
	deleteRoomUseCase := usecase.NewDeleteRoomUseCase(roomRepository)
	roomHandler := web.NewRoomHandler(createRoomUseCase, searchRoomUseCase, findRoomUseCase, updateRoomUseCase, deleteRoomUseCase)
	messageRepository := database.NewMessageRepository(sqlDB)
	connection := event.BrokerConnection(broker)
	messageGateway := event.NewMessageGateway(connection)
	createMessageUseCase := usecase.NewCreateMessageUseCase(roomRepository, messageRepository, messageGateway)
	messageHandler := web.NewMessageHandler(createMessageUseCase)
	engine := web.ApiRouter(api, roomHandler, messageHandler)
	return engine
}

// wire.go:

// Repositories
var setRoomRepositoryInterface = wire.NewSet(database.NewRoomRepository, wire.Bind(new(repository.RoomRepositoryInterface), new(*database.RoomRepository)))

var setMessageRepositoryInterface = wire.NewSet(database.NewMessageRepository, wire.Bind(new(repository.MessageRepositoryInterface), new(*database.MessageRepository)))

// Gateways
var setMessageGatewayInterface = wire.NewSet(event.NewMessageGateway, wire.Bind(new(gateway.MessageGatewayInterface), new(*event.MessageGateway)))

// Use Cases
var setCreateRoomUseCaseInterface = wire.NewSet(usecase.NewCreateRoomUseCase, wire.Bind(new(usecase.CreateRoomUseCaseInterface), new(*usecase.CreateRoomUseCase)))

var setSearchRoomUseCaseInterface = wire.NewSet(usecase.NewSearchRoomUseCase, wire.Bind(new(usecase.SearchRoomUseCaseInterface), new(*usecase.SearchRoomUseCase)))

var setFindRoomUseCaseInterface = wire.NewSet(usecase.NewFindRoomUseCase, wire.Bind(new(usecase.FindRoomUseCaseInterface), new(*usecase.FindRoomUseCase)))

var setUpdateRoomUseCaseInterface = wire.NewSet(usecase.NewUpdateRoomUseCase, wire.Bind(new(usecase.UpdateRoomUseCaseInterface), new(*usecase.UpdateRoomUseCase)))

var setDeleteRoomUseCaseInterface = wire.NewSet(usecase.NewDeleteRoomUseCase, wire.Bind(new(usecase.DeleteRoomUseCaseInterface), new(*usecase.DeleteRoomUseCase)))

var setCreateMessageUseCaseInterface = wire.NewSet(usecase.NewCreateMessageUseCase, wire.Bind(new(usecase.CreateMessageUseCaseInterface), new(*usecase.CreateMessageUseCase)))

// Handlers
var setRoomHandlerInterface = wire.NewSet(web.NewRoomHandler, wire.Bind(new(web.RoomHandlerInterface), new(*web.RoomHandler)))

var setMessageHandlerInterface = wire.NewSet(web.NewMessageHandler, wire.Bind(new(web.MessageHandlerInterface), new(*web.MessageHandler)))
