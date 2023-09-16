package store

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
)

type StoreDto struct {
	ID            int64                  `json:"id"`
	Code          string                 `json:"code"`
	Name          string                 `json:"name"`
	Location      string                 `json:"location"`
	StoreTypes    []string               `json:"store_types"`
	Owner         *user.UserMinimalDto   `json:"owner"`
	Status        customtypes.StatusEnum `json:"status"`
	CreatedAt     customtypes.NullTime   `json:"created_at"`
	UpdatedAt     customtypes.NullTime   `json:"updated_at"`
	LastUpdatedBy customtypes.NullString `json:"last_updated_by"`
}
