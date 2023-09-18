package container

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ContainerFilterDto struct {
	Query         string                     `form:"query" json:"query"`
	ID            int64                      `form:"id" json:"id"`
	StoreID       int64                      `form:"store_id" json:"store_id"`
	ContainerType *customtypes.ContainerType `form:"container_type" binding:"omitempty,oneof=PALLET BIN RACK" json:"container_type"`
	Code          string                     `form:"code" json:"code"`
	Name          string                     `form:"name" json:"name"`
	Status        customtypes.StatusEnum     `form:"status" json:"status"`
}
