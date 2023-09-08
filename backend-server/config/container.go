package config

import (
	"sync"

	"github.com/vamika-digital/wms-api-server/core/business/admin/repository"
	"github.com/vamika-digital/wms-api-server/core/business/admin/service"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type Container struct {
	RoleService service.RoleService
	UserService service.UserService
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

func (c *Container) initialize(db drivers.Connection) {
	permissionRepo := repository.NewSQLPermissionRepository(db)
	roleRepo := repository.NewSQLRoleRepository(db)
	userRepo := repository.NewSQLUserRepository(db)

	c.RoleService = service.NewRoleService(roleRepo, permissionRepo)
	c.UserService = service.NewUserService(userRepo, roleRepo)
}
