package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ShipperLabel struct {
	ID               int64                  `db:"id" json:"id"`
	PackedAt         time.Time              `db:"packed_at" json:"packed_at"`
	ShipperNumber    string                 `db:"shipper_number" json:"shipper_number"`
	CustomerName     string                 `db:"customer_name" json:"customer_name"`
	ProductCode      string                 `db:"product_code" json:"product_code"`
	ProductName      string                 `db:"product_name" json:"product_name"`
	BatchNo          string                 `db:"batch_no" json:"batch_no"`
	PackedQty        string                 `db:"packed_qty" json:"packed_qty"`
	OutwardRequestID int64                  `db:"outwardrequest_id" json:"outwardrequest_id"`
	Status           customtypes.StatusEnum `db:"status" json:"status"`
	CreatedAt        time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt        *time.Time             `db:"updated_at" json:"updated_at"`
	LastUpdatedBy    customtypes.NullString `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt        *time.Time             `db:"deleted_at" json:"deleted_at,omitempty"`
}
