package user

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/role"
)

type UserDto struct {
	ID             int64                  `json:"id"`
	Name           string                 `json:"name"`
	Username       customtypes.NullString `json:"username"`
	StaffID        customtypes.NullString `json:"staff_id"`
	EMail          customtypes.NullString `json:"email"`
	EMailConfirmed bool                   `json:"email_confirmed"`
	Status         customtypes.StatusEnum `json:"status"`
	CreatedAt      customtypes.NullTime   `json:"created_at"`
	UpdatedAt      customtypes.NullTime   `json:"updated_at"`
	LastUpdatedBy  customtypes.NullString `json:"last_updated_by"`
	Roles          []role.RoleMinimalDto  `json:"roles"`
}
