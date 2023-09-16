package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLCustomerRepository struct {
	DB *sqlx.DB
}

func NewSQLCustomerRepository(conn drivers.Connection) (CustomerRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLCustomerRepository{DB: conn.GetDB()}, nil
}

func (r *SQLCustomerRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *customer.CustomerFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (code LIKE :query OR name LIKE :query OR contact_person LIKE :query)")
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
	if filter.Code != "" {
		queryBuffer.WriteString(" AND code LIKE :code")
		args["code"] = "%" + filter.Code + "%"
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

func (r *SQLCustomerRepository) GetTotalCount(filter *customer.CustomerFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM customers WHERE deleted_at IS NULL")
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

func (r *SQLCustomerRepository) GetAll(page int, pageSize int, sort string, filter *customer.CustomerFilterDto) ([]*domain.Customer, error) {
	customers := make([]*domain.Customer, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM customers WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&customers, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return customers, nil
}

func (r *SQLCustomerRepository) Create(customer *domain.Customer) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM customers WHERE code = ?", customer.Code)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("a customer with the code %s already exists", customer.Code)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := `INSERT INTO customers (code, name, contact_person, billing_address_address1, billing_address_address2, billing_address_state, billing_address_country, billing_address_pincode, shipping_address_address1, shipping_address_address2, shipping_address_state, shipping_address_country, shipping_address_pincode, status, last_updated_by) VALUES(:code, :name, :contact_person, :billing_address_address1, :billing_address_address2, :billing_address_state, :billing_address_country, :billing_address_pincode, :shipping_address_address1, :shipping_address_address2, :shipping_address_state, :shipping_address_country, :shipping_address_pincode, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, customer)
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
	customer.ID = id

	return tx.Commit()
}

func (r *SQLCustomerRepository) GetById(customerID int64) (*domain.Customer, error) {
	customer := &domain.Customer{}
	err := r.DB.Get(customer, "SELECT * FROM customers WHERE id = ? AND deleted_at IS NULL", customerID)
	return customer, err
}

func (r *SQLCustomerRepository) GetByIds(customerIDs []int64) ([]*domain.Customer, error) {
	var customers []*domain.Customer
	query, args, err := sqlx.In("SELECT * FROM customers WHERE id IN (?) AND deleted_at IS NULL", customerIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	query = r.DB.Rebind(query)
	err = r.DB.Select(&customers, query, args...)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return customers, nil
}

func (r *SQLCustomerRepository) Update(customer *domain.Customer) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM customers WHERE code = ? AND id != ?", customer.Code, customer.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("a customer with the code %s already exists", customer.Code)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := "UPDATE customers SET code=:code, name=:name, contact_person=:contact_person, billing_address_address1=:billing_address_address1, billing_address_address2=:billing_address_address2, billing_address_state=:billing_address_state, billing_address_country=:billing_address_country, billing_address_pincode=:billing_address_pincode, shipping_address_address1=:shipping_address_address1, shipping_address_address2=:shipping_address_address2, shipping_address_state=:shipping_address_state, shipping_address_country=:shipping_address_country, shipping_address_pincode=:shipping_address_pincode, status=:status, last_updated_by=:last_updated_by WHERE id=:id"
	_, err = tx.NamedExec(query, customer)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLCustomerRepository) Delete(customerID int64) error {
	_, err := r.DB.Exec("UPDATE customers SET deleted_at = NOW() WHERE id = ?", customerID)
	return err
}

func (r *SQLCustomerRepository) DeleteByIDs(customerIDs []int64) error {
	if len(customerIDs) == 0 {
		return nil
	}

	query := "UPDATE customers SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, customerIDs)

	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}
