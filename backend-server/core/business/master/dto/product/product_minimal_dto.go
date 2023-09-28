package product

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ProductMinimalDto struct {
	ID             int64                      `json:"id"`
	ProductType    customtypes.ProductType    `json:"product_type"`
	ProductSubType string                     `json:"product_subtype"`
	Code           string                     `json:"code"`
	LinkCode       string                     `json:"link_code"`
	Name           string                     `json:"name"`
	Description    string                     `json:"description"`
	UnitType       customtypes.UnitType       `json:"unit_type"`
	UnitWeight     float64                    `json:"unit_weight"`
	UnitWeightType customtypes.UnitWeightType `json:"unit_weight_type"`
	Status         customtypes.StatusEnum     `json:"status"`
}
