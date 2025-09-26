package handler

import (
	"room-planner/app"

	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db          *gorm.DB
	minioClient *minio.Client
	redisClient *redis.Client
}

func New(app *app.Application) *Handler {
	return &Handler{
		db:          app.Db,
		minioClient: app.MinioClient,
		redisClient: app.RedisClient,
	}
}
