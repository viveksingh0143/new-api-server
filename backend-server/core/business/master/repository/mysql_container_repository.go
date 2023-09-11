package repository

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
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
		queryBuffer.WriteString(" AND code LIKE :code")
		args["code"] = "%" + filter.Code + "%"
	}

	if filter.Name != "" {
		queryBuffer.WriteString(" AND name LIKE :name")
		args["name"] = "%" + filter.Name + "%"
	}

	if filter.ContainerType.IsValid() {
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
		return 0, err
	}

	err = namedQuery.Get(&count, args)
	if err != nil {
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
		return nil, err
	}

	err = namedQuery.Select(&containers, args)
	if err != nil {
		return nil, err
	}

	return containers, nil
}

func (r *SQLContainerRepository) Create(container *domain.Container) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM containers WHERE code = ?", container.Code)
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

	// Insert container
	query := `INSERT INTO containers (container_type, code, name, address, status, last_updated_by) VALUES(:container_type, :code, :name, :address, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, container)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	container.ID = id

	return tx.Commit()
}

func (r *SQLContainerRepository) GetById(containerID int64) (*domain.Container, error) {
	container := &domain.Container{}
	err := r.DB.Get(container, "SELECT * FROM containers WHERE id = ? AND deleted_at IS NULL", containerID)
	return container, err
}

func (r *SQLContainerRepository) GetByCode(containerCode string) (*domain.Container, error) {
	container := &domain.Container{}
	err := r.DB.Get(container, "SELECT * FROM containers WHERE code = ? AND deleted_at IS NULL", containerCode)
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
