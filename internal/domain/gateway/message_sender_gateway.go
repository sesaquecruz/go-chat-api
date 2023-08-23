package gateway

import (
	"context"

	"github.com/sesaquecruz/go-chat-api/internal/domain/entity"
)

type MessageSenderGatewayInterface interface {
	Send(ctx context.Context, message *entity.Message) error
}
