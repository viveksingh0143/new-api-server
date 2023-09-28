package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/outwardrequest"
	"github.com/vamika-digital/wms-api-server/core/business/master/reports"
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

func (r *SQLOutwardRequestRepository) GetShipperLabels(requestID int64, requestName string) ([]*reports.OutwardRequestShipperReport, error) {
	var shipperReports []*reports.OutwardRequestShipperReport

	query := "SELECT report.*, p.code as product_code, p.name as product_name, s.shipper_number as shipper_number, s.packed_at as shipper_packed_at FROM ( SELECT request_id, request_name, batch_no, product_id, shipperlabel_id as shipper_id, count(*) as package_count from wms.stocks where request_id = :request_id AND request_name = :request_name AND status = 'STOCK-OUT' group by request_id, request_name, batch_no, product_id, shipperlabel_id) as report LEFT JOIN products p ON report.product_id = p.id LEFT JOIN shipperlabels s ON report.shipper_id = s.id"
	args := map[string]interface{}{
		"request_id":   requestID,
		"request_name": requestName,
	}

	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&shipperReports, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return shipperReports, nil
}

func (r *SQLOutwardRequestRepository) GenerateShipperLabels(shipperLabelForm *domain.ShipperLabel, batchNo string, productID int64) error {
	var shipperReports []*reports.OutwardRequestShipperReport
	var shippers []*domain.ShipperLabel

	query := "SELECT * from shipperlabels s WHERE s.outwardrequest_id = :outwardrequest_id AND batch_no = :batch_no ORDER BY created_at desc"
	args := map[string]interface{}{
		"outwardrequest_id": shipperLabelForm.OutwardRequestID,
		"batch_no":          batchNo,
	}

	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	err = namedQuery.Select(&shippers, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query = "SELECT report.*, p.name as product_name, p.code as product_code, s.shipper_number as shipper_number, s.packed_at as shipper_packed_at FROM ( SELECT request_id, request_name, batch_no, product_id, shipperlabel_id as shipper_id, count(*) as package_count from wms.stocks where request_id=:request_id AND request_name=:request_name  AND batch_no=:batch_no AND product_id=:product_id AND status = 'STOCK-OUT' AND shipperlabel_id IS NULL group by request_id, request_name, batch_no, product_id, shipperlabel_id) as report LEFT JOIN products p ON report.product_id = p.id LEFT JOIN shipperlabels s ON report.shipper_id = s.id"
	args = map[string]interface{}{
		"request_id":   shipperLabelForm.OutwardRequestID,
		"request_name": "*domain.OutwardRequest",
		"batch_no":     batchNo,
		"product_id":   productID,
	}

	namedQuery, err = r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	err = namedQuery.Select(&shipperReports, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	if len(shipperReports) <= 0 {
		return fmt.Errorf("no data found with the given request")
	}
	if len(shipperReports) > 1 {
		return fmt.Errorf("somthing wrong with the given request")
	}

	shipperCount := len(shippers)
	shipperReport := shipperReports[0]
	shipperLabelForm.PackedAt = time.Now()
	shipperLabelForm.ShipperNumber = fmt.Sprintf("%s%s%d", shipperReport.BatchNo, "00", shipperCount+1)
	shipperLabelForm.ProductCode = shipperReport.ProductCode
	shipperLabelForm.ProductName = shipperReport.ProductName
	shipperLabelForm.BatchNo = shipperReport.BatchNo
	shipperLabelForm.PackedQty = fmt.Sprintf("%dNos", shipperReport.PackageCount)

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query = `INSERT INTO shipperlabels ( shipper_number, customer_name, product_code, product_name, batch_no, packed_qty, packed_at, outwardrequest_id ) VALUES(:shipper_number, :customer_name, :product_code, :product_name, :batch_no, :packed_qty, :packed_at, :outwardrequest_id)`
	res, err := tx.NamedExec(query, shipperLabelForm)
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
	shipperLabelForm.ID = id
	err = tx.Commit()
	if err != nil {
		log.Printf("%+v\n", err)
		_ = tx.Rollback()
		return err
	}

	tx, err = r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	_, err = r.DB.Exec("UPDATE stocks SET shipperlabel_id = ? WHERE request_id = ? AND  request_name = ? AND batch_no = ? AND product_id = ? AND shipperlabel_id IS NULL", shipperLabelForm.ID, shipperLabelForm.OutwardRequestID, "*domain.OutwardRequest", shipperLabelForm.BatchNo, productID)
	if err != nil {
		log.Printf("%+v\n", err)
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
