package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/outwardrequest"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLOutwardRequestRepository struct {
	DB *sqlx.DB
}

func NewSQLOutwardRequestRepository(conn drivers.Connection) (OutwardRequestRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLOutwardRequestRepository{DB: conn.GetDB()}, nil
}

func (r *SQLOutwardRequestRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *outwardrequest.OutwardRequestFilterDto) (string, map[string]interface{}) {
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

func (r *SQLOutwardRequestRepository) GetTotalCount(filter *outwardrequest.OutwardRequestFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM outwardrequests WHERE deleted_at IS NULL")
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

func (r *SQLOutwardRequestRepository) GetAll(page int, pageSize int, sort string, filter *outwardrequest.OutwardRequestFilterDto) ([]*domain.OutwardRequest, error) {
	outwardrequests := make([]*domain.OutwardRequest, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM outwardrequests WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&outwardrequests, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return outwardrequests, nil
}

func (r *SQLOutwardRequestRepository) Create(outwardrequest *domain.OutwardRequest) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM outwardrequests WHERE order_no = ?", outwardrequest.OrderNo)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("a outwardrequest with the order number %s already exists", outwardrequest.OrderNo)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := `INSERT INTO outwardrequests (issued_date, order_no, customer_id, status, last_updated_by) VALUES(:issued_date, :order_no, :customer_id, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, outwardrequest)
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
	outwardrequest.ID = id

	for _, item := range outwardrequest.Items {
		item.OutwardRequestID = outwardrequest.ID
		query = "INSERT INTO outwardrequest_items (outwardrequest_id, product_id, quantity) VALUES(:outwardrequest_id, :product_id, :quantity)"
		if _, err := tx.NamedExec(query, item); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLOutwardRequestRepository) GetById(outwardrequestID int64) (*domain.OutwardRequest, error) {
	outwardrequest := &domain.OutwardRequest{}
	err := r.DB.Get(outwardrequest, "SELECT * FROM outwardrequests WHERE id = ? AND deleted_at IS NULL", outwardrequestID)
	return outwardrequest, err
}

func (r *SQLOutwardRequestRepository) GetByOrderNo(outwardrequestNumber string) (*domain.OutwardRequest, error) {
	outwardrequest := &domain.OutwardRequest{}
	err := r.DB.Get(outwardrequest, "SELECT * FROM outwardrequests WHERE order_no = ? AND deleted_at IS NULL", outwardrequestNumber)
	return outwardrequest, err
}

func (r *SQLOutwardRequestRepository) Update(outwardrequest *domain.OutwardRequest) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM outwardrequests WHERE order_no = ? AND id != ?", outwardrequest.OrderNo, outwardrequest.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("a outwardrequest with the order number %s already exists", outwardrequest.OrderNo)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := "UPDATE outwardrequests SET issued_date=:issued_date, order_no=:order_no, customer_id=:customer_id, last_updated_by=:last_updated_by WHERE id=:id"
	_, err = tx.NamedExec(query, outwardrequest)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	query = "DELETE FROM outwardrequest_items WHERE outwardrequest_id=?"
	_, err = tx.Exec(query, outwardrequest.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	for _, item := range outwardrequest.Items {
		item.OutwardRequestID = outwardrequest.ID
		query = "INSERT INTO outwardrequest_items (outwardrequest_id, product_id, quantity) VALUES(:outwardrequest_id, :product_id, :quantity)"
		if _, err := tx.NamedExec(query, item); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLOutwardRequestRepository) Delete(outwardrequestID int64) error {
	_, err := r.DB.Exec("UPDATE outwardrequests SET deleted_at = NOW() WHERE id = ?", outwardrequestID)
	return err
}

func (r *SQLOutwardRequestRepository) DeleteByIDs(outwardrequestIDs []int64) error {
	if len(outwardrequestIDs) == 0 {
		return nil
	}

	query := "UPDATE outwardrequests SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, outwardrequestIDs)

	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}

func (r *SQLOutwardRequestRepository) GetItemsForOutwardRequests(orderIDs []int64) (map[int64][]*domain.OutwardRequestItem, error) {
	itemsMap := make(map[int64][]*domain.OutwardRequestItem)

	query := `SELECT * FROM outwardrequest_items WHERE outwardrequest_id IN (?)`
	query, args, err := sqlx.In(query, orderIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	query = r.DB.Rebind(query)
	var items []*domain.OutwardRequestItem

	if err := r.DB.Select(&items, query, args...); err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	for _, item := range items {
		itemsMap[item.OutwardRequestID] = append(itemsMap[item.OutwardRequestID], item)
	}
	return itemsMap, nil
}

func (r *SQLOutwardRequestRepository) GetItemsForOutwardRequest(orderID int64) ([]*domain.OutwardRequestItem, error) {
	var outwardrequestItems []*domain.OutwardRequestItem

	query := "SELECT * FROM outwardrequest_items WHERE outwardrequest_id=:outwardrequest_id"
	args := map[string]interface{}{
		"outwardrequest_id": orderID,
	}

	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&outwardrequestItems, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return outwardrequestItems, nil
}
