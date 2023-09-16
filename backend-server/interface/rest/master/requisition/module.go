package requisition

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type RequisitionRestModule struct {
	Handler *RequisitionRestHandler
}

func NewRequisitionRestModule(requisitionService service.RequisitionService) *RequisitionRestModule {
	requisitionHandler := NewRequisitionHandler(requisitionService)
	return &RequisitionRestModule{Handler: requisitionHandler}
}

func (m *RequisitionRestModule) RegisterRoutes(r *gin.RouterGroup) {
	requisitionGroup := r.Group("/requisitions")
	{
		requisitionGroup.POST("", m.Handler.CreateRequisition)
		requisitionGroup.GET("", m.Handler.GetAllRequisitions)
		requisitionGroup.POST("/bulkdelete", m.Handler.DeleteRequisitionByIDs)
		requisitionGroup.GET("/:id", m.Handler.GetRequisitionByID)
		requisitionGroup.PUT("/:id", m.Handler.UpdateRequisition)
		requisitionGroup.DELETE("/:id", m.Handler.DeleteRequisition)
	}
}
