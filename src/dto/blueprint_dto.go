package dto

import (
	"room-planner/model"
	"time"
)

type BlueprintDto struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"userId"`
	CreatedAt time.Time `json:"createAT"`
}

func FromBlueprintModel(bp model.Blueprint) *BlueprintDto {
	return &BlueprintDto{
		Id:        bp.Id,
		UserId:    bp.UserId,
		CreatedAt: bp.CreatedAt,
	}
}
