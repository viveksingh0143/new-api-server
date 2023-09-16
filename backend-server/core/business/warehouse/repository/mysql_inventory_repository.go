package repository

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	masterDomain "github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/inventory"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/reports"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLInventoryRepository struct {
	DB *sqlx.DB
}

func NewSQLInventoryRepository(conn drivers.Connection) (InventoryRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLInventoryRepository{DB: conn.GetDB()}, nil
}

func (r *SQLInventoryRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *inventory.InventoryFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (p.code LIKE :query OR p.name LIKE :query)")
		args["query"] = "%" + filter.Query + "%"
	}

	if filter.ProductTypes != "" {
		productTypesFilter := strings.Split(filter.ProductTypes, ",")
		if len(productTypesFilter) > 0 {
			placeholders := make([]string, len(productTypesFilter))
			for i, pt := range productTypesFilter {
				key := fmt.Sprintf("product_type%d", i+1)
				placeholders[i] = fmt.Sprintf(":%s", key)
				args[key] = pt
			}
			queryBuffer.WriteString(fmt.Sprintf(" AND p.product_type IN (%s)", strings.Join(placeholders, ",")))
		}
	}

	if filter.StoreID != 0 {
		queryBuffer.WriteString(" AND s.store_id = :store_id")
		args["store_id"] = filter.StoreID
	}

	if filter.ProductID != 0 {
		queryBuffer.WriteString(" AND p.id = :id")
		args["id"] = filter.ProductID
	}

	if filter.ProductCode != "" {
		queryBuffer.WriteString(" AND p.code LIKE :code")
		args["code"] = "%" + filter.ProductCode + "%"
	}

	if filter.ContainerID != 0 && filter.ContainerType != "" {
		cType, err := customtypes.GetContainerTypeFromString(filter.ContainerType)
		if err == nil {
			if cType.IsPalletType() {
				queryBuffer.WriteString(" AND s.pallet_id=:container_id")
				args["container_id"] = filter.ContainerID
			} else if cType.IsBinType() {
				queryBuffer.WriteString(" AND s.bin_id=:container_id")
				args["container_id"] = filter.ContainerID
			} else if cType.IsRackType() {
				queryBuffer.WriteString(" AND s.rack_id=:container_id")
				args["container_id"] = filter.ContainerID
			}
		}
	}

	// if filter.Status.IsValid() {
	// 	queryBuffer.WriteString(" AND status = :status")
	// 	args["status"] = filter.Status
	// }

	return queryBuffer.String(), args
}

