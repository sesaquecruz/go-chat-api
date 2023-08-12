package router

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

type WebApiTestSuite struct {
	suite.Suite
	ctx         context.Context
	postgres    *test.PostgresContainer
	auth0       *test.Auth0Server
	roomGateway gateway.RoomGatewayInterface
	router      *gin.Engine
}

func (s *WebApiTestSuite) SetupTest() {
	ctx := context.Background()
	postgres, err := test.NewPostgresContainer(ctx, "file://../../../../migrations")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("postgres", postgres.DSN)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	postgres.ClearDB()

	auth0 := test.NewAuth0Server()
	go func() {
		auth0.Run()
	}()

	cfg := config.APIConfig{
		Port:         "",
		Path:         "/api/v1",
		Mode:         "release",
		AllowOrigins: "*",
		JwtIssuer:    auth0.GetIssuer(),
		JwtAudience:  auth0.GetAudience(),
	}

	roomPostgresGateway := database.NewRoomPostgresGateway(db)
	createRoomUseCase := usecase.NewCreateRoomUseCase(roomPostgresGateway)
	roomFindRoomUseCase := usecase.NewFindRoomUseCase(roomPostgresGateway)
	roomHandler := handler.NewRoomHandler(
		createRoomUseCase,
		roomFindRoomUseCase,
	)

	apiRouter := ApiRouter(&cfg, roomHandler)

	s.ctx = ctx
	s.postgres = postgres
	s.auth0 = auth0
	s.roomGateway = roomPostgresGateway
	s.router = apiRouter
}

func (s *WebApiTestSuite) TearDownSuite() {
	if err := s.postgres.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating postgres container: %s", err)
	}

	if err := s.auth0.Stop(s.ctx); err != nil {
		s.T().Fatalf("error terminating auth0 server: %s", err)
	}
}

func TestApiRouterTestSuite(t *testing.T) {
	suite.Run(t, new(WebApiTestSuite))
}
