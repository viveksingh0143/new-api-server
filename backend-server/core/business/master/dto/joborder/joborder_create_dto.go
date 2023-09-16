package joborder

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
)

type JobOrderCreateDto struct {
	IssuedDate    time.Time                    `json:"issued_date"`
	OrderNo       string                       `json:"order_no"`
	POCategory    string                       `json:"po_category"`
	CustomerID    int64                        `json:"customer_id"`
	LastUpdatedBy customtypes.NullString       `json:"last_updated_by"`
	Items         []*JobOrderItemDto           `json:"items"`
	Customer      *customer.CustomerMinimalDto `json:"customer"`
}
