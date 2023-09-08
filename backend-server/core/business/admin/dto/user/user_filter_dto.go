package user

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type UserFilterDto struct {
	Query    string                 `json:"query"`
	ID       int64                  `json:"id"`
	Name     string                 `json:"name"`
	Username string                 `json:"username"`
	StaffID  string                 `json:"staff_id"`
	EMail    string                 `json:"email"`
	Status   customtypes.StatusEnum `json:"status"`
}
