package model

import (
	"time"
)

type Blueprint struct {
	Id         int64
	UserId     int64
	CreatedAt  time.Time
	ModifiedAt time.Time
	Corners    []*Corner
	Walls      []*Wall
	Items      []*Item
}

func (Blueprint) TableName() string {
	return "blueprint_t"
}
