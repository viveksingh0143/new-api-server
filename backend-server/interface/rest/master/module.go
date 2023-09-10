package master

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/customer"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/machine"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/store"
)

type MasterRestModule struct {
	MachineModule  *machine.MachineRestModule
	CustomerModule *customer.CustomerRestModule
	StoreModule    *store.StoreRestModule
}

func NewMasterRestModule(
	machineService service.MachineService,
	customerService service.CustomerService,
	storeService service.StoreService,
) *MasterRestModule {
	machineModule := machine.NewMachineRestModule(machineService)
	customerModule := customer.NewCustomerRestModule(customerService)
	storeModule := store.NewStoreRestModule(storeService)
	return &MasterRestModule{MachineModule: machineModule, CustomerModule: customerModule, StoreModule: storeModule}
}

func (m *MasterRestModule) RegisterRoutes(r *gin.RouterGroup) {
	masterGroup := r.Group("/master")
	m.MachineModule.RegisterRoutes(masterGroup)
	m.CustomerModule.RegisterRoutes(masterGroup)
	m.StoreModule.RegisterRoutes(masterGroup)
}
