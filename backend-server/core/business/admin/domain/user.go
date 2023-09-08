package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type User struct {
	ID                  int64                  `db:"id" json:"id"`
	Name                string                 `db:"name" json:"name"`
	Username            string                 `db:"username" json:"username"`
	Password            string                 `db:"password" json:"password"`
	StaffID             string                 `db:"staff_id" json:"staff_id"`
	EMail               string                 `db:"email" json:"email"`
	EMailConfirmationAt *time.Time             `db:"email_confirmation_at" json:"email_confirmation_at,omitempty"`
	Status              customtypes.StatusEnum `db:"status" json:"status"`
	CreatedAt           time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt           *time.Time             `db:"updated_at" json:"updated_at"`
	LastUpdatedBy       customtypes.NullString `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt           *time.Time             `db:"deleted_at" json:"deleted_at,omitempty"`
	Roles               []*Role                `json:"roles"`
}

func NewUserWithDefaults() User {
	return User{
		Status: customtypes.Enable,
	}
}
