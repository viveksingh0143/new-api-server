package role

import "github.com/vamika-digital/wms-api-server/common/types"

type RoleUpdateDto struct {
	Name   string           `json:"name"`
	Status types.StatusEnum `json:"status"`
}
