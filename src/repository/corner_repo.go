package repository

import (
	"context"
	"room-planner/model"

	"gorm.io/gorm"
)

type CornerRepository interface {
	WithTx(tx *gorm.DB) CornerRepository
	Create(ctx context.Context, corners *model.Corner) error
	CreateMultiple(ctx context.Context, corners []*model.Corner) error
	DeleteByBlueprintId(ctx context.Context, id int64) error
}

func (r *cornerRepository) DeleteByBlueprintId(ctx context.Context, id int64) error {
	sql := `DELETE FROM corner_t 
			WHERE blueprint_id = ?`

	return r.db.Exec(sql, id).Error
}

type cornerRepository struct {
	db *gorm.DB
}

func NewCornerRepository(db *gorm.DB) CornerRepository {
	return &cornerRepository{db: db}
}

func (r *cornerRepository) WithTx(tx *gorm.DB) CornerRepository {
	return NewCornerRepository(tx)
}

func (r *cornerRepository) Create(ctx context.Context, corners *model.Corner) error {
	err := r.db.WithContext(ctx).Create(&corners).Error
	return err
}

func (r *cornerRepository) CreateMultiple(ctx context.Context, corners []*model.Corner) error {
	err := r.db.WithContext(ctx).Create(&corners).Error
	return err
}
