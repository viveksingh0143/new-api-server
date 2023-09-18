package outwardrequest

import "github.com/vamika-digital/wms-api-server/core/business/master/dto/product"

type OutwardRequestItemDto struct {
	ID               int64                      `json:"id"`
	OutwardRequestID int64                      `json:"outwardrequest_id"`
	ProductID        int64                      `json:"product_id"`
	Quantity         int64                      `json:"quantity"`
	LockedQuantity   int64                      `json:"locked_quantity"`
	Product          *product.ProductMinimalDto `json:"product"`
}
