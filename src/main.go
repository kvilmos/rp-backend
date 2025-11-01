package main

import (
	"context"
	"log"
	"room-planner/handler"
	"room-planner/middleware"
	"room-planner/observer"
	"room-planner/repository"
	"room-planner/router"
	"room-planner/service"
	"room-planner/storage"
	"room-planner/token"

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
	secretKey := "d4mqw2lvfivh7fcr32igzf5q12345678" // min 32
	// env

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

	jwtMaker := token.NewJWTMaker(secretKey)

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	userService := service.NewUserService(userRepo, sessionRepo, db, jwtMaker)

	furnitureRepo := repository.NewFurnitureRepository(db)
	fileRepo := repository.NewFileRepository(minioClient)
	cacheRepo := repository.NewCacheRepository(redisClient)
	locker := repository.NewLocker(redisClient)
	queue := repository.NewQueue(redisClient)
	furnitureService := service.NewFurnitureService(furnitureRepo, fileRepo, cacheRepo, queue, locker)

	furnitureCategoryRepo := repository.NewFurnitureCategoryRepository(db)
	furnitureCategoryService := service.NewFurnitureCategoryService(furnitureCategoryRepo)

	bpRepo := repository.NewBlueprintRepository(db)
	cornerRepo := repository.NewCornerRepository(db)
	wallRepo := repository.NewWallRepository(db)
	itemRepo := repository.NewItemRepository(db)
	bpService := service.NewBlueprintService(bpRepo, cornerRepo, wallRepo, itemRepo)

	authMiddleware := middleware.NewAuthMiddleware(userService, jwtMaker)
	apiHandler := handler.NewHandler(userService, furnitureService, furnitureCategoryService, bpService)

	observer := observer.New(redisClient, *furnitureService)
	observer.Start(context.Background())

	server := echo.New()

	server.Use(middleware.ErrorMiddleware)

	router.SetupRoutes(server, apiHandler, authMiddleware)
	err = server.Start(":4747")
	if err != nil {
		log.Fatal(err)
	}
}
