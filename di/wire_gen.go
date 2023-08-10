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
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/router"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
)

// Injectors from wire.go:

func NewApiRouter(db *config.DatabaseConfig, api *config.APIConfig) *gin.Engine {
	sqlDB := database.PostgresDb(db)
	roomPostgresGateway := database.NewRoomPostgresGateway(sqlDB)
	createRoomUseCase := usecase.NewCreateRoomUseCase(roomPostgresGateway)
	findRoomUseCase := usecase.NewFindRoomUseCase(roomPostgresGateway)
	roomHandler := handler.NewRoomHandler(createRoomUseCase, findRoomUseCase)
	engine := router.ApiRouter(api, roomHandler)
	return engine
}

// wire.go:

var setRoomGatewayInterface = wire.NewSet(database.NewRoomPostgresGateway, wire.Bind(new(gateway.RoomGatewayInterface), new(*database.RoomPostgresGateway)))

var setCreateRoomUseCaseInterface = wire.NewSet(usecase.NewCreateRoomUseCase, wire.Bind(new(usecase.CreateRoomUseCaseInterface), new(*usecase.CreateRoomUseCase)))

var setFindRoomUseCaseInterface = wire.NewSet(usecase.NewFindRoomUseCase, wire.Bind(new(usecase.FindRoomUseCaseInterface), new(*usecase.FindRoomUseCase)))

var setRoomHandlerInterface = wire.NewSet(handler.NewRoomHandler, wire.Bind(new(handler.RoomHandlerInterface), new(*handler.RoomHandler)))
