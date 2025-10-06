package request

type NewFurnitureRequest struct {
	UserId int64
	Name   string `json:"name" validate:"required,min=4,max=100"`
}
