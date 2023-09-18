package service

import (
	"log"

	"github.com/vamika-digital/wms-api-server/core/base/helpers"
	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/requisition"
	"github.com/vamika-digital/wms-api-server/core/business/master/repository"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/reports"
	warehouseRepository "github.com/vamika-digital/wms-api-server/core/business/warehouse/repository"
)

type RequisitionServiceImpl struct {
	RequisitionRepo      repository.RequisitionRepository
	InventoryRepo        warehouseRepository.InventoryRepository
	StoreRepo            repository.StoreRepository
	ProductRepo          repository.ProductRepository
	RequisitionConverter *converter.RequisitionConverter
}

func NewRequisitionService(requisitionRepo repository.RequisitionRepository, inventoryRepo warehouseRepository.InventoryRepository, storeRepo repository.StoreRepository, productRepo repository.ProductRepository, requisitionConverter *converter.RequisitionConverter) RequisitionService {
	return &RequisitionServiceImpl{RequisitionRepo: requisitionRepo, InventoryRepo: inventoryRepo, StoreRepo: storeRepo, ProductRepo: productRepo, RequisitionConverter: requisitionConverter}
}

func (s *RequisitionServiceImpl) GetAllRequisitions(page int16, pageSize int16, sort string, filter *requisition.RequisitionFilterDto) ([]*requisition.RequisitionDto, int64, error) {
	totalCount, err := s.RequisitionRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainRequisitions, err := s.RequisitionRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}

	if len(domainRequisitions) > 0 {
		jobOrderIds := make([]int64, 0, len(domainRequisitions))
		for _, jobOrder := range domainRequisitions {
			jobOrderIds = append(jobOrderIds, jobOrder.ID)
		}

		jobOrderItemsMap, err := s.RequisitionRepo.GetItemsForRequisitions(jobOrderIds)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, 0, err
		}
		for i, jobOrder := range domainRequisitions {
			if roles, ok := jobOrderItemsMap[jobOrder.ID]; ok {
				domainRequisitions[i].Items = roles
			}
		}
	}

	// Convert domain requisitions to DTOs. You can do this based on your requirements.
	var requisitionDtos []*requisition.RequisitionDto = s.RequisitionConverter.ToDtoSlice(domainRequisitions)
	return requisitionDtos, int64(totalCount), nil
}

func (s *RequisitionServiceImpl) CreateRequisition(requisitionDto *requisition.RequisitionCreateDto) error {
	var newRequisition *domain.Requisition = s.RequisitionConverter.ToDomain(requisitionDto)
	err := s.RequisitionRepo.Create(newRequisition)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *RequisitionServiceImpl) GetRequisitionByID(requisitionID int64) (*requisition.RequisitionDto, error) {
	domainRequisition, err := s.RequisitionRepo.GetById(requisitionID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	domainStore, err := s.StoreRepo.GetById(domainRequisition.StoreID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainRequisition.Store = domainStore

	items, err := s.RequisitionRepo.GetItemsForRequisition(requisitionID)
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
	domainRequisition.Items = items

	return s.RequisitionConverter.ToDto(domainRequisition), nil
}

func (s *RequisitionServiceImpl) GetMinimalRequisitionByID(requisitionID int64) (*requisition.RequisitionMinimalDto, error) {
	domainRequisition, err := s.RequisitionRepo.GetById(requisitionID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.RequisitionConverter.ToMinimalDto(domainRequisition), nil
}

func (s *RequisitionServiceImpl) UpdateRequisition(requisitionID int64, requisitionDto *requisition.RequisitionUpdateDto) error {
	existingRequisition, err := s.RequisitionRepo.GetById(requisitionID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	s.RequisitionConverter.ToUpdateDomain(existingRequisition, requisitionDto)
	if err := s.RequisitionRepo.Update(existingRequisition); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *RequisitionServiceImpl) DeleteRequisition(requisitionID int64) error {
	if err := s.RequisitionRepo.Delete(requisitionID); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *RequisitionServiceImpl) DeleteRequisitionByIDs(requisitionIDs []int64) error {
	if err := s.RequisitionRepo.DeleteByIDs(requisitionIDs); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *RequisitionServiceImpl) GetRequisitionByCode(requisitionCode string) (*requisition.RequisitionDto, []*reports.InventoryRackStatusDetail, []*reports.InventoryPalletStatusDetail, error) {
	domainRequisition, err := s.RequisitionRepo.GetByOrderNo(requisitionCode)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}

	domainStore, err := s.StoreRepo.GetById(domainRequisition.StoreID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}
	domainRequisition.Store = domainStore

	items, err := s.RequisitionRepo.GetItemsForRequisition(domainRequisition.ID)
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
	domainRequisition.Items = items

	var productIds []int64
	for _, item := range domainRequisition.Items {
		productIds = append(productIds, item.ProductID)
	}

	// Get Locked Quantity
	requestLockedReports, err := s.InventoryRepo.GetLockedInventoryDetailForRequest(domainRequisition.ID, helpers.GetNameOfTheVariable(domainRequisition))
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}

	// Attached Locked Quantity
	for _, requisitionItem := range domainRequisition.Items {
		for i := 0; i < len(requestLockedReports); i++ {
			if requestLockedReports[i].ProductID == requisitionItem.ProductID {
				requisitionItem.LockedQuantity = requestLockedReports[i].LockCount
			}
		}
	}

	palletItems, err := s.InventoryRepo.GetLockedInventoryStocksWithPalletForRequest(domainRequisition.ID, helpers.GetNameOfTheVariable(domainRequisition))
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, nil, err
	}

	// Attached Stockout Quantity, Required Quantity
	for _, requisitionItem := range domainRequisition.Items {
		requiredQty := requisitionItem.Quantity
		for i := 0; i < len(palletItems) && requiredQty > 0; i++ {
			if palletItems[i].ProductID == requisitionItem.ProductID {
				requiredQty -= palletItems[i].StockOutCount
				stockOutQty := min(requiredQty, palletItems[i].LockCount)
				palletItems[i].RequiredStocks = stockOutQty
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
	for _, requisitionItem := range domainRequisition.Items {
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
	return s.RequisitionConverter.ToDto(domainRequisition), reports, palletItems, nil
}
