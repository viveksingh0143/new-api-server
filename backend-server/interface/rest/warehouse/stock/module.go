package stock

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/service"
)

type StockRestModule struct {
	Handler *StockRestHandler
}

func NewStockRestModule(stock service.StockService) *StockRestModule {
	labelstickerHandler := NewStockHandler(stock)
	return &StockRestModule{Handler: labelstickerHandler}
}

func (m *StockRestModule) RegisterRoutes(r *gin.RouterGroup) {
	labelstickerGroup := r.Group("/stocks")
	{
		labelstickerGroup.GET("", m.Handler.GetAllInventories)
		labelstickerGroup.GET("/:id", m.Handler.GetStockByID)
	}
}
