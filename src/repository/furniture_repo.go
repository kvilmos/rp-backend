package repository

import (
	"context"
	"room-planner/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FurnitureRepository interface {
	Create(ctx context.Context, furniture *model.Furniture) error
	IsExistByFileId(fileId uuid.UUID) (bool, error)
}

type furnitureRepository struct {
	db *gorm.DB
}

func NewFurnitureRepository(db *gorm.DB) FurnitureRepository {
	return &furnitureRepository{db: db}
}

func (r *furnitureRepository) Create(ctx context.Context, furniture *model.Furniture) error {
	err := r.db.WithContext(ctx).Create(&furniture).Error
	return err
}

func (r *furnitureRepository) IsExistByFileId(fileId uuid.UUID) (bool, error) {
	var isExist bool
	sql := `SELECT COUNT(0) 
			FROM furniture_t 
			WHERE file_name = ?`

	err := r.db.Raw(sql, fileId).Scan(&isExist).Error
	return isExist, err
}
