package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	domain_event "github.com/sesaquecruz/go-chat-api/internal/domain/event"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/event"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/dto"
	room_handler "github.com/sesaquecruz/go-chat-api/internal/infra/web/handler/impl/room"
	usecase "github.com/sesaquecruz/go-chat-api/internal/usecase/impl"
	"github.com/sesaquecruz/go-chat-api/pkg/health"
	"github.com/sesaquecruz/go-chat-api/test/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var db, _ = services.NewPostgresContainer(context.Background(), "file://../../../../")
var broker, _ = services.NewRabbitmqContainer(context.Background(), "../../../../")
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

	createRoomUseCase := usecase.NewCreateRoomUseCase(roomRepository)
	findRoomUseCase := usecase.NewFindRoomUseCase(roomRepository)
	searchRoomUseCase := usecase.NewSearchRoomUseCase(roomRepository)
	updateRoomUsecase := usecase.NewUpdateRoomUseCase(roomRepository)
	deleteRoomUseCase := usecase.NewDeleteRoomUseCase(roomRepository)
	createMessageUseCase := usecase.NewSendMessageUseCase(roomRepository, messageRepository, messageEventGateway)

	health := health.NewHealthCheck(db, conn)

	roomHandler := room_handler.NewRoomHandler(
		createRoomUseCase,
		searchRoomUseCase,
		findRoomUseCase,
		updateRoomUsecase,
		deleteRoomUseCase,
		createMessageUseCase,
	)

	router := ApiRouter(&config.ApiConfig{
		Port:         "",
		Path:         "/api/v1",
		Mode:         "release",
		AllowOrigins: "*",
		JwtIssuer:    auth.GetIssuer(),
		JwtAudience:  auth.GetAudience(),
	},
		health,
		roomHandler,
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

func createARoom(adminId, name, category string) *entity.Room {
	roomAdminId, _ := valueobject.NewUserIdWith(adminId)
	roomName, _ := valueobject.NewRoomNameWith(name)
	roomCategory, _ := valueobject.NewRoomCategoryWith(category)
	room := entity.NewRoom(roomAdminId, roomName, roomCategory)
	return room
}

func (s *RouterTestSuite) TestShouldReturnUnauthorizedWhenUnauthenticated() {
	defer db.Clear()
	t := s.T()
	r := s.router

	testCases := []struct {
		test   string
		method string
		url    string
	}{
		{
			"post room",
			http.MethodPost,
			"/api/v1/rooms",
		},
		{
			"get room",
			http.MethodGet,
			"/api/v1/rooms",
		},
		{
			"get room with id",
			http.MethodGet,
			"/api/v1/rooms/id",
		},
		{
			"put room with id",
			http.MethodPut,
			"/api/v1/rooms/id",
		},
		{
			"delete room with id",
			http.MethodDelete,
			"/api/v1/rooms/id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tc.method, tc.url, nil)

			r.ServeHTTP(w, req)
			res := w.Result()

			assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		})
	}
}

func (s *RouterTestSuite) TestShouldCreateARoom() {
	defer db.Clear()
	t := s.T()
	r := s.router

	sub := auth.GenerateSub()
	jwt, _ := auth.GenerateJWT(sub)

	payload := struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}{
		"A Game",
		"Game",
	}

	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rooms", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	location := res.Header.Get("Location")
	assert.NotEmpty(t, location)

	id, err := valueobject.NewIdWith(location[len("/api/v1/rooms/"):])
	assert.NotNil(t, id)
	assert.Nil(t, err)

	room, err := s.roomRepository.FindById(s.ctx, id)
	assert.NotNil(t, room)
	assert.Nil(t, err)
	assert.Equal(t, sub, room.AdminId().Value())
	assert.Equal(t, payload.Name, room.Name().Value())
	assert.Equal(t, payload.Category, room.Category().Value())
}

