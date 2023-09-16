package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type JobOrder struct {
	ID            int64                  `db:"id" json:"id"`
	IssuedDate    time.Time              `db:"issued_date" json:"issued_date"`
	OrderNo       string                 `db:"order_no" json:"order_no"`
	POCategory    string                 `db:"po_category" json:"po_category"`
	CustomerID    int64                  `db:"customer_id" json:"customer_id"`
	Status        customtypes.StatusEnum `db:"status" json:"status"`
	CreatedAt     time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time             `db:"updated_at" json:"updated_at"`
	LastUpdatedBy customtypes.NullString `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt     *time.Time             `db:"deleted_at" json:"deleted_at,omitempty"`
	Items         []*JobOrderItem        `db:"_" json:"items,omitempty"`
	Customer      *Customer              `db:"_" json:"customer"`
}

type JobOrderItem struct {
	ID         int64    `db:"id" json:"id"`
	JobOrderID int64    `db:"joborder_id" json:"joborder_id"`
	ProductID  int64    `db:"product_id" json:"product_id"`
	Quantity   int64    `db:"quantity" json:"quantity"`
	Product    *Product `db:"_" json:"product"`
}
