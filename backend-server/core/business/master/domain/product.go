package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type Product struct {
	ID             int64                      `db:"id" json:"id"`
	ProductType    customtypes.ProductType    `db:"product_type" json:"product_type"`
	ProductSubType string                     `db:"product_subtype" json:"product_subtype"`
	Code           string                     `db:"code" json:"code"`
	LinkCode       string                     `db:"link_code" json:"link_code"`
	Name           string                     `db:"name" json:"name"`
	Description    string                     `db:"description" json:"description"`
	UnitType       customtypes.UnitType       `db:"unit_type" json:"unit_type"`
	UnitWeight     float64                    `db:"unit_weight" json:"unit_weight"`
	UnitWeightType customtypes.UnitWeightType `db:"unit_weight_type" json:"unit_weight_type"`
	Status         customtypes.StatusEnum     `db:"status" json:"status"`
	CreatedAt      time.Time                  `db:"created_at" json:"created_at"`
	UpdatedAt      *time.Time                 `db:"updated_at" json:"updated_at"`
	LastUpdatedBy  customtypes.NullString     `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt      *time.Time                 `db:"deleted_at" json:"deleted_at,omitempty"`
}
