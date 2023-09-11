package service

import "github.com/vamika-digital/wms-api-server/core/business/admin/dto/role"

type RoleService interface {
	GetAllRoles(page int16, pageSize int16, sort string, filter *role.RoleFilterDto) ([]*role.RoleDto, int64, error)
	CreateRole(roleDto *role.RoleCreateDto) error
	GetRoleByID(roleID int64) (*role.RoleDto, error)
	UpdateRole(roleID int64, role *role.RoleUpdateDto) error
	DeleteRole(roleID int64) error
	DeleteRoleByIDs(roleIDs []int64) error
}