func (s *RouterTestSuite) TestShouldReturnRoomPages() {
	defer db.Clear()
	t := s.T()
	r := s.router

	sub := auth.GenerateSub()
	jwt, _ := auth.GenerateJWT(sub)

	room1 := createARoom(sub, "Rust", "Tech")
	room2 := createARoom(sub, "Go", "Tech")
	room3 := createARoom(sub, "Java", "Tech")
	room4 := createARoom(sub, "Need for Speed Undergroud", "Game")
	room5 := createARoom(sub, "Need for Speed Most Wanted", "Game")

	s.roomRepository.Save(s.ctx, room1)
	s.roomRepository.Save(s.ctx, room2)
	s.roomRepository.Save(s.ctx, room3)
	s.roomRepository.Save(s.ctx, room4)
	s.roomRepository.Save(s.ctx, room5)

	testCases := []struct {
		query string
		page  dto.RoomPage
	}{
		{
			query: "?page=0&size=2&sort=asc&search=tech",
			page: dto.RoomPage{
				Page:  0,
				Size:  2,
				Total: 3,
				Rooms: []*dto.RoomResponse{
					{
						Id:       room2.Id().Value(),
						Name:     room2.Name().Value(),
						Category: room2.Category().Value(),
					},
					{
						Id:       room3.Id().Value(),
						Name:     room3.Name().Value(),
						Category: room3.Category().Value(),
					},
				},
			},
		},
		{
			query: "?page=1&size=2&sort=asc&search=tech",
			page: dto.RoomPage{
				Page:  1,
				Size:  2,
				Total: 3,
				Rooms: []*dto.RoomResponse{
					{
						Id:       room1.Id().Value(),
						Name:     room1.Name().Value(),
						Category: room1.Category().Value(),
					},
				},
			},
		},
		{
			query: "?page=0&size=3&sort=desc&search=speed",
			page: dto.RoomPage{
				Page:  0,
				Size:  3,
				Total: 2,
				Rooms: []*dto.RoomResponse{
					{
						Id:       room4.Id().Value(),
						Name:     room4.Name().Value(),
						Category: room4.Category().Value(),
					},
					{
						Id:       room5.Id().Value(),
						Name:     room5.Name().Value(),
						Category: room5.Category().Value(),
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms"+tc.query, nil)
		req.Header.Set("Authorization", "Bearer "+jwt)

		r.ServeHTTP(w, req)
		res := w.Result()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		assert.Nil(t, err)

		var page dto.RoomPage
		err = json.Unmarshal(body, &page)
		assert.Nil(t, err)

		assert.Equal(t, tc.page.Page, page.Page)
		assert.Equal(t, tc.page.Size, page.Size)
		assert.Equal(t, tc.page.Total, page.Total)
		assert.Equal(t, len(tc.page.Rooms), len(page.Rooms))

		for i := 0; i < len(tc.page.Rooms); i++ {
			expectedItem := tc.page.Rooms[i]
			actualItem := page.Rooms[i]

			assert.Equal(t, expectedItem.Id, actualItem.Id)
			assert.Equal(t, expectedItem.Name, actualItem.Name)
			assert.Equal(t, expectedItem.Category, actualItem.Category)
		}
	}
}

func (s *RouterTestSuite) TestShouldReturnARoom() {
	defer db.Clear()
	t := s.T()
	r := s.router

	sub := auth.GenerateSub()
	jwt, _ := auth.GenerateJWT(sub)

	adminId, _ := valueobject.NewUserIdWith(sub)
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room := entity.NewRoom(adminId, name, category)
	s.roomRepository.Save(s.ctx, room)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms/"+room.Id().Value(), nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	assert.Nil(t, err)

	var result dto.RoomResponse
	err = json.Unmarshal(body, &result)
	assert.Nil(t, err)

	assert.Equal(t, room.Id().Value(), result.Id)
	assert.Equal(t, room.Name().Value(), result.Name)
	assert.Equal(t, room.Category().Value(), result.Category)
}

func (s *RouterTestSuite) TestShouldUpdateARoom() {
	defer db.Clear()
	t := s.T()
	r := s.router

	sub := auth.GenerateSub()
	jwt, _ := auth.GenerateJWT(sub)

	adminId, _ := valueobject.NewUserIdWith(sub)
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room := entity.NewRoom(adminId, name, category)
	s.roomRepository.Save(s.ctx, room)

	payload := struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}{
		"Need for Speed",
		"Game",
	}

	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/rooms/"+room.Id().Value(), bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	savedRoom, err := s.roomRepository.FindById(s.ctx, room.Id())
	assert.NotNil(t, savedRoom)
	assert.Nil(t, err)
	assert.Equal(t, room.AdminId().Value(), savedRoom.AdminId().Value())
	assert.Equal(t, payload.Name, savedRoom.Name().Value())
	assert.Equal(t, payload.Category, savedRoom.Category().Value())
	assert.True(t, room.CreatedAt().Time().Equal(savedRoom.CreatedAt().Time()))
	assert.True(t, room.UpdatedAt().Time().Before(savedRoom.UpdatedAt().Time()))
}

func (s *RouterTestSuite) TestShouldDeleteARoom() {
	defer db.Clear()
	t := s.T()
	r := s.router

	sub := auth.GenerateSub()
	jwt, _ := auth.GenerateJWT(sub)

	adminId, _ := valueobject.NewUserIdWith(sub)
	name, _ := valueobject.NewRoomNameWith("A Game")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room := entity.NewRoom(adminId, name, category)

	s.roomRepository.Save(s.ctx, room)
	assert.False(t, room.IsDeleted())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/rooms/"+room.Id().Value(), nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	savedRoom, err := s.roomRepository.FindById(s.ctx, room.Id())
	assert.NotNil(t, savedRoom)
	assert.Nil(t, err)
	assert.True(t, savedRoom.IsDeleted())
}

func (s *RouterTestSuite) TestCreateMessage_ShouldCreateAMessage() {
	defer db.Clear()
	t := s.T()
	r := s.router

	userName := auth.GetNickname()
	userId := auth.GenerateSub()
	jwt, _ := auth.GenerateJWT(userId)

	room := createARoom(userId, "A Game", "Game")
	s.roomRepository.Save(s.ctx, room)

	payload := struct {
		Text string `json:"text"`
	}{
		"A text",
	}

	url := fmt.Sprintf("/api/v1/rooms/%s/send", room.Id().Value())
	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	msgs := make(chan *domain_event.MessageEvent)
	defer close(msgs)

	go func() {
		err := s.messageEventGateway.Receive(s.ctx, msgs)
		if err != nil {
			t.Error(err)
		}
	}()

	select {
	case msg := <-msgs:
		assert.Equal(t, room.Id().Value(), msg.RoomId)
		assert.Equal(t, userId, msg.SenderId)
		assert.Equal(t, userName, msg.SenderName)
		assert.Equal(t, payload.Text, msg.Text)
	case <-time.After(30 * time.Second):
		t.Fail()
	}
}
