package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
)

type PermissionRepository interface {
	GetAllByRoleId(roleID int64) ([]*domain.Permission, error)
}
