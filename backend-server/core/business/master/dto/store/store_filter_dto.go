package store

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
)

type StoreFilterDto struct {
	Query  string                 `json:"query"`
	ID     int64                  `json:"id"`
	Code   string                 `json:"code"`
	Name   string                 `json:"name"`
	Status customtypes.StatusEnum `json:"status"`
	Owner  *user.UserMinimalDto   `json:"owner"`
}
