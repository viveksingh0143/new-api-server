package container

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
)

type ContainerCreateDto struct {
	ContainerType customtypes.ContainerType `json:"container_type" validate:"ContainerType=PALLET BIN RACK"`
	Code          string                    `json:"code" validate:"required"`
	Name          customtypes.NullString    `json:"name" validate:"required"`
	Address       string                    `json:"address"`
	Status        customtypes.StatusEnum    `json:"status"`
	LastUpdatedBy customtypes.NullString    `json:"last_updated_by"`
	Store         *store.StoreMinimalDto    `json:"store"`
}
