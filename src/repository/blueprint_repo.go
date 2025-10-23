package repository

import (
	"context"
	"room-planner/model"

	"gorm.io/gorm"
)

type BlueprintRepository interface {
	Create(ctx context.Context, blueprint *model.Blueprint) error
	GetById(ctx context.Context, id int64) (*model.Blueprint, error)
	Update(ctx context.Context, id int64) error
}

type blueprintRepository struct {
	db *gorm.DB
}

func NewBlueprintRepository(db *gorm.DB) BlueprintRepository {
	return &blueprintRepository{db: db}
}

func (r *blueprintRepository) Create(ctx context.Context, blueprint *model.Blueprint) error {
	err := r.db.WithContext(ctx).Omit("id").Create(&blueprint).Error
	return err
}

func (r *blueprintRepository) GetById(ctx context.Context, id int64) (*model.Blueprint, error) {
	var bp *model.Blueprint
	sql := `SELECT * 
			FROM blueprint_t
			WHERE id = ?`
	err := r.db.Raw(sql, id).Scan(&bp).Error

	return bp, err
}

func (r *blueprintRepository) Update(ctx context.Context, id int64) error {
	var bp *model.Blueprint
	sql := `UPDATE blueprint_t 
			SET modified_at = NOW()
			WHERE id = ?`
	return r.db.Raw(sql, id).Scan(&bp).Error
}

func (r *blueprintRepository) List() error {
	panic("unimplemented")
}
