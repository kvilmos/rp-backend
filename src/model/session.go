package model

import "time"

type Session struct {
	Id           string
	UserEmail    string
	RefreshToken string
	IsRevoked    bool
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

func (Session) TableName() string {
	return "session_t"
}
