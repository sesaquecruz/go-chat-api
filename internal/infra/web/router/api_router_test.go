package router

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/sesaquecruz/go-chat-api/internal/infra/web/handler"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

func (s *WebApiTestSuite) TestReturnAllowOrigins() {
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodOptions, "/api/v1", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
}

func (s *WebApiTestSuite) Test_CreateRoom_ShouldReturnUnauthorizedWhenUnauthenticated() {
	t := s.T()
	r := s.router

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/rooms", nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *WebApiTestSuite) Test_CreateRoom_ShouldCreateANewRoom() {
	s.postgres.ClearDB()

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

func (s *WebApiTestSuite) Test_FindRoom_ShouldReturnUnauthorizedWhenUnauthenticated() {
	t := s.T()
	r := s.router
	id := valueobject.NewID().Value()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/rooms/"+id, nil)

	r.ServeHTTP(w, req)
	res := w.Result()

	assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func (s *WebApiTestSuite) Test_FindRoom_ShouldReturnNotFoundWhenRoomIdDoesNotExist() {
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

func (s *WebApiTestSuite) Test_FindRoom_ShouldReturnARoomWhenIdExists() {
	s.postgres.ClearDB()

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

	var data handler.RoomResponse
	err = json.Unmarshal(body, &data)
	assert.Nil(t, err)

	assert.Equal(t, room.Id().Value(), data.Id)
	assert.Equal(t, room.Name().Value(), data.Name)
	assert.Equal(t, room.Category().Value(), data.Category)
}
