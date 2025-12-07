package repository

import (
	"context"
	"encoding/json"
	"room-planner/common/constant"
	"room-planner/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type FurnitureUploadStatusRepository interface {
	Get(ctx context.Context, key string) (*model.FurnitureUploadStatus, error)
	Set(ctx context.Context, key string, status *model.FurnitureUploadStatus, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}
type furnitureUploadStatusRepository struct {
	redisClient *redis.Client
}

func NewUploadStatusRepository(client *redis.Client) FurnitureUploadStatusRepository {
	return &furnitureUploadStatusRepository{redisClient: client}
}

func (r *furnitureUploadStatusRepository) Get(ctx context.Context, key string) (*model.FurnitureUploadStatus, error) {
	uploadStatus := model.FurnitureUploadStatus{}
	val, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(val), &uploadStatus)
	if err != nil {
		return nil, err
	}

	return &uploadStatus, nil
}

func (r *furnitureUploadStatusRepository) Set(ctx context.Context, key string, status *model.FurnitureUploadStatus, ttl time.Duration) error {
	statusJson, err := json.Marshal(&status)
	if err != nil {
		return err
	}

	return r.redisClient.Set(ctx, key, statusJson, constant.UPLOAD_TTL).Err()
}

func (r *furnitureUploadStatusRepository) Delete(ctx context.Context, key string) error {
	return r.redisClient.Del(ctx, key).Err()
}
