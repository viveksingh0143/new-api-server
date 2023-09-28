package outwardrequest

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type OutwardRequestRestModule struct {
	Handler *OutwardRequestRestHandler
}

func NewOutwardRequestRestModule(outwardrequestService service.OutwardRequestService) *OutwardRequestRestModule {
	outwardrequestHandler := NewOutwardRequestHandler(outwardrequestService)
	return &OutwardRequestRestModule{Handler: outwardrequestHandler}
}

func (m *OutwardRequestRestModule) RegisterRoutes(r *gin.RouterGroup) {
	outwardrequestGroup := r.Group("/outwardrequests")
	{
		outwardrequestGroup.POST("", m.Handler.CreateOutwardRequest)
		outwardrequestGroup.GET("", m.Handler.GetAllOutwardRequests)
		outwardrequestGroup.POST("/bulkdelete", m.Handler.DeleteOutwardRequestByIDs)
		outwardrequestGroup.GET("/:id", m.Handler.GetOutwardRequestByID)
		outwardrequestGroup.PUT("/:id", m.Handler.UpdateOutwardRequest)
		outwardrequestGroup.DELETE("/:id", m.Handler.DeleteOutwardRequest)
		outwardrequestGroup.GET("/:id/shipperlabels", m.Handler.GetShipperLabelsByID)
		outwardrequestGroup.POST("/:id/shipperlabels", m.Handler.GenerateShipperLabelsByID)
	}
}
