package cache

import (
	"context"
	"testing"

	"github.com/sesaquecruz/go-chat-api/config"
	"github.com/sesaquecruz/go-chat-api/test/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var redisQueue, _ = services.NewRedisContainer(context.Background())

type RedisQueueTestSuite struct {
	suite.Suite
	ctx        context.Context
	redisQueue Queue
}

func (s *RedisQueueTestSuite) SetupSuite() {
	rdb := RedisConnection(&config.CacheConfig{
		Host: redisQueue.Host,
		Port: redisQueue.Port,
	})

	s.ctx = context.Background()
	s.redisQueue = NewRedisQueue(rdb)
}

func (s *RedisQueueTestSuite) TearDownSuite() {
	if err := redisQueue.Terminate(s.ctx); err != nil {
		s.T().Fatalf("error terminating redis container: %s", err)
	}
}

func TestRedisQueueTestSuite(t *testing.T) {
	suite.Run(t, new(RedisQueueTestSuite))
}

func (s *RedisQueueTestSuite) TestShouldAddAndRemoveData() {
	t := s.T()

	data := []struct {
		queueName string
		value     string
	}{
		{
			"queue 1",
			"data 1",
		},
		{
			"queue 1",
			"data 2",
		},
		{
			"queue 2",
			"data 3",
		},
		{
			"queue 2",
			"data 4",
		},
	}

	for _, d := range data {
		err := s.redisQueue.Push(s.ctx, d.queueName, []byte(d.value))
		assert.Nil(t, err)
	}

	for _, d := range data {
		value, err := s.redisQueue.Pop(s.ctx, d.queueName)
		assert.NotNil(t, value)
		assert.Nil(t, err)
		assert.Equal(t, d.value, string(value))
	}

	value, err := s.redisQueue.Pop(s.ctx, "queue 1")
	assert.Nil(t, value)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrEmptyQueue)

	value, err = s.redisQueue.Pop(s.ctx, "queue 2")
	assert.Nil(t, value)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrEmptyQueue)
}
