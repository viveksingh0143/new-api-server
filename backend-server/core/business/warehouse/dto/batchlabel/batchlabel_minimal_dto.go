package batchlabel

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
)

type BatchLabelMinimalDto struct {
	ID         int64                        `json:"id"`
	BatchDate  customtypes.NullTime         `json:"batch_date" validate:"required"`
	BatchNo    string                       `json:"batch_no" validate:"required"`
	SoNumber   string                       `json:"so_number"`
	PoCategory string                       `json:"po_category" validate:"required"`
	Status     customtypes.StatusEnum       `json:"status" validate:"required"`
	Customer   *customer.CustomerMinimalDto `json:"customer" validate:"required"`
	Product    *product.ProductMinimalDto   `json:"product" validate:"required"`
}
