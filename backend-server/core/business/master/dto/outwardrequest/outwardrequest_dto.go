package outwardrequest

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
)

type OutwardRequestDto struct {
	ID            int64                        `json:"id"`
	IssuedDate    customtypes.NullTime         `json:"issued_date"`
	OrderNo       string                       `json:"order_no"`
	CustomerID    int64                        `json:"customer_id"`
	Status        customtypes.StatusEnum       `json:"status"`
	CreatedAt     customtypes.NullTime         `json:"created_at"`
	UpdatedAt     customtypes.NullTime         `json:"updated_at"`
	LastUpdatedBy customtypes.NullString       `json:"last_updated_by"`
	Items         []*OutwardRequestItemDto     `json:"items"`
	Customer      *customer.CustomerMinimalDto `json:"customer"`
}
