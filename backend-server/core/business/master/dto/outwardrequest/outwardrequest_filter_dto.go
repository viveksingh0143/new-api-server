package outwardrequest

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
)

type OutwardRequestFilterDto struct {
	Query    string                       `form:"query"`
	ID       int64                        `form:"id"`
	OrderNo  string                       `form:"order_no"`
	Status   customtypes.StatusEnum       `form:"status"`
	Customer *customer.CustomerMinimalDto `form:"customer"`
}
