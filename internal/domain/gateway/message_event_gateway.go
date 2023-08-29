package gateway

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/event"
)

type MessageEventGateway interface {
	Send(ctx context.Context, messageEvent *event.MessageEvent) error
	Receive(ctx context.Context, messageEvents chan<- *event.MessageEvent) error
}
