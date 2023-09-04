package dto

import (
	"github.com/vamika-digital/wms-api-server/config/types/data"
	"github.com/vamika-digital/wms-api-server/config/types/status"
)

type UserCreateDto struct {
	Username          data.NullString   `json:"username"`
	Password          string            `json:"-"`
	Name              string            `json:"name"`
	StaffID           data.NullString   `json:"staff_id"`
	Email             string            `json:"email"`
	EmailConfirmation data.NullBool     `json:"email_confirmation"`
	Status            status.StatusType `json:"status"`
	Roles             []RoleDto         `json:"roles"`
}
