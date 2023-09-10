package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/admin/converter"
	adminDomain "github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
)

type StoreConverter struct {
	userConverter *converter.UserConverter
}

func NewStoreConverter(userConv *converter.UserConverter) *StoreConverter {
	return &StoreConverter{userConverter: userConv}
}

func (c *StoreConverter) ToMinimalDto(domainStore *domain.Store) *store.StoreMinimalDto {
	storeDto := &store.StoreMinimalDto{
		ID:     domainStore.ID,
		Code:   domainStore.Code,
		Name:   domainStore.Name,
		Status: domainStore.Status,
	}
	return storeDto
}

func (c *StoreConverter) ToDto(domainStore *domain.Store) *store.StoreDto {
	storeDto := &store.StoreDto{
		ID:            domainStore.ID,
		Code:          domainStore.Code,
		Name:          domainStore.Name,
		Location:      domainStore.Location,
		Status:        domainStore.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainStore.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainStore.UpdatedAt),
		LastUpdatedBy: domainStore.LastUpdatedBy,
	}
	storeDto.Owner = c.userConverter.ToMinimalDto(domainStore.Owner)
	return storeDto
}

func (c *StoreConverter) ToDtoSlice(domainStores []*domain.Store) []*store.StoreDto {
	var storeDtos = make([]*store.StoreDto, 0)
	for _, domainStore := range domainStores {
		storeDtos = append(storeDtos, c.ToDto(domainStore))
	}
	return storeDtos
}

func (c *StoreConverter) ToDomain(storeDto *store.StoreCreateDto) *domain.Store {
	domainStore := &domain.Store{
		Code:          storeDto.Code,
		Name:          storeDto.Name,
		Location:      storeDto.Location,
		Status:        storeDto.Status,
		LastUpdatedBy: storeDto.LastUpdatedBy,
	}
	if storeDto.Owner != nil && storeDto.Owner.ID != 0 {
		domainStore.Owner = &adminDomain.User{ID: storeDto.Owner.ID}
	}
	return domainStore
}

func (c *StoreConverter) ToUpdateDomain(domainStore *domain.Store, storeDto *store.StoreUpdateDto) {
	domainStore.Code = storeDto.Code
	domainStore.Name = storeDto.Name
	domainStore.Location = storeDto.Location
	if storeDto.Status.IsValid() {
		domainStore.Status = storeDto.Status
	}
	domainStore.LastUpdatedBy = storeDto.LastUpdatedBy
	if storeDto.Owner != nil && storeDto.Owner.ID > 0 {
		domainStore.Owner = &adminDomain.User{ID: storeDto.Owner.ID}
	}
}
