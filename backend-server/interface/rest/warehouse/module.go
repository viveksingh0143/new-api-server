package warehouse

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/service"
	"github.com/vamika-digital/wms-api-server/interface/rest/warehouse/batchlabel"
	"github.com/vamika-digital/wms-api-server/interface/rest/warehouse/inventory"
	"github.com/vamika-digital/wms-api-server/interface/rest/warehouse/labelsticker"
	"github.com/vamika-digital/wms-api-server/interface/rest/warehouse/stock"
)

type WarehouseRestModule struct {
	BatchLabelModule    *batchlabel.BatchLabelRestModule
	LabelStickerModule  *labelsticker.LabelStickerRestModule
	InventoryRestModule *inventory.InventoryRestModule
	StockRestModule     *stock.StockRestModule
}

func NewWarehouseRestModule(
	batchlabelService service.BatchLabelService,
	labelstickerService service.LabelStickerService,
	inventoryService service.InventoryService,
	stockService service.StockService,
) *WarehouseRestModule {
	batchlabelModule := batchlabel.NewBatchLabelRestModule(batchlabelService)
	labelstickerModule := labelsticker.NewLabelStickerRestModule(labelstickerService)
	inventoryModule := inventory.NewInventoryRestModule(inventoryService)
	stockModule := stock.NewStockRestModule(stockService)
	return &WarehouseRestModule{BatchLabelModule: batchlabelModule, LabelStickerModule: labelstickerModule, InventoryRestModule: inventoryModule, StockRestModule: stockModule}
}

func (m *WarehouseRestModule) RegisterRoutes(r *gin.RouterGroup) {
	warehouseGroup := r.Group("/warehouse")
	m.BatchLabelModule.RegisterRoutes(warehouseGroup)
	m.LabelStickerModule.RegisterRoutes(warehouseGroup)
	m.InventoryRestModule.RegisterRoutes(warehouseGroup)
	m.StockRestModule.RegisterRoutes(warehouseGroup)
}
