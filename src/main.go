package main

import (
	"context"
	"log"
	"room-planner/app"
	"room-planner/observer"
	"room-planner/router"
	"room-planner/storage"
)

func main() {
	minioClient, err := storage.NewMinioClient()
	if err != nil {
		log.Fatal(err)
	}

	redisClient, err := storage.NewRedisClient()
	if err != nil {
		log.Fatal(err)
	}

	db, err := storage.NewMySQLConnection()
	if err != nil {
		log.Fatal(err)
	}

	application := app.New(db, minioClient, redisClient)

	observer := observer.NewObserver(application)
	observer.Start(context.Background())

	server := router.New(application)
	err = server.Start(":4747")
	if err != nil {
		log.Fatal(err)
	}
}
