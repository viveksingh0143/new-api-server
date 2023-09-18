package reports

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type InventoryPalletStatusDetail struct {
	PalletID              int64                `db:"pallet_id" json:"pallet_id"`
	PalletCode            string               `db:"pallet_code" json:"pallet_code"`
	PalletName            string               `db:"pallet_name" json:"pallet_name"`
	ProductID             int64                `db:"product_id" json:"product_id"`
	ProductName           string               `db:"product_name" json:"product_name"`
	ProductCode           string               `db:"product_code" json:"product_code"`
	LockCount             int64                `db:"lock_count" json:"lock_count"`
	StockDispatchingCount int64                `db:"stockdispatching_count" json:"stockdispatching_count"`
	StockOutCount         int64                `db:"stockout_count" json:"stockout_count"`
	Stockin_at            customtypes.NullTime `db:"stockin_at" json:"stockin_at"`
	RequiredStocks        int64                `db:"_" json:"required_stocks"`
}
