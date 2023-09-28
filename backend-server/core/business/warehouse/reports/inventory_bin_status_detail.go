package reports

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type InventoryBinStatusDetail struct {
	BinID                 int64                `db:"bin_id" json:"bin_id"`
	BinCode               string               `db:"bin_code" json:"bin_code"`
	BinName               string               `db:"bin_name" json:"bin_name"`
	ProductID             int64                `db:"product_id" json:"product_id"`
	ProductName           string               `db:"product_name" json:"product_name"`
	ProductCode           string               `db:"product_code" json:"product_code"`
	LockCount             int64                `db:"lock_count" json:"lock_count"`
	StockInCount          int64                `db:"stockin_count" json:"stockin_count"`
	StockDispatchingCount int64                `db:"stockdispatching_count" json:"stockdispatching_count"`
	StockOutCount         int64                `db:"stockout_count" json:"stockout_count"`
	Stockin_at            customtypes.NullTime `db:"stockin_at" json:"stockin_at"`
	RequiredStocks        int64                `db:"_" json:"required_stocks"`
}
