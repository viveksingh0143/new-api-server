package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	adminDomain "github.com/vamika-digital/wms-api-server/core/business/admin/domain"
)

type Store struct {
	ID            int64                  `db:"id" json:"id"`
	Code          string                 `db:"code" json:"code"`
	Name          string                 `db:"name" json:"name"`
	Location      string                 `db:"location" json:"location"`
	OwnerID       int64                  `db:"owner_id" json:"owner_id"`
	Owner         *adminDomain.User      `db:"_" json:"owner"`
	Status        customtypes.StatusEnum `db:"status" json:"status"`
	CreatedAt     time.Time              `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time             `db:"updated_at" json:"updated_at"`
	LastUpdatedBy customtypes.NullString `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt     *time.Time             `db:"deleted_at" json:"deleted_at,omitempty"`
}

func NewStoreWithDefaults() Store {
	return Store{
		Status: customtypes.Enable,
	}
}
