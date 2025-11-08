package handler

import (
	"room-planner/service"
)

type Handler struct {
	UserService         *service.UserService
	FurnitureService    *service.FurnitureService
	FurnitureCatService *service.FurnitureCategoryService
	BpService           *service.BlueprintService
}

func NewHandler(us *service.UserService, fs *service.FurnitureService, fcs *service.FurnitureCategoryService, bps *service.BlueprintService) *Handler {
	return &Handler{
		UserService:         us,
		FurnitureService:    fs,
		FurnitureCatService: fcs,
		BpService:           bps,
	}
}
