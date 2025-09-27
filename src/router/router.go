package router

import (
	"room-planner/app"
	"room-planner/handler"

	"github.com/labstack/echo/v4"
)

type server struct {
	echo *echo.Echo
}

func New(app *app.Application) *server {
	e := echo.New()
	server := server{echo: e}
	initRoutes(server.echo, app)

	return &server
}

func initRoutes(e *echo.Echo, app *app.Application) {
	handler := handler.New(app)

	e.POST("minio/upload-hook", handler.UploadFurnitureFile)
	e.POST("/furniture", handler.CreateFurniture)
}

func (s server) Start(port string) error {
	err := s.echo.Start(port)
	if err != nil {
		return err
	}

	return nil
}
