package service

import (
	"github.com/vamika-digital/wms-api-server/internal/auth/dto/role"
	"github.com/vamika-digital/wms-api-server/internal/auth/repository"
)

type RoleService interface {
	GetAllRoles(page int64, pageSize int64, sort string, filter repository.RoleFilterOptions) ([]role.RoleDto, int64, error)
	CreateRole(user role.RoleCreateDto) error
	GetRoleByID(userID int64) (role.RoleDto, error)
	UpdateRole(user role.RoleUpdateDto) error
	DeleteRole(userID int64) error
}
