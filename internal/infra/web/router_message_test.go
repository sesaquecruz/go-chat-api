package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/domain/event"
	"github.com/sesaquecruz/go-chat-api/internal/domain/valueobject"

	"github.com/stretchr/testify/assert"
)

func (s *RouterTestSuite) TestCreateMessage_ShouldCreateAMessage() {
	defer db.Clear()
	t := s.T()
	r := s.router

	user := auth.GetNickname()
	sub := auth.GenerateSub()
	jwt, _ := auth.GenerateJWT(sub)

	room := createARoom(sub, "A Game", "Game")
	s.roomRepository.Save(s.ctx, room)

	payload := struct {
		RoomId string `json:"room_id"`
		Text   string `json:"text"`
	}{
		room.Id().Value(),
		"A text",
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
	assert.NotNil(t, message)
	assert.Nil(t, err)
	assert.Equal(t, payload.RoomId, message.RoomId().Value())
	assert.Equal(t, sub, message.SenderId().Value())
	assert.Equal(t, user, message.SenderName().Value())
	assert.Equal(t, payload.Text, message.Text().Value())

	msgs := make(chan *event.MessageEvent)
	defer close(msgs)

	go func() {
		err = s.messageEventGateway.Receive(s.ctx, msgs)
		if err != nil {
			t.Error(err)
		}
	}()

	select {
	case msg := <-msgs:
		assert.Equal(t, message.Id().Value(), msg.Id)
		assert.Equal(t, message.RoomId().Value(), msg.RoomId)
		assert.Equal(t, message.SenderId().Value(), msg.SenderId)
		assert.Equal(t, message.SenderName().Value(), msg.SenderName)
		assert.Equal(t, message.Text().Value(), msg.Text)
		assert.Equal(t, message.CreatedAt().Value(), msg.CreatedAt)
	case <-time.After(10 * time.Second):
		t.Fail()
	}
}
