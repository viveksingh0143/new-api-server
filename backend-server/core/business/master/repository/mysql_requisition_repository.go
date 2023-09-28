package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/requisition"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLRequisitionRepository struct {
	DB *sqlx.DB
}

func NewSQLRequisitionRepository(conn drivers.Connection) (RequisitionRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLRequisitionRepository{DB: conn.GetDB()}, nil
}

func (r *SQLRequisitionRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *requisition.RequisitionFilterDto) (string, map[string]interface{}) {
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

	if filter.IsApproved.Valid {
		queryBuffer.WriteString(" AND approved = :approved")
		args["approved"] = filter.IsApproved.Bool
	}

	if filter.Status.IsValid() {
		queryBuffer.WriteString(" AND status = :status")
		args["status"] = filter.Status
	}

	if filter.Store != nil && filter.Store.ID > 0 {
		queryBuffer.WriteString(" AND store_id = :store_id")
		args["store_id"] = filter.Store.ID
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

func (r *SQLRequisitionRepository) GetTotalCount(filter *requisition.RequisitionFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM requisitions WHERE deleted_at IS NULL")
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

func (r *SQLRequisitionRepository) GetAll(page int, pageSize int, sort string, filter *requisition.RequisitionFilterDto) ([]*domain.Requisition, error) {
	requisitions := make([]*domain.Requisition, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM requisitions WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&requisitions, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return requisitions, nil
}

func (r *SQLRequisitionRepository) Create(requisition *domain.Requisition) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM requisitions WHERE order_no = ?", requisition.OrderNo)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if count > 0 {
		return fmt.Errorf("a requisition with the order number %s already exists", requisition.OrderNo)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := `INSERT INTO requisitions (issued_date, order_no, department, store_id, approved, status, last_updated_by) VALUES(:issued_date, :order_no, :department, :store_id, :approved, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, requisition)
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
	requisition.ID = id

	for _, item := range requisition.Items {
		item.RequisitionID = requisition.ID
		query = "INSERT INTO requisition_items (requisition_id, product_id, quantity) VALUES(:requisition_id, :product_id, :quantity)"
		if _, err := tx.NamedExec(query, item); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLRequisitionRepository) GetById(requisitionID int64) (*domain.Requisition, error) {
	requisition := &domain.Requisition{}
	err := r.DB.Get(requisition, "SELECT * FROM requisitions WHERE id = ? AND deleted_at IS NULL", requisitionID)
	return requisition, err
}

func (r *SQLRequisitionRepository) GetByOrderNo(requisitionNumber string) (*domain.Requisition, error) {
	requisition := &domain.Requisition{}
	err := r.DB.Get(requisition, "SELECT * FROM requisitions WHERE order_no = ? AND deleted_at IS NULL", requisitionNumber)
	return requisition, err
}

func (r *SQLRequisitionRepository) Update(requisition *domain.Requisition) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM requisitions WHERE order_no = ? AND id != ?", requisition.OrderNo, requisition.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if count > 0 {
		return fmt.Errorf("a requisition with the order number %s already exists", requisition.OrderNo)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := "UPDATE requisitions SET issued_date=:issued_date, order_no=:order_no, department=:department, store_id=:store_id, last_updated_by=:last_updated_by WHERE id=:id"
	_, err = tx.NamedExec(query, requisition)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	query = "DELETE FROM requisition_items WHERE requisition_id=?"
	_, err = tx.Exec(query, requisition.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	for _, item := range requisition.Items {
		item.RequisitionID = requisition.ID
		query = "INSERT INTO requisition_items (requisition_id, product_id, quantity) VALUES(:requisition_id, :product_id, :quantity)"
		if _, err := tx.NamedExec(query, item); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *SQLRequisitionRepository) Delete(requisitionID int64) error {
	_, err := r.DB.Exec("UPDATE requisitions SET deleted_at = NOW() WHERE id = ?", requisitionID)
	return err
}

func (r *SQLRequisitionRepository) DeleteByIDs(requisitionIDs []int64) error {
	if len(requisitionIDs) == 0 {
		return nil
	}

	query := "UPDATE requisitions SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, requisitionIDs)

	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}

func (r *SQLRequisitionRepository) GetItemsForRequisitions(orderIDs []int64) (map[int64][]*domain.RequisitionItem, error) {
	itemsMap := make(map[int64][]*domain.RequisitionItem)

	query := `SELECT * FROM requisition_items WHERE requisition_id IN (?)`
	query, args, err := sqlx.In(query, orderIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	query = r.DB.Rebind(query)
	var items []*domain.RequisitionItem

	if err := r.DB.Select(&items, query, args...); err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	for _, item := range items {
		itemsMap[item.RequisitionID] = append(itemsMap[item.RequisitionID], item)
	}
	return itemsMap, nil
}

func (r *SQLRequisitionRepository) GetItemsForRequisition(orderID int64) ([]*domain.RequisitionItem, error) {
	var requisitionItems []*domain.RequisitionItem

	query := "SELECT * FROM requisition_items WHERE requisition_id=:requisition_id"
	args := map[string]interface{}{
		"requisition_id": orderID,
	}

	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&requisitionItems, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return requisitionItems, nil
}

func (r *SQLRequisitionRepository) ApproveRequisitionByIDs(requisitionIDs []int64) error {
	if len(requisitionIDs) == 0 {
		return nil
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	query := "UPDATE requisitions SET approved = ? WHERE id IN (?)"
	query, args, err := sqlx.In(query, 1, requisitionIDs)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = tx.Rebind(query)

	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
