package user

import (
	"github.com/vamika-digital/wms-api-server/common/types"
	"github.com/vamika-digital/wms-api-server/internal/auth/dto/role"
)

type UserUpdateDto struct {
	Username          types.NullString `json:"username"`
	Password          types.NullString `json:"-"`
	Name              types.NullString `json:"name"`
	StaffID           types.NullString `json:"staff_id"`
	Email             types.NullString `json:"email"`
	EmailConfirmation types.NullBool   `json:"email_confirmation"`
	Status            types.StatusEnum `json:"status"`
	LastUpdatedBy     types.NullString `json:"last_updated_by"`
	Roles             []role.RoleDto   `json:"roles"`
}
