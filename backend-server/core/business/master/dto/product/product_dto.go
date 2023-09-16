package product

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ProductDto struct {
	ID             int64                      `json:"id"`
	ProductType    customtypes.ProductType    `json:"product_type"`
	Code           string                     `json:"code"`
	LinkCode       string                     `json:"link_code"`
	Name           string                     `json:"name"`
	Description    string                     `json:"description"`
	UnitType       customtypes.UnitType       `json:"unit_type"`
	UnitWeight     float64                    `json:"unit_weight"`
	UnitWeightType customtypes.UnitWeightType `json:"unit_weight_type"`
	Status         customtypes.StatusEnum     `json:"status"`
	CreatedAt      customtypes.NullTime       `json:"created_at"`
	UpdatedAt      customtypes.NullTime       `json:"updated_at"`
	LastUpdatedBy  customtypes.NullString     `json:"last_updated_by"`
}
