package web

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository"
	"github.com/sesaquecruz/go-chat-api/internal/domain/repository/pagination"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/infra/event"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

var dbApi, _ = test.NewPostgresContainer(context.Background(), "file://../../../")
var brokerApi, _ = test.NewRabbitmqContainer(context.Background(), "../../../")
var authApi = test.NewAuth0Server()

type ApiRouterTestSuite struct {
	suite.Suite
	ctx                    context.Context
	roomRepository         repository.RoomRepositoryInterface
	messageRepository      repository.MessageRepositoryInterface
	messageSenderGateway   gateway.MessageSenderGatewayInterface
	messageRecevierGateway gateway.MessageReceiverGatewayInterface
	router                 *gin.Engine
}

func (s *ApiRouterTestSuite) SetupTest() {
	dbApi.Clear()

	db := database.DbConnection(&config.DatabaseConfig{
		Host:     dbApi.Host,
		Port:     dbApi.Port,
		User:     dbApi.User,
		Password: dbApi.Password,
		Name:     dbApi.Name,
	})

	conn := event.BrokerConnection(&config.BrokerConfig{
		Host:     brokerApi.Host,
		Port:     brokerApi.Port,
		User:     brokerApi.User,
		Password: brokerApi.Password,
	})

	roomRepository := database.NewRoomRepository(db)
	messageRepository := database.NewMessageRepository(db)

	messageSenderGateway := event.NewMessageGateway(conn)
	messageReceiverGateway := event.NewMessageGateway(conn)

	createRoomUseCase := usecase.NewCreateRoomUseCase(roomRepository)
	findRoomUseCase := usecase.NewFindRoomUseCase(roomRepository)
	searchRoomUseCase := usecase.NewSearchRoomUseCase(roomRepository)
	updateRoomUsecase := usecase.NewUpdateRoomUseCase(roomRepository)
	deleteRoomUseCase := usecase.NewDeleteRoomUseCase(roomRepository)
	createMessageUseCase := usecase.NewCreateMessageUseCase(roomRepository, messageRepository, messageSenderGateway)

	roomHandler := NewRoomHandler(
		createRoomUseCase,
		searchRoomUseCase,
		findRoomUseCase,
		updateRoomUsecase,
		deleteRoomUseCase,
	)

	messageHandler := NewMessageHandler(
		createMessageUseCase,
	)

	apiRouter := ApiRouter(&config.APIConfig{
		Port:         "",
		Path:         "/api/v1",
		Mode:         "release",
		AllowOrigins: "*",
		JwtIssuer:    authApi.GetIssuer(),
		JwtAudience:  authApi.GetAudience(),
	},
		roomHandler,
		messageHandler,
	)

	s.ctx = context.Background()
	s.roomRepository = roomRepository
	s.messageRepository = messageRepository
	s.messageSenderGateway = messageSenderGateway
	s.messageRecevierGateway = messageReceiverGateway
	s.router = apiRouter
}

func (s *ApiRouterTestSuite) TearDownSuite() {
	if err := brokerApi.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating broker container: %s", err)
	}

	if err := dbApi.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating database container: %s", err)
	}

	if err := authApi.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating auth0 server: %s", err)
	}
}

func TestApiRouterTestSuite(t *testing.T) {
	suite.Run(t, new(ApiRouterTestSuite))
}

func (s *ApiRouterTestSuite) TestReturnAllowOrigins() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/api/v1", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
}

func (s *ApiRouterTestSuite) TestCreateRoom_ShouldReturnUnauthorizedWhenUnauthenticated() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rooms", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestCreateRoom_ShouldCreateARoom() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	sub := authApi.GenerateSub()
	jwt, _ := authApi.GenerateJWT(sub)

	payload := struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}{
		"Need for Speed",
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
	assert.Nil(t, err)

	room, err := s.roomRepository.FindById(s.ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, sub, room.AdminId().Value())
	assert.Equal(t, payload.Name, room.Name().Value())
	assert.Equal(t, payload.Category, room.Category().Value())
}

