package master

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/customer"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/machine"
)

type MasterRestModule struct {
	MachineModule  *machine.MachineRestModule
	CustomerModule *customer.CustomerRestModule
}

func NewMasterRestModule(machineService service.MachineService, customerService service.CustomerService) *MasterRestModule {
	machineModule := machine.NewMachineRestModule(machineService)
	customerModule := customer.NewCustomerRestModule(customerService)
	return &MasterRestModule{MachineModule: machineModule, CustomerModule: customerModule}
}

func (m *MasterRestModule) RegisterRoutes(r *gin.RouterGroup) {
	masterGroup := r.Group("/master")
	m.MachineModule.RegisterRoutes(masterGroup)
	m.CustomerModule.RegisterRoutes(masterGroup)
}
