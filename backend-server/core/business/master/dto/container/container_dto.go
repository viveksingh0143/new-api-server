package container

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
)

type ContainerDto struct {
	ID               int64                         `json:"id"`
	ContainerType    customtypes.ContainerType     `json:"container_type"`
	Code             string                        `json:"code"`
	Name             string                        `json:"name"`
	Address          string                        `json:"address"`
	IsApproved       bool                          `json:"approved"`
	Status           customtypes.StatusEnum        `json:"status"`
	CreatedAt        customtypes.NullTime          `json:"created_at"`
	UpdatedAt        customtypes.NullTime          `json:"updated_at"`
	LastUpdatedBy    customtypes.NullString        `json:"last_updated_by"`
	MinCapacity      int64                         `json:"min_capacity"`
	MaxCapacity      int64                         `json:"max_capacity"`
	CanContains      []customtypes.ContainableType `json:"can_contains"`
	Level            customtypes.StockLevel        `json:"stock_level"`
	StoreID          customtypes.NullInt64         `json:"store_id"`
	Store            *store.StoreMinimalDto        `json:"store"`
	OtherInfo        *domain.ContainerInfo         `json:"other_info"`
	ResourceID       customtypes.NullInt64         `json:"resource_id"`
	ResourceName     customtypes.NullString        `json:"resource_name"`
	ItemsCount       int64                         `json:"items_count"`
	ContainerItemDto *ContainerItemDto             `json:"item_info"`
}
