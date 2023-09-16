package requisition

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
)

type RequisitionCreateDto struct {
	IssuedDate    time.Time              `json:"issued_date"`
	OrderNo       string                 `json:"order_no"`
	Department    string                 `json:"department"`
	StoreID       int64                  `json:"store_id"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
	Items         []*RequisitionItemDto  `json:"items"`
	Store         *store.StoreMinimalDto `json:"store"`
}
