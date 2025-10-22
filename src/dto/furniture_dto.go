package dto

import (
	"room-planner/model"
	"time"

	"github.com/google/uuid"
)

type FurnitureDto struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	UserId       int64     `json:"userId"`
	SizeX        float64   `json:"sizeX"`
	SizeY        float64   `json:"sizeY"`
	SizeZ        float64   `json:"sizeZ"`
	CategoryId   int8      `json:"categoryId"`
	FileName     uuid.UUID `json:"fileName"`
	ThumbnailUrl string    `json:"thumbnailUrl"`
	ObjectUrl    string    `json:"objectUrl"`
	CreatedAt    time.Time `json:"createdAt"`
}

func FromFurnitureModel(furniture *model.Furniture) *FurnitureDto {
	return &FurnitureDto{
		Id:         furniture.Id,
		Name:       furniture.Name,
		UserId:     furniture.UserId,
		SizeX:      furniture.SizeX,
		SizeY:      furniture.SizeY,
		SizeZ:      furniture.SizeZ,
		CategoryId: furniture.CategoryId,
		FileName:   furniture.FileName,
		CreatedAt:  furniture.CreatedAt,
	}
}

type FurniturePaginationDto struct {
	NextPage   int            `json:"nextPage"`
	PrevPage   int            `json:"prevPage"`
	CurrPage   int            `json:"currPage"`
	TotalPages int            `json:"totalPages"`
	List       []FurnitureDto `json:"furniture"`
}
