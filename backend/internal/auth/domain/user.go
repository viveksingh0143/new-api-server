package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/config/types/data"
	"github.com/vamika-digital/wms-api-server/config/types/status"
)

type User struct {
	ID                int64             `json:"id"`
	Username          data.NullString   `json:"username"`
	Password          string            `json:"-"`
	Name              string            `json:"name"`
	StaffID           data.NullString   `json:"staff_id"`
	Email             string            `json:"email"`
	EmailConfirmation bool              `json:"email_confirmation"`
	Status            status.StatusType `json:"status"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         data.NullTime     `json:"updated_at"`
	LastUpdatedBy     data.NullString   `json:"last_updated_by"`
	Roles             []Role            `json:"roles"`
}

func NewUserWithDefaults() *User {
	return &User{
		EmailConfirmation: false,
		Status:            status.Enabled,
		CreatedAt:         time.Now(),
	}
}
