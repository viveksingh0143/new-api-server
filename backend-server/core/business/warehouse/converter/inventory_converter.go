package converter

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/inventory"
)

type InventoryConverter struct{}

func NewInventoryConverter() *InventoryConverter {
	return &InventoryConverter{}
}

func (c *InventoryConverter) ToDto(domainInventory *domain.Inventory) *inventory.InventoryDto {
	inventoryDto := &inventory.InventoryDto{
		ID:             domainInventory.ID,
		ProductType:    domainInventory.ProductType,
		Code:           domainInventory.Code,
		Name:           domainInventory.Name,
		UnitType:       domainInventory.UnitType,
		UnitWeight:     domainInventory.UnitWeight,
		UnitWeightType: domainInventory.UnitWeightType,
		StockCount:     domainInventory.TotalStockCount,
		StockinAt:      customtypes.GetNullTime(domainInventory.LastStockinAt),
	}
	return inventoryDto
}

func (c *InventoryConverter) ToDtoSlice(domainInventorys []*domain.Inventory) []*inventory.InventoryDto {
	var inventoryDtos = make([]*inventory.InventoryDto, 0)
	for _, domainInventory := range domainInventorys {
		inventoryDtos = append(inventoryDtos, c.ToDto(domainInventory))
	}
	return inventoryDtos
}

func (c *InventoryConverter) ToStockDomain(rmStockForm *inventory.InventoryRMStockCreateDto) *domain.Stock {
	return &domain.Stock{
		ProductID:     rmStockForm.ProductID,
		StoreID:       rmStockForm.StoreID,
		Quantity:      rmStockForm.Quantity,
		StockInAt:     time.Now(),
		Status:        customtypes.STOCK_IN,
		LastUpdatedBy: rmStockForm.LastUpdatedBy,
	}
}
