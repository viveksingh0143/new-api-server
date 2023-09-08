package repository

import (
	"bytes"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/role"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLRoleRepository struct {
	DB *sqlx.DB
}

func NewSQLRoleRepository(conn drivers.Connection) *SQLRoleRepository {
	return &SQLRoleRepository{DB: conn.GetDB()}
}

func (r *SQLRoleRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter role.RoleFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.ID != 0 {
		queryBuffer.WriteString(" AND id = :id")
		args["id"] = filter.ID
	}

	if filter.Name != "" {
		queryBuffer.WriteString(" AND name LIKE :name")
		args["name"] = "%" + filter.Name + "%"
	}

	if filter.Status.IsValid() {
		queryBuffer.WriteString(" AND status = :status")
		args["status"] = filter.Status
	}

	if sort != "" {
		queryBuffer.WriteString(fmt.Sprintf(" ORDER BY %s", sort))
	}

	if page > 0 {
		queryBuffer.WriteString(" LIMIT :start, :end")
		args["start"] = (page - 1) * pageSize
		args["end"] = pageSize
	}

	return queryBuffer.String(), args
}

func (r *SQLRoleRepository) GetTotalCount(filter role.RoleFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM roles WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(0, 0, "", filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		return 0, err
	}

	err = namedQuery.Get(&count, args)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SQLRoleRepository) GetAll(page int, pageSize int, sort string, filter role.RoleFilterDto) ([]*domain.Role, error) {
	roles := make([]*domain.Role, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM roles WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	err = namedQuery.Select(&roles, args)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *SQLRoleRepository) Create(role *domain.Role) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM roles WHERE name = ?", role.Name)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("a role with the name %s already exists", role.Name)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	// Insert role
	query := `INSERT INTO roles (name, status, last_updated_by) VALUES (:name, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, role)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	role.ID = id

	// Insert permissions
	for _, perm := range role.Permissions {
		perm.RoleID = role.ID
		query = "INSERT INTO permissions (role_id, module_name, create_perm, read_perm, update_perm, delete_perm, export_perm, import_perm) VALUES (:role_id, :module_name, :create_perm, :read_perm, :update_perm, :delete_perm, :export_perm, :import_perm)"
		if _, err := tx.NamedExec(query, perm); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLRoleRepository) GetById(roleID int64) (*domain.Role, error) {
	role := &domain.Role{}
	err := r.DB.Get(role, "SELECT * FROM roles WHERE id = ? AND deleted_at IS NULL", roleID)
	return role, err
}

func (r *SQLRoleRepository) GetByName(roleName string) (*domain.Role, error) {
	role := &domain.Role{}
	err := r.DB.Get(role, "SELECT * FROM roles WHERE name = ? AND deleted_at IS NULL", roleName)
	return role, err
}

func (r *SQLRoleRepository) Update(role *domain.Role) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM roles WHERE name = ? AND id != ?", role.Name, role.ID)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("a role with the name %s already exists", role.Name)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	query := "UPDATE roles SET name=:name, status=:status, last_updated_by=:last_updated_by WHERE id=:id"

	_, err = tx.NamedExec(query, role)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = "DELETE FROM permissions WHERE role_id=?"
	_, err = tx.Exec(query, role.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, permission := range role.Permissions {
		permission.RoleID = role.ID
		query = "INSERT INTO permissions (role_id, module_name, create_perm, read_perm, update_perm, delete_perm, export_perm, import_perm) VALUES (:role_id, :module_name, :create_perm, :read_perm, :update_perm, :delete_perm, :export_perm, :import_perm)"
		_, err := tx.NamedExec(query, &permission)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLRoleRepository) Delete(roleID int64) error {
	_, err := r.DB.Exec("UPDATE roles SET deleted_at = NOW() WHERE id = ?", roleID)
	return err
}

func (r *SQLRoleRepository) GetRolesForUser(userID int64) ([]*domain.Role, error) {
	roles := make([]*domain.Role, 0)
	query := `SELECT roles.* FROM roles JOIN roles_users ON roles.id = roles_users.role_id WHERE roles_users.user_id = ?`
	err := r.DB.Select(&roles, query, userID)
	if err != nil {
		log.Printf("Error fetching roles for user %d: %v\n", userID, err)
		return nil, err
	}
	return roles, nil
}

func (r *SQLRoleRepository) GetRolesForUsers(userIDs []int64) (map[int64][]*domain.Role, error) {
	userRolesMap := make(map[int64][]*domain.Role)

	query := `SELECT ur.user_id, r.id, r.name, r.status FROM roles_users ur JOIN roles r ON ur.role_id = r.id WHERE ur.user_id IN (?)`
	query, args, err := sqlx.In(query, userIDs)
	if err != nil {
		return nil, err
	}

	query = r.DB.Rebind(query)
	var roles []struct {
		UserID     int64                  `db:"user_id"`
		RoleID     int64                  `db:"id"`
		RoleName   string                 `db:"name"`
		RoleStatus customtypes.StatusEnum `db:"status"`
	}

	if err := r.DB.Select(&roles, query, args...); err != nil {
		return nil, err
	}

	for _, roleInfo := range roles {
		role := &domain.Role{
			ID:     roleInfo.RoleID,
			Name:   roleInfo.RoleName,
			Status: roleInfo.RoleStatus,
		}
		userRolesMap[roleInfo.UserID] = append(userRolesMap[roleInfo.UserID], role)
	}
	return userRolesMap, nil
}
