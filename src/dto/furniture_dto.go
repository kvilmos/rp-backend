package dto

import "time"

type FurnitureListDto struct {
	Furniture []FurnitureDto `json:"Furniture"`
}

type FurnitureDto struct {
	Id           int64     `json:"id"`
	Name         string    `json:"name"`
	CategoryId   int8      `json:"categoryId"`
	ObjectUrl    string    `json:"objectUrl"`
	ThumbnailUrl string    `json:"thumbnailUrl"`
	CreatedAt    time.Time `json:"createdAt"`
}
