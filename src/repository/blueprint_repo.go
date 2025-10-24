package repository

import (
	"context"
	"room-planner/common/constant"
	"room-planner/model"

	"gorm.io/gorm"
)

type BlueprintRepository interface {
	Create(ctx context.Context, blueprint *model.Blueprint) error
	List(ctx context.Context) ([]*model.Blueprint, error)
	Page(ctx context.Context, page int) ([]*model.Blueprint, error)
	GetById(ctx context.Context, id int64) (*model.Blueprint, error)
	Update(ctx context.Context, id int64) error
	ListComplete(ctx context.Context) ([]*model.Blueprint, error)
	GetCompleteById(ctx context.Context, id int64) (*model.Blueprint, error)
	Count(ctx context.Context) (int, error)
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

func (r *blueprintRepository) Page(ctx context.Context, page int) ([]*model.Blueprint, error) {
	var blueprints []*model.Blueprint
	sql := `SELECT * 
			FROM blueprint_t
			LIMIT ?
			OFFSET ?`
	err := r.db.Raw(sql, constant.PAGE_LIMIT, constant.PAGE_LIMIT*(page-1)).Scan(&blueprints).Error

	return blueprints, err
}

func (r *blueprintRepository) GetById(ctx context.Context, id int64) (*model.Blueprint, error) {
	var blueprint *model.Blueprint
	sql := `SELECT * 
			FROM blueprint_t
			WHERE id = ?`
	err := r.db.Raw(sql, id).Scan(&blueprint).Error

	return blueprint, err
}

func (r *blueprintRepository) List(ctx context.Context) ([]*model.Blueprint, error) {
	var blueprints []*model.Blueprint
	sql := `SELECT * 
			FROM blueprint_t`
	err := r.db.Raw(sql).Scan(&blueprints).Error

	return blueprints, err
}

func (r *blueprintRepository) Update(ctx context.Context, id int64) error {
	var blueprint *model.Blueprint
	sql := `UPDATE blueprint_t 
			SET modified_at = NOW()
			WHERE id = ?`
	return r.db.Raw(sql, id).Scan(&blueprint).Error
}

func (r *blueprintRepository) ListComplete(ctx context.Context) ([]*model.Blueprint, error) {
	var blueprints []*model.Blueprint
	err := r.db.Preload("Corners").Preload("Items").Preload("Walls").Find(&blueprints).Error

	return blueprints, err
}

func (r *blueprintRepository) GetCompleteById(ctx context.Context, id int64) (*model.Blueprint, error) {
	var blueprint *model.Blueprint
	err := r.db.Where("id = ?", id).Preload("Corners").Preload("Items.Furniture").Preload("Walls").Find(&blueprint).Error

	return blueprint, err
}

func (r *blueprintRepository) Count(ctx context.Context) (int, error) {
	var count int
	sql := `SELECT COUNT(1) 
			FROM blueprint_t`
	err := r.db.Raw(sql).Scan(&count).Error

	return count, err
}
