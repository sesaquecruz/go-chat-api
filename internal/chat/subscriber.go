package chat

import "context"

type subscriber struct {
	ctx    context.Context
	roomId string
	ch     chan []byte
}
