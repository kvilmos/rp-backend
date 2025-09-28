package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TODO STORY-201 Handler error
var ErrLockNotAcquired = errors.New("lock not acquired")

type Locker interface {
	Acquire(ctx context.Context, key string, ttl time.Duration) (string, error)
	Release(ctx context.Context, key string, value string) error
}

type locker struct {
	redisClient *redis.Client
}

func NewLocker(client *redis.Client) Locker {
	return &locker{
		redisClient: client,
	}
}

func (r *locker) Acquire(ctx context.Context, key string, ttl time.Duration) (string, error) {
	lockValue := uuid.New().String()

	ok, err := r.redisClient.SetNX(ctx, key, lockValue, ttl).Result()
	if err != nil {
		return "", err
	}

	if !ok {
		return "", ErrLockNotAcquired
	}

	return lockValue, nil
}

func (r *locker) Release(ctx context.Context, key string, value string) error {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end`

	_, err := r.redisClient.Eval(ctx, script, []string{key}, value).Result()
	return err
}
