package outwardrequest

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type OutwardRequestDetailDto struct {
	ID            int64                          `json:"id"`
	IssuedDate    customtypes.NullTime           `json:"issued_date"`
	OrderNo       string                         `json:"order_no"`
	POCategory    string                         `json:"po_category"`
	CustomerID    int64                          `json:"customer_id"`
	Status        customtypes.StatusEnum         `json:"status"`
	CreatedAt     customtypes.NullTime           `json:"created_at"`
	UpdatedAt     customtypes.NullTime           `json:"updated_at"`
	LastUpdatedBy customtypes.NullString         `json:"last_updated_by"`
	Items         []*OutwardRequestItemDetailDto `json:"items"`
}
