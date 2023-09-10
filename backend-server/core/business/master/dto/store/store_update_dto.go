package store

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
)

type StoreUpdateDto struct {
	Code          string                 `json:"code" validate:"required"`
	Name          string                 `json:"name" validate:"required"`
	Location      string                 `json:"location"`
	Status        customtypes.StatusEnum `json:"status" validate:"required"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
	Owner         *user.UserMinimalDto   `json:"owner"`
}
