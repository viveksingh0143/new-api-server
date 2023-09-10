package role

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/permission"
)

type RoleDto struct {
	ID            int64                       `json:"id"`
	Name          string                      `json:"name"`
	Status        customtypes.StatusEnum      `json:"status"`
	CreatedAt     customtypes.NullTime        `json:"created_at"`
	UpdatedAt     customtypes.NullTime        `json:"updated_at"`
	LastUpdatedBy customtypes.NullString      `json:"last_updated_by"`
	Permissions   []*permission.PermissionDto `json:"permissions" validate:"required"`
}
