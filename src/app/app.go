package app

import (
	"github.com/labstack/echo/v4"
	"github.com/minio/minio-go/v7"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Server *echo.Echo
	Db     *gorm.DB
	MinIO  *minio.Client
	Redis  *redis.Client
}
