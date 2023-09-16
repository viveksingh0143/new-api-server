package requisition

import "github.com/vamika-digital/wms-api-server/core/business/master/dto/product"

type RequisitionItemDto struct {
	ID            int64                      `json:"id"`
	RequisitionID int64                      `json:"requisition_id"`
	ProductID     int64                      `json:"product_id"`
	Quantity      int64                      `json:"quantity"`
	Product       *product.ProductMinimalDto `json:"product"`
}
