package chat

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/gateway"
	"github.com/sesaquecruz/go-chat-api/internal/infra/cache"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

var chat *ChatCached

type ChatCached struct {
	ctx      context.Context
	cancel   context.CancelFunc
	listener *listener
	rooms    map[string]*broadcaster
	waitList chan *subscriber
	logger   *log.Logger
}

func NewChatCached(
	gateway gateway.MessageEventGateway,
	cache cache.Queue,
) *ChatCached {
	if chat == nil {
		ctx, cancel := context.WithCancel(context.Background())

		chat = &ChatCached{
			ctx:      ctx,
			cancel:   cancel,
			listener: newListener(ctx, gateway, cache),
			rooms:    make(map[string]*broadcaster),
			waitList: make(chan *subscriber),
			logger:   log.NewLogger("Chat"),
		}

		go chat.start()
	}

	return chat
}

func (c *ChatCached) start() {
	c.logger.Info("chat started")
	defer c.cancel()

	for sub := range c.waitList {
		if _, ok := c.rooms[sub.roomId]; !ok {
			c.rooms[sub.roomId] = newBroadcaster(c.ctx, sub.roomId, c.listener.cache)
		}

		c.rooms[sub.roomId].waitList <- sub
	}

	c.logger.Info("chat finished")
}

func (c *ChatCached) Subscribe(ctx context.Context, roomId string) (<-chan []byte, error) {
	select {
	case <-c.ctx.Done():
		return nil, ErrChatClosed
	default:
		sub := &subscriber{
			ctx:    ctx,
			roomId: roomId,
			ch:     make(chan []byte),
		}

		c.waitList <- sub

		return sub.ch, nil
	}
}
