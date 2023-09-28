package requisition

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
)

type RequisitionFilterDto struct {
	Query      string                 `form:"query"`
	ID         int64                  `form:"id"`
	OrderNo    string                 `form:"order_no"`
	IsApproved customtypes.NullBool   `form:"approved"`
	Status     customtypes.StatusEnum `form:"status"`
	Store      *store.StoreMinimalDto `form:"store"`
}
