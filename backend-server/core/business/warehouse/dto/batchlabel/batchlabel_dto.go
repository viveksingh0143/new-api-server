package batchlabel

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/machine"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
)

type BatchLabelDto struct {
	ID              int64                        `json:"id"`
	BatchDate       customtypes.NullTime         `json:"batch_date"`
	BatchNo         string                       `json:"batch_no"`
	SoNumber        string                       `json:"so_number"`
	TargetQuantity  int64                        `json:"target_quantity"`
	PackageQuantity int64                        `json:"package_quantity"`
	PoCategory      string                       `json:"po_category"`
	UnitWeight      float32                      `json:"unit_weight"`
	UnitWeightType  customtypes.UnitWeightType   `json:"unit_weight_type"`
	Status          customtypes.StatusEnum       `json:"status"`
	CreatedAt       customtypes.NullTime         `json:"created_at"`
	UpdatedAt       customtypes.NullTime         `json:"updated_at"`
	LastUpdatedBy   customtypes.NullString       `json:"last_updated_by"`
	Customer        *customer.CustomerMinimalDto `json:"customer"`
	Product         *product.ProductMinimalDto   `json:"product"`
	Machine         *machine.MachineMinimalDto   `json:"machine"`

	LabelsToPrint int64 `json:"labels_to_print"`
	TotalPrinted  int64 `json:"total_printed"`
}
