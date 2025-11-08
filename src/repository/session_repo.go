package repository

import (
	"room-planner/model"

	"gorm.io/gorm"
)

type SessionRepository interface {
	WithTx(tx *gorm.DB) SessionRepository
	Create(session *model.Session) (*model.Session, error)
	GetById(id string) (*model.Session, error)
	Revoke(id string) error
	Delete(id string) error
}

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) WithTx(tx *gorm.DB) SessionRepository {
	return NewSessionRepository(tx)
}

func (r *sessionRepository) Create(session *model.Session) (*model.Session, error) {
	res := r.db.Create(&session)

	return session, res.Error
}

func (r *sessionRepository) GetById(id string) (*model.Session, error) {
	session := new(model.Session)
	sql := `SELECT id, user_email, refresh_token, is_revoked, created_at, expires_at 
			FROM session_t 
			WHERE id = ?`

	err := r.db.Raw(sql, id).Scan(&session).Error
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (r *sessionRepository) Revoke(id string) error {
	var newRevoke = 1
	sql := `UPDATE session_t 
			SET is_revoked = ?
			WHERE id = ?`

	return r.db.Exec(sql, newRevoke, id).Error
}

func (r *sessionRepository) Delete(id string) error {
	sql := `DELETE FROM session_t 
			WHERE id = ?`

	return r.db.Exec(sql, id).Error
}
