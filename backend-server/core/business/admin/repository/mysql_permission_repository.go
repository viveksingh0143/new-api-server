package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLPermissionRepository struct {
	DB *sqlx.DB
}

func NewSQLPermissionRepository(conn drivers.Connection) *SQLPermissionRepository {
	return &SQLPermissionRepository{DB: conn.GetDB()}
}

func (r *SQLPermissionRepository) GetAllByRoleId(roleID int64) ([]*domain.Permission, error) {
	var permissions []*domain.Permission

	query := "SELECT * FROM permissions WHERE role_id=:role_id"
	args := map[string]interface{}{
		"role_id": roleID,
	}

	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	err = namedQuery.Select(&permissions, args)
	if err != nil {
		return nil, err
	}
	return permissions, nil
}
