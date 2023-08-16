package web

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/domain/search"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/infra/database"
	"github.com/sesaquecruz/go-chat-api/internal/usecase"
	"github.com/sesaquecruz/go-chat-api/test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	_ "github.com/lib/pq"
)

type ApiRouterTestSuite struct {
	suite.Suite
	ctx         context.Context
	postgres    *test.PostgresContainer
	auth0       *test.Auth0Server
	roomGateway gateway.RoomGatewayInterface
	router      *gin.Engine
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

	roomGateway := database.NewRoomGateway(db)
	createRoomUseCase := usecase.NewCreateRoomUseCase(roomGateway)
	findRoomUseCase := usecase.NewFindRoomUseCase(roomGateway)
	searchRoomUseCase := usecase.NewSearchRoomUseCase(roomGateway)
	updateRoomUsecase := usecase.NewUpdateRoomUseCase(roomGateway)
	deleteRoomUseCase := usecase.NewDeleteRoomUseCase(roomGateway)
	roomHandler := NewRoomHandler(
		createRoomUseCase,
		searchRoomUseCase,
		findRoomUseCase,
		updateRoomUsecase,
		deleteRoomUseCase,
	)

	apiRouter := ApiRouter(&cfg, roomHandler)

	s.ctx = ctx
	s.postgres = postgres
	s.auth0 = auth0
	s.roomGateway = roomGateway
	s.router = apiRouter
}

func (s *ApiRouterTestSuite) TearDownSuite() {
	if err := s.postgres.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating postgres container: %s", err)
	}

	if err := s.auth0.Stop(s.ctx); err != nil {
		s.T().Fatalf("error terminating auth0 server: %s", err)
	}
}

func TestApiRouterTestSuite(t *testing.T) {
	suite.Run(t, new(ApiRouterTestSuite))
}

func (s *ApiRouterTestSuite) TestReturnAllowOrigins() {
	defer s.postgres.ClearDB()
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
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rooms", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestCreateRoom_ShouldCreateANewRoom() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	sub := s.auth0.GenerateSub()
	jwt, _ := s.auth0.GenerateJWT(sub)

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

	id, err := valueobject.NewIDWith(location[len("/api/v1/rooms/"):])
	assert.Nil(t, err)

	room, err := s.roomGateway.FindById(s.ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, sub, room.AdminId().Value())
	assert.Equal(t, payload.Name, room.Name().Value())
	assert.Equal(t, payload.Category, room.Category().Value())
}

