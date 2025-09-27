package app

import (
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Db          *gorm.DB
	MinioClient *minio.Client
	RedisClient *redis.Client
}

func New(db *gorm.DB, minio *minio.Client, redis *redis.Client) *Application {
	return &Application{
		Db:          db,
		MinioClient: minio,
		RedisClient: redis,
	}
}
