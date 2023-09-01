package chat

import "context"

const ErrChatClosed = ChatError("chat closed")

type Chat interface {
	Subscribe(ctx context.Context, roomId string) (<-chan []byte, error)
}
