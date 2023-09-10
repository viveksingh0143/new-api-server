package machine

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type MachineRestModule struct {
	Handler *MachineRestHandler
}

func NewMachineRestModule(machineService service.MachineService) *MachineRestModule {
	machineHandler := NewMachineHandler(machineService)
	return &MachineRestModule{Handler: machineHandler}
}

func (m *MachineRestModule) RegisterRoutes(r *gin.RouterGroup) {
	machineGroup := r.Group("/machines")
	{
		machineGroup.POST("", m.Handler.CreateMachine)
		machineGroup.GET("", m.Handler.GetAllMachines)
		machineGroup.POST("/bulkdelete", m.Handler.DeleteMachineByIDs)
		machineGroup.GET("/:id", m.Handler.GetMachineByID)
		machineGroup.PUT("/:id", m.Handler.UpdateMachine)
		machineGroup.DELETE("/:id", m.Handler.DeleteMachine)
	}
}
