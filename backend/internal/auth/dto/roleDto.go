package dto

import (
	"github.com/vamika-digital/wms-api-server/config/types/data"
	"github.com/vamika-digital/wms-api-server/config/types/status"
)

type RoleDto struct {
	ID            int64             `json:"id"`
	Name          data.NullString   `json:"name"`
	Status        status.StatusType `json:"status"`
	CreatedAt     data.NullTime     `json:"created_at"`
	UpdatedAt     data.NullTime     `json:"updated_at"`
	LastUpdatedBy data.NullString   `json:"last_updated_by"`
}
