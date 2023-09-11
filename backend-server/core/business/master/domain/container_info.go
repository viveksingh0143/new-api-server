package domain

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ContainerInfo struct {
	Type        customtypes.ContainerType     `db:"type" json:"type"`
	MinCapacity int64                         `db:"min_capacity" json:"min_capacity"`
	MaxCapacity int64                         `db:"max_capacity" json:"max_capacity"`
	CanContains []customtypes.ContainableType `db:"can_contains" json:"can_contains"`
}

var PalletContainerInfo ContainerInfo = ContainerInfo{
	Type:        customtypes.PALLET_TYPE,
	MinCapacity: 0,
	MaxCapacity: 0,
	CanContains: []customtypes.ContainableType{customtypes.RAW_MATERIAL_CONTAINABLE, customtypes.FINISHED_GOOD_CONTAINABLE},
}
var BinContainerInfo ContainerInfo = ContainerInfo{
	Type:        customtypes.BIN_TYPE,
	MinCapacity: 0,
	MaxCapacity: 0,
	CanContains: []customtypes.ContainableType{customtypes.PALLET_CONTAINABLE, customtypes.RAW_MATERIAL_CONTAINABLE, customtypes.FINISHED_GOOD_CONTAINABLE},
}
var RackContainerInfo ContainerInfo = ContainerInfo{
	Type:        customtypes.RACK_TYPE,
	MinCapacity: 0,
	MaxCapacity: 0,
	CanContains: []customtypes.ContainableType{customtypes.BIN_CONTAINABLE},
}
