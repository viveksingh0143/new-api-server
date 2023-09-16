package batchlabel

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/machine"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
)

type BatchLabelCreateDto struct {
	BatchDate       time.Time                    `json:"batch_date" validate:"required"`
	BatchNo         string                       `json:"batch_no" validate:"required"`
	SoNumber        string                       `json:"so_number"`
	TargetQuantity  int64                        `json:"target_quantity" validate:"required"`
	PackageQuantity int64                        `json:"package_quantity" validate:"required"`
	PoCategory      string                       `json:"po_category" validate:"required"`
	UnitWeight      float32                      `json:"unit_weight" validate:"required"`
	UnitWeightType  customtypes.UnitWeightType   `json:"unit_weight_type" validate:"required"`
	Status          customtypes.StatusEnum       `json:"status" validate:"required"`
	LastUpdatedBy   customtypes.NullString       `json:"last_updated_by"`
	Customer        *customer.CustomerMinimalDto `json:"customer" validate:"required"`
	Product         *product.ProductMinimalDto   `json:"product" validate:"required"`
	Machine         *machine.MachineMinimalDto   `json:"machine" validate:"required"`
}
