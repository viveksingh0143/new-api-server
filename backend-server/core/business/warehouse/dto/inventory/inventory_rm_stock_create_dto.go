package inventory

import "github.com/vamika-digital/wms-api-server/core/base/customtypes"

type InventoryRMStockCreateDto struct {
	StoreID       int64                  `json:"store_id" validate:"required,min=1"`
	ProductID     int64                  `json:"product_id" validate:"required,min=1"`
	Quantity      int64                  `json:"quantity" validate:"required,min=1"`
	PalletCode    string                 `json:"pallet" validate:"required"`
	CreatePallet  customtypes.NullBool   `json:"create_pallet"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
}
