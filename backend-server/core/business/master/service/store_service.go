package service

import "github.com/vamika-digital/wms-api-server/core/business/master/dto/store"

type StoreService interface {
	GetAllStores(page int64, pageSize int64, sort string, filter *store.StoreFilterDto) ([]*store.StoreDto, int64, error)
	CreateStore(storeDto *store.StoreCreateDto) error
	GetStoreByID(storeID int64) (*store.StoreDto, error)
	GetStoreByCode(storeCode string) (*store.StoreDto, error)
	UpdateStore(storeID int64, store *store.StoreUpdateDto) error
	DeleteStore(storeID int64) error
	DeleteStoreByIDs(storeIDs []int64) error
}
