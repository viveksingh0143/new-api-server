package batchlabel

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/machine"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
)

type BatchLabelFilterDto struct {
	Query     string                       `form:"query"`
	ID        int64                        `form:"id"`
	BatchDate customtypes.NullTime         `form:"batch_date"`
	BatchNo   string                       `form:"batch_no"`
	SoNumber  string                       `form:"so_number"`
	Status    customtypes.StatusEnum       `form:"status"`
	Customer  *customer.CustomerMinimalDto `form:"customer"`
	Product   *product.ProductMinimalDto   `form:"product"`
	Machine   *machine.MachineMinimalDto   `form:"machine"`
}
