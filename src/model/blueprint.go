package model

import (
	"time"
)

type Blueprint struct {
	Id        int64
	UserId    int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Corners   []*Corner
	Walls     []*Wall
	Items     []*Item
}

func (Blueprint) TableName() string {
	return "blueprint_t"
}

type BlueprintWithTotal struct {
	Blueprint Blueprint `gorm:"embedded"`
	TotalRows int
}
