package main

import (
	"context"
	"log"
	"room-planner/handler"
	"room-planner/middleware"
	"room-planner/observer"
	"room-planner/repository"
	"room-planner/route"
	"room-planner/service"
	"room-planner/storage"

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

	APIHandler     *handler.Handler
	AuthMiddleware *middleware.AuthMiddleware
}

func main() {
	db, err := storage.NewMySQLConnection()
	if err != nil {
		log.Fatal(err)
	}
	minioClient, err := storage.NewMinioClient()
	if err != nil {
		log.Fatal(err)
	}
	redisClient, err := storage.NewRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo, db)

	furnitureRepo := repository.NewFurnitureRepository(db)
	fileRepo := repository.NewFileRepository(minioClient)
	cacheRepo := repository.NewCacheRepository(redisClient)
	locker := repository.NewLocker(redisClient)
	queue := repository.NewQueue(redisClient)
	furnitureService := service.NewFurnitureService(furnitureRepo, fileRepo, cacheRepo, queue, locker)

	authMiddleware := middleware.NewAuthMiddleware(userService)
	apiHandler := handler.NewHandler(userService, furnitureService)

	observer := observer.New(redisClient, *furnitureService)
	observer.Start(context.Background())

	server := echo.New()
	route.SetupRoutes(server, apiHandler, authMiddleware)

	err = server.Start(":4747")
	if err != nil {
		log.Fatal(err)
	}
}
