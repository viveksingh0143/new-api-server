package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type Stock struct {
	ID            int64                   `db:"id"`
	ProductID     int64                   `db:"product_id"`
	StoreID       int64                   `db:"store_id"`
	BinID         customtypes.NullInt64   `db:"bin_id"`
	PalletID      customtypes.NullInt64   `db:"pallet_id"`
	RackID        customtypes.NullInt64   `db:"rack_id"`
	BatchLabelID  customtypes.NullInt64   `db:"batchlabel_id"`
	Barcode       string                  `db:"barcode"`
	BatchNo       string                  `db:"batch_no"`
	UnitWeight    float64                 `db:"unit_weight"`
	Quantity      int64                   `db:"quantity"`
	MachineCode   string                  `db:"machine_code"`
	StockInAt     time.Time               `db:"stockin_at"`
	StockOutAt    *time.Time              `db:"stockout_at"`
	Status        customtypes.StockStatus `db:"status"`
	CreatedAt     time.Time               `db:"created_at"`
	UpdatedAt     *time.Time              `db:"updated_at"`
	LastUpdatedBy customtypes.NullString  `db:"last_updated_by"`
}
