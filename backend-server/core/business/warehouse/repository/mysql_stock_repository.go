package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/stock"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLStockRepository struct {
	DB *sqlx.DB
}

func NewSQLStockRepository(conn drivers.Connection) (StockRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLStockRepository{DB: conn.GetDB()}, nil
}

func (r *SQLStockRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *stock.StockFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (barcode LIKE :query OR batch_no LIKE :query)")
		args["query"] = "%" + filter.Query + "%"
	}

	if filter.ID != 0 {
		queryBuffer.WriteString(" AND id = :id")
		args["id"] = filter.ID
	}

	if filter.ProductID != 0 {
		queryBuffer.WriteString(" AND product_id = :product_id")
		args["product_id"] = filter.ProductID
	}
	if filter.StoreID != 0 {
		queryBuffer.WriteString(" AND store_id = :store_id")
		args["store_id"] = filter.StoreID
	}
	if filter.BinID.Valid {
		queryBuffer.WriteString(" AND bin_id = :bin_id")
		args["bin_id"] = filter.BinID.Int64
	}
	if filter.PalletID.Valid {
		queryBuffer.WriteString(" AND pallet_id = :pallet_id")
		args["pallet_id"] = filter.PalletID.Int64
	}
	if filter.RackID.Valid {
		queryBuffer.WriteString(" AND rack_id = :rack_id")
		args["rack_id"] = filter.RackID.Int64
	}
	if filter.BatchLabelID.Valid {
		queryBuffer.WriteString(" AND batchlabel_id = :batchlabel_id")
		args["batchlabel_id"] = filter.BatchLabelID.Int64
	}

	if filter.Barcode != "" {
		queryBuffer.WriteString(" AND barcode=:barcode")
		args["barcode"] = filter.Barcode
	}

	if filter.BatchNo != "" {
		queryBuffer.WriteString(" AND batch_no=:batch_no")
		args["batch_no"] = filter.BatchNo
	}

	if filter.Status.IsValid() {
		queryBuffer.WriteString(" AND status=:status")
		args["status"] = filter.Status
	}

	if sort != "" {
		queryBuffer.WriteString(fmt.Sprintf(" ORDER BY %s", sort))
	}

	return queryBuffer.String(), args
}

func (r *SQLStockRepository) GetTotalCount(filter *stock.StockFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM stocks WHERE 1=1")
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

func (r *SQLStockRepository) GetAll(page int, pageSize int, sort string, filter *stock.StockFilterDto) ([]*domain.Stock, error) {
	stocks := make([]*domain.Stock, 0)
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("SELECT * FROM stocks WHERE 1=1")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	if page > 0 {
		queryBuffer.WriteString(" LIMIT :start, :end")
		args["start"] = (page - 1) * pageSize
		args["end"] = pageSize
	}

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	err = namedQuery.Select(&stocks, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return stocks, nil
}

func (r *SQLStockRepository) GetById(stockID int64) (*domain.Stock, error) {
	stock := &domain.Stock{}
	err := r.DB.Get(stock, "SELECT * FROM stocks WHERE id = ?", stockID)
	return stock, err
}

func (r *SQLStockRepository) GetByBarcode(stockBarcode string) (*domain.Stock, error) {
	stock := &domain.Stock{}
	err := r.DB.Get(stock, "SELECT * FROM stocks WHERE barcode = ?", stockBarcode)
	return stock, err
}
