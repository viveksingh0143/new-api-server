package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/common/types"
)

type Role struct {
	ID            int64            `json:"id"`
	Name          string           `json:"name"`
	Status        types.StatusEnum `json:"status"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     types.NullTime   `json:"updated_at"`
	LastUpdatedBy types.NullString `json:"last_updated_by"`
}

func NewRoleWithDefaults() Role {
	return Role{
		Status:    types.EnabledStatus,
		CreatedAt: time.Now(),
	}
}
