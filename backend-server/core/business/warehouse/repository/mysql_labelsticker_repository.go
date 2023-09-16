package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/labelsticker"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLLabelStickerRepository struct {
	DB *sqlx.DB
}

func NewSQLLabelStickerRepository(conn drivers.Connection) (LabelStickerRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLLabelStickerRepository{DB: conn.GetDB()}, nil
}

func (r *SQLLabelStickerRepository) Create(labelstickers []*domain.LabelSticker) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	for _, labelsticker := range labelstickers {
		query := `INSERT INTO labelstickers (uuid, packet_no, print_count, shift, product_line, batch_no, unit_weight, quantity, machine_no, batchlabel_id, last_updated_by) VALUES(:uuid, :packet_no, :print_count, :shift, :product_line, :batch_no, :unit_weight, :quantity, :machine_no, :batchlabel_id, :last_updated_by)`
		res, err := tx.NamedExec(query, labelsticker)
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
		labelsticker.ID = id
	}

	return tx.Commit()
}

func (r *SQLLabelStickerRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *labelsticker.LabelStickerFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (uuid LIKE :query OR packet_no LIKE :query OR last_updated_by LIKE :query)")
		args["query"] = "%" + filter.Query + "%"
	}
	if filter.ID != 0 {
		queryBuffer.WriteString(" AND id = :id")
		args["id"] = filter.ID
	}
	if filter.UUIDCode != "" {
		queryBuffer.WriteString(" AND uuid LIKE :uuid")
		args["uuid"] = "%" + filter.UUIDCode + "%"
	}
	if filter.LastUpdatedBy.Valid {
		queryBuffer.WriteString(" AND last_updated_by LIKE :last_updated_by")
		args["last_updated_by"] = "%" + filter.LastUpdatedBy.String + "%"
	}
	if filter.BatchID > 0 {
		queryBuffer.WriteString(" AND batchlabel_id = :batchlabel_id")
		args["batchlabel_id"] = filter.BatchID
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

func (r *SQLLabelStickerRepository) GetTotalCount(filter *labelsticker.LabelStickerFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM labelstickers WHERE 1=1")
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

func (r *SQLLabelStickerRepository) GetAll(page int, pageSize int, sort string, filter *labelsticker.LabelStickerFilterDto) ([]*domain.LabelSticker, error) {
	labelstickers := make([]*domain.LabelSticker, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM labelstickers WHERE 1=1")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&labelstickers, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return labelstickers, nil
}

func (r *SQLLabelStickerRepository) GetById(labelstickerID int64) (*domain.LabelSticker, error) {
	labelsticker := &domain.LabelSticker{}
	err := r.DB.Get(labelsticker, "SELECT * FROM labelstickers WHERE id = ? AND 1=1", labelstickerID)
	return labelsticker, err
}

func (r *SQLLabelStickerRepository) GetStickerCountByIds(batchIDs []int64) (map[int64]int64, error) {
	batchIdStickerCountMap := make(map[int64]int64)

	query := `SELECT batchlabel_id, COUNT(*) FROM labelstickers WHERE batchlabel_id IN (?) GROUP BY batchlabel_id;`
	query, args, err := sqlx.In(query, batchIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	query = r.DB.Rebind(query)
	rows, err := r.DB.Queryx(query, args...)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var batchId int64
		var count int64
		if err := rows.Scan(&batchId, &count); err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
		batchIdStickerCountMap[batchId] = count
	}

	if err := rows.Err(); err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return batchIdStickerCountMap, nil
}
