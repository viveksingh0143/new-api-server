package registry

import (
	"sync"

	adminConverter "github.com/vamika-digital/wms-api-server/core/business/admin/converter"
	adminRepository "github.com/vamika-digital/wms-api-server/core/business/admin/repository"
	adminService "github.com/vamika-digital/wms-api-server/core/business/admin/service"
	authService "github.com/vamika-digital/wms-api-server/core/business/auth/service"
	masterConverter "github.com/vamika-digital/wms-api-server/core/business/master/converter"
	masterRepository "github.com/vamika-digital/wms-api-server/core/business/master/repository"
	masterService "github.com/vamika-digital/wms-api-server/core/business/master/service"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type Container struct {
	RoleService      adminService.RoleService
	UserService      adminService.UserService
	AuthService      authService.AuthService
	MachineService   masterService.MachineService
	CustomerService  masterService.CustomerService
	StoreService     masterService.StoreService
	ContainerService masterService.ContainerService
	ProductService   masterService.ProductService
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
		return err
	}
	roleRepo, err := adminRepository.NewSQLRoleRepository(db)
	if err != nil {
		return err
	}
	userRepo, err := adminRepository.NewSQLUserRepository(db)
	if err != nil {
		return err
	}

	customerRepo, err := masterRepository.NewSQLCustomerRepository(db)
	if err != nil {
		return err
	}

	machineRepo, err := masterRepository.NewSQLMachineRepository(db)
	if err != nil {
		return err
	}

	storeRepo, err := masterRepository.NewSQLStoreRepository(db)
	if err != nil {
		return err
	}

	containerRepo, err := masterRepository.NewSQLContainerRepository(db)
	if err != nil {
		return err
	}

	productRepo, err := masterRepository.NewSQLProductRepository(db)
	if err != nil {
		return err
	}

	roleConverter := adminConverter.NewRoleConverter()
	userConverter := adminConverter.NewUserConverter(roleConverter)

	machineConverter := masterConverter.NewMachineConverter()
	customerConverter := masterConverter.NewCustomerConverter()
	storeConverter := masterConverter.NewStoreConverter(userConverter)
	containerConverter := masterConverter.NewContainerConverter()
	productConverter := masterConverter.NewProductConverter()

	c.RoleService = adminService.NewRoleService(roleRepo, permissionRepo, roleConverter)
	c.UserService = adminService.NewUserService(userRepo, roleRepo, userConverter)
	c.AuthService = authService.NewAuthService(userRepo, *userConverter)
	c.CustomerService = masterService.NewCustomerService(customerRepo, customerConverter)
	c.MachineService = masterService.NewMachineService(machineRepo, machineConverter)
	c.StoreService = masterService.NewStoreService(storeRepo, userRepo, storeConverter)
	c.ContainerService = masterService.NewContainerService(containerRepo, containerConverter)
	c.ProductService = masterService.NewProductService(productRepo, productConverter)
	return nil
}
