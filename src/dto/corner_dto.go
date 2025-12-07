package dto

import (
	"room-planner/model"

	"github.com/google/uuid"
)

type CornerDto struct {
	Id          uuid.UUID `json:"id"`
	BlueprintId int64     `json:"blueprintId"`
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
}

func FromCornerModel(corner model.Corner) *CornerDto {
	return &CornerDto{
		Id:          corner.Id,
		BlueprintId: corner.BlueprintId,
		X:           corner.X,
		Y:           corner.Y,
	}
}

func FromCornersModel(corners []*model.Corner) []CornerDto {
	var cornersDto []CornerDto
	for _, corner := range corners {
		dto := FromCornerModel(*corner)
		cornersDto = append(cornersDto, *dto)
	}

	return cornersDto
}
