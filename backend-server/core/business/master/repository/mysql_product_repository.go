package repository

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLProductRepository struct {
	DB *sqlx.DB
}

func NewSQLProductRepository(conn drivers.Connection) (ProductRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLProductRepository{DB: conn.GetDB()}, nil
}

func (r *SQLProductRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *product.ProductFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (code LIKE :query OR link_code LIKE :query OR name LIKE :query OR description LIKE :query)")
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

	if filter.LinkCode != "" {
		queryBuffer.WriteString(" AND link_code LIKE :code")
		args["link_code"] = "%" + filter.LinkCode + "%"
	}

	if filter.Name != "" {
		queryBuffer.WriteString(" AND name LIKE :name")
		args["name"] = "%" + filter.Name + "%"
	}

	if filter.ProductType.IsValid() {
		queryBuffer.WriteString(" AND product_type = :product_type")
		args["product_type"] = filter.ProductType
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

func (r *SQLProductRepository) GetTotalCount(filter *product.ProductFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM products WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(0, 0, "", filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		return 0, err
	}

	err = namedQuery.Get(&count, args)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *SQLProductRepository) GetAll(page int, pageSize int, sort string, filter *product.ProductFilterDto) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM products WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	err = namedQuery.Select(&products, args)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (r *SQLProductRepository) Create(product *domain.Product) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM products WHERE code = ?", product.Code)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("a product with the code %s already exists", product.Code)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	// Insert product
	query := `INSERT INTO products (product_type, code, link_code, name, description, unit_type, unit_weight, unit_weight_type, status, last_updated_by) VALUES (:product_type, :code, :link_code, :name, :description, :unit_type, :unit_weight, :unit_weight_type, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, product)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	product.ID = id

	return tx.Commit()
}

func (r *SQLProductRepository) GetById(productID int64) (*domain.Product, error) {
	product := &domain.Product{}
	err := r.DB.Get(product, "SELECT * FROM products WHERE id = ? AND deleted_at IS NULL", productID)
	return product, err
}

func (r *SQLProductRepository) GetByCode(productCode string) (*domain.Product, error) {
	product := &domain.Product{}
	err := r.DB.Get(product, "SELECT * FROM products WHERE code = ? AND deleted_at IS NULL", productCode)
	return product, err
}

func (r *SQLProductRepository) Update(product *domain.Product) error {
	var count int
	err := r.DB.Get(&count, "SELECT COUNT(*) FROM products WHERE code = ? AND id != ?", product.Code, product.ID)
	if err != nil {
		return err
	}

	if count > 0 {
		return fmt.Errorf("a product with the code %s already exists", product.Code)
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	query := "UPDATE products SET product_type=:product_type, link_code=:link_code, name=:name, description=:description, unit_type=:unit_type, unit_weight=:unit_weight, unit_weight_type=:unit_weight_type, status=:status, last_updated_by=:last_updated_by WHERE id=:id"
	_, err = tx.NamedExec(query, product)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLProductRepository) Delete(productID int64) error {
	_, err := r.DB.Exec("UPDATE products SET deleted_at = NOW() WHERE id = ?", productID)
	return err
}

func (r *SQLProductRepository) DeleteByIDs(productIDs []int64) error {
	if len(productIDs) == 0 {
		return nil
	}

	query := "UPDATE products SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, productIDs)

	if err != nil {
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}
