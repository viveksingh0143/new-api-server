package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type Role struct {
	ID            int64                  `db:"id" json:"id"`
	Name          string                 `db:"name" json:"name"`
	Status        customtypes.StatusEnum `db:"status" json:"status"`
	CreatedAt     time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time             `db:"updated_at" json:"updated_at"`
	LastUpdatedBy customtypes.NullString `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt     *time.Time             `db:"deleted_at" json:"deleted_at,omitempty"`
	Permissions   []*Permission          `json:"permissions"`
}

func NewRoleWithDefaults() Role {
	return Role{
		Status: customtypes.Enable,
	}
}
