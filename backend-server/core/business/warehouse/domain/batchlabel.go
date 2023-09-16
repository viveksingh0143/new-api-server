package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
)

type BatchLabel struct {
	ID              int64                      `db:"id" json:"id"`
	BatchDate       time.Time                  `db:"batch_date" json:"batch_date"`
	BatchNo         string                     `db:"batch_no" json:"batch_no"`
	SoNumber        string                     `db:"so_number" json:"so_number"`
	TargetQuantity  int64                      `db:"target_quantity" json:"target_quantity"`
	PackageQuantity int64                      `db:"package_quantity" json:"package_quantity"`
	PoCategory      string                     `db:"po_category" json:"po_category"`
	UnitWeight      float32                    `db:"unit_weight" json:"unit_weight"`
	UnitWeightType  customtypes.UnitWeightType `db:"unit_weight_type" json:"unit_weight_type"`
	Status          customtypes.StatusEnum     `db:"status" json:"status"`
	CreatedAt       time.Time                  `db:"created_at" json:"created_at"`
	UpdatedAt       *time.Time                 `db:"updated_at" json:"updated_at"`
	LastUpdatedBy   customtypes.NullString     `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt       *time.Time                 `db:"deleted_at" json:"deleted_at,omitempty"`
	CustomerID      int64                      `db:"customer_id" json:"customer_id,string"`
	ProductID       int64                      `db:"product_id" json:"product_id,string"`
	MachineID       int64                      `db:"machine_id" json:"machine_id,string"`
	Customer        *domain.Customer           `db:"-" json:"customer"`
	Product         *domain.Product            `db:"-" json:"product"`
	Machine         *domain.Machine            `db:"-" json:"machine"`
	TotalPrinted    int64                      `json:"total_printed"`
}

func (b *BatchLabel) GetStickerCountToPrint() int64 {
	if b.PackageQuantity <= 0 {
		return 0
	}
	return b.TargetQuantity / b.PackageQuantity
}
