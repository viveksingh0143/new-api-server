package customer

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type CustomerFilterDto struct {
	Query         string                 `json:"query"`
	ID            int64                  `json:"id"`
	Code          string                 `json:"code"`
	Name          string                 `json:"name"`
	ContactPerson string                 `json:"contact_person"`
	Status        customtypes.StatusEnum `json:"status"`
}
