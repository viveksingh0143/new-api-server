package user

import (
	"github.com/vamika-digital/wms-api-server/common/types"
	"github.com/vamika-digital/wms-api-server/internal/auth/dto/role"
)

type UserCreateDto struct {
	Username          types.NullString `json:"username"`
	Password          string           `json:"-"`
	Name              string           `json:"name"`
	StaffID           types.NullString `json:"staff_id"`
	Email             string           `json:"email"`
	EmailConfirmation types.NullBool   `json:"email_confirmation"`
	Status            types.StatusEnum `json:"status"`
	Roles             []role.RoleDto   `json:"roles"`
}
