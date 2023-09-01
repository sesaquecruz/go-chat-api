package cache

import (
	"context"
	"errors"

	"github.com/sesaquecruz/go-chat-api/pkg/log"

	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	rdb    *redis.Client
	logger *log.Logger
}

func NewRedisQueue(rdb *redis.Client) *RedisQueue {
	return &RedisQueue{
		rdb:    rdb,
		logger: log.NewLogger("RedisQueue"),
	}
}

func (c *RedisQueue) Push(ctx context.Context, queueName string, data []byte) error {
	err := c.rdb.RPush(ctx, queueName, data).Err()
	if err != nil {
		c.logger.Error(err)
		return err
	}

	return nil
}

func (c *RedisQueue) Pop(ctx context.Context, queueName string) ([]byte, error) {
	msg, err := c.rdb.LPop(ctx, queueName).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrEmptyQueue
		}

		c.logger.Error(err)
		return nil, err
	}

	return msg, nil
}
