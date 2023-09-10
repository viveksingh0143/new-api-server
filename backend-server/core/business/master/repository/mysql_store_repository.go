package repository

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
	"github.com/vamika-digital/wms-api-server/global/drivers"
)

type SQLStoreRepository struct {
	DB *sqlx.DB
}

func NewSQLStoreRepository(conn drivers.Connection) (StoreRepository, error) {
	db := conn.GetDB()
	if db == nil {
		return nil, errors.New("failed to get database connection")
	}
	return &SQLStoreRepository{DB: conn.GetDB()}, nil
}

func (r *SQLStoreRepository) getFilterQueryWithArgs(page int, pageSize int, sort string, filter *store.StoreFilterDto) (string, map[string]interface{}) {
	var queryBuffer bytes.Buffer
	args := make(map[string]interface{})

	if filter.Query != "" {
		queryBuffer.WriteString(" AND (code LIKE :query OR name LIKE :query)")
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
	if filter.Name != "" {
		queryBuffer.WriteString(" AND name LIKE :name")
		args["name"] = "%" + filter.Name + "%"
	}
	if filter.Status.IsValid() {
		queryBuffer.WriteString(" AND status = :status")
		args["status"] = filter.Status
	}
	if filter.Owner != nil && filter.Owner.ID > 0 {
		queryBuffer.WriteString(" AND owner_id = :owner_id")
		args["owner_id"] = filter.Owner.ID
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

func (r *SQLStoreRepository) GetTotalCount(filter *store.StoreFilterDto) (int, error) {
	var count int
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT COUNT(*) FROM stores WHERE deleted_at IS NULL")
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

func (r *SQLStoreRepository) GetAll(page int, pageSize int, sort string, filter *store.StoreFilterDto) ([]*domain.Store, error) {
	stores := make([]*domain.Store, 0)
	var queryBuffer bytes.Buffer

	queryBuffer.WriteString("SELECT * FROM stores WHERE deleted_at IS NULL")
	filterQuery, args := r.getFilterQueryWithArgs(page, pageSize, sort, filter)
	queryBuffer.WriteString(filterQuery)

	query := queryBuffer.String()
	namedQuery, err := r.DB.PrepareNamed(query)
	if err != nil {
		return nil, err
	}

	err = namedQuery.Select(&stores, args)
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (r *SQLStoreRepository) Create(store *domain.Store) error {
	if store.Code != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM stores WHERE code = ?", store.Code)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a store with the code %s already exists", store.Code)
		}
	}

	if store.Name != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM stores WHERE name = ?", store.Name)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a store with the name %s already exists", store.Name)
		}
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	if store.Owner != nil && store.Owner.ID > 0 {
		store.OwnerID = store.Owner.ID
	}

	query := `INSERT INTO stores (code, name, location, owner_id, status, last_updated_by) VALUES(:code, :name, :location, :owner_id, :status, :last_updated_by)`
	res, err := tx.NamedExec(query, store)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	store.ID = id
	return tx.Commit()
}

func (r *SQLStoreRepository) GetById(storeID int64) (*domain.Store, error) {
	store := &domain.Store{}
	err := r.DB.Get(store, "SELECT * FROM stores WHERE id = ? AND deleted_at IS NULL", storeID)
	return store, err
}

func (r *SQLStoreRepository) GetByCode(storeCode string) (*domain.Store, error) {
	store := &domain.Store{}
	err := r.DB.Get(store, "SELECT * FROM stores WHERE code = ? AND deleted_at IS NULL", storeCode)
	return store, err
}

func (r *SQLStoreRepository) Update(store *domain.Store) error {
	if store.Code != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM stores WHERE code = ? AND id != ?", store.Code, store.ID)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a store with the code %s already exists", store.Code)
		}
	}

	if store.Name != "" {
		var count int
		err := r.DB.Get(&count, "SELECT COUNT(*) FROM stores WHERE name = ? AND id != ?", store.Name, store.ID)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("a store with the name %s already exists", store.Name)
		}
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	if store.Owner != nil && store.Owner.ID > 0 {
		store.OwnerID = store.Owner.ID
	}
	query := "UPDATE stores SET code=:code, name=:name, location=:location, owner_id=:owner_id, status=:status, last_updated_by=:last_updated_by WHERE id=:id"
	_, err = tx.NamedExec(query, store)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *SQLStoreRepository) Delete(storeID int64) error {
	_, err := r.DB.Exec("UPDATE stores SET deleted_at = NOW() WHERE id = ?", storeID)
	return err
}

func (r *SQLStoreRepository) DeleteByIDs(storeIDs []int64) error {
	if len(storeIDs) == 0 {
		return nil
	}

	query := "UPDATE stores SET deleted_at = NOW() WHERE id IN (?)"
	query, args, err := sqlx.In(query, storeIDs)

	if err != nil {
		return err
	}

	query = r.DB.Rebind(query)

	_, err = r.DB.Exec(query, args...)
	return err
}
