package container

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ContainerCodeInfoDto struct {
	ContainerType customtypes.ContainerType `db:"container_type" json:"container_type"`
	Code          string                    `db:"code" json:"code"`
}
