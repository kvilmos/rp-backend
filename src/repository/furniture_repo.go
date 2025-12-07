package repository

import (
	"context"
	"fmt"
	"room-planner/common/constant"
	"room-planner/model"
	"room-planner/request"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FurnitureRepository interface {
	Create(ctx context.Context, furniture *model.Furniture) error
	IsExistByFileId(ctx context.Context, fileId uuid.UUID) (bool, error)
	Page(ctx context.Context, filter request.FurnitureFilter) ([]*model.Furniture, int, error)
	GetById(ctx context.Context, id int64) (*model.Furniture, error)
	DeleteForUser(ctx context.Context, furnitureId int64, userId int64) error
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

func (r *furnitureRepository) IsExistByFileId(ctx context.Context, fileId uuid.UUID) (bool, error) {
	var isExist bool
	sql := `SELECT COUNT(0) 
			FROM furniture_t 
			WHERE file_name = ?`

	err := r.db.WithContext(ctx).Raw(sql, fileId).Scan(&isExist).Error
	return isExist, err
}

func (r *furnitureRepository) GetById(ctx context.Context, id int64) (*model.Furniture, error) {
	var furniture *model.Furniture
	sql := `SELECT * 
			FROM furniture_t 
			WHERE id = ?`
	err := r.db.WithContext(ctx).Raw(sql, id).Scan(&furniture).Error

	return furniture, err
}

func (r *furnitureRepository) Page(ctx context.Context, filter request.FurnitureFilter) ([]*model.Furniture, int, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(`SELECT *,
							COUNT(1) OVER() as total_rows
							FROM furniture_t`)

	var params []interface{}
	var conditions []string
	if filter.CreatorId != 0 {
		conditions = append(conditions, "user_id = ?")
		params = append(params, filter.CreatorId)
	}

	if filter.CategoryId != nil {
		conditions = append(conditions, "category_id = ?")
		params = append(params, filter.CategoryId)
	}

	if len(conditions) > 0 {
		sqlBuilder.WriteString(" WHERE ")
		sqlBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	sortRule, ok := constant.FurnitureOrders[constant.OrderOption(filter.Order)]
	if !ok {
		sortRule = constant.FurnitureOrders[constant.LATEST_CREATED]
	}
	sqlBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s ", sortRule.Column, sortRule.Direction))

	sqlBuilder.WriteString(" LIMIT ? OFFSET ? ")
	params = append(params, constant.PAGE_LIMIT)
	params = append(params, constant.PAGE_LIMIT*(filter.Page-1))

	sqlString := sqlBuilder.String()

	var result []*model.FurnitureWithTotal
	err := r.db.WithContext(ctx).Raw(sqlString, params...).Scan(&result).Error
	if err != nil {
		return nil, 0, err
	}
	if len(result) == 0 {
		return []*model.Furniture{}, 0, nil
	}
	furnitureList := make([]*model.Furniture, len(result))
	for i, item := range result {
		furnitureList[i] = &item.Furniture
	}
	totalRows := result[0].TotalRows

	return furnitureList, totalRows, nil
}

func (r *furnitureRepository) DeleteForUser(ctx context.Context, userId int64, furnitureId int64) error {
	sql := `DELETE FROM furniture_t
			WHERE id  = ?
			AND user_id = ?`

	result := r.db.WithContext(ctx).Exec(sql, furnitureId, userId)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
