package main

import (
	"log"
	minIOService "room-planner/pkg/minio"

	"github.com/labstack/echo/v4"
)

func main() {
	err := minIOService.InitMinIO()
	if err != nil {
		log.Fatal(err)
	}

	server := echo.New()
	server.Logger.Fatal(server.Start(":4747"))
}
