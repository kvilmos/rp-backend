package model

import (
	"github.com/google/uuid"
)

type Corner struct {
	Id          uuid.UUID
	BlueprintId int64
	X           float64
	Y           float64
}

func (Corner) TableName() string {
	return "corner_t"
}
