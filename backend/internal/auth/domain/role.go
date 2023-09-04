package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/config/types/data"
	"github.com/vamika-digital/wms-api-server/config/types/status"
)

type Role struct {
	ID            int64             `json:"id"`
	Name          string            `json:"name"`
	Status        status.StatusType `json:"status"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     data.NullTime     `json:"updated_at"`
	LastUpdatedBy data.NullString   `json:"last_updated_by"`
}

func NewRoleWithDefaults() User {
	return User{
		Status:    status.Enabled,
		CreatedAt: time.Now(),
	}
}
