package role

import (
	"github.com/vamika-digital/wms-api-server/common/types"
)

type RoleDto struct {
	ID            int64            `json:"id"`
	Name          types.NullString `json:"name"`
	Status        types.StatusEnum `json:"status"`
	CreatedAt     types.NullTime   `json:"created_at"`
	UpdatedAt     types.NullTime   `json:"updated_at"`
	LastUpdatedBy types.NullString `json:"last_updated_by"`
}