func (s *ApiRouterTestSuite) TestSearchRoom_ShouldReturnUnauthorizedWhenUnauthenticated() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestSearchRoom_ShouldReturnRoomPages() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	sub := authApi.GenerateSub()
	jwt, _ := authApi.GenerateJWT(sub)

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
		page  pagination.Page[RoomResponseDto]
	}{
		{
			query: "?page=0&size=2&sort=asc&search=tech",
			page: pagination.Page[RoomResponseDto]{
				Page:  0,
				Size:  2,
				Total: 3,
				Items: []RoomResponseDto{
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
			page: pagination.Page[RoomResponseDto]{
				Page:  1,
				Size:  2,
				Total: 3,
				Items: []RoomResponseDto{
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
			page: pagination.Page[RoomResponseDto]{
				Page:  0,
				Size:  3,
				Total: 2,
				Items: []RoomResponseDto{
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

	for _, test := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms"+test.query, nil)
		req.Header.Set("Authorization", "Bearer "+jwt)

		r.ServeHTTP(w, req)
		res := w.Result()

		assert.Equal(t, http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		assert.Nil(t, err)

		var data pagination.Page[RoomResponseDto]
		err = json.Unmarshal(body, &data)
		assert.Nil(t, err)

		assert.Equal(t, test.page.Page, data.Page)
		assert.Equal(t, test.page.Size, data.Size)
		assert.Equal(t, test.page.Total, data.Total)
		assert.Equal(t, len(test.page.Items), len(data.Items))

		for i := 0; i < len(test.page.Items); i++ {
			testRoom := test.page.Items[i]
			dataRoom := data.Items[i]

			assert.Equal(t, testRoom.Id, dataRoom.Id)
			assert.Equal(t, testRoom.Name, dataRoom.Name)
			assert.Equal(t, testRoom.Category, dataRoom.Category)
		}
	}
}

func (s *ApiRouterTestSuite) TestFindRoom_ShouldReturnUnauthorizedWhenUnauthenticated() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router

	id := valueobject.NewId().Value()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms/"+id, nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestFindRoom_ShouldReturnNotFoundWhenRoomIdDoesNotExist() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	jwt, _ := authApi.GenerateJWT(authApi.GenerateSub())

	testCases := []struct {
		id      string
		expCode int
	}{
		{
			id:      "dfsoifdsiuroewrdf",
			expCode: http.StatusNotFound,
		},
		{
			id:      valueobject.NewId().Value(),
			expCode: http.StatusNotFound,
		},
	}

	for _, test := range testCases {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms/"+test.id, nil)
		req.Header.Set("Authorization", "Bearer "+jwt)

		r.ServeHTTP(w, req)
		res := w.Result()

		assert.Equal(t, test.expCode, res.StatusCode)
	}
}

func (s *ApiRouterTestSuite) TestFindRoom_ShouldReturnARoomWhenIdExists() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	jwt, _ := authApi.GenerateJWT(authApi.GenerateSub())

	adminId, _ := valueobject.NewUserIdWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room, _ := entity.NewRoom(adminId, name, category)
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

	var data RoomResponseDto
	err = json.Unmarshal(body, &data)
	assert.Nil(t, err)

	assert.Equal(t, room.Id().Value(), data.Id)
	assert.Equal(t, room.Name().Value(), data.Name)
	assert.Equal(t, room.Category().Value(), data.Category)
}

func (s *ApiRouterTestSuite) TestUpdateRoom_ShouldReturnUnauthorizedWhenUnauthenticated() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router

	id := valueobject.NewId().Value()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/rooms/"+id, nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestUpdateRoom_ShouldReturnNotFoundWhenRoomIdDoesNotExist() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	sub := authApi.GenerateSub()
	jwt, _ := authApi.GenerateJWT(sub)

	id := valueobject.NewId().Value()

	payload := struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}{
		"Need for Speed",
		"Game",
	}

	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/rooms/"+id, bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestUpdateRoom_ShouldReturnForbiddenWhenAdminIdIsInvalid() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	sub := authApi.GenerateSub()
	jwt, _ := authApi.GenerateJWT(sub)

	adminId, _ := valueobject.NewUserIdWith(authApi.GenerateSub())
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := entity.NewRoom(adminId, name, category)
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

	assert.Equal(t, http.StatusForbidden, res.StatusCode)

	savedRoom, err := s.roomRepository.FindById(s.ctx, room.Id())
	assert.NotNil(t, savedRoom)
	assert.Nil(t, err)
	assert.Equal(t, room.AdminId().Value(), savedRoom.AdminId().Value())
	assert.Equal(t, room.Name().Value(), savedRoom.Name().Value())
	assert.Equal(t, room.Category().Value(), savedRoom.Category().Value())
	assert.True(t, room.CreatedAt().Value().Equal(savedRoom.CreatedAt().Value()))
	assert.True(t, room.UpdatedAt().Value().Equal(savedRoom.UpdatedAt().Value()))
}

func (s *ApiRouterTestSuite) TestUpdateRoom_ShouldUpdateARoomWhenIdExists() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	sub := authApi.GenerateSub()
	jwt, _ := authApi.GenerateJWT(sub)

	adminId, _ := valueobject.NewUserIdWith(sub)
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := entity.NewRoom(adminId, name, category)
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
	assert.Nil(t, err)
	assert.Equal(t, room.AdminId().Value(), savedRoom.AdminId().Value())
	assert.Equal(t, payload.Name, savedRoom.Name().Value())
	assert.Equal(t, payload.Category, savedRoom.Category().Value())
	assert.True(t, room.CreatedAt().Value().Equal(savedRoom.CreatedAt().Value()))
	assert.True(t, room.UpdatedAt().Value().Before(savedRoom.UpdatedAt().Value()))
}

func (s *ApiRouterTestSuite) TestDeleteRoom_ShouldReturnUnauthorizedWhenUnauthenticated() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router

	id := valueobject.NewId().Value()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/rooms/"+id, nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestDeleteRoom_ShouldReturnNotFoundWhenRoomIdDoesNotExist() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	sub := authApi.GenerateSub()
	jwt, _ := authApi.GenerateJWT(sub)

	id := valueobject.NewId().Value()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/rooms/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestDeleteRoom_ShouldReturnForbiddenWhenAdminIdIsInvalid() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	sub := authApi.GenerateSub()
	jwt, _ := authApi.GenerateJWT(sub)

	adminId, _ := valueobject.NewUserIdWith(authApi.GenerateSub())
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := entity.NewRoom(adminId, name, category)
	s.roomRepository.Save(s.ctx, room)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/rooms/"+room.Id().Value(), nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusForbidden, res.StatusCode)

	savedRoom, err := s.roomRepository.FindById(s.ctx, room.Id())
	assert.NotNil(t, savedRoom)
	assert.Nil(t, err)
}

func (s *ApiRouterTestSuite) TestDeleteRoom_ShouldDeleteARoomWhenIdExists() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	sub := authApi.GenerateSub()
	jwt, _ := authApi.GenerateJWT(sub)

	adminId, _ := valueobject.NewUserIdWith(sub)
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := entity.NewRoom(adminId, name, category)
	s.roomRepository.Save(s.ctx, room)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/rooms/"+room.Id().Value(), nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	savedRoom, err := s.roomRepository.FindById(s.ctx, room.Id())
	assert.Nil(t, savedRoom)
	assert.NotNil(t, err)
}