func (r *SQLInventoryRepository) GetTotalCount(filter *inventory.InventoryFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT count(DISTINCT p.id) FROM products p LEFT JOIN stocks s ON p.id = s.product_id WHERE p.deleted_at IS NULL")
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

func (r *SQLInventoryRepository) GetAll(page int, pageSize int, sort string, filter *inventory.InventoryFilterDto) ([]*domain.Inventory, error) {
	products := make([]*domain.Inventory, 0)
	var queryBuffer bytes.Buffer
	queryBuffer.WriteString("SELECT p.*, p.product_type, p.code, p.name, SUM(CASE WHEN s.status = 'STOCK-IN' THEN s.quantity ELSE 0 END) AS total_stock_count, MAX(CASE WHEN s.status = 'STOCK-IN' THEN s.stockin_at ELSE NULL END) AS last_stockin_at FROM products p LEFT JOIN stocks s ON p.id = s.product_id WHERE p.deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	queryBuffer.WriteString(" GROUP BY p.id")

	if sort != "" {
		if strings.HasPrefix(sort, "code") {
			queryBuffer.WriteString(fmt.Sprintf(" ORDER BY p.%s", sort))
		} else if strings.HasPrefix(sort, "name") {
			queryBuffer.WriteString(fmt.Sprintf(" ORDER BY p.%s", sort))
		} else if strings.HasPrefix(sort, "last_stockin_at") {
			queryBuffer.WriteString(fmt.Sprintf(" ORDER BY %s", sort))
		} else if strings.HasPrefix(sort, "stock_count") {
			newSort := strings.Replace(sort, "stock_count", "last_stockin_at", -1)
			queryBuffer.WriteString(fmt.Sprintf(" ORDER BY %s", newSort))
		} else {
			queryBuffer.WriteString(" ORDER BY last_stockin_at desc")
		}
	} else {
		queryBuffer.WriteString(" ORDER BY last_stockin_at desc")
	}

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

	err = namedQuery.Select(&products, args)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return products, nil
}

func (r *SQLInventoryRepository) GetById(productID int64) (*domain.Inventory, error) {
	product := &domain.Inventory{}
	err := r.DB.Get(product, "SELECT p.*, p.product_type, p.code, p.name, SUM(CASE WHEN s.status = 'STOCK-IN' THEN s.quantity ELSE 0 END) AS total_stock_count, MAX(CASE WHEN s.status = 'STOCK-IN' THEN s.stockin_at ELSE NULL END) AS last_stockin_at FROM products p LEFT JOIN stocks s ON p.id = s.product_id WHERE p.id = ? AND p.deleted_at IS NULL GROUP BY p.id, p.product_type, p.code, p.name ORDER BY p.id", productID)
	return product, err
}

func (r *SQLInventoryRepository) GetByCode(productCode string) (*domain.Inventory, error) {
	product := &domain.Inventory{}
	err := r.DB.Get(product, "SELECT p.*, p.product_type, p.code, p.name, SUM(CASE WHEN s.status = 'STOCK-IN' THEN s.quantity ELSE 0 END) AS total_stock_count, MAX(CASE WHEN s.status = 'STOCK-IN' THEN s.stockin_at ELSE NULL END) AS last_stockin_at FROM products p LEFT JOIN stocks s ON p.id = s.product_id WHERE p.code = ? AND p.deleted_at IS NULL GROUP BY p.id, p.product_type, p.code, p.name ORDER BY p.id", productCode)
	return product, err
}

func (r *SQLInventoryRepository) CreateRawMaterialStock(stockForm *domain.Stock, containerForm *masterDomain.Container) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	query := `INSERT INTO stocks (product_id, store_id, bin_id, pallet_id, rack_id, batchlabel_id, barcode, batch_no, unit_weight, quantity, machine_code, stockin_at, stockout_at, status, last_updated_by ) VALUES (:product_id, :store_id, :bin_id, :pallet_id, :rack_id, :batchlabel_id, :barcode, :batch_no, :unit_weight, :quantity, :machine_code, :stockin_at, :stockout_at, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, stockForm)
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
	stockForm.ID = id

	containerQuery := "UPDATE containers SET stock_level=:stock_level, resource_id=:resource_id, resource_name=:resource_name, items_count=:items_count WHERE id=:id"
	_, err = tx.NamedExec(containerQuery, containerForm)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLInventoryRepository) CreateFinishedStocks(fdStocks []*domain.Stock, stickers []*domain.LabelSticker, containerForm *masterDomain.Container) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	for _, fdStock := range fdStocks {
		query := `INSERT INTO stocks (product_id, store_id, bin_id, pallet_id, rack_id, batchlabel_id, barcode, batch_no, unit_weight, quantity, machine_code, stockin_at, stockout_at, status, last_updated_by ) VALUES (:product_id, :store_id, :bin_id, :pallet_id, :rack_id, :batchlabel_id, :barcode, :batch_no, :unit_weight, :quantity, :machine_code, :stockin_at, :stockout_at, :status, :last_updated_by)`
		res, err := tx.NamedExec(query, fdStock)
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
		fdStock.ID = id
	}

	for _, sticker := range stickers {
		stickerQuery := "UPDATE labelstickers SET is_used=:is_used WHERE id=:id"
		_, err = tx.NamedExec(stickerQuery, sticker)
		if err != nil {
			log.Printf("%+v\n", err)
			tx.Rollback()
			return err
		}
	}

	containerQuery := "UPDATE containers SET stock_level=:stock_level, resource_id=:resource_id, resource_name=:resource_name, items_count=:items_count WHERE id=:id"
	_, err = tx.NamedExec(containerQuery, containerForm)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLInventoryRepository) AttachContainer(destinationContainer *masterDomain.Container, sourceContainer *masterDomain.Container) error {
	tx, err := r.DB.Beginx()
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	containerQuery := "UPDATE containers SET stock_level=:stock_level, resource_id=:resource_id, resource_name=:resource_name, items_count=:items_count WHERE id=:id"
	_, err = tx.NamedExec(containerQuery, destinationContainer)
	if err != nil {
		log.Printf("%+v\n", err)
		tx.Rollback()
		return err
	}

	updateQuery := "UPDATE stocks SET DESTINATIONID=? WHERE SOURCEID=?"
	if destinationContainer.ContainerType.IsBinType() {
		updateQuery = strings.Replace(updateQuery, "DESTINATIONID", "bin_id", 1)
	} else if destinationContainer.ContainerType.IsPalletType() {
		updateQuery = strings.Replace(updateQuery, "DESTINATIONID", "pallet_id", 1)
	} else if destinationContainer.ContainerType.IsRackType() {
		updateQuery = strings.Replace(updateQuery, "DESTINATIONID", "rack_id", 1)
	}

	if sourceContainer.ContainerType.IsBinType() {
		updateQuery = strings.Replace(updateQuery, "SOURCEID", "bin_id", 1)
	} else if sourceContainer.ContainerType.IsPalletType() {
		updateQuery = strings.Replace(updateQuery, "SOURCEID", "pallet_id", 1)
	} else if sourceContainer.ContainerType.IsRackType() {
		updateQuery = strings.Replace(updateQuery, "SOURCEID", "rack_id", 1)
	}

	_, err = tx.Exec(updateQuery, destinationContainer.ID, sourceContainer.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLInventoryRepository) GetInventoryDetailForProductIds(productIds []int64) ([]*reports.InventoryStatusDetail, error) {
	result := make([]*reports.InventoryStatusDetail, 0)

	query := "SELECT c.id AS rack_id, c.name AS rack_name, p.id AS product_id, p.name AS product_name, p.code AS product_code, SUM(CASE WHEN s.status = 'STOCK-IN' THEN s.quantity ELSE 0 END) AS stockin_count, SUM(CASE WHEN s.status IN ('STOCK-DISPATCHING','STOCK-OUT','STOCK-REJECT') THEN s.quantity ELSE 0 END) AS stockout_count, MAX(CASE WHEN s.status = 'STOCK-IN' THEN s.stockin_at ELSE NULL END) AS stockin_at FROM stocks s JOIN containers c ON s.rack_id = c.id JOIN products p ON s.product_id = p.id WHERE s.rack_id IS NOT NULL AND s.product_id IN (?) GROUP BY c.id, c.name, p.id, p.name, p.code  HAVING stockin_count - stockout_count  > 0 ORDER BY stockin_at"
	query, args, err := sqlx.In(query, productIds)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	query = r.DB.Rebind(query)

	if err := r.DB.Select(&result, query, args...); err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	return result, nil
}
