package model

import "github.com/google/uuid"

type Wall struct {
	Id            int64
	BlueprintId   int64
	StartCornerId uuid.UUID
	EndCornerId   uuid.UUID
}

func (Wall) TableName() string {
	return "wall_t"
}
