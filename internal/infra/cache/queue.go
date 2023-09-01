package cache

import (
	"context"
	"errors"
)

var ErrEmptyQueue = errors.New("empty queue")

type Queue interface {
	Push(ctx context.Context, queueName string, data []byte) error
	Pop(ctx context.Context, queueName string) ([]byte, error)
}
