package dto

import (
	"room-planner/model"
	"time"
)

type BlueprintDto struct {
	Id         int64     `json:"id"`
	UserId     int64     `json:"userId"`
	CreatedAt  time.Time `json:"createAt"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

type BlueprintCompleteDto struct {
	Id         int64          `json:"id"`
	UserId     int64          `json:"userId"`
	CreatedAt  time.Time      `json:"createAt"`
	ModifiedAt time.Time      `json:"modifiedAt"`
	Corners    []CornerDto    `json:"corners"`
	Walls      []WallDto      `json:"walls"`
	Items      []ItemDto      `json:"items"`
	Furniture  []FurnitureDto `json:"furniture"`
}

type BlueprintPageDto struct {
	NextPage   int            `json:"nextPage"`
	PrevPage   int            `json:"prevPage"`
	CurrPage   int            `json:"currPage"`
	TotalPages int            `json:"totalPages"`
	List       []BlueprintDto `json:"blueprints"`
}

func FromBlueprintModel(bp model.Blueprint) *BlueprintDto {
	return &BlueprintDto{
		Id:         bp.Id,
		UserId:     bp.UserId,
		CreatedAt:  bp.CreatedAt,
		ModifiedAt: bp.ModifiedAt,
	}
}

func FromBlueprintsModel(blueprints []*model.Blueprint) []BlueprintDto {
	var blueprintsDto []BlueprintDto
	for _, bp := range blueprints {
		dto := FromBlueprintModel(*bp)
		blueprintsDto = append(blueprintsDto, *dto)
	}

	return blueprintsDto
}

func FromCompleteBlueprintModel(bp model.Blueprint) *BlueprintCompleteDto {
	return &BlueprintCompleteDto{
		Id:        bp.Id,
		UserId:    bp.UserId,
		CreatedAt: bp.CreatedAt,
		Corners:   FromCornersModel(bp.Corners),
		Walls:     FromWallsModel(bp.Walls),
		Items:     FromItemsModel(bp.Items),
	}
}

func FromCompleteBlueprintsModel(blueprints []*model.Blueprint) []BlueprintCompleteDto {
	var blueprintsDto []BlueprintCompleteDto
	for _, bp := range blueprints {
		dto := FromCompleteBlueprintModel(*bp)
		blueprintsDto = append(blueprintsDto, *dto)
	}

	return blueprintsDto
}
