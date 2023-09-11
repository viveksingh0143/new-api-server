package container

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ContainerFilterDto struct {
	Query         string                     `form:"query" json:"query"`
	ID            int64                      `form:"id" json:"id"`
	ContainerType *customtypes.ContainerType `form:"container_type" binding:"oneof=PALLET BIN RACK" json:"container_type"`
	Code          string                     `form:"code" json:"code"`
	Name          string                     `form:"name" json:"name"`
	Status        customtypes.StatusEnum     `form:"status" json:"status"`
}
