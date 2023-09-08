package role

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/permission"
)

type RoleUpdateDto struct {
	Name          string                     `json:"name"`
	Status        customtypes.StatusEnum     `json:"status"`
	LastUpdatedBy customtypes.NullString     `json:"last_updated_by"`
	Permissions   []permission.PermissionDto `json:"permissions"`
}
