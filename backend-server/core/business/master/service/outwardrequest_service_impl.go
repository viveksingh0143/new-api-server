package service

import (
	"log"

	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/outwardrequest"
	"github.com/vamika-digital/wms-api-server/core/business/master/repository"
)

type OutwardRequestServiceImpl struct {
	OutwardRequestRepo      repository.OutwardRequestRepository
	CustomerRepo            repository.CustomerRepository
	ProductRepo             repository.ProductRepository
	OutwardRequestConverter *converter.OutwardRequestConverter
}

func NewOutwardRequestService(outwardrequestRepo repository.OutwardRequestRepository, customerRepo repository.CustomerRepository, productRepo repository.ProductRepository, outwardrequestConverter *converter.OutwardRequestConverter) OutwardRequestService {
	return &OutwardRequestServiceImpl{OutwardRequestRepo: outwardrequestRepo, CustomerRepo: customerRepo, ProductRepo: productRepo, OutwardRequestConverter: outwardrequestConverter}
}

func (s *OutwardRequestServiceImpl) GetAllOutwardRequests(page int16, pageSize int16, sort string, filter *outwardrequest.OutwardRequestFilterDto) ([]*outwardrequest.OutwardRequestDto, int64, error) {
	totalCount, err := s.OutwardRequestRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainOutwardRequests, err := s.OutwardRequestRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}

	if len(domainOutwardRequests) > 0 {
		outwardrequestIds := make([]int64, 0, len(domainOutwardRequests))
		for _, outwardrequest := range domainOutwardRequests {
			outwardrequestIds = append(outwardrequestIds, outwardrequest.ID)
		}

		outwardrequestItemsMap, err := s.OutwardRequestRepo.GetItemsForOutwardRequests(outwardrequestIds)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, 0, err
		}
		for i, outwardrequest := range domainOutwardRequests {
			if roles, ok := outwardrequestItemsMap[outwardrequest.ID]; ok {
				domainOutwardRequests[i].Items = roles
			}
		}
	}

	// Convert domain outwardrequests to DTOs. You can do this based on your requirements.
	var outwardrequestDtos []*outwardrequest.OutwardRequestDto = s.OutwardRequestConverter.ToDtoSlice(domainOutwardRequests)
	return outwardrequestDtos, int64(totalCount), nil
}

func (s *OutwardRequestServiceImpl) CreateOutwardRequest(outwardrequestDto *outwardrequest.OutwardRequestCreateDto) error {
	var newOutwardRequest *domain.OutwardRequest = s.OutwardRequestConverter.ToDomain(outwardrequestDto)
	err := s.OutwardRequestRepo.Create(newOutwardRequest)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *OutwardRequestServiceImpl) GetOutwardRequestByID(outwardrequestID int64) (*outwardrequest.OutwardRequestDto, error) {
	domainOutwardRequest, err := s.OutwardRequestRepo.GetById(outwardrequestID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	domainCustomer, err := s.CustomerRepo.GetById(domainOutwardRequest.CustomerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainOutwardRequest.Customer = domainCustomer

	items, err := s.OutwardRequestRepo.GetItemsForOutwardRequest(domainOutwardRequest.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	for _, item := range items {
		productDomain, err := s.ProductRepo.GetById(item.ProductID)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
		item.Product = productDomain
	}
	domainOutwardRequest.Items = items

	return s.OutwardRequestConverter.ToDto(domainOutwardRequest), nil
}

func (s *OutwardRequestServiceImpl) GetMinimalOutwardRequestByID(outwardrequestID int64) (*outwardrequest.OutwardRequestMinimalDto, error) {
	domainOutwardRequest, err := s.OutwardRequestRepo.GetById(outwardrequestID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.OutwardRequestConverter.ToMinimalDto(domainOutwardRequest), nil
}

func (s *OutwardRequestServiceImpl) GetOutwardRequestByCode(outwardrequestCode string) (*outwardrequest.OutwardRequestDto, error) {
	domainOutwardRequest, err := s.OutwardRequestRepo.GetByOrderNo(outwardrequestCode)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.OutwardRequestConverter.ToDto(domainOutwardRequest), nil
}

func (s *OutwardRequestServiceImpl) UpdateOutwardRequest(outwardrequestID int64, outwardrequestDto *outwardrequest.OutwardRequestUpdateDto) error {
	existingOutwardRequest, err := s.OutwardRequestRepo.GetById(outwardrequestID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	s.OutwardRequestConverter.ToUpdateDomain(existingOutwardRequest, outwardrequestDto)
	if err := s.OutwardRequestRepo.Update(existingOutwardRequest); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *OutwardRequestServiceImpl) DeleteOutwardRequest(outwardrequestID int64) error {
	if err := s.OutwardRequestRepo.Delete(outwardrequestID); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *OutwardRequestServiceImpl) DeleteOutwardRequestByIDs(outwardrequestIDs []int64) error {
	if err := s.OutwardRequestRepo.DeleteByIDs(outwardrequestIDs); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}
