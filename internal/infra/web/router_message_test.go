package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/domain/event"

	"github.com/stretchr/testify/assert"
)

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

	msgs := make(chan *event.MessageEvent)
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
	case <-time.After(10 * time.Second):
		t.Fail()
	}
}
