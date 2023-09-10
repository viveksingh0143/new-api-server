package store

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type StoreRestModule struct {
	Handler *StoreRestHandler
}

func NewStoreRestModule(storeService service.StoreService) *StoreRestModule {
	storeHandler := NewStoreHandler(storeService)
	return &StoreRestModule{Handler: storeHandler}
}

func (m *StoreRestModule) RegisterRoutes(r *gin.RouterGroup) {
	storeGroup := r.Group("/stores")
	{
		storeGroup.POST("", m.Handler.CreateStore)
		storeGroup.GET("", m.Handler.GetAllStores)
		storeGroup.POST("/bulkdelete", m.Handler.DeleteStoreByIDs)
		storeGroup.GET("/:id", m.Handler.GetStoreByID)
		storeGroup.PUT("/:id", m.Handler.UpdateStore)
		storeGroup.DELETE("/:id", m.Handler.DeleteStore)
	}
}
