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
	if domainStore == nil {
		return nil
	}
	storeDto := &store.StoreMinimalDto{
		ID:         domainStore.ID,
		Code:       domainStore.Code,
		Name:       domainStore.Name,
		StoreTypes: domainStore.GetStoreTypes(),
		Status:     domainStore.Status,
	}
	return storeDto
}

func (c *StoreConverter) ToDto(domainStore *domain.Store) *store.StoreDto {
	storeDto := &store.StoreDto{
		ID:            domainStore.ID,
		Code:          domainStore.Code,
		Name:          domainStore.Name,
		Location:      domainStore.Location,
		StoreTypes:    domainStore.GetStoreTypes(),
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

func (c *StoreConverter) ToMinimalDtoSlice(domainStores []*domain.Store) []*store.StoreMinimalDto {
	var storeDtos = make([]*store.StoreMinimalDto, 0)
	for _, domainStore := range domainStores {
		storeDtos = append(storeDtos, c.ToMinimalDto(domainStore))
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
	domainStore.SetStoreTypes(storeDto.StoreTypes)
	if storeDto.Owner != nil && storeDto.Owner.ID != 0 {
		domainStore.Owner = &adminDomain.User{ID: storeDto.Owner.ID}
	}
	return domainStore
}

func (c *StoreConverter) ToUpdateDomain(domainStore *domain.Store, storeDto *store.StoreUpdateDto) {
	domainStore.Code = storeDto.Code
	domainStore.Name = storeDto.Name
	domainStore.Location = storeDto.Location
	domainStore.SetStoreTypes(storeDto.StoreTypes)
	if storeDto.Status.IsValid() {
		domainStore.Status = storeDto.Status
	}
	domainStore.LastUpdatedBy = storeDto.LastUpdatedBy
	if storeDto.Owner != nil && storeDto.Owner.ID > 0 {
		domainStore.Owner = &adminDomain.User{ID: storeDto.Owner.ID}
	}
}
