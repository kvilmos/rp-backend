package repository

import (
	"context"
	"fmt"
	"room-planner/common/constant"
	"room-planner/model"
	"room-planner/request"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BlueprintRepository interface {
	Create(ctx context.Context, blueprint *model.Blueprint) error
	Page(ctx context.Context, filter request.BlueprintFilter) ([]*model.Blueprint, int, error)
	GetById(ctx context.Context, id int64) (*model.Blueprint, error)
	Update(ctx context.Context, blueprint *model.Blueprint, id int64) error
	GetCompleteById(ctx context.Context, id int64) (*model.Blueprint, error)
}

type blueprintRepository struct {
	db *gorm.DB
}

func NewBlueprintRepository(db *gorm.DB) BlueprintRepository {
	return &blueprintRepository{db: db}
}

func (r *blueprintRepository) Create(ctx context.Context, blueprint *model.Blueprint) error {
	err := r.db.WithContext(ctx).Omit("id", "name").Create(&blueprint).Error
	return err
}

func (r *blueprintRepository) Page(ctx context.Context, filter request.BlueprintFilter) ([]*model.Blueprint, int, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(`SELECT *, 
							COUNT(1) OVER() as total_rows
							FROM blueprint_t`)

	var params []interface{}
	var conditions []string
	if filter.CreatorId != 0 {
		conditions = append(conditions, "user_id = ?")
		params = append(params, filter.CreatorId)
	}
	if len(conditions) > 0 {
		sqlBuilder.WriteString(" WHERE ")
		sqlBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	sortRule, ok := constant.BlueprintOrders[constant.OrderOption(filter.Order)]
	if !ok {
		sortRule = constant.BlueprintOrders[constant.RECENTLY_MODIFIED]
	}
	sqlBuilder.WriteString(fmt.Sprintf(" ORDER BY %s %s ", sortRule.Column, sortRule.Direction))

	sqlBuilder.WriteString(" LIMIT ? OFFSET ? ")
	params = append(params, constant.PAGE_LIMIT)
	params = append(params, constant.PAGE_LIMIT*(filter.Page-1))

	sqlString := sqlBuilder.String()

	var result []*model.BlueprintWithTotal
	err := r.db.WithContext(ctx).Raw(sqlString, params...).Scan(&result).Error
	if err != nil {
		return nil, 0, err
	}
	if len(result) == 0 {
		return []*model.Blueprint{}, 0, nil
	}
	blueprints := make([]*model.Blueprint, len(result))
	for i, item := range result {
		blueprints[i] = &item.Blueprint
	}
	totalRows := result[0].TotalRows

	return blueprints, totalRows, nil
}

func (r *blueprintRepository) GetById(ctx context.Context, id int64) (*model.Blueprint, error) {
	var blueprint *model.Blueprint
	sql := `SELECT * 
			FROM blueprint_t
			WHERE id = ?`
	err := r.db.Raw(sql, id).Scan(&blueprint).Error

	return blueprint, err
}

func (r *blueprintRepository) Update(ctx context.Context, blueprint *model.Blueprint, id int64) error {
	blueprint.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Model(blueprint).Updates(blueprint).Error
}

func (r *blueprintRepository) GetCompleteById(ctx context.Context, id int64) (*model.Blueprint, error) {
	var blueprint *model.Blueprint
	err := r.db.Where("id = ?", id).Preload("Corners").Preload("Items.Furniture").Preload("Walls").Find(&blueprint).Error

	return blueprint, err
}
