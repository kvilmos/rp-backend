package route

import (
	"room-planner/handler"
	"room-planner/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, handler *handler.Handler, authMiddleware *middleware.AuthMiddleware) {
	auth := e.Group("/auth")
	{
		auth.POST("/register", handler.HandlerRegisterUser)
		auth.POST("/login", handler.HandlerLoginUser)
		auth.POST("/token", handler.HandlerRenewToken)
		auth.POST("/logout", handler.HandleLogoutUser)
	}

	verified := e.Group("/verify")
	verified.Use(authMiddleware.Authenticate)
	{
		verified.GET("/me", handler.HandleVerifyUser)
	}

	e.POST("/furniture", handler.NewFurniture)
	e.POST("/minio/upload-hook", handler.HandleUploadNotification)

	protected := e.Group("/protected")
	protected.Use(authMiddleware.Authenticate)
	{
		protected.GET("", handler.VerifyUser)
	}

}
