package repository

import (
	"context"
	"room-planner/model"

	"gorm.io/gorm"
)

type FurnitureCategoryRepository interface {
	List(ctx context.Context) ([]*model.FurnitureCategory, error)
}

type furnitureCategoryRepository struct {
	db *gorm.DB
}

func NewFurnitureCategoryRepository(db *gorm.DB) FurnitureCategoryRepository {
	return &furnitureCategoryRepository{db: db}
}

func (r *furnitureCategoryRepository) List(ctx context.Context) ([]*model.FurnitureCategory, error) {
	sql := `SELECT * 
			FROM category_t`

	var list []*model.FurnitureCategory
	err := r.db.WithContext(ctx).Raw(sql).Scan(&list).Error
	if err != nil {
		return nil, err
	}

	return list, nil
}
