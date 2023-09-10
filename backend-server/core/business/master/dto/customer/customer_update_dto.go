package customer

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type CustomerUpdateDto struct {
	Code            string                 `json:"code" validate:"required"`
	Name            string                 `json:"name" validate:"required"`
	ContactPerson   customtypes.NullString `json:"contact_person"`
	BillingAddress  BillingAddress         `json:"billing_address"`
	ShippingAddress ShippingAddress        `json:"shipping_address"`
	Status          customtypes.StatusEnum `json:"status" validate:"required"`
	LastUpdatedBy   customtypes.NullString `json:"last_updated_by"`
}
