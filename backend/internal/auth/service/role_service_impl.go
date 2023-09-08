package service

import (
	"github.com/vamika-digital/wms-api-server/internal/auth/dto/role"
	"github.com/vamika-digital/wms-api-server/internal/auth/repository"
)

type RoleServiceImpl struct {
}

// CreateRole implements RoleService.
func (*RoleServiceImpl) CreateRole(user role.RoleCreateDto) error {
	panic("unimplemented")
}

// DeleteRole implements RoleService.
func (*RoleServiceImpl) DeleteRole(userID int64) error {
	panic("unimplemented")
}

// GetAllRoles implements RoleService.
func (*RoleServiceImpl) GetAllRoles(page int64, pageSize int64, sort string, filter repository.RoleFilterOptions) ([]role.RoleDto, int64, error) {
	panic("unimplemented")
}

// GetRoleByID implements RoleService.
func (*RoleServiceImpl) GetRoleByID(userID int64) (role.RoleDto, error) {
	panic("unimplemented")
}

// UpdateRole implements RoleService.
func (*RoleServiceImpl) UpdateRole(user role.RoleUpdateDto) error {
	panic("unimplemented")
}

func NewRoleService() RoleService {
	return &RoleServiceImpl{}
}
