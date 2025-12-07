package dto

import "room-planner/model"

type ItemDto struct {
	Id          int64        `json:"id"`
	BlueprintId int64        `json:"blueprintId"`
	FurnitureId int64        `json:"furnitureId"`
	PosX        float64      `json:"posX"`
	PosY        float64      `json:"posY"`
	PosZ        float64      `json:"posZ"`
	Rot         float64      `json:"rot"`
	Furniture   FurnitureDto `json:"furniture"`
}

func FromItemModel(item model.Item) *ItemDto {
	return &ItemDto{
		Id:          item.Id,
		BlueprintId: item.BlueprintId,
		FurnitureId: item.FurnitureId,
		PosX:        item.PosX,
		PosY:        item.PosY,
		PosZ:        item.PosZ,
		Rot:         item.Rot,
		Furniture:   *FromFurnitureModel(item.Furniture),
	}
}

func FromItemsModel(items []*model.Item) []ItemDto {
	var itemsDto []ItemDto
	for _, item := range items {
		dto := FromItemModel(*item)
		itemsDto = append(itemsDto, *dto)
	}

	return itemsDto
}
