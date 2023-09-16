package stock

import (
	"github.com/gin-gonic/gin"
	masterService "github.com/vamika-digital/wms-api-server/core/business/master/service"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/service"
)

type StockRestModule struct {
	Handler *StockRestHandler
}

func NewStockRestModule(inventoryService service.InventoryService, requisitionService masterService.RequisitionService, outwardrequestService masterService.OutwardRequestService) *StockRestModule {
	containerHandler := NewStockHandler(inventoryService, requisitionService, outwardrequestService)
	return &StockRestModule{Handler: containerHandler}
}

func (m *StockRestModule) RegisterRoutes(r *gin.RouterGroup) {
	containerGroup := r.Group("/inventories")
	{
		containerGroup.POST("raw-material", m.Handler.CreateRawMaterialStock)
		containerGroup.POST("finished-goods", m.Handler.CreateFinishedStocks)
		containerGroup.POST("finished-good", m.Handler.CreateFinishedStock)
		containerGroup.POST("attach-container", m.Handler.AttachContainer)
		containerGroup.GET("container-stocks", m.Handler.GetAllContainerStocks)
		containerGroup.GET("find-requisition", m.Handler.GetRequisitionByCode)
	}
}
