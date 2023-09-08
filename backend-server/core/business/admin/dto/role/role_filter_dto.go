package role

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type RoleFilterDto struct {
	Query  string                 `json:"query"`
	ID     int64                  `json:"id"`
	Name   string                 `json:"name"`
	Status customtypes.StatusEnum `json:"status"`
}
