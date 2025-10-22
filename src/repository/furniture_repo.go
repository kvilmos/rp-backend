package repository

import (
	"context"
	"room-planner/common/constant"
	"room-planner/model"
	"room-planner/request"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FurnitureRepository interface {
	Create(ctx context.Context, furniture *model.Furniture) error
	IsExistByFileId(fileId uuid.UUID) (bool, error)
	GetById(id int64) (*model.Furniture, error)
	List() ([]*model.Furniture, error)
	Page(page int, filter request.FurnitureFilter) ([]*model.Furniture, error)
	Count() (int, error)
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

func (r *furnitureRepository) List() ([]*model.Furniture, error) {
	var list []*model.Furniture
	sql := `SELECT * 
			FROM furniture_t `
	err := r.db.Raw(sql).Scan(&list).Error

	return list, err
}

func (r *furnitureRepository) GetById(id int64) (*model.Furniture, error) {
	var furniture *model.Furniture
	sql := `SELECT * 
			FROM furniture_t 
			WHERE id = ?`
	err := r.db.Raw(sql, id).Scan(&furniture).Error

	return furniture, err
}

func (r *furnitureRepository) Page(page int, filter request.FurnitureFilter) ([]*model.Furniture, error) {
	sql := `SELECT * 
			FROM furniture_t `

	switch filter.SortByDate {
	case "latest":
		sql += ` ORDER BY created_at DESC `
	case "oldest":
		sql += ` ORDER BY created_at ASC `
	default:
		sql += ` ORDER BY created_at DESC `
	}

	sql += `LIMIT ?
			OFFSET ?`

	var list []*model.Furniture
	err := r.db.Raw(sql, constant.PAGE_LIMIT, constant.PAGE_LIMIT*(page-1)).Scan(&list).Error

	return list, err
}

func (r *furnitureRepository) Count() (int, error) {
	var count int
	sql := `SELECT COUNT(1) 
			FROM furniture_t`
	err := r.db.Raw(sql).Scan(&count).Error
	return count, err
}
