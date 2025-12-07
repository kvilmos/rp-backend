package repository

import (
	"room-planner/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	WithTx(tx *gorm.DB) UserRepository
	Create(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetById(id int64) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) WithTx(tx *gorm.DB) UserRepository {
	return NewUserRepository(tx)
}

func (r *userRepository) Create(user *model.User) error {
	err := r.db.Omit("id").Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	user := new(model.User)
	sql := `SELECT id, username, email, password, created_at 
			FROM user_t 
			WHERE email = ?`

	err := r.db.Raw(sql, email).Scan(&user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

func (r *userRepository) GetById(id int64) (*model.User, error) {
	var user model.User
	sql := `SELECT id, username, email, password, created_at 
			FROM user_t 
			WHERE id = ?`

	err := r.db.Raw(sql, id).Scan(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, err
}
