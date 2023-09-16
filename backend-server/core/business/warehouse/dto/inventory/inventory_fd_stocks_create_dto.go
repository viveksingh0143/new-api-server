package inventory

import "github.com/vamika-digital/wms-api-server/core/base/customtypes"

type InventoryFDSingleStockCreateDto struct {
	StoreID       int64                  `json:"store_id" validate:"required,min=1"`
	BinCode       string                 `json:"bin" validate:"required"`
	Barcode       string                 `json:"barcode" validate:"required"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
}
