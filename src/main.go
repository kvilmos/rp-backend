package main

import (
	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()
	server.Logger.Fatal(server.Start(":4747"))
}
