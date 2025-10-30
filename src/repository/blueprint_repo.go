package repository

import (
	"context"
	"room-planner/common/constant"
	"room-planner/model"
	"room-planner/request"
	"strings"
	"time"

	"gorm.io/gorm"
)

type BlueprintRepository interface {
	Create(ctx context.Context, blueprint *model.Blueprint) error
	List(ctx context.Context) ([]*model.Blueprint, error)
	Page(ctx context.Context, page int) ([]*model.Blueprint, error)
	PageByFilter(ctx context.Context, page int, filter request.BlueprintFilter) ([]*model.Blueprint, error)
	GetById(ctx context.Context, id int64) (*model.Blueprint, error)
	Update(ctx context.Context, blueprint *model.Blueprint, id int64) error
	ListComplete(ctx context.Context) ([]*model.Blueprint, error)
	GetCompleteById(ctx context.Context, id int64) (*model.Blueprint, error)
	Count(ctx context.Context) (int, error)
	CountByUserId(ctx context.Context, id int64) (int, error)
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

func (r *blueprintRepository) Page(ctx context.Context, page int) ([]*model.Blueprint, error) {
	var blueprints []*model.Blueprint
	sql := `SELECT * 
			FROM blueprint_t
			LIMIT ?
			OFFSET ?`
	err := r.db.Raw(sql, constant.PAGE_LIMIT, constant.PAGE_LIMIT*(page-1)).Scan(&blueprints).Error

	return blueprints, err
}

func (r *blueprintRepository) PageByFilter(ctx context.Context, page int, filter request.BlueprintFilter) ([]*model.Blueprint, error) {
	var sqlBuilder strings.Builder
	sqlBuilder.WriteString(`SELECT * FROM blueprint_t`)
	var params []interface{}
	var conditions []string
	conditions = append(conditions, "user_id LIKE ?")
	params = append(params, filter.CreatorId)

	if len(conditions) > 0 {
		sqlBuilder.WriteString(" WHERE ")
		sqlBuilder.WriteString(strings.Join(conditions, " AND "))
	}

	switch filter.OrderBy {
	case "latest":
		sqlBuilder.WriteString(" ORDER BY created_at DESC ")
	case "oldest":
		sqlBuilder.WriteString(" ORDER BY created_at ASC ")
	default:
		sqlBuilder.WriteString(" ORDER BY created_at DESC ")
	}

	sqlBuilder.WriteString(" LIMIT ? OFFSET ? ")
	params = append(params, constant.PAGE_LIMIT)
	params = append(params, constant.PAGE_LIMIT*(page-1))

	var list []*model.Blueprint
	sqlString := sqlBuilder.String()
	err := r.db.Raw(sqlString, params...).Scan(&list).Error

	return list, err
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

func (r *blueprintRepository) Update(ctx context.Context, blueprint *model.Blueprint, id int64) error {
	blueprint.UpdatedAt = time.Now()
	return r.db.WithContext(ctx).Model(blueprint).Updates(blueprint).Error
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

func (r *blueprintRepository) CountByUserId(ctx context.Context, id int64) (int, error) {
	var count int
	sql := `SELECT COUNT(1) 
			FROM blueprint_t
			WHERE user_id = ?`
	err := r.db.Raw(sql, id).Scan(&count).Error

	return count, err
}
