package chat

import (
	"context"
	"encoding/json"

	"github.com/sesaquecruz/go-chat-api/internal/domain/event"
	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/infra/cache"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

const listenerChannelSize = 10

type listener struct {
	ctx     context.Context
	gateway gateway.MessageEventGateway
	cache   cache.Queue
	logger  *log.Logger
}

func newListener(
	ctx context.Context,
	gateway gateway.MessageEventGateway,
	cache cache.Queue,
) *listener {
	listener := &listener{
		ctx:     ctx,
		gateway: gateway,
		cache:   cache,
		logger:  log.NewLogger("Listener"),
	}

	go listener.start()

	return listener
}

func (w *listener) start() {
	w.logger.Info("listener started")

	msgs := make(chan *event.MessageEvent, listenerChannelSize)
	defer close(msgs)

	go func() {
		for msg := range msgs {
			queueName := msg.RoomId

			data, err := json.Marshal(&message{User: msg.SenderName, Text: msg.Text, Time: msg.CreatedAt})
			if err != nil {
				w.logger.Error(err)
				continue
			}

			err = w.cache.Push(w.ctx, queueName, data)
			if err != nil {
				w.logger.Error(err)
			}
		}
	}()

	err := w.gateway.Receive(w.ctx, msgs)
	if err != nil {
		w.logger.Error(err)
	}

	w.logger.Info("listener finished")
}
