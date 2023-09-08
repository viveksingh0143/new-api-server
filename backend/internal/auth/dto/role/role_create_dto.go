package role

import "github.com/vamika-digital/wms-api-server/common/types"

type RoleCreateDto struct {
	Name   types.NullString     `json:"name" validate:"required"`
	Status types.NullStatusEnum `json:"status" validate:"required"`
}
