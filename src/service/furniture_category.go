package service

import (
	"context"
	"room-planner/model"
	"room-planner/repository"
)

type FurnitureCategoryService struct {
	CategoryRepo repository.FurnitureCategoryRepository
}

func NewFurnitureCategoryService(fcr repository.FurnitureCategoryRepository) *FurnitureCategoryService {
	return &FurnitureCategoryService{
		CategoryRepo: fcr,
	}
}

func (s FurnitureCategoryService) GetCategories(ctx context.Context) ([]*model.FurnitureCategory, error) {
	return s.CategoryRepo.List(ctx)
}
