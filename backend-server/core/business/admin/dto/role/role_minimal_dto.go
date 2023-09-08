package role

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type RoleMinimalDto struct {
	ID     int64                  `json:"id"`
	Name   string                 `json:"name"`
	Status customtypes.StatusEnum `json:"status"`
}
