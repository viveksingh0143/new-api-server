package product

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ProductUpdateDto struct {
	ProductType    customtypes.ProductType    `json:"product_type" validate:"ProductType,required"`
	Code           string                     `json:"code" validate:"required"`
	LinkCode       string                     `json:"link_code"`
	Name           string                     `json:"name" validate:"required"`
	Description    string                     `json:"description"`
	UnitType       customtypes.UnitType       `json:"unit_type" validate:"UnitType,required"`
	UnitWeight     float32                    `json:"unit_weight"`
	UnitWeightType customtypes.UnitWeightType `json:"unit_weight_type" validate:"UnitWeightType,required"`
	Status         customtypes.StatusEnum     `json:"status"`
	LastUpdatedBy  customtypes.NullString     `json:"last_updated_by"`
}
