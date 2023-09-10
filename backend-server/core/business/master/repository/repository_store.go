package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
)

type StoreRepository interface {
	Create(store *domain.Store) error
	Update(store *domain.Store) error
	Delete(storeID int64) error
	DeleteByIDs(storeIDs []int64) error
	GetById(storeID int64) (*domain.Store, error)
	GetByCode(storeCode string) (*domain.Store, error)
	GetTotalCount(filter *store.StoreFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *store.StoreFilterDto) ([]*domain.Store, error)
}
