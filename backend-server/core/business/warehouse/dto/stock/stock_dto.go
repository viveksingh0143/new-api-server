package stock

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type StockDto struct {
	ID            int64                   `json:"id"`
	ProductID     int64                   `json:"product_id"`
	StoreID       int64                   `json:"store_id"`
	BinID         customtypes.NullInt64   `json:"bin_id"`
	PalletID      customtypes.NullInt64   `json:"pallet_id"`
	RackID        customtypes.NullInt64   `json:"rack_id"`
	BatchLabelID  customtypes.NullInt64   `json:"batchlabel_id"`
	Barcode       string                  `json:"barcode"`
	BatchNo       string                  `json:"batch_no"`
	UnitWeight    float64                 `json:"unit_weight"`
	Quantity      int64                   `json:"quantity"`
	MachineCode   string                  `json:"machine_code"`
	StockInAt     time.Time               `json:"stockin_at"`
	StockOutAt    *time.Time              `json:"stockout_at"`
	Status        customtypes.StockStatus `json:"status"`
	CreatedAt     time.Time               `json:"created_at"`
	UpdatedAt     *time.Time              `json:"updated_at"`
	LastUpdatedBy customtypes.NullString  `json:"last_updated_by"`
}
