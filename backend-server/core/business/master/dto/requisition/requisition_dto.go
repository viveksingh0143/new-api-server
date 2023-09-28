package requisition

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
)

type RequisitionDto struct {
	ID            int64                  `json:"id"`
	IssuedDate    customtypes.NullTime   `json:"issued_date"`
	OrderNo       string                 `json:"order_no"`
	Department    string                 `json:"department"`
	StoreID       int64                  `json:"store_id"`
	Store         *store.StoreMinimalDto `json:"store"`
	IsApproved    bool                   `json:"approved"`
	Status        customtypes.StatusEnum `json:"status"`
	CreatedAt     customtypes.NullTime   `json:"created_at"`
	UpdatedAt     customtypes.NullTime   `json:"updated_at"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
	Items         []*RequisitionItemDto  `json:"items"`
}
