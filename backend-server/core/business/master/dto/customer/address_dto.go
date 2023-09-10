package customer

import "github.com/vamika-digital/wms-api-server/core/base/customtypes"

type BillingAddress struct {
	Address1 customtypes.NullString `json:"address1"`
	Address2 customtypes.NullString `json:"address2"`
	State    customtypes.NullString `json:"state"`
	Country  customtypes.NullString `json:"country"`
	Pincode  customtypes.NullString `json:"pincode"`
}

type ShippingAddress struct {
	Address1 customtypes.NullString `json:"address1"`
	Address2 customtypes.NullString `json:"address2"`
	State    customtypes.NullString `json:"state"`
	Country  customtypes.NullString `json:"country"`
	Pincode  customtypes.NullString `json:"pincode"`
}
