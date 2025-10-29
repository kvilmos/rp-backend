package dto

import (
	"room-planner/model"
	"time"
)

type BlueprintDto struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"userId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type BlueprintCompleteDto struct {
	Id        int64          `json:"id"`
	UserId    int64          `json:"userId"`
	Name      string         `json:"name"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	Corners   []CornerDto    `json:"corners"`
	Walls     []WallDto      `json:"walls"`
	Items     []ItemDto      `json:"items"`
	Furniture []FurnitureDto `json:"furniture"`
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
		Id:        bp.Id,
		UserId:    bp.UserId,
		Name:      bp.Name,
		CreatedAt: bp.CreatedAt,
		UpdatedAt: bp.UpdatedAt,
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
		Name:      bp.Name,
		CreatedAt: bp.CreatedAt,
		UpdatedAt: bp.UpdatedAt,
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
