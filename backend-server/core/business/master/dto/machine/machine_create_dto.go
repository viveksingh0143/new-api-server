package machine

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type MachineCreateDto struct {
	Code          string                 `json:"code" validate:"required"`
	Name          customtypes.NullString `json:"name" validate:"required"`
	Status        customtypes.StatusEnum `json:"status"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
}
