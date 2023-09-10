package machine

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type MachineFilterDto struct {
	Query  string                 `json:"query"`
	ID     int64                  `json:"id"`
	Code   string                 `json:"code"`
	Name   string                 `json:"name"`
	Status customtypes.StatusEnum `json:"status"`
}