func (s *ApiRouterTestSuite) TestCreateMessage_ShouldReturnUnauthorizedWhenUnauthenticated() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/messages", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestCreateMessage_ShouldCreateAMessage() {
	defer dbApi.Clear()
	t := s.T()
	r := s.router
	sub := authApi.GenerateSub()
	user := authApi.GetNickname()
	jwt, _ := authApi.GenerateJWT(sub)

	adminId, _ := valueobject.NewUserIdWith(sub)
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := entity.NewRoom(adminId, name, category)
	s.roomRepository.Save(s.ctx, room)

	payload := struct {
		RoomId string `json:"room_id"`
		Text   string `json:"text"`
	}{
		room.Id().Value(),
		"A simple text",
	}

	body, _ := json.Marshal(payload)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/messages", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusCreated, res.StatusCode)

	location := res.Header.Get("Location")
	assert.NotEmpty(t, location)

	id, err := valueobject.NewIdWith(location[len("/api/v1/messages/"):])
	assert.Nil(t, err)

	message, err := s.messageRepository.FindById(s.ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, payload.RoomId, message.RoomId().Value())
	assert.Equal(t, sub, message.SenderId().Value())
	assert.Equal(t, user, message.SenderName().Value())
	assert.Equal(t, payload.Text, message.Text().Value())

	msgs, err := s.messageRecevierGateway.Receive()
	assert.Nil(t, err)

	var data event.MessageEvent

	select {
	case msg := <-msgs:
		err := json.Unmarshal(msg.Body, &data)
		if err != nil {
			t.Fatal(err)
		}
		err = msg.Ack(true)
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second * 5):
		t.FailNow()
	}

	assert.Equal(t, message.Id().Value(), data.Id)
	assert.Equal(t, message.RoomId().Value(), data.RoomId)
	assert.Equal(t, message.SenderId().Value(), data.SenderId)
	assert.Equal(t, message.SenderName().Value(), data.SenderName)
	assert.Equal(t, message.Text().Value(), data.Text)
	assert.Equal(t, message.CreatedAt().String(), data.CreatedAt)
}

func createARoom(adminId, name, category string) *entity.Room {
	roomAdminId, _ := valueobject.NewUserIdWith(adminId)
	roomName, _ := valueobject.NewRoomNameWith(name)
	roomCategory, _ := valueobject.NewRoomCategoryWith(category)
	room, _ := entity.NewRoom(roomAdminId, roomName, roomCategory)
	return room
}
