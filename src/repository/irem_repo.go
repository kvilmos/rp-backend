package repository

import (
	"context"
	"room-planner/model"

	"gorm.io/gorm"
)

type ItemRepository interface {
	Create(ctx context.Context, item *model.Item) error
	CreateMultiple(ctx context.Context, items []*model.Item) error
	DeleteByBlueprintId(ctx context.Context, id int64) error
}

func (r *itemRepository) DeleteByBlueprintId(ctx context.Context, id int64) error {
	sql := `DELETE FROM item_t 
			WHERE blueprint_id = ?`

	return r.db.Exec(sql, id).Error
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(ctx context.Context, item *model.Item) error {
	err := r.db.WithContext(ctx).Create(&item).Error
	return err
}

func (r *itemRepository) CreateMultiple(ctx context.Context, items []*model.Item) error {
	err := r.db.WithContext(ctx).Create(&items).Error
	return err
}
