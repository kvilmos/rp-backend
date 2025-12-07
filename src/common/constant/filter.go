package constant

type OrderOption string

const (
	LATEST_CREATED    OrderOption = "latest_created"
	OLDEST_CREATED    OrderOption = "oldest_created"
	RECENTLY_MODIFIED OrderOption = "recently_modified"
	OLDEST_MODIFIED   OrderOption = "oldest_modified"
)

type SortRule struct {
	Column    string
	Direction string
}

var BlueprintOrders = map[OrderOption]SortRule{
	LATEST_CREATED:    {Column: "created_at", Direction: "DESC"},
	OLDEST_CREATED:    {Column: "created_at", Direction: "ASC"},
	RECENTLY_MODIFIED: {Column: "updated_at", Direction: "DESC"},
	OLDEST_MODIFIED:   {Column: "updated_at", Direction: "ASC"},
}

var FurnitureOrders = map[OrderOption]SortRule{
	LATEST_CREATED: {Column: "created_at", Direction: "DESC"},
	OLDEST_CREATED: {Column: "created_at", Direction: "ASC"},
}

type FilterParm string

const (
	PAGE        FilterParm = "page"
	ORDER       FilterParm = "order"
	CATEGORY_ID FilterParm = "category"
)
