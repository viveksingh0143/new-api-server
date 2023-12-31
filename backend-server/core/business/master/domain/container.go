package domain

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
)

type Container struct {
	ID            int64                     `db:"id" json:"id"`
	ContainerType customtypes.ContainerType `db:"container_type" json:"container_type"`
	Code          string                    `db:"code" json:"code"`
	Name          string                    `db:"name" json:"name"`
	Address       string                    `db:"address" json:"address"`
	IsApproved    bool                      `db:"approved" json:"approved"`
	Status        customtypes.StatusEnum    `db:"status" json:"status"`
	CreatedAt     time.Time                 `db:"created_at" json:"created_at"`
	UpdatedAt     *time.Time                `db:"updated_at" json:"updated_at"`
	LastUpdatedBy customtypes.NullString    `db:"last_updated_by" json:"last_updated_by"`
	DeletedAt     *time.Time                `db:"deleted_at" json:"deleted_at,omitempty"`
	Level         customtypes.StockLevel    `db:"stock_level" json:"stock_level"`
	ResourceID    customtypes.NullInt64     `db:"resource_id" json:"resource_id"`
	ResourceName  customtypes.NullString    `db:"resource_name" json:"resource_name"`
	ItemsCount    int64                     `db:"items_count" json:"items_count"`
	StoreID       customtypes.NullInt64     `db:"store_id" json:"store_id"`
	Store         *Store                    `db:"store" json:"store"`
}

func NewContainerWithDefaults() Container {
	return Container{
		Status: customtypes.Enable,
	}
}

func (c *Container) IncreamentStock(quantity int64, rID int64, rName string) {
	c.ItemsCount = c.ItemsCount + quantity
	if c.Info().MaxCapacity != -1 && c.ItemsCount >= c.Info().MaxCapacity {
		c.Level = customtypes.FULL_STOCK
	} else {
		c.Level = customtypes.PARTIAL_STOCK
	}
	c.ResourceID = customtypes.NewValidNullInt64(rID)
	c.ResourceName = customtypes.NewValidNullString(rName)
}

func (c *Container) DecreamentStock(quantity int64) {
	c.ItemsCount = c.ItemsCount - quantity
	if c.ItemsCount <= 0 {
		c.Level = customtypes.EMPTY_STOCK
		c.ResourceID = customtypes.NewInvalidNullInt64()
		c.ResourceName = customtypes.NewInvalidNullString()
	} else {
		c.Level = customtypes.PARTIAL_STOCK
	}
}

func (c *Container) Clear() {
	c.ItemsCount = 0
	c.Level = customtypes.EMPTY_STOCK
	c.ResourceID = customtypes.NewInvalidNullInt64()
	c.ResourceName = customtypes.NewInvalidNullString()
}

func (c *Container) IsFull() bool {
	if c.Level == customtypes.FULL_STOCK {
		return true
	}
	cInfo := GetContainerInfo(c.ContainerType)
	return c.Info().MaxCapacity != -1 && c.ItemsCount >= cInfo.MaxCapacity
}

func (c *Container) Info() *ContainerInfo {
	return GetContainerInfo(c.ContainerType)
}

func (c *Container) IsResourceMatch(rID int64, rName string) bool {
	return c.ResourceID.Int64 == rID && c.ResourceName.String == rName
}

func (c *Container) IsConnectedWithResource() bool {
	return c.ResourceID.Valid && c.ResourceName.Valid
}