func (s *ApiRouterTestSuite) TestSearchRoom_ShouldReturnUnauthorizedWhenUnauthenticated() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestSearchRoom_ShouldReturnRoomPages() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	sub := s.auth0.GenerateSub()
	jwt, _ := s.auth0.GenerateJWT(sub)

	room1 := createARoom(sub, "Rust", "Tech")
	room2 := createARoom(sub, "Go", "Tech")
	room3 := createARoom(sub, "Java", "Tech")
	room4 := createARoom(sub, "Need for Speed Undergroud", "Game")
	room5 := createARoom(sub, "Need for Speed Most Wanted", "Game")
	s.roomGateway.Save(s.ctx, room1)
	s.roomGateway.Save(s.ctx, room2)
	s.roomGateway.Save(s.ctx, room3)
	s.roomGateway.Save(s.ctx, room4)
	s.roomGateway.Save(s.ctx, room5)

	testCases := []struct {
		query string
		page  search.Page[RoomResponseDto]
	}{
		{
			query: "?page=0&size=2&sort=asc&search=tech",
			page: search.Page[RoomResponseDto]{
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
			page: search.Page[RoomResponseDto]{
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
			page: search.Page[RoomResponseDto]{
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

		var data search.Page[RoomResponseDto]
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
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router

	id := valueobject.NewID().Value()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms/"+id, nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestFindRoom_ShouldReturnNotFoundWhenRoomIdDoesNotExist() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	jwt, _ := s.auth0.GenerateJWT(s.auth0.GenerateSub())

	testCases := []struct {
		id      string
		expCode int
	}{
		{
			id:      "dfsoifdsiuroewrdf",
			expCode: http.StatusNotFound,
		},
		{
			id:      valueobject.NewID().Value(),
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
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	jwt, _ := s.auth0.GenerateJWT(s.auth0.GenerateSub())

	adminId, _ := valueobject.NewAuth0IDWith("auth0|64c8457bb160e37c8c34533b")
	name, _ := valueobject.NewRoomNameWith("Need for Speed")
	category, _ := valueobject.NewRoomCategoryWith("Game")
	room, _ := entity.NewRoom(adminId, name, category)
	s.roomGateway.Save(s.ctx, room)

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
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router

	id := valueobject.NewID().Value()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPut, "/api/v1/rooms/"+id, nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestUpdateRoom_ShouldReturnNotFoundWhenRoomIdDoesNotExist() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	sub := s.auth0.GenerateSub()
	jwt, _ := s.auth0.GenerateJWT(sub)

	id := valueobject.NewID().Value()

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
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	sub := s.auth0.GenerateSub()
	jwt, _ := s.auth0.GenerateJWT(sub)

	adminId, _ := valueobject.NewAuth0IDWith(s.auth0.GenerateSub())
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := entity.NewRoom(adminId, name, category)
	s.roomGateway.Save(s.ctx, room)

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

	savedRoom, err := s.roomGateway.FindById(s.ctx, room.Id())
	assert.NotNil(t, savedRoom)
	assert.Nil(t, err)
	assert.Equal(t, room.AdminId().Value(), savedRoom.AdminId().Value())
	assert.Equal(t, room.Name().Value(), savedRoom.Name().Value())
	assert.Equal(t, room.Category().Value(), savedRoom.Category().Value())
	assert.True(t, room.CreatedAt().Time().Equal(savedRoom.CreatedAt().Time()))
	assert.True(t, room.UpdatedAt().Time().Equal(savedRoom.UpdatedAt().Time()))
}

func (s *ApiRouterTestSuite) TestUpdateRoom_ShouldUpdateARoomWhenIdExists() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	sub := s.auth0.GenerateSub()
	jwt, _ := s.auth0.GenerateJWT(sub)

	adminId, _ := valueobject.NewAuth0IDWith(sub)
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := entity.NewRoom(adminId, name, category)
	s.roomGateway.Save(s.ctx, room)

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

	savedRoom, err := s.roomGateway.FindById(s.ctx, room.Id())
	assert.Nil(t, err)
	assert.Equal(t, room.AdminId().Value(), savedRoom.AdminId().Value())
	assert.Equal(t, payload.Name, savedRoom.Name().Value())
	assert.Equal(t, payload.Category, savedRoom.Category().Value())
	assert.True(t, room.CreatedAt().Time().Equal(savedRoom.CreatedAt().Time()))
	assert.True(t, room.UpdatedAt().Time().Before(savedRoom.UpdatedAt().Time()))
}

func (s *ApiRouterTestSuite) TestDeleteRoom_ShouldReturnUnauthorizedWhenUnauthenticated() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router

	id := valueobject.NewID().Value()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/rooms/"+id, nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestDeleteRoom_ShouldReturnNotFoundWhenRoomIdDoesNotExist() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	sub := s.auth0.GenerateSub()
	jwt, _ := s.auth0.GenerateJWT(sub)

	id := valueobject.NewID().Value()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/rooms/"+id, nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
}

func (s *ApiRouterTestSuite) TestDeleteRoom_ShouldReturnForbiddenWhenAdminIdIsInvalid() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	sub := s.auth0.GenerateSub()
	jwt, _ := s.auth0.GenerateJWT(sub)

	adminId, _ := valueobject.NewAuth0IDWith(s.auth0.GenerateSub())
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := entity.NewRoom(adminId, name, category)
	s.roomGateway.Save(s.ctx, room)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/rooms/"+room.Id().Value(), nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusForbidden, res.StatusCode)

	savedRoom, err := s.roomGateway.FindById(s.ctx, room.Id())
	assert.NotNil(t, savedRoom)
	assert.Nil(t, err)
}

func (s *ApiRouterTestSuite) TestDeleteRoom_ShouldDeleteARoomWhenIdExists() {
	defer s.postgres.ClearDB()
	t := s.T()
	r := s.router
	sub := s.auth0.GenerateSub()
	jwt, _ := s.auth0.GenerateJWT(sub)

	adminId, _ := valueobject.NewAuth0IDWith(sub)
	name, _ := valueobject.NewRoomNameWith("Rust")
	category, _ := valueobject.NewRoomCategoryWith("Tech")
	room, _ := entity.NewRoom(adminId, name, category)
	s.roomGateway.Save(s.ctx, room)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/api/v1/rooms/"+room.Id().Value(), nil)
	req.Header.Set("Authorization", "Bearer "+jwt)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNoContent, res.StatusCode)

	savedRoom, err := s.roomGateway.FindById(s.ctx, room.Id())
	assert.Nil(t, savedRoom)
	assert.NotNil(t, err)
}

func createARoom(adminId, name, category string) *entity.Room {
	roomAdminId, _ := valueobject.NewAuth0IDWith(adminId)
	roomName, _ := valueobject.NewRoomNameWith(name)
	roomCategory, _ := valueobject.NewRoomCategoryWith(category)
	room, _ := entity.NewRoom(roomAdminId, roomName, roomCategory)
	return room
}
