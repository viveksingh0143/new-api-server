package reports

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type InventoryStatusDetail struct {
	RackID        int64                `db:"rack_id" json:"rack_id"`
	RackCode      string               `db:"rack_code" json:"rack_code"`
	RackName      string               `db:"rack_name" json:"rack_name"`
	RackAddress   string               `db:"rack_address" json:"rack_address"`
	ProductID     int64                `db:"product_id" json:"product_id"`
	ProductName   string               `db:"product_name" json:"product_name"`
	ProductCode   string               `db:"product_code" json:"product_code"`
	StockinCount  int64                `db:"stockin_count" json:"stockin_count"`
	StockoutCount int64                `db:"stockout_count" json:"stockout_count"`
	Stockin_at    customtypes.NullTime `db:"stockin_at" json:"stockin_at"`
}
