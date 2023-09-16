package repository

import (
	masterDomain "github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/inventory"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/reports"
)

type InventoryRepository interface {
	GetById(productID int64) (*domain.Inventory, error)
	GetTotalCount(filter *inventory.InventoryFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *inventory.InventoryFilterDto) ([]*domain.Inventory, error)
	CreateRawMaterialStock(*domain.Stock, *masterDomain.Container) error
	CreateFinishedStocks([]*domain.Stock, []*domain.LabelSticker, *masterDomain.Container) error
	AttachContainer(*masterDomain.Container, *masterDomain.Container) error
	GetInventoryDetailForProductIds(productIds []int64) ([]*reports.InventoryStatusDetail, error)
}
