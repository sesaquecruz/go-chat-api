package web

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"

	_ "github.com/lib/pq"
)

const roomPath = "/rooms"

func (s *ApiRouterTestSuite) TestRoom_ShouldReturnUnauthorizedWhenBearerIsEmpty() {
	t := s.T()
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		req, err := http.NewRequest(method, s.apiUrl+roomPath, http.NoBody)
		assert.Nil(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.Nil(t, err)
		defer res.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
	}
}

func (s *ApiRouterTestSuite) TestRoom_ShouldCreateANewRoom() {
	t := s.T()

	payload := struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}{
		"Need for Speed",
		"Game",
	}

	body, err := json.Marshal(payload)
	assert.Nil(t, err)

	req, err := http.NewRequest(http.MethodPost, s.apiUrl+roomPath, bytes.NewReader(body))
	assert.Nil(t, err)

	sub := s.iam.GenerateSub()
	token, err := s.iam.GenerateJWT(sub)
	assert.Nil(t, err)

	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)

	location := res.Header.Get("Location")
	assert.NotEmpty(t, location)

	id, err := valueobject.NewIDWith(location[len(roomPath)+1:])
	assert.Nil(t, err)

	room, err := s.roomGateway.FindById(s.ctx, id)
	assert.Nil(t, err)
	assert.Equal(t, sub, room.AdminId().Value())
	assert.Equal(t, payload.Name, room.Name().Value())
	assert.Equal(t, payload.Category, room.Category().Value())
}
