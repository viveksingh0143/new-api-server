package dto

import (
	"github.com/vamika-digital/wms-api-server/config/types/data"
	"github.com/vamika-digital/wms-api-server/config/types/status"
)

type UserUpdateDto struct {
	Username          data.NullString   `json:"username"`
	Password          data.NullString   `json:"-"`
	Name              data.NullString   `json:"name"`
	StaffID           data.NullString   `json:"staff_id"`
	Email             data.NullString   `json:"email"`
	EmailConfirmation data.NullBool     `json:"email_confirmation"`
	Status            status.StatusType `json:"status"`
	LastUpdatedBy     data.NullString   `json:"last_updated_by"`
	Roles             []RoleDto         `json:"roles"`
}
