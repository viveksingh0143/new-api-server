package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type Requisition struct {
	ID            int64                  `db:"id" json:"id"`
	IssuedDate    time.Time              `db:"issued_date" json:"issued_date"`
	OrderNo       string                 `db:"order_no" json:"order_no"`
	Department    string                 `db:"department" json:"department"`
	StoreID       int64                  `db:"store_id" json:"store_id"`
	Status        customtypes.StatusEnum `db:"status" json:"status"`
	CreatedAt     time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time             `db:"updated_at" json:"updated_at"`
	LastUpdatedBy customtypes.NullString `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt     *time.Time             `db:"deleted_at" json:"deleted_at,omitempty"`
	Items         []*RequisitionItem     `db:"_" json:"items,omitempty"`
	Store         *Store                 `db:"_" json:"store"`
}

type RequisitionItem struct {
	ID            int64    `db:"id" json:"id"`
	RequisitionID int64    `db:"requisition_id" json:"requisition_id"`
	ProductID     int64    `db:"product_id" json:"product_id"`
	Quantity      int64    `db:"quantity" json:"quantity"`
	Product       *Product `db:"_" json:"product"`
}
