package outwardrequest

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
)

type OutwardRequestUpdateDto struct {
	IssuedDate    time.Time                    `json:"issued_date"`
	OrderNo       string                       `json:"order_no"`
	CustomerID    int64                        `json:"customer_id"`
	Status        customtypes.StatusEnum       `json:"status"`
	LastUpdatedBy customtypes.NullString       `json:"last_updated_by"`
	Customer      *customer.CustomerMinimalDto `json:"customer"`
	Items         []*OutwardRequestItemDto     `json:"items"`
}
