package mobile

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
	warehouseService "github.com/vamika-digital/wms-api-server/core/business/warehouse/service"
	"github.com/vamika-digital/wms-api-server/interface/rest/mobile/container"
	"github.com/vamika-digital/wms-api-server/interface/rest/mobile/stock"
)

type MobileRestModule struct {
	ContainerModule *container.ContainerRestModule
	StockModule     *stock.StockRestModule
}

func NewMobileRestModule(
	containerService service.ContainerService,
	inventoryService warehouseService.InventoryService,
	requisitionService service.RequisitionService,
	outwardrequestService service.OutwardRequestService,
) *MobileRestModule {
	containerModule := container.NewContainerRestModule(containerService)
	stockModule := stock.NewStockRestModule(inventoryService, requisitionService, outwardrequestService)
	return &MobileRestModule{ContainerModule: containerModule, StockModule: stockModule}
}

func (m *MobileRestModule) RegisterRoutes(r *gin.RouterGroup) {
	deviceGroup := r.Group("/device/v1")
	m.ContainerModule.RegisterRoutes(deviceGroup)
	m.StockModule.RegisterRoutes(deviceGroup)
}
