package user

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type UserMinimalDto struct {
	ID      int64                  `json:"id"`
	Name    string                 `json:"name"`
	StaffID customtypes.NullString `json:"staff_id"`
	EMail   customtypes.NullString `json:"email"`
	Status  customtypes.StatusEnum `json:"status"`
}
