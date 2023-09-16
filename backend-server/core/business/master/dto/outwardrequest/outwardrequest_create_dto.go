package outwardrequest

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
)

type OutwardRequestCreateDto struct {
	IssuedDate    time.Time                    `json:"issued_date"`
	OrderNo       string                       `json:"order_no"`
	CustomerID    int64                        `json:"customer_id"`
	LastUpdatedBy customtypes.NullString       `json:"last_updated_by"`
	Items         []*OutwardRequestItemDto     `json:"items"`
	Customer      *customer.CustomerMinimalDto `json:"customer"`
}
