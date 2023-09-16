package repository

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/base/customerrors"
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/container"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLContainerRepository struct {
	DB *sqlx.DB
}

func NewSQLContainerRepository(conn drivers.Connection) (ContainerRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLContainerRepository{DB: conn.GetDB()}, nil
}

func (r *SQLContainerRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *container.ContainerFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (code LIKE :query OR name LIKE :query)")
		args["query"] = "%" + filter.Query + "%"
	}

	if filter.ID != 0 {
		queryBuffer.WriteString(" AND id = :id")
		args["id"] = filter.ID
	}

	if filter.Code != "" {
		queryBuffer.WriteString(" AND code=:code")
		args["code"] = filter.Code
	}

	if filter.Name != "" {
		queryBuffer.WriteString(" AND name LIKE :name")
		args["name"] = "%" + filter.Name + "%"
	}

	if filter.ContainerType != nil && filter.ContainerType.IsValid() {
		queryBuffer.WriteString(" AND container_type = :container_type")
		args["container_type"] = filter.ContainerType
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

func (r *SQLContainerRepository) GetTotalCount(filter *container.ContainerFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM containers WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(0, 0, "", filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return 0, err
	}

	err = namedQuery.Get(&count, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return 0, err
	}

	return count, nil
}

func (r *SQLContainerRepository) GetAll(page int, pageSize int, sort string, filter *container.ContainerFilterDto) ([]*domain.Container, error) {
	containers := make([]*domain.Container, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM containers WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&containers, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return containers, nil
}

func (r *SQLContainerRepository) Create(container *domain.Container) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM containers WHERE code = ?", container.Code)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("a container with the code %s already exists", container.Code)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	// Insert container
	query := `INSERT INTO containers (container_type, code, name, address, status, last_updated_by) VALUES(:container_type, :code, :name, :address, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, container)
	if err != nil {
		log.Printf("%+v\n", err)
		_ = tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("%+v\n", err)
		_ = tx.Rollback()
		return err
	}
	container.ID = id

	return tx.Commit()
}

func (r *SQLContainerRepository) GetById(containerID int64) (*domain.Container, error) {
	container := &domain.Container{}
	err := r.DB.Get(container, "SELECT * FROM containers WHERE id = ? AND deleted_at IS NULL", containerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.NewNotFoundError(fmt.Sprintf("No container found with ID: %d", containerID))
		}
		return nil, err
	}
	return container, err
}

func (r *SQLContainerRepository) GetByIds(containerIDs []int64) ([]*domain.Container, error) {
	var containers []*domain.Container
	query, args, err := sqlx.In("SELECT * FROM containers WHERE id IN (?) AND deleted_at IS NULL", containerIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	query = r.DB.Rebind(query)
	err = r.DB.Select(&containers, query, args...)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return containers, nil
}

func (r *SQLContainerRepository) GetByCode(containerCode string) (*domain.Container, error) {
	container := &domain.Container{}
	err := r.DB.Get(container, "SELECT * FROM containers WHERE code = ? AND deleted_at IS NULL", containerCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.NewNotFoundError(fmt.Sprintf("No container found with Code: %s", containerCode))
		}
		return nil, err
	}
	return container, err
}

func (r *SQLContainerRepository) Update(container *domain.Container) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM containers WHERE code = ? AND id != ?", container.Code, container.ID)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("a container with the code %s already exists", container.Code)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	query := "UPDATE containers SET container_type=:container_type, code=:code, name=:name, address=:address, status=:status, last_updated_by=:last_updated_by WHERE id=:id"
	_, err = tx.NamedExec(query, container)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLContainerRepository) Delete(containerID int64) error {
	_, err := r.DB.Exec("UPDATE containers SET deleted_at = NOW() WHERE id = ?", containerID)
	return err
}

func (r *SQLContainerRepository) DeleteByIDs(containerIDs []int64) error {
	if len(containerIDs) == 0 {
		return nil
	}

	query := "UPDATE containers SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, containerIDs)

	if err != nil {
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}

func (r *SQLContainerRepository) GetContainerCodeInfoDto() ([]*container.ContainerCodeInfoDto, error) {
	containerInfo := make([]*container.ContainerCodeInfoDto, 0)
	err := r.DB.Select(&containerInfo, "SELECT c1.container_type, c1.code FROM containers c1 JOIN (SELECT container_type, MAX(id) as max_id FROM containers GROUP BY container_type) c2 ON c1.container_type = c2.container_type AND c1.id = c2.max_id")
	return containerInfo, err
}

func (r *SQLContainerRepository) GetOneContainerByCodeAndType(code string, containerType customtypes.ContainerType) (*domain.Container, error) {
	container := &domain.Container{}
	err := r.DB.Get(container, "SELECT * FROM containers WHERE status = 1 AND code = ? AND container_type = ? AND deleted_at IS NULL", code, containerType)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.NewNotFoundError(fmt.Sprintf("No %s found with Code: %s", strings.ToLower(containerType.String()), code))
		}
		return nil, err
	}
	return container, err
}

func (r *SQLContainerRepository) MarkContainerFullById(id int64) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	query := "UPDATE containers SET stock_level='FULL' WHERE id=?"
	_, err = tx.Exec(query, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *SQLContainerRepository) AttachedCount(resourceId int64, resourceName string) (int, error) {
	var count int
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	queryBuffer.WriteString("SELECT COUNT(*) FROM containers WHERE deleted_at IS NULL AND resource_id=:resource_id AND resource_name=:resource_name")
	args["resource_id"] = resourceId
	args["resource_name"] = resourceName

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return 0, err
	}

	err = namedQuery.Get(&count, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return 0, err
	}

	return count, nil
}
