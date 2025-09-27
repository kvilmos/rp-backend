package observer

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var ErrLockNotAcquired = errors.New("lock not acquired")

// TODO STORY-201 Handler error

func acquireLock(ctx context.Context, client *redis.Client, lockKey string, ttl time.Duration) (string, error) {
	lockValue := uuid.New().String()

	ok, err := client.SetNX(ctx, lockKey, lockValue, ttl).Result()
	if err != nil {
		return "", err
	}

	if !ok {
		return "", ErrLockNotAcquired
	}

	return lockValue, nil
}

func releaseLock(ctx context.Context, client *redis.Client, lockKey string, lockValue string) error {
	script := `
	if redis.call("get", KEYS[1]) == ARGV[1] then
		return redis.call("del", KEYS[1])
	else
		return 0
	end`

	_, err := client.Eval(ctx, script, []string{lockKey}, lockValue).Result()
	return err
}
