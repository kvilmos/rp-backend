package storage

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient() (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return redisClient, nil
}
