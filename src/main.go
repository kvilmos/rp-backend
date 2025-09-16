package main

import (
	"room-planner/api"

	"github.com/labstack/echo/v4"
)

func main() {
	webServer := echo.New()
	api.InitRoutes(webServer)
	webServer.Logger.Fatal(webServer.Start(":4747"))
}
