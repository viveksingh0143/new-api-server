package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/batchlabel"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLBatchLabelRepository struct {
	DB *sqlx.DB
}

func NewSQLBatchLabelRepository(conn drivers.Connection) (BatchLabelRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLBatchLabelRepository{DB: conn.GetDB()}, nil
}

func (r *SQLBatchLabelRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *batchlabel.BatchLabelFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (batch_no LIKE :query OR so_number LIKE :query)")
		args["query"] = "%" + filter.Query + "%"
	}

	if filter.ID != 0 {
		queryBuffer.WriteString(" AND id = :id")
		args["id"] = filter.ID
	}
	if filter.BatchDate.Valid {
		queryBuffer.WriteString(" AND batch_date = :batch_date")
		args["batch_date"] = filter.BatchDate
	}
	if filter.BatchNo != "" {
		queryBuffer.WriteString(" AND batch_no LIKE :batch_no")
		args["batch_no"] = "%" + filter.BatchNo + "%"
	}
	if filter.SoNumber != "" {
		queryBuffer.WriteString(" AND so_number LIKE :so_number")
		args["so_number"] = "%" + filter.SoNumber + "%"
	}
	if filter.Status.IsValid() {
		queryBuffer.WriteString(" AND status = :status")
		args["status"] = filter.Status
	}
	if filter.Product != nil && filter.Product.ID > 0 {
		queryBuffer.WriteString(" AND product_id = :product_id")
		args["product_id"] = filter.Product.ID
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

func (r *SQLBatchLabelRepository) GetTotalCount(filter *batchlabel.BatchLabelFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM batchlabels WHERE deleted_at IS NULL")
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

func (r *SQLBatchLabelRepository) GetAll(page int, pageSize int, sort string, filter *batchlabel.BatchLabelFilterDto) ([]*domain.BatchLabel, error) {
	batchlabels := make([]*domain.BatchLabel, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM batchlabels WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&batchlabels, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return batchlabels, nil
}

func (r *SQLBatchLabelRepository) Create(batchlabel *domain.BatchLabel) error {
	if batchlabel.BatchNo != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM batchlabels WHERE batch_no = ?", batchlabel.BatchNo)
		if err != nil {
			log.Printf("%+v\n", err)
			return err
		}
		if count > 0 {
			return fmt.Errorf("a batch label with the batch number %s already exists", batchlabel.BatchNo)
		}
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if batchlabel.Customer != nil && batchlabel.Customer.ID > 0 {
		batchlabel.CustomerID = batchlabel.Customer.ID
	}

	if batchlabel.Product != nil && batchlabel.Product.ID > 0 {
		batchlabel.ProductID = batchlabel.Product.ID
	}

	if batchlabel.Machine != nil && batchlabel.Machine.ID > 0 {
		batchlabel.MachineID = batchlabel.Machine.ID
	}
	query := `INSERT INTO batchlabels (batch_date, batch_no, so_number, target_quantity, po_category, customer_id, product_id, machine_id, unit_weight, unit_weight_type, package_quantity, status, last_updated_by ) VALUES(:batch_date, :batch_no, :so_number, :target_quantity, :po_category, :customer_id, :product_id, :machine_id, :unit_weight, :unit_weight_type, :package_quantity, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, batchlabel)
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
	batchlabel.ID = id
	return tx.Commit()
}

func (r *SQLBatchLabelRepository) GetById(batchlabelID int64) (*domain.BatchLabel, error) {
	batchlabel := &domain.BatchLabel{}
	err := r.DB.Get(batchlabel, "SELECT * FROM batchlabels WHERE id = ? AND deleted_at IS NULL", batchlabelID)
	return batchlabel, err
}

func (r *SQLBatchLabelRepository) GetByIds(batchlabelIDs []int64) ([]*domain.BatchLabel, error) {
	var batchlabels []*domain.BatchLabel
	query, args, err := sqlx.In("SELECT * FROM batchlabels WHERE id IN (?) AND deleted_at IS NULL", batchlabelIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	query = r.DB.Rebind(query)
	err = r.DB.Select(&batchlabels, query, args...)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return batchlabels, nil
}

func (r *SQLBatchLabelRepository) GetByBatchNumber(batchlabelBatchNumber string) (*domain.BatchLabel, error) {
	batchlabel := &domain.BatchLabel{}
	err := r.DB.Get(batchlabel, "SELECT * FROM batchlabels WHERE batch_no = ? AND deleted_at IS NULL", batchlabelBatchNumber)
	return batchlabel, err
}

func (r *SQLBatchLabelRepository) Update(batchlabel *domain.BatchLabel) error {
	if batchlabel.BatchNo != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM batchlabels WHERE batch_no = ? AND id != ?", batchlabel.BatchNo, batchlabel.ID)
		if err != nil {
			log.Printf("%+v\n", err)
			return err
		}
		if count > 0 {
			return fmt.Errorf("a batch label with the batch number %s already exists", batchlabel.BatchNo)
		}
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	if batchlabel.Customer != nil && batchlabel.Customer.ID > 0 {
		batchlabel.CustomerID = batchlabel.Customer.ID
	}

	if batchlabel.Product != nil && batchlabel.Product.ID > 0 {
		batchlabel.ProductID = batchlabel.Product.ID
	}

	if batchlabel.Machine != nil && batchlabel.Machine.ID > 0 {
		batchlabel.MachineID = batchlabel.Machine.ID
	}

	query := "UPDATE batchlabels SET batch_date=:batch_date, batch_no=:batch_no, so_number=:so_number, target_quantity=:target_quantity, po_category=:po_category, customer_id=:customer_id, product_id=:product_id, machine_id=:machine_id, unit_weight=:unit_weight, unit_weight_type=:unit_weight_type, package_quantity=:package_quantity, status=:status, last_updated_by=:last_updated_by WHERE id=:id"
	_, err = tx.NamedExec(query, batchlabel)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLBatchLabelRepository) Delete(batchlabelID int64) error {
	_, err := r.DB.Exec("UPDATE batchlabels SET deleted_at = NOW() WHERE id = ?", batchlabelID)
	return err
}

func (r *SQLBatchLabelRepository) DeleteByIDs(batchlabelIDs []int64) error {
	if len(batchlabelIDs) == 0 {
		return nil
	}

	query := "UPDATE batchlabels SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, batchlabelIDs)

	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}

func (r *SQLBatchLabelRepository) GetTotalStickers(batchlabelID int64) (int64, error) {
	var count int64
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	queryBuffer.WriteString("SELECT COUNT(*) FROM labelstickers WHERE batchlabel_id=:batchlabel_id")
	args["batchlabel_id"] = batchlabelID

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

func (r *SQLBatchLabelRepository) GetTotalStickersForShift(batchlabelID int64, shift string, createdAt time.Time) (int64, error) {
	var count int64
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	queryBuffer.WriteString("SELECT COUNT(*) FROM labelstickers WHERE batchlabel_id=:batchlabel_id AND shift=:shift AND DATE(created_at)=DATE(:created_at)")
	args["batchlabel_id"] = batchlabelID
	args["shift"] = shift
	args["created_at"] = createdAt

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

func (r *SQLBatchLabelRepository) GetByBarcode(barcode string) (*domain.BatchLabel, *domain.LabelSticker, error) {
	batchlabel := &domain.BatchLabel{}
	err := r.DB.Get(batchlabel, "SELECT b.* from labelstickers l INNER JOIN batchlabels b ON l.batchlabel_id = b.id WHERE l.uuid = ?", barcode)

	if err != nil {
		return nil, nil, err
	}

	labelSticker := &domain.LabelSticker{}
	err = r.DB.Get(labelSticker, "SELECT * FROM labelstickers WHERE uuid = ?", barcode)

	return batchlabel, labelSticker, err
}
