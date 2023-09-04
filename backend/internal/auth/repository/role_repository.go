package repository

import (
	"github.com/vamika-digital/wms-api-server/config/types/status"
	"github.com/vamika-digital/wms-api-server/internal/auth/domain"
)

type RoleRepository interface {
	Create(role *domain.Role) error
	Update(role *domain.Role) error
	Delete(roleID int64) error
	GetById(roleID int64) (*domain.Role, error)
	GetTotalCount(filter RoleFilterOptions) (int, error)
	GetAll(page int, pageSize int, sort string, filter RoleFilterOptions) ([]*domain.Role, error)
}

type RoleFilterOptions struct {
	Name   string
	Status status.StatusType // Filter by role status, using an exact match
}
