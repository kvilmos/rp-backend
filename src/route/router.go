package route

import (
	"room-planner/handler"
	"room-planner/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, handler *handler.Handler, auth *middleware.AuthMiddleware) {
	e.POST("/signup", handler.Register)
	e.POST("/login", handler.Login)

	e.POST("/furniture", handler.NewFurniture)
	e.POST("/minio/upload-hook", handler.HandleUploadNotification)

	authRequired := e.Group("/v1")
	authRequired.Use(auth.Authenticate)
	{
		authRequired.GET("/authrequired", handler.TestAuth)
	}
}
