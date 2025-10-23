package model

type Item struct {
	Id          int64
	BlueprintId int64
	FurnitureId int64
	PosX        float64
	PosY        float64
	PosZ        float64
	Rot         float64
}

func (Item) TableName() string {
	return "item_t"
}
