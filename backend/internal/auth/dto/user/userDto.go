package user

import (
	"github.com/vamika-digital/wms-api-server/common/types"
	"github.com/vamika-digital/wms-api-server/internal/auth/dto/role"
)

type UserDto struct {
	ID                int64            `json:"id"`
	Username          types.NullString `json:"username"`
	Name              string           `json:"name"`
	StaffID           types.NullString `json:"staff_id"`
	Email             string           `json:"email"`
	EmailConfirmation bool             `json:"email_confirmation"`
	Status            types.StatusEnum `json:"status"`
	CreatedAt         types.NullTime   `json:"created_at"`
	UpdatedAt         types.NullTime   `json:"updated_at"`
	LastUpdatedBy     types.NullString `json:"last_updated_by"`
	Roles             []role.RoleDto   `json:"roles"`
}
