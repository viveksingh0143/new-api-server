package joborder

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
)

type JobOrderUpdateDto struct {
	IssuedDate    time.Time                    `json:"issued_date"`
	OrderNo       string                       `json:"order_no"`
	POCategory    string                       `json:"po_category"`
	CustomerID    int64                        `json:"customer_id"`
	Status        customtypes.StatusEnum       `json:"status"`
	LastUpdatedBy customtypes.NullString       `json:"last_updated_by"`
	Customer      *customer.CustomerMinimalDto `json:"customer"`
	Items         []*JobOrderItemDto           `json:"items"`
}
