package inventory

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/service"
)

type InventoryRestModule struct {
	Handler *InventoryRestHandler
}

func NewInventoryRestModule(inventory service.InventoryService) *InventoryRestModule {
	labelstickerHandler := NewInventoryHandler(inventory)
	return &InventoryRestModule{Handler: labelstickerHandler}
}

func (m *InventoryRestModule) RegisterRoutes(r *gin.RouterGroup) {
	labelstickerGroup := r.Group("/inventories")
	{
		labelstickerGroup.GET("", m.Handler.GetAllInventories)
		labelstickerGroup.GET("/:id", m.Handler.GetInventoryByID)
		labelstickerGroup.POST("/raw-material", m.Handler.CreateRawMaterialStock)
	}
}
