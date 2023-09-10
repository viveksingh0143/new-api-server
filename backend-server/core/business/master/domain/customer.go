package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type Customer struct {
	ID               int64                  `db:"id" json:"id"`
	Code             string                 `db:"code" json:"code"`
	Name             string                 `db:"name" json:"name"`
	ContactPerson    customtypes.NullString `db:"contact_person" json:"contact_person"`
	BillingAddress1  customtypes.NullString `db:"billing_address_address1" json:"billing_address_address1"`
	BillingAddress2  customtypes.NullString `db:"billing_address_address2" json:"billing_address_address2"`
	BillingState     customtypes.NullString `db:"billing_address_state" json:"billing_address_state"`
	BillingCountry   customtypes.NullString `db:"billing_address_country" json:"billing_address_country"`
	BillingPincode   customtypes.NullString `db:"billing_address_pincode" json:"billing_address_pincode"`
	ShippingAddress1 customtypes.NullString `db:"shipping_address_address1" json:"shipping_address_address1"`
	ShippingAddress2 customtypes.NullString `db:"shipping_address_address2" json:"shipping_address_address2"`
	ShippingState    customtypes.NullString `db:"shipping_address_state" json:"shipping_address_state"`
	ShippingCountry  customtypes.NullString `db:"shipping_address_country" json:"shipping_address_country"`
	ShippingPincode  customtypes.NullString `db:"shipping_address_pincode" json:"shipping_address_pincode"`
	Status           customtypes.StatusEnum `db:"status" json:"status"`
	CreatedAt        time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt        *time.Time             `db:"updated_at" json:"updated_at"`
	LastUpdatedBy    customtypes.NullString `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt        *time.Time             `db:"deleted_at" json:"deleted_at,omitempty"`
}

func NewCustomerWithDefaults() Customer {
	return Customer{
		Status: customtypes.Enable,
	}
}
