package stock

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/container"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/batchlabel"
)

type StockDetailDto struct {
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

	Product *product.ProductMinimalDto `json:"product"`
	Store   *store.StoreMinimalDto     `json:"store"`

	Bin        *container.ContainerMinimalDto   `json:"bin"`
	Pallet     *container.ContainerMinimalDto   `json:"pallet"`
	Rack       *container.ContainerMinimalDto   `json:"rack"`
	BatchLabel *batchlabel.BatchLabelMinimalDto `json:"batchlabel"`
}
