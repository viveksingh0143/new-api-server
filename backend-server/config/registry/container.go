package registry

import (
	"log"
	"sync"

	adminConverter "github.com/vamika-digital/wms-api-server/core/business/admin/converter"
	adminRepository "github.com/vamika-digital/wms-api-server/core/business/admin/repository"
	adminService "github.com/vamika-digital/wms-api-server/core/business/admin/service"
	authService "github.com/vamika-digital/wms-api-server/core/business/auth/service"
	masterConverter "github.com/vamika-digital/wms-api-server/core/business/master/converter"
	masterRepository "github.com/vamika-digital/wms-api-server/core/business/master/repository"
	masterService "github.com/vamika-digital/wms-api-server/core/business/master/service"
	warehouseConverter "github.com/vamika-digital/wms-api-server/core/business/warehouse/converter"
	warehouseRepository "github.com/vamika-digital/wms-api-server/core/business/warehouse/repository"
	warehouseService "github.com/vamika-digital/wms-api-server/core/business/warehouse/service"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type Container struct {
	RoleService adminService.RoleService
	UserService adminService.UserService
	AuthService authService.AuthService

	MachineService        masterService.MachineService
	CustomerService       masterService.CustomerService
	StoreService          masterService.StoreService
	ContainerService      masterService.ContainerService
	ProductService        masterService.ProductService
	JobOrderService       masterService.JobOrderService
	OutwardRequestService masterService.OutwardRequestService
	RequisitionService    masterService.RequisitionService

	BatchLabelService   warehouseService.BatchLabelService
	LabelStickerService warehouseService.LabelStickerService
	InventoryService    warehouseService.InventoryService
	StockService        warehouseService.StockService
}

var (
	once      sync.Once
	container *Container
)

func GetContainerInstance(db drivers.Connection) *Container {
	once.Do(func() {
		container = &Container{}
		container.initialize(db)
	})
	return container
}

func (c *Container) initialize(db drivers.Connection) error {
	permissionRepo, err := adminRepository.NewSQLPermissionRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	roleRepo, err := adminRepository.NewSQLRoleRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	userRepo, err := adminRepository.NewSQLUserRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	customerRepo, err := masterRepository.NewSQLCustomerRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	machineRepo, err := masterRepository.NewSQLMachineRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	storeRepo, err := masterRepository.NewSQLStoreRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	containerRepo, err := masterRepository.NewSQLContainerRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	productRepo, err := masterRepository.NewSQLProductRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	joborderRepo, err := masterRepository.NewSQLJobOrderRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	outwardrequestRepo, err := masterRepository.NewSQLOutwardRequestRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	requisitionRepo, err := masterRepository.NewSQLRequisitionRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	batchlabelRepo, err := warehouseRepository.NewSQLBatchLabelRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	labelstickerRepo, err := warehouseRepository.NewSQLLabelStickerRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	inventoryRepo, err := warehouseRepository.NewSQLInventoryRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	stockRepo, err := warehouseRepository.NewSQLStockRepository(db)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	roleConverter := adminConverter.NewRoleConverter()
	userConverter := adminConverter.NewUserConverter(roleConverter)

	machineConverter := masterConverter.NewMachineConverter()
	customerConverter := masterConverter.NewCustomerConverter()
	storeConverter := masterConverter.NewStoreConverter(userConverter)
	containerConverter := masterConverter.NewContainerConverter()
	productConverter := masterConverter.NewProductConverter()
	joborderConverter := masterConverter.NewJobOrderConverter(*productConverter, *customerConverter)
	outwardrequestConverter := masterConverter.NewOutwardRequestConverter(*productConverter, *customerConverter)
	requisitionConverter := masterConverter.NewRequisitionConverter(*productConverter, *storeConverter)

	batchlabelConverter := warehouseConverter.NewBatchLabelConverter(userConverter, productConverter, customerConverter, machineConverter)
	labelstickerConverter := warehouseConverter.NewLabelStickerConverter(batchlabelConverter)
	inventoryConverter := warehouseConverter.NewInventoryConverter()
	stockConverter := warehouseConverter.NewStockConverter(*productConverter, *storeConverter, *containerConverter, *batchlabelConverter)

	c.RoleService = adminService.NewRoleService(roleRepo, permissionRepo, roleConverter)
	c.UserService = adminService.NewUserService(userRepo, roleRepo, userConverter)
	c.AuthService = authService.NewAuthService(userRepo, *userConverter)

	c.CustomerService = masterService.NewCustomerService(customerRepo, customerConverter)
	c.MachineService = masterService.NewMachineService(machineRepo, machineConverter)
	c.StoreService = masterService.NewStoreService(storeRepo, userRepo, storeConverter)
	c.ContainerService = masterService.NewContainerService(containerRepo, containerConverter)
	c.ProductService = masterService.NewProductService(productRepo, productConverter)
	c.JobOrderService = masterService.NewJobOrderService(joborderRepo, customerRepo, productRepo, joborderConverter)
	c.OutwardRequestService = masterService.NewOutwardRequestService(outwardrequestRepo, customerRepo, productRepo, outwardrequestConverter)
	c.RequisitionService = masterService.NewRequisitionService(requisitionRepo, inventoryRepo, storeRepo, productRepo, requisitionConverter)

	c.BatchLabelService = warehouseService.NewBatchLabelService(batchlabelRepo, labelstickerRepo, productRepo, machineRepo, customerRepo, batchlabelConverter, labelstickerConverter)
	c.LabelStickerService = warehouseService.NewLabelStickerService(labelstickerRepo, batchlabelRepo, labelstickerConverter)
	c.InventoryService = warehouseService.NewInventoryService(c.BatchLabelService, inventoryRepo, productRepo, storeRepo, containerRepo, inventoryConverter)
	c.StockService = warehouseService.NewStockService(stockRepo, productRepo, storeRepo, containerRepo, batchlabelRepo, stockConverter)
	return nil
}
