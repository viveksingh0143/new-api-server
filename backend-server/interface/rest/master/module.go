package master

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/container"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/customer"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/joborder"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/machine"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/outwardrequest"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/product"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/requisition"
	"github.com/vamika-digital/wms-api-server/interface/rest/master/store"
)

type MasterRestModule struct {
	MachineModule        *machine.MachineRestModule
	CustomerModule       *customer.CustomerRestModule
	StoreModule          *store.StoreRestModule
	ContainerModule      *container.ContainerRestModule
	ProductModule        *product.ProductRestModule
	JobOrderModule       *joborder.JobOrderRestModule
	OutwardRequestModule *outwardrequest.OutwardRequestRestModule
	RequisitionModule    *requisition.RequisitionRestModule
}

func NewMasterRestModule(
	machineService service.MachineService,
	customerService service.CustomerService,
	storeService service.StoreService,
	containerService service.ContainerService,
	productService service.ProductService,
	joborderService service.JobOrderService,
	outwardrequestService service.OutwardRequestService,
	requisitionService service.RequisitionService,
) *MasterRestModule {
	machineModule := machine.NewMachineRestModule(machineService)
	customerModule := customer.NewCustomerRestModule(customerService)
	storeModule := store.NewStoreRestModule(storeService)
	containerModule := container.NewContainerRestModule(containerService)
	productModule := product.NewProductRestModule(productService)
	joborderModule := joborder.NewJobOrderRestModule(joborderService)
	outwardrequestModule := outwardrequest.NewOutwardRequestRestModule(outwardrequestService)
	requisitionModule := requisition.NewRequisitionRestModule(requisitionService)
	return &MasterRestModule{MachineModule: machineModule, CustomerModule: customerModule, StoreModule: storeModule, ContainerModule: containerModule, ProductModule: productModule, JobOrderModule: joborderModule, OutwardRequestModule: outwardrequestModule, RequisitionModule: requisitionModule}
}

func (m *MasterRestModule) RegisterRoutes(r *gin.RouterGroup) {
	masterGroup := r.Group("/master")
	m.MachineModule.RegisterRoutes(masterGroup)
	m.CustomerModule.RegisterRoutes(masterGroup)
	m.StoreModule.RegisterRoutes(masterGroup)
	m.ContainerModule.RegisterRoutes(masterGroup)
	m.ProductModule.RegisterRoutes(masterGroup)
	m.JobOrderModule.RegisterRoutes(masterGroup)
	m.OutwardRequestModule.RegisterRoutes(masterGroup)
	m.RequisitionModule.RegisterRoutes(masterGroup)
}
