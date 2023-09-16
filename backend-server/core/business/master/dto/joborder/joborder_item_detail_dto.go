package joborder

import "github.com/vamika-digital/wms-api-server/core/business/master/dto/product"

type JobOrderItemDetailDto struct {
	ID         int64 `json:"id"`
	JobOrderID int64 `json:"joborder_id"`
	ProductID  int64 `json:"product_id"`
	Quantity   int64 `json:"quantity"`

	Product product.ProductMinimalDto `json:"product"`
}
