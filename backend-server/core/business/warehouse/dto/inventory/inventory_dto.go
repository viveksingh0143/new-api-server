package inventory

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type InventoryDto struct {
	ID             int64                      `json:"id"`
	ProductType    customtypes.ProductType    `json:"product_type"`
	Code           string                     `json:"code"`
	Name           string                     `json:"name"`
	UnitType       customtypes.UnitType       `json:"unit_type"`
	UnitWeight     float32                    `json:"unit_weight"`
	UnitWeightType customtypes.UnitWeightType `json:"unit_weight_type"`
	StockCount     int64                      `json:"stock_count"`
	StockinAt      customtypes.NullTime       `json:"stockin_at"`
}
