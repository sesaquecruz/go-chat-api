package web

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/event"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler/impl/message"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler/impl/room"
	"github.com/sesaquecruz/go-chat-api/internal/usecase/impl"
	"github.com/sesaquecruz/go-chat-api/pkg/health"
	"github.com/sesaquecruz/go-chat-api/test/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var db, _ = services.NewPostgresContainer(context.Background(), "file://../../../")
var broker, _ = services.NewRabbitmqContainer(context.Background(), "../../../")
var auth = services.NewAuth0Server()

type RouterTestSuite struct {
	suite.Suite
	ctx                 context.Context
	roomRepository      repository.RoomRepository
	messageRepository   repository.MessageRepository
	messageEventGateway gateway.MessageEventGateway
	router              *gin.Engine
}

func (s *RouterTestSuite) SetupTest() {
	db.Clear()

	db := database.PostgresConnection(&config.DatabaseConfig{
		Host:     db.Host,
		Port:     db.Port,
		User:     db.User,
		Password: db.Password,
		Name:     db.Name,
	})

	conn := event.RabbitMqConnection(&config.BrokerConfig{
		Host:     broker.Host,
		Port:     broker.Port,
		User:     broker.User,
		Password: broker.Password,
	})

	roomRepository := database.NewRoomPostgresRepository(db)
	messageRepository := database.NewMessagePostgresRepository(db)

	messageEventGateway := event.NewMessageEventRabbitMqGateway(conn)

	createRoomUseCase := impl.NewCreateRoomUseCase(roomRepository)
	findRoomUseCase := impl.NewFindRoomUseCase(roomRepository)
	searchRoomUseCase := impl.NewSearchRoomUseCase(roomRepository)
	updateRoomUsecase := impl.NewUpdateRoomUseCase(roomRepository)
	deleteRoomUseCase := impl.NewDeleteRoomUseCase(roomRepository)
	createMessageUseCase := impl.NewCreateMessageUseCase(roomRepository, messageRepository, messageEventGateway)

	health := health.NewHealthCheck(db, conn)

	roomHandler := room.NewRoomHandler(
		createRoomUseCase,
		searchRoomUseCase,
		findRoomUseCase,
		updateRoomUsecase,
		deleteRoomUseCase,
	)

	messageHandler := message.NewMessageHandler(
		createMessageUseCase,
		findRoomUseCase,
	)

	router := Router(&config.ApiConfig{
		Port:         "",
		Path:         "/api/v1",
		Mode:         "release",
		AllowOrigins: "*",
		JwtIssuer:    auth.GetIssuer(),
		JwtAudience:  auth.GetAudience(),
	},
		health,
		roomHandler,
		messageHandler,
	)

	s.ctx = context.Background()
	s.roomRepository = roomRepository
	s.messageRepository = messageRepository
	s.messageEventGateway = messageEventGateway
	s.router = router
}

func (s *RouterTestSuite) TearDownSuite() {
	if err := db.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating db container: %s", err)
	}

	if err := broker.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating broker container: %s", err)
	}

	if err := auth.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating auth server: %s", err)
	}
}

func TestRouterTestSuite(t *testing.T) {
	suite.Run(t, new(RouterTestSuite))
}

func (s *RouterTestSuite) TestHealth() {
	defer db.Clear()
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/healthz", nil)

	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func (s *RouterTestSuite) TestAllowOrigins() {
	defer db.Clear()
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/api/v1", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
}
