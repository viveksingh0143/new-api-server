package role

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/permission"
)

type RoleCreateDto struct {
	Name          customtypes.NullString     `json:"name" validate:"required"`
	Status        customtypes.NullStatusEnum `json:"status" validate:"required"`
	LastUpdatedBy customtypes.NullString     `json:"last_updated_by"`
	Permissions   []permission.PermissionDto `json:"permissions"`
}
