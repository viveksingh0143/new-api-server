package service

import (
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/inventory"
)

type InventoryService interface {
	GetAllProductsWithStockCounts(page int16, pageSize int16, sort string, filter *inventory.InventoryFilterDto) ([]*inventory.InventoryDto, int64, error)
	GetInventoryByID(inventoryID int64) (*inventory.InventoryDto, error)
	CreateRawMaterialStock(rmStockForm *inventory.InventoryRMStockCreateDto) error
	CreateFinishedGoodsStock(fdStockForm *inventory.InventoryFDStockCreateDto) error
	AttachContainer(string, string) error
}
