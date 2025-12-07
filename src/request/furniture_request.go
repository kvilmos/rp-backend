package request

type NewFurnitureRequest struct {
	UserId     int64
	Name       string  `json:"name" validate:"required,min=4,max=100"`
	CategoryId int64   `json:"categoryId"`
	SizeX      float64 `json:"sizeX"`
	SizeY      float64 `json:"sizeY"`
	SizeZ      float64 `json:"sizeZ"`
}
