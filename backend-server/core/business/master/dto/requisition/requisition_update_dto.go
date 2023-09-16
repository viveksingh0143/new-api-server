package requisition

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
)

type RequisitionUpdateDto struct {
	IssuedDate    time.Time              `json:"issued_date"`
	OrderNo       string                 `json:"order_no"`
	Department    string                 `json:"department"`
	StoreID       int64                  `json:"store_id"`
	Status        customtypes.StatusEnum `json:"status"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
	Store         *store.StoreMinimalDto `json:"store"`
	Items         []*RequisitionItemDto  `json:"items"`
}
