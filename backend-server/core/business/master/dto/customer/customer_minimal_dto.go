package customer

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type CustomerMinimalDto struct {
	ID     int64                  `json:"id"`
	Code   string                 `json:"code"`
	Name   string                 `json:"name"`
	Status customtypes.StatusEnum `json:"status"`
}
