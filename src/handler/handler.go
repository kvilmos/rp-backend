package handler

import "room-planner/service"

type Handler struct {
	UserService      *service.UserService
	FurnitureService *service.FurnitureService
}

func NewHandler(us *service.UserService, fs *service.FurnitureService) *Handler {
	return &Handler{
		UserService:      us,
		FurnitureService: fs,
	}
}
