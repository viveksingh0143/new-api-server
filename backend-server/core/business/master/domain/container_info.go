package domain

import (
	"strings"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type ContainerInfo struct {
	Type        customtypes.ContainerType     `db:"type" json:"type"`
	MinCapacity int64                         `db:"min_capacity" json:"min_capacity"`
	MaxCapacity int64                         `db:"max_capacity" json:"max_capacity"`
	CanContains []customtypes.ContainableType `db:"can_contains" json:"can_contains"`
}

func (c *ContainerInfo) ContainsType(cType string) bool {
	uCType := strings.ToUpper(cType)
	var containerType customtypes.ContainableType
	switch uCType {
	case "PALLET":
		containerType = customtypes.PALLET_CONTAINABLE
	case "BIN":
		containerType = customtypes.BIN_CONTAINABLE
	case "RACK":
		containerType = customtypes.RACK_CONTAINABLE
	case "RAW MATERIAL":
		containerType = customtypes.RAW_MATERIAL_CONTAINABLE
	case "FINISHED_GOOD":
		containerType = customtypes.FINISHED_GOOD_CONTAINABLE
	}

	for _, item := range c.CanContains {
		if item == containerType {
			return true
		}
	}
	return false
}

func GetContainerInfo(containerType customtypes.ContainerType) *ContainerInfo {
	if containerType == customtypes.PALLET_TYPE {
		return &PalletContainerInfo
	} else if containerType == customtypes.BIN_TYPE {
		return &BinContainerInfo
	} else if containerType == customtypes.RACK_TYPE {
		return &RackContainerInfo
	}
	return nil
}

var PalletContainerInfo ContainerInfo = ContainerInfo{
	Type:        customtypes.PALLET_TYPE,
	MinCapacity: 1,
	MaxCapacity: 1,
	CanContains: []customtypes.ContainableType{customtypes.RAW_MATERIAL_CONTAINABLE},
}
var BinContainerInfo ContainerInfo = ContainerInfo{
	Type:        customtypes.BIN_TYPE,
	MinCapacity: -1,
	MaxCapacity: -1,
	CanContains: []customtypes.ContainableType{customtypes.FINISHED_GOOD_CONTAINABLE},
}
var RackContainerInfo ContainerInfo = ContainerInfo{
	Type:        customtypes.RACK_TYPE,
	MinCapacity: 1,
	MaxCapacity: 1,
	CanContains: []customtypes.ContainableType{customtypes.PALLET_CONTAINABLE, customtypes.BIN_CONTAINABLE},
}
