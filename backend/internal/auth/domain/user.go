package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/common/types"
)

type User struct {
	ID                int64            `json:"id"`
	Username          types.NullString `json:"username"`
	Password          string           `json:"-"`
	Name              string           `json:"name"`
	StaffID           types.NullString `json:"staff_id"`
	Email             string           `json:"email"`
	EmailConfirmation bool             `json:"email_confirmation"`
	Status            types.StatusEnum `json:"status"`
	CreatedAt         time.Time        `json:"created_at"`
	UpdatedAt         types.NullTime   `json:"updated_at"`
	LastUpdatedBy     types.NullString `json:"last_updated_by"`
	Roles             []Role           `json:"roles"`
}

func NewUserWithDefaults() *User {
	return &User{
		EmailConfirmation: false,
		Status:            types.EnabledStatus,
		CreatedAt:         time.Now(),
	}
}
