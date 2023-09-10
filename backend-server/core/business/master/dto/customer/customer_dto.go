package customer

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type CustomerDto struct {
	ID              int64                  `json:"id"`
	Code            string                 `json:"code"`
	Name            string                 `json:"name"`
	ContactPerson   customtypes.NullString `json:"contact_person"`
	BillingAddress  BillingAddress         `json:"billing_address"`
	ShippingAddress ShippingAddress        `json:"shipping_address"`
	Status          customtypes.StatusEnum `json:"status"`
	CreatedAt       customtypes.NullTime   `json:"created_at"`
	UpdatedAt       customtypes.NullTime   `json:"updated_at"`
	LastUpdatedBy   customtypes.NullString `json:"last_updated_by"`
}
