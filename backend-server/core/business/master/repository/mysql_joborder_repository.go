package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/joborder"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLJobOrderRepository struct {
	DB *sqlx.DB
}

func NewSQLJobOrderRepository(conn drivers.Connection) (JobOrderRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLJobOrderRepository{DB: conn.GetDB()}, nil
}

func (r *SQLJobOrderRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *joborder.JobOrderFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (order_no LIKE :query)")
		args["query"] = "%" + filter.Query + "%"
	}

	if filter.ID != 0 {
		queryBuffer.WriteString(" AND id = :id")
		args["id"] = filter.ID
	}

	if filter.OrderNo != "" {
		queryBuffer.WriteString(" AND order_no=:order_no")
		args["order_no"] = filter.OrderNo
	}

	if filter.Status.IsValid() {
		queryBuffer.WriteString(" AND status = :status")
		args["status"] = filter.Status
	}

	if filter.Customer != nil && filter.Customer.ID > 0 {
		queryBuffer.WriteString(" AND customer_id = :customer_id")
		args["customer_id"] = filter.Customer.ID
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

func (r *SQLJobOrderRepository) GetTotalCount(filter *joborder.JobOrderFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM joborders WHERE deleted_at IS NULL")
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

func (r *SQLJobOrderRepository) GetAll(page int, pageSize int, sort string, filter *joborder.JobOrderFilterDto) ([]*domain.JobOrder, error) {
	joborders := make([]*domain.JobOrder, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM joborders WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&joborders, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return joborders, nil
}

func (r *SQLJobOrderRepository) Create(joborder *domain.JobOrder) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM joborders WHERE order_no = ?", joborder.OrderNo)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("a joborder with the order number %s already exists", joborder.OrderNo)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := `INSERT INTO joborders (issued_date ,order_no ,po_category ,customer_id ,status ,last_updated_by) VALUES(:issued_date, :order_no, :po_category, :customer_id, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, joborder)
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
	joborder.ID = id

	for _, item := range joborder.Items {
		item.JobOrderID = joborder.ID
		query = "INSERT INTO joborder_items (joborder_id, product_id, quantity) VALUES(:joborder_id, :product_id, :quantity)"
		if _, err := tx.NamedExec(query, item); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLJobOrderRepository) GetById(joborderID int64) (*domain.JobOrder, error) {
	joborder := &domain.JobOrder{}
	err := r.DB.Get(joborder, "SELECT * FROM joborders WHERE id = ? AND deleted_at IS NULL", joborderID)
	return joborder, err
}

func (r *SQLJobOrderRepository) GetByOrderNo(joborderNumber string) (*domain.JobOrder, error) {
	joborder := &domain.JobOrder{}
	err := r.DB.Get(joborder, "SELECT * FROM joborders WHERE order_no = ? AND deleted_at IS NULL", joborderNumber)
	return joborder, err
}

func (r *SQLJobOrderRepository) Update(joborder *domain.JobOrder) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM joborders WHERE order_no = ? AND id != ?", joborder.OrderNo, joborder.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("a joborder with the order number %s already exists", joborder.OrderNo)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := "UPDATE joborders SET issued_date=:issued_date, order_no=:order_no, po_category=:po_category, customer_id=:customer_id, last_updated_by=:last_updated_by WHERE id=:id"
	_, err = tx.NamedExec(query, joborder)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	query = "DELETE FROM joborder_items WHERE joborder_id=?"
	_, err = tx.Exec(query, joborder.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	for _, item := range joborder.Items {
		item.JobOrderID = joborder.ID
		query = "INSERT INTO joborder_items (joborder_id, product_id, quantity) VALUES(:joborder_id, :product_id, :quantity)"
		if _, err := tx.NamedExec(query, item); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLJobOrderRepository) Delete(joborderID int64) error {
	_, err := r.DB.Exec("UPDATE joborders SET deleted_at = NOW() WHERE id = ?", joborderID)
	return err
}

func (r *SQLJobOrderRepository) DeleteByIDs(joborderIDs []int64) error {
	if len(joborderIDs) == 0 {
		return nil
	}

	query := "UPDATE joborders SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, joborderIDs)

	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}

func (r *SQLJobOrderRepository) GetItemsForJobOrders(orderIDs []int64) (map[int64][]*domain.JobOrderItem, error) {
	itemsMap := make(map[int64][]*domain.JobOrderItem)

	query := `SELECT * FROM joborder_items WHERE joborder_id IN (?)`
	query, args, err := sqlx.In(query, orderIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	query = r.DB.Rebind(query)
	var items []*domain.JobOrderItem

	if err := r.DB.Select(&items, query, args...); err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	for _, item := range items {
		itemsMap[item.JobOrderID] = append(itemsMap[item.JobOrderID], item)
	}
	return itemsMap, nil
}

func (r *SQLJobOrderRepository) GetItemsForJobOrder(orderID int64) ([]*domain.JobOrderItem, error) {
	var joborderItems []*domain.JobOrderItem

	query := "SELECT * FROM joborder_items WHERE joborder_id=:joborder_id"
	args := map[string]interface{}{
		"joborder_id": orderID,
	}

	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&joborderItems, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return joborderItems, nil
}
