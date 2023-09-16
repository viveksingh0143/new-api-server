package stock

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type StockFilterDto struct {
	Query        string                  `form:"query" json:"query"`
	ID           int64                   `form:"id" json:"id"`
	ProductID    int64                   `form:"product_id" json:"product_id"`
	StoreID      int64                   `form:"store_id" json:"store_id"`
	BinID        customtypes.NullInt64   `form:"bin_id" json:"bin_id"`
	PalletID     customtypes.NullInt64   `form:"pallet_id" json:"pallet_id"`
	RackID       customtypes.NullInt64   `form:"rack_id" json:"rack_id"`
	BatchLabelID customtypes.NullInt64   `form:"batchlabel_id" json:"batchlabel_id"`
	Barcode      string                  `form:"barcode" json:"barcode"`
	BatchNo      string                  `form:"batch_no" json:"batch_no"`
	StockInAt    customtypes.NullTime    `form:"stockin_at" json:"stockin_at"`
	StockOutAt   customtypes.NullTime    `form:"stockout_at" json:"stockout_at"`
	Status       customtypes.StockStatus `form:"status" json:"status"`
}
