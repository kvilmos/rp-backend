package model

import (
	"time"
)

type Blueprint struct {
	Id        int64
	UserId    int64
	CreatedAt time.Time
}

func (Blueprint) TableName() string {
	return "blueprint_t"
}
