package chat

import (
	"context"
	"errors"
	"time"

	"github.com/sesaquecruz/go-chat-api/internal/infra/cache"
	"github.com/sesaquecruz/go-chat-api/pkg/log"
)

const broadcasterWaitListChannelSize = 10
const broadcasterEmptyQueueWaitTimeMilli = 100

type broadcaster struct {
	ctx         context.Context
	roomId      string
	cache       cache.Queue
	subscribers map[*subscriber]struct{}
	waitList    chan *subscriber
	removeList  chan *subscriber
	logger      *log.Logger
}

func newBroadcaster(
	ctx context.Context,
	roomId string,
	cache cache.Queue,
) *broadcaster {
	broadcaster := &broadcaster{
		ctx:         ctx,
		roomId:      roomId,
		cache:       cache,
		subscribers: make(map[*subscriber]struct{}),
		waitList:    make(chan *subscriber, broadcasterWaitListChannelSize),
		logger:      log.NewLogger("Broadcaster"),
	}

	go broadcaster.start()

	return broadcaster
}

func (b *broadcaster) start() {
	b.logger.Infof("broadcast started to room %s\n", b.roomId)

	for {
		select {
		case <-b.ctx.Done():
			for sub := range b.subscribers {
				close(sub.ch)
			}

			for sub := range b.waitList {
				close(sub.ch)
			}

			b.subscribers = make(map[*subscriber]struct{})

			b.logger.Infof("broadcast finished of room %s\n", b.roomId)
			return

		case sub := <-b.waitList:
			if _, ok := b.subscribers[sub]; !ok {
				b.subscribers[sub] = struct{}{}
			} else {
				b.logger.Warning("can not register duplicated subscriber")
			}

		case sub := <-b.removeList:
			b.logger.Info("removing sub")
			delete(b.subscribers, sub)

		default:
			msg, err := b.cache.Pop(b.ctx, b.roomId)
			if err != nil {
				if errors.Is(err, cache.ErrEmptyQueue) {
					<-time.After(broadcasterEmptyQueueWaitTimeMilli * time.Millisecond)
				} else {
					b.logger.Error(err)
				}

				continue
			}

			removeSubs := make([]*subscriber, 0)

			for sub := range b.subscribers {
				select {
				case <-sub.ctx.Done():
					removeSubs = append(removeSubs, sub)
				default:
					sub.ch <- msg
				}
			}

			for _, sub := range removeSubs {
				delete(b.subscribers, sub)
				close(sub.ch)
			}
		}
	}
}
