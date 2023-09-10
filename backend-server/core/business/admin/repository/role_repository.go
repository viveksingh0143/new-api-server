package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/role"
)

type RoleRepository interface {
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	Delete(roleID int64) error
	DeleteByIDs(roleIDs []int64) error
	GetById(roleID int64) (*domain.Role, error)
	GetByName(roleName string) (*domain.Role, error)
	GetTotalCount(filter *role.RoleFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *role.RoleFilterDto) ([]*domain.Role, error)
	GetRolesForUser(userID int64) ([]*domain.Role, error)
	GetRolesForUsers(userIDs []int64) (map[int64][]*domain.Role, error)
}
