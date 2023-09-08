package user

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/role"
)

type UserCreateDto struct {
	Name          string                 `json:"name" validate:"required"`
	Username      customtypes.NullString `json:"username"`
	Password      string                 `json:"password" validate:"required"`
	StaffID       customtypes.NullString `json:"staff_id"`
	EMail         customtypes.NullString `json:"email" validate:"required"`
	Status        customtypes.StatusEnum `json:"status"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
	Roles         []role.RoleMinimalDto  `json:"roles"`
}
