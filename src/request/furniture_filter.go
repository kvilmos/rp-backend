package request

type FurnitureFilter struct {
	Page       int
	Order      string
	CategoryId *int64
	CreatorId  int64
}
