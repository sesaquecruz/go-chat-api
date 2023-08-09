package web

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/test"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

type ApiRouterTestSuite struct {
	suite.Suite
	ctx         context.Context
	postgres    *test.PostgresContainer
	iam         *test.IamServer
	cfg         *config.APIConfig
	roomGateway gateway.RoomGatewayInterface
	api         *ApiRouter
	apiUrl      string
}

func (s *ApiRouterTestSuite) SetupTest() {
	ctx := context.Background()
	postgres, err := test.NewPostgresContainer(ctx, "file://../../../migrations")
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

	iam := test.NewIamServer()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		log.Fatal(err)
	}

	cfg := config.APIConfig{
		Port:         strconv.Itoa(port),
		Path:         "/api/v1",
		GinMode:      "release",
		AllowOrigins: "*",
		JwtIssuer:    iam.GetIssuer(),
		JwtAudience:  iam.GetAudience(),
	}

	roomPostgresGateway := database.NewRoomPostgresGateway(db)
	createRoomUseCase := usecase.NewCreateRoomUseCase(roomPostgresGateway)
	roomHandler := handler.NewRoomHandler(createRoomUseCase)

	s.ctx = ctx
	s.postgres = postgres
	s.iam = iam
	s.cfg = &cfg
	s.roomGateway = roomPostgresGateway
	s.api = NewApiRouter(&cfg, roomHandler)
	s.apiUrl = fmt.Sprintf("http://127.0.0.1:%d/api/v1", port)

	go func() {
		s.iam.Run()
	}()

	go func() {
		s.api.Run()
	}()

	time.Sleep(5 * time.Second)
}

func (s *ApiRouterTestSuite) TearDownSuite() {
	if err := s.postgres.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating postgres container: %s", err)
	}

	if err := s.iam.Stop(s.ctx); err != nil {
		s.T().Fatalf("error terminating iam server: %s", err)
	}

	if err := s.api.Stop(s.ctx); err != nil {
		s.T().Fatalf("error terminating api server: %s", err)
	}
}

func (s *ApiRouterTestSuite) TestShouldReturnAllowOrigins() {
	t := s.T()

	req, err := http.NewRequest(http.MethodOptions, s.apiUrl, http.NoBody)
	assert.Nil(t, err)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	defer res.Body.Close()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)
	assert.Equal(t, s.cfg.AllowOrigins, res.Header.Get("Access-Control-Allow-Origin"))
}

func TestApiRouterTestSuite(t *testing.T) {
	suite.Run(t, new(ApiRouterTestSuite))
}
