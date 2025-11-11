package main

import (
	"context"
	"log"
	"os"
	"room-planner/handler"
	"room-planner/middleware"
	"room-planner/observer"
	"room-planner/repository"
	"room-planner/router"
	"room-planner/service"
	"room-planner/storage"
	"room-planner/token"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

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

	jwtMaker := token.NewJWTMaker(os.Getenv("JWT_SECRET"))

	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	userService := service.NewUserService(db, userRepo, sessionRepo, jwtMaker)

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
	bpService := service.NewBlueprintService(db, bpRepo, cornerRepo, wallRepo, itemRepo)

	authMiddleware := middleware.NewAuthMiddleware(userService, jwtMaker)
	apiHandler := handler.NewHandler(userService, furnitureService, furnitureCategoryService, bpService)

	observer := observer.New(redisClient, *furnitureService)
	observer.Start(context.Background())

	server := echo.New()

	server.Use(middleware.ErrorMiddleware)

	router.SetupRoutes(server, apiHandler, authMiddleware)
	err = server.Start(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal(err)
	}
}
