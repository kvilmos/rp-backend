package dto

import "room-planner/model"

type FurnitureCategoryDto struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func FromFurnitureCategoryModel(category *model.FurnitureCategory) *FurnitureCategoryDto {
	return &FurnitureCategoryDto{
		Id:   category.Id,
		Name: category.Name,
	}
}

func FromFurnitureCategoryModels(categories []*model.FurnitureCategory) []*FurnitureCategoryDto {
	var categoriesDto []*FurnitureCategoryDto
	for _, category := range categories {
		dto := FromFurnitureCategoryModel(category)
		categoriesDto = append(categoriesDto, dto)
	}

	return categoriesDto
}
