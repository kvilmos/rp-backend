package router

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

	furniture := e.Group("/furniture")
	furniture.Use(authMiddleware.Authenticate)
	{
		furniture.POST("", handler.HandleNewFurniture)
		furniture.GET("", handler.HandleGetFurnitureList)
		furniture.GET("/:id", handler.HandleGetFurnitureById)
		furniture.DELETE("/:id", handler.HandleDeleteFurniture)
	}

	category := e.Group("/furniture-category")
	category.Use(authMiddleware.Authenticate)
	{
		category.GET("", handler.HandleGetFurnitureCategory)
	}

	blueprint := e.Group("/blueprint")
	blueprint.Use(authMiddleware.Authenticate)
	{
		blueprint.POST("", handler.HandleCreateBlueprint)
		blueprint.PUT("/:id", handler.HandleSaveBlueprint)
		blueprint.DELETE("/:id", handler.HandlerDeleteBlueprint)
		blueprint.GET("/complete/:id", handler.HandleGetCompleteBlueprintById)
	}

	profile := e.Group("/profile")
	profile.Use(authMiddleware.Authenticate)
	{
		profile.GET("/furniture", handler.HandlerGetProfileFurniture)
		profile.GET("/blueprint", handler.HandlerGetProfileBlueprints)
	}

	minio := e.Group("/minio")
	{
		minio.POST("/upload-hook", handler.HandleUploadNotification)
	}

}
