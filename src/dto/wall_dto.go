package dto

import (
	"room-planner/model"

	"github.com/google/uuid"
)

type WallDto struct {
	Id            int64     `json:"id"`
	BlueprintId   int64     `json:"blueprintId"`
	StartCornerId uuid.UUID `json:"startCornerId"`
	EndCornerId   uuid.UUID `json:"endCornerId"`
}

func FromWallModel(wall model.Wall) *WallDto {
	return &WallDto{
		Id:            wall.Id,
		BlueprintId:   wall.BlueprintId,
		StartCornerId: wall.StartCornerId,
		EndCornerId:   wall.EndCornerId,
	}
}

func FromWallsModel(corners []*model.Wall) []WallDto {
	var wallsDto []WallDto
	for _, wall := range corners {
		dto := FromWallModel(*wall)
		wallsDto = append(wallsDto, *dto)
	}

	return wallsDto
}
