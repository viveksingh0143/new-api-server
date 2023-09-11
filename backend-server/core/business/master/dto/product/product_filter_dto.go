package product

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ProductFilterDto struct {
	Query       string                  `form:"query" json:"query"`
	ID          int64                   `form:"id" json:"id"`
	ProductType customtypes.ProductType `form:"product_type" json:"product_type"`
	LinkCode    string                  `form:"link_code" json:"link_code"`
	Code        string                  `form:"code" json:"code"`
	Name        string                  `form:"name" json:"name"`
	Status      customtypes.StatusEnum  `form:"status" json:"status"`
}
