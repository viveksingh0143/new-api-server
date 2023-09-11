package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type Container struct {
	ID            int64                     `db:"id" json:"id"`
	ContainerType customtypes.ContainerType `db:"container_type" json:"container_type"`
	Code          string                    `db:"code" json:"code"`
	Name          string                    `db:"name" json:"name"`
	Address       string                    `db:"address" json:"address"`
	Status        customtypes.StatusEnum    `db:"status" json:"status"`
	CreatedAt     time.Time                 `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time                `db:"updated_at" json:"updated_at"`
	LastUpdatedBy customtypes.NullString    `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt     *time.Time                `db:"deleted_at" json:"deleted_at,omitempty"`
}

func NewContainerWithDefaults() Container {
	return Container{
		Status: customtypes.Enable,
	}
}
