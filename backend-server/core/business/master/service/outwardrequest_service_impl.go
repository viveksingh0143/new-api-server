package service

import (
	"log"

	"github.com/vamika-digital/wms-api-server/core/base/helpers"
	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/outwardrequest"
	"github.com/vamika-digital/wms-api-server/core/business/master/repository"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/reports"
	warehouseRepository "github.com/vamika-digital/wms-api-server/core/business/warehouse/repository"
)

type OutwardRequestServiceImpl struct {
	OutwardRequestRepo      repository.OutwardRequestRepository
	InventoryRepo           warehouseRepository.InventoryRepository
	CustomerRepo            repository.CustomerRepository
	ProductRepo             repository.ProductRepository
	OutwardRequestConverter *converter.OutwardRequestConverter
}

func NewOutwardRequestService(outwardrequestRepo repository.OutwardRequestRepository, inventoryRepo warehouseRepository.InventoryRepository, customerRepo repository.CustomerRepository, productRepo repository.ProductRepository, outwardrequestConverter *converter.OutwardRequestConverter) OutwardRequestService {
	return &OutwardRequestServiceImpl{OutwardRequestRepo: outwardrequestRepo, InventoryRepo: inventoryRepo, CustomerRepo: customerRepo, ProductRepo: productRepo, OutwardRequestConverter: outwardrequestConverter}
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

func (s *OutwardRequestServiceImpl) GetOutwardRequestByCode(outwardrequestCode string) (*outwardrequest.OutwardRequestDto, []*reports.InventoryRackStatusDetail, []*reports.InventoryBinStatusDetail, error) {
	domainOutwardRequest, err := s.OutwardRequestRepo.GetByOrderNo(outwardrequestCode)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}

	domainCustomer, err := s.CustomerRepo.GetById(domainOutwardRequest.CustomerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}
	domainOutwardRequest.Customer = domainCustomer

	items, err := s.OutwardRequestRepo.GetItemsForOutwardRequest(domainOutwardRequest.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}

	for _, item := range items {
		productDomain, err := s.ProductRepo.GetById(item.ProductID)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, nil, nil, err
		}
		item.Product = productDomain
	}
	domainOutwardRequest.Items = items

	var productIds []int64
	for _, item := range domainOutwardRequest.Items {
		productIds = append(productIds, item.ProductID)
	}

	// Get Locked Quantity
	requestLockedReports, err := s.InventoryRepo.GetLockedInventoryDetailForRequest(domainOutwardRequest.ID, helpers.GetNameOfTheVariable(domainOutwardRequest))
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}

	// Attached Locked Quantity
	for _, requisitionItem := range domainOutwardRequest.Items {
		for i := 0; i < len(requestLockedReports); i++ {
			if requestLockedReports[i].ProductID == requisitionItem.ProductID {
				requisitionItem.LockedQuantity = requestLockedReports[i].LockCount
			}
		}
	}

	binItems, err := s.InventoryRepo.GetLockedInventoryStocksWithBinForRequest(domainOutwardRequest.ID, helpers.GetNameOfTheVariable(domainOutwardRequest))
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}

	// Attached Stockout Quantity, Required Quantity
	for _, requisitionItem := range domainOutwardRequest.Items {
		requiredQty := requisitionItem.Quantity
		for i := 0; i < len(binItems) && requiredQty > 0; i++ {
			if binItems[i].ProductID == requisitionItem.ProductID {
				requiredQty -= binItems[i].StockOutCount
				stockOutQty := min(requiredQty, binItems[i].LockCount)
				binItems[i].RequiredStocks = stockOutQty
				requiredQty -= stockOutQty
			}
		}
		requisitionItem.PendingQuantity = requiredQty
	}

	// Get Locked Rack Quantity
	reports, err := s.InventoryRepo.GetInventoryDetailForProductIds(productIds)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}

	// Attached Stockout Quantity, Required Quantity
	for _, requisitionItem := range domainOutwardRequest.Items {
		requiredQty := requisitionItem.PendingQuantity - requisitionItem.LockedQuantity
		for i := 0; i < len(reports) && requiredQty > 0; i++ {
			if reports[i].ProductID == requisitionItem.ProductID {
				stockOutQty := min(requiredQty, reports[i].StockCount)
				reports[i].RequiredStocks = stockOutQty
				requiredQty -= stockOutQty
			}
		}

		if requiredQty == 0 {
			for j := len(reports) - 1; j >= 0; j-- {
				if reports[j].ProductID == requisitionItem.ProductID {
					if reports[j].RequiredStocks <= 0 {
						reports = append(reports[:j], reports[j+1:]...)
					}
				}
			}
		}
	}
	return s.OutwardRequestConverter.ToDto(domainOutwardRequest), reports, binItems, nil
}
