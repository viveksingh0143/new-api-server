package machine

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type MachineUpdateDto struct {
	Code          string                 `json:"code" validate:"required"`
	Name          string                 `json:"name" validate:"required"`
	Status        customtypes.StatusEnum `json:"status"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
}
