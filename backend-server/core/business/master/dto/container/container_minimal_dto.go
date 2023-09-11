package container

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ContainerMinimalDto struct {
	ID            int64                     `json:"id"`
	ContainerType customtypes.ContainerType `json:"container_type"`
	Code          string                    `json:"code"`
	Name          string                    `json:"name"`
	Status        customtypes.StatusEnum    `json:"status"`
}
