package service

import (
	adminDomain "github.com/vamika-digital/wms-api-server/core/business/admin/domain"
	adminRepository "github.com/vamika-digital/wms-api-server/core/business/admin/repository"
	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/store"
	masterRepository "github.com/vamika-digital/wms-api-server/core/business/master/repository"
)

type StoreServiceImpl struct {
	StoreRepo      masterRepository.StoreRepository
	UserRepo       adminRepository.UserRepository
	StoreConverter *converter.StoreConverter
}

func NewStoreService(storeRepo masterRepository.StoreRepository, userRepo adminRepository.UserRepository, storeConverter *converter.StoreConverter) StoreService {
	return &StoreServiceImpl{StoreRepo: storeRepo, UserRepo: userRepo, StoreConverter: storeConverter}
}

func (s *StoreServiceImpl) GetAllStores(page int64, pageSize int64, sort string, filter *store.StoreFilterDto) ([]*store.StoreDto, int64, error) {
	totalCount, err := s.StoreRepo.GetTotalCount(filter)
	if err != nil {
		return nil, 0, err
	}
	domainStores, err := s.StoreRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		return nil, 0, err
	}
	if len(domainStores) > 0 {
		uniqueOwnerIDsMap := make(map[int64]bool)
		var uniqueOwnerIDs []int64

		for _, store := range domainStores {
			if _, exists := uniqueOwnerIDsMap[store.OwnerID]; !exists {
				uniqueOwnerIDsMap[store.OwnerID] = true
				uniqueOwnerIDs = append(uniqueOwnerIDs, store.OwnerID)
			}
		}

		owners, err := s.UserRepo.GetByIds(uniqueOwnerIDs)
		if err != nil {
			return nil, 0, err
		}

		ownerMap := make(map[int64]*adminDomain.User)
		for _, owner := range owners {
			ownerMap[owner.ID] = owner
		}

		for i, store := range domainStores {
			if owner, ok := ownerMap[store.OwnerID]; ok {
				domainStores[i].Owner = owner
			}
		}
	}

	var storeDtos []*store.StoreDto = s.StoreConverter.ToDtoSlice(domainStores)
	return storeDtos, int64(totalCount), nil
}

func (s *StoreServiceImpl) CreateStore(storeDto *store.StoreCreateDto) error {
	var newStore *domain.Store = s.StoreConverter.ToDomain(storeDto)
	if err := s.StoreRepo.Create(newStore); err != nil {
		return err
	}
	return nil
}

func (s *StoreServiceImpl) GetStoreByID(storeID int64) (*store.StoreDto, error) {
	domainStore, err := s.StoreRepo.GetById(storeID)
	if err != nil {
		return nil, err
	}

	owner, err := s.UserRepo.GetById(domainStore.OwnerID)
	if err != nil {
		return nil, err
	}
	domainStore.Owner = owner
	return s.StoreConverter.ToDto(domainStore), nil
}

func (s *StoreServiceImpl) GetStoreByCode(storeCode string) (*store.StoreDto, error) {
	domainStore, err := s.StoreRepo.GetByCode(storeCode)
	if err != nil {
		return nil, err
	}

	owner, err := s.UserRepo.GetById(domainStore.OwnerID)
	if err != nil {
		return nil, err
	}
	domainStore.Owner = owner
	return s.StoreConverter.ToDto(domainStore), nil
}

func (s *StoreServiceImpl) UpdateStore(storeID int64, storeDto *store.StoreUpdateDto) error {
	existingStore, err := s.StoreRepo.GetById(storeID)
	if err != nil {
		return err
	}

	s.StoreConverter.ToUpdateDomain(existingStore, storeDto)
	if err := s.StoreRepo.Update(existingStore); err != nil {
		return err
	}
	return nil
}

func (s *StoreServiceImpl) DeleteStore(storeID int64) error {
	if err := s.StoreRepo.Delete(storeID); err != nil {
		return err
	}
	return nil
}

func (s *StoreServiceImpl) DeleteStoreByIDs(storeIDs []int64) error {
	if err := s.StoreRepo.DeleteByIDs(storeIDs); err != nil {
		return err
	}
	return nil
}
