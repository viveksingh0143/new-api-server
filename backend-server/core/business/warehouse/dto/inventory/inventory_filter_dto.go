package inventory

type InventoryFilterDto struct {
	Query         string `form:"query" json:"query"`
	ProductID     int64  `form:"product_id" json:"product_id"`
	ProductCode   string `form:"product_code" json:"product_code"`
	ProductTypes  string `form:"product_types" json:"product_types"`
	StoreID       int64  `form:"store_id" json:"store_id"`
	ContainerID   int64  `form:"container_id" json:"container_id"`
	ContainerType string `form:"container_type" json:"container_type"`
	ContainerCode string `form:"container_code" json:"container_code"`
}
