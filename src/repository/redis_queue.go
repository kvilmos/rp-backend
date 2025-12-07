package repository

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Queue interface {
	Enqueue(ctx context.Context, queueName string, value any) error
}

type queue struct {
	redisClient *redis.Client
}

func NewQueue(client *redis.Client) Queue {
	return &queue{
		redisClient: client,
	}
}

func (r *queue) Enqueue(ctx context.Context, queueName string, value any) error {
	message, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.redisClient.LPush(ctx, queueName, message).Err()
}
