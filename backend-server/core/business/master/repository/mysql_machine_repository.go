package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/machine"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLMachineRepository struct {
	DB *sqlx.DB
}

func NewSQLMachineRepository(conn drivers.Connection) (MachineRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLMachineRepository{DB: conn.GetDB()}, nil
}

func (r *SQLMachineRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *machine.MachineFilterDto) (string, map[string]interface{}) {
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

func (r *SQLMachineRepository) GetTotalCount(filter *machine.MachineFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM machines WHERE deleted_at IS NULL")
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

func (r *SQLMachineRepository) GetAll(page int, pageSize int, sort string, filter *machine.MachineFilterDto) ([]*domain.Machine, error) {
	machines := make([]*domain.Machine, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM machines WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&machines, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return machines, nil
}

func (r *SQLMachineRepository) Create(machine *domain.Machine) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM machines WHERE code = ?", machine.Code)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("a machine with the code %s already exists", machine.Code)
	}

	err = r.DB.Get(&count, "SELECT COUNT(*) FROM machines WHERE name = ?", machine.Name)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("a machine with the name %s already exists", machine.Name)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	// Insert machine
	query := `INSERT INTO machines (code, name, status, last_updated_by) VALUES (:code, :name, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, machine)
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
	machine.ID = id

	return tx.Commit()
}

func (r *SQLMachineRepository) GetById(machineID int64) (*domain.Machine, error) {
	machine := &domain.Machine{}
	err := r.DB.Get(machine, "SELECT * FROM machines WHERE id = ? AND deleted_at IS NULL", machineID)
	return machine, err
}

func (r *SQLMachineRepository) GetByIds(machineIDs []int64) ([]*domain.Machine, error) {
	var machines []*domain.Machine
	query, args, err := sqlx.In("SELECT * FROM machines WHERE id IN (?) AND deleted_at IS NULL", machineIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	query = r.DB.Rebind(query)
	err = r.DB.Select(&machines, query, args...)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return machines, nil
}

func (r *SQLMachineRepository) GetByCode(machineCode string) (*domain.Machine, error) {
	machine := &domain.Machine{}
	err := r.DB.Get(machine, "SELECT * FROM machines WHERE code = ? AND deleted_at IS NULL", machineCode)
	return machine, err
}

func (r *SQLMachineRepository) GetByName(machineName string) (*domain.Machine, error) {
	machine := &domain.Machine{}
	err := r.DB.Get(machine, "SELECT * FROM machines WHERE name = ? AND deleted_at IS NULL", machineName)
	return machine, err
}

func (r *SQLMachineRepository) Update(machine *domain.Machine) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM machines WHERE code = ? AND id != ?", machine.Code, machine.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("a machine with the code %s already exists", machine.Code)
	}

	err = r.DB.Get(&count, "SELECT COUNT(*) FROM machines WHERE name = ? AND id != ?", machine.Name, machine.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("a machine with the name %s already exists", machine.Name)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := "UPDATE machines SET code=:code, name=:name, status=:status, last_updated_by=:last_updated_by WHERE id=:id"

	_, err = tx.NamedExec(query, machine)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLMachineRepository) Delete(machineID int64) error {
	_, err := r.DB.Exec("UPDATE machines SET deleted_at = NOW() WHERE id = ?", machineID)
	return err
}

func (r *SQLMachineRepository) DeleteByIDs(machineIDs []int64) error {
	if len(machineIDs) == 0 {
		return nil
	}

	query := "UPDATE machines SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, machineIDs)

	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}
