package joborder

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type JobOrderRestModule struct {
	Handler *JobOrderRestHandler
}

func NewJobOrderRestModule(joborderService service.JobOrderService) *JobOrderRestModule {
	joborderHandler := NewJobOrderHandler(joborderService)
	return &JobOrderRestModule{Handler: joborderHandler}
}

func (m *JobOrderRestModule) RegisterRoutes(r *gin.RouterGroup) {
	joborderGroup := r.Group("/joborders")
	{
		joborderGroup.POST("", m.Handler.CreateJobOrder)
		joborderGroup.GET("", m.Handler.GetAllJobOrders)
		joborderGroup.POST("/bulkdelete", m.Handler.DeleteJobOrderByIDs)
		joborderGroup.GET("/:id", m.Handler.GetJobOrderByID)
		joborderGroup.PUT("/:id", m.Handler.UpdateJobOrder)
		joborderGroup.DELETE("/:id", m.Handler.DeleteJobOrder)
	}
}
