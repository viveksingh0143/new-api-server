package repository

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLUserRepository struct {
	DB *sqlx.DB
}

func NewSQLUserRepository(conn drivers.Connection) (UserRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLUserRepository{DB: conn.GetDB()}, nil
}

func (r *SQLUserRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *user.UserFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (name LIKE :query OR username LIKE :query OR email LIKE :query OR staff_id LIKE :query)")
		args["query"] = "%" + filter.Query + "%"
	}

	if filter.ID != 0 {
		queryBuffer.WriteString(" AND id = :id")
		args["id"] = filter.ID
	}
	if filter.Name != "" {
		queryBuffer.WriteString(" AND name LIKE :name")
		args["name"] = "%" + filter.Name + "%"
	}
	if filter.Username != "" {
		queryBuffer.WriteString(" AND username LIKE :username")
		args["username"] = "%" + filter.Username + "%"
	}
	if filter.StaffID != "" {
		queryBuffer.WriteString(" AND staff_id LIKE :staff_id")
		args["staff_id"] = "%" + filter.StaffID + "%"
	}
	if filter.EMail != "" {
		queryBuffer.WriteString(" AND email LIKE :email")
		args["email"] = "%" + filter.EMail + "%"
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

func (r *SQLUserRepository) GetTotalCount(filter *user.UserFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM users WHERE deleted_at IS NULL")
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

func (r *SQLUserRepository) GetAll(page int, pageSize int, sort string, filter *user.UserFilterDto) ([]*domain.User, error) {
	users := make([]*domain.User, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM users WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	err = namedQuery.Select(&users, args)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *SQLUserRepository) Create(user *domain.User) error {
	if user.Username != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE username = ?", user.Username)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a user with the username %s already exists", user.Username)
		}
	}

	if user.EMail != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email = ?", user.EMail)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a user with the email %s already exists", user.EMail)
		}
	}

	if user.StaffID != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE staff_id = ?", user.StaffID)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a user with the staff ID %s already exists", user.StaffID)
		}
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	query := `INSERT INTO users (name, username, password, staff_id, email, status, last_updated_by) VALUES(:name, :username, :password, :staff_id, :email, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, user)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	user.ID = id

	// Attach roles
	for _, role := range user.Roles {
		query = "INSERT INTO roles_users (role_id, user_id) VALUES (:role_id, :user_id)"
		args := map[string]interface{}{
			"role_id": role.ID,
			"user_id": user.ID,
		}
		if _, err := tx.NamedExec(query, args); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLUserRepository) GetById(userID int64) (*domain.User, error) {
	user := &domain.User{}
	err := r.DB.Get(user, "SELECT * FROM users WHERE id = ? AND deleted_at IS NULL", userID)
	return user, err
}

func (r *SQLUserRepository) GetByIds(userIDs []int64) ([]*domain.User, error) {
	var users []*domain.User
	query, args, err := sqlx.In("SELECT * FROM users WHERE id IN (?) AND deleted_at IS NULL", userIDs)
	if err != nil {
		return nil, err
	}

	query = r.DB.Rebind(query)
	err = r.DB.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *SQLUserRepository) Update(user *domain.User) error {
	if user.Username != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE username = ? AND id != ?", user.Username, user.ID)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a user with the username %s already exists", user.Username)
		}
	}

	if user.EMail != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email = ? AND id != ?", user.EMail, user.ID)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a user with the email %s already exists", user.EMail)
		}
	}

	if user.StaffID != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE staff_id = ? AND id != ?", user.StaffID, user.ID)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a user with the staff ID %s already exists", user.StaffID)
		}
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	query := "UPDATE users SET name=:name, username=:username, password=:password, staff_id=:staff_id, email=:email, email_confirmation_at=:email_confirmation_at, status=:status, last_updated_by=:last_updated_by WHERE id=:id"
	_, err = tx.NamedExec(query, user)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = "DELETE FROM roles_users WHERE user_id=?"
	_, err = tx.Exec(query, user.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Attach roles
	// Attach roles
	for _, role := range user.Roles {
		query = "INSERT INTO roles_users (role_id, user_id) VALUES (:role_id, :user_id)"
		args := map[string]interface{}{
			"role_id": role.ID,
			"user_id": user.ID,
		}
		if _, err := tx.NamedExec(query, args); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLUserRepository) Delete(userID int64) error {
	_, err := r.DB.Exec("UPDATE users SET deleted_at = NOW() WHERE id = ?", userID)
	return err
}

func (r *SQLUserRepository) GetByUsername(username string) (*domain.User, error) {
	user := &domain.User{}
	err := r.DB.Get(user, "SELECT * FROM users WHERE username = ? AND deleted_at IS NULL", username)
	return user, err
}
func (r *SQLUserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	err := r.DB.Get(user, "SELECT * FROM users WHERE email = ? AND deleted_at IS NULL", email)
	return user, err
}
func (r *SQLUserRepository) GetByStaffID(staffID string) (*domain.User, error) {
	user := &domain.User{}
	err := r.DB.Get(user, "SELECT * FROM users WHERE staff_id = ? AND deleted_at IS NULL", staffID)
	return user, err
}

func (r *SQLUserRepository) DeleteByIDs(userIDs []int64) error {
	if len(userIDs) == 0 {
		return nil
	}

	query := "UPDATE users SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, userIDs)

	if err != nil {
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}
