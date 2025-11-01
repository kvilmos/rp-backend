package handler

import (
	"room-planner/service"
)

type Handler struct {
	UserService      *service.UserService
	FurnitureService *service.FurnitureService
	BpService        *service.BlueprintService
}

func NewHandler(us *service.UserService, fs *service.FurnitureService, bps *service.BlueprintService) *Handler {
	return &Handler{
		UserService:      us,
		FurnitureService: fs,
		BpService:        bps,
	}
}
