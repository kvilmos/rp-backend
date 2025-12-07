package repository

import (
	"context"
	"room-planner/model"

	"gorm.io/gorm"
)

type WallRepository interface {
	WithTx(tx *gorm.DB) WallRepository
	CreateMultiple(ctx context.Context, wall []*model.Wall) error
	DeleteByBlueprintId(ctx context.Context, bpId int64) error
}

type wallRepository struct {
	db *gorm.DB
}

func NewWallRepository(db *gorm.DB) WallRepository {
	return &wallRepository{db: db}
}
func (r *wallRepository) WithTx(tx *gorm.DB) WallRepository {
	return NewWallRepository(tx)
}

func (r *wallRepository) CreateMultiple(ctx context.Context, wall []*model.Wall) error {
	err := r.db.WithContext(ctx).Create(&wall).Error
	return err
}

func (r *wallRepository) DeleteByBlueprintId(ctx context.Context, id int64) error {
	sql := `DELETE FROM wall_t
			WHERE start_corner_id IN 
				(SELECT id FROM corner_t WHERE blueprint_id = ?)`

	return r.db.WithContext(ctx).Exec(sql, id).Error
}
