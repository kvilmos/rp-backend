package request

type NewBlueprintRequest struct {
	Id      int64
	UserId  int64
	Corners []NewCornerRequest
	Walls   []NewWallRequest
	Items   []NewItemRequest
}

type NewCornerRequest struct {
	Id string
	X  float64
	Y  float64
}

type NewWallRequest struct {
	StartCornerId string
	EndCornerId   string
}

type NewItemRequest struct {
	FurnitureId int64
	PosX        float64
	PosY        float64
	PosZ        float64
	Rot         float64
}
