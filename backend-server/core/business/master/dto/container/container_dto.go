package container

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ContainerDto struct {
	ID            int64                         `json:"id"`
	ContainerType customtypes.ContainerType     `json:"container_type"`
	Code          string                        `json:"code"`
	Name          string                        `json:"name"`
	Address       string                        `json:"address"`
	Status        customtypes.StatusEnum        `json:"status"`
	CreatedAt     customtypes.NullTime          `json:"created_at"`
	UpdatedAt     customtypes.NullTime          `json:"updated_at"`
	LastUpdatedBy customtypes.NullString        `json:"last_updated_by"`
	MinCapacity   int64                         `json:"min_capacity"`
	MaxCapacity   int64                         `json:"max_capacity"`
	CanContains   []customtypes.ContainableType `json:"can_contains"`
}
