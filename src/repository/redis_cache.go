package repository

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, value any) error
	Delete(ctx context.Context, key string) error
}

type cacheRepository struct {
	redisClient *redis.Client
}

func NewCacheRepository(client *redis.Client) Cache {
	return &cacheRepository{
		redisClient: client,
	}
}

func (r *cacheRepository) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	valueByte, err := json.Marshal(value)
	if err != nil {

		return err
	}

	return r.redisClient.Set(ctx, key, valueByte, ttl).Err()
}

func (r *cacheRepository) Get(ctx context.Context, key string, value any) error {
	if reflect.ValueOf(value).Kind() != reflect.Ptr {
		return errors.New("not a pointer")
	}

	result, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(result), &value)
	if err != nil {
		return err
	}

	return nil
}

func (r *cacheRepository) Delete(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
