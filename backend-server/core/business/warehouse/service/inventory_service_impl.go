package service

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customerrors"
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/base/helpers"
	masterDomain "github.com/vamika-digital/wms-api-server/core/business/master/domain"
	masterRepository "github.com/vamika-digital/wms-api-server/core/business/master/repository"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/converter"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/inventory"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/reports"
	warehouseRepository "github.com/vamika-digital/wms-api-server/core/business/warehouse/repository"
)

type InventoryServiceImpl struct {
	BatchLabelService  BatchLabelService
	RequisitionRepo    masterRepository.RequisitionRepository
	OutwardRequestRepo masterRepository.OutwardRequestRepository
	StockRepo          warehouseRepository.StockRepository
	InventoryRepo      warehouseRepository.InventoryRepository
	ProductRepo        masterRepository.ProductRepository
	StoreRepo          masterRepository.StoreRepository
	ContainerRepo      masterRepository.ContainerRepository
	WarehouseRepo      warehouseRepository.BatchLabelRepository
	InventoryConverter *converter.InventoryConverter
}

func NewInventoryService(batchLabelService BatchLabelService, requisitionRepo masterRepository.RequisitionRepository, outwardRequestRepo masterRepository.OutwardRequestRepository, stockRepo warehouseRepository.StockRepository, inventoryRepo warehouseRepository.InventoryRepository, productRepo masterRepository.ProductRepository, storeRepo masterRepository.StoreRepository, containerRepo masterRepository.ContainerRepository, inventoryConverter *converter.InventoryConverter) InventoryService {
	return &InventoryServiceImpl{BatchLabelService: batchLabelService, RequisitionRepo: requisitionRepo, OutwardRequestRepo: outwardRequestRepo, StockRepo: stockRepo, InventoryRepo: inventoryRepo, ProductRepo: productRepo, StoreRepo: storeRepo, ContainerRepo: containerRepo, InventoryConverter: inventoryConverter}
}

func (s *InventoryServiceImpl) GetAllProductsWithStockCounts(page int16, pageSize int16, sort string, filter *inventory.InventoryFilterDto) ([]*inventory.InventoryDto, int64, error) {
	if filter.ContainerCode != "" {
		containerDomain, err := s.ContainerRepo.GetByCode(filter.ContainerCode)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, 0, err
		}
		filter.ContainerID = containerDomain.ID
		filter.ContainerType = containerDomain.ContainerType.String()
	}

	totalCount, err := s.InventoryRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainInventorys, err := s.InventoryRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	// Convert domain inventorys to DTOs. You can do this based on your requirements.
	var inventoryDtos []*inventory.InventoryDto = s.InventoryConverter.ToDtoSlice(domainInventorys)
	return inventoryDtos, int64(totalCount), nil
}

func (s *InventoryServiceImpl) GetInventoryByID(inventoryID int64) (*inventory.InventoryDto, error) {
	domainInventory, err := s.InventoryRepo.GetById(inventoryID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.InventoryConverter.ToDto(domainInventory), nil
}

func (s *InventoryServiceImpl) CreateRawMaterialStock(rmStockForm *inventory.InventoryRMStockCreateDto) error {
	domainProduct, err := s.ProductRepo.GetById(rmStockForm.ProductID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	domainStore, err := s.StoreRepo.GetById(rmStockForm.StoreID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	domainPallet, err := s.ContainerRepo.GetByCode(rmStockForm.PalletCode)
	{
		needToCreate := false
		if err != nil {
			switch err.(type) {
			case *customerrors.NotFoundError:
				needToCreate = rmStockForm.CreatePallet.Valid && rmStockForm.CreatePallet.Bool
			default:
				log.Printf("%+v\n", err)
				return err
			}
		}

		if needToCreate {
			domainPallet = &masterDomain.Container{
				ContainerType: customtypes.PALLET_TYPE,
				Code:          rmStockForm.PalletCode,
				Name:          rmStockForm.PalletCode,
				Status:        customtypes.Enable,
				LastUpdatedBy: rmStockForm.LastUpdatedBy,
			}
			err = s.ContainerRepo.Create(domainPallet)
			if err != nil {
				return err
			}
			domainPallet, err = s.ContainerRepo.GetByCode(rmStockForm.PalletCode)
			if err != nil {
				return err
			}
		}

		if err != nil {
			log.Printf("%+v\n", err)
			return err
		}

		if !domainPallet.ContainerType.IsPalletType() {
			return errors.New("given code is not for pallet")
		}
	}

	if domainPallet.IsFull() {
		return fmt.Errorf("pallet reached to it's maximum capacity of %d", domainPallet.Info().MaxCapacity)
	}

	resourceType := helpers.GetNameOfTheVariable(domainProduct)
	if domainPallet.IsConnectedWithResource() && !domainPallet.IsResourceMatch(domainProduct.ID, resourceType) {
		return errors.New("pallet already attached with different item")
	}
	var stockForm *domain.Stock = s.InventoryConverter.ToStockDomain(rmStockForm)
	stockForm.PalletID = customtypes.NewValidNullInt64(domainPallet.ID)
	stockForm.ProductID = domainProduct.ID
	stockForm.StoreID = domainStore.ID
	stockForm.Barcode = helpers.GenerateBarcode(domainProduct.Code)
	stockForm.UnitWeight = domainProduct.UnitWeight
	stockForm.PackageQuantity = 1
	stockForm.Status = customtypes.STOCK_IN_APPROVAL

	domainPallet.IsApproved = false
	domainPallet.IncreamentStock(rmStockForm.Quantity, domainProduct.ID, resourceType)

	err = s.InventoryRepo.CreateRawMaterialStock(stockForm, domainPallet)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *InventoryServiceImpl) CreateFinishedGoodsStock(fdStockForm *inventory.InventoryFDStockCreateDto) error {
	domainStore, err := s.StoreRepo.GetById(fdStockForm.StoreID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	domainBin, err := s.ContainerRepo.GetByCode(fdStockForm.BinCode)
	if err != nil {
		return err
	}
	if !domainBin.ContainerType.IsBinType() {
		return errors.New("given code is not for bin")
	}
	if domainBin.IsFull() {
		return fmt.Errorf("bin reached to it's maximum capacity of %d", domainBin.Info().MaxCapacity)
	}

	var fdStocks = make([]*domain.Stock, 0)
	var stickers = make([]*domain.LabelSticker, 0)
	for _, barcode := range fdStockForm.Barcodes {
		batchLabel, stickerDto, err := s.BatchLabelService.GetBatchLabelByBarcode(barcode)
		if err != nil {
			return fmt.Errorf("barcode '%s' sticker not found", barcode)
		}
		if stickerDto.IsUsed {
			return fmt.Errorf("barcode '%s' sticker already used", stickerDto.UUIDCode)
		}
		resourceType := helpers.GetNameOfTheVariable(&masterDomain.Product{})
		if domainBin.IsConnectedWithResource() && !domainBin.IsResourceMatch(batchLabel.Product.ID, resourceType) {
			return errors.New("bin already attached with different item")
		}

		fdStock := &domain.Stock{
			ProductID:       batchLabel.Product.ID,
			StoreID:         domainStore.ID,
			BinID:           customtypes.NewValidNullInt64(domainBin.ID),
			BatchLabelID:    customtypes.NewValidNullInt64(batchLabel.ID),
			Barcode:         barcode,
			BatchNo:         batchLabel.BatchNo,
			UnitWeight:      float64(batchLabel.UnitWeight),
			Quantity:        1,
			PackageQuantity: batchLabel.PackageQuantity,
			MachineCode:     batchLabel.Machine.Code,
			StockInAt:       time.Now(),
			Status:          customtypes.STOCK_IN,
			LastUpdatedBy:   fdStockForm.LastUpdatedBy,
		}
		domainBin.IncreamentStock(batchLabel.PackageQuantity, batchLabel.Product.ID, resourceType)
		stickers = append(stickers, &domain.LabelSticker{
			ID:     stickerDto.ID,
			IsUsed: true,
		})
		fdStocks = append(fdStocks, fdStock)
	}

	err = s.InventoryRepo.CreateFinishedStocks(fdStocks, stickers, domainBin)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *InventoryServiceImpl) AttachContainer(sourceCode string, destinationCode string) error {
	destinationContainer, err := s.ContainerRepo.GetByCode(destinationCode)
	if err != nil {
		return err
	}
	if destinationContainer.IsFull() {
		return fmt.Errorf("%s reached to it's maximum capacity of %d", strings.ToLower(destinationContainer.ContainerType.String()), destinationContainer.Info().MaxCapacity)
	}

	sourceContainer, err := s.ContainerRepo.GetByCode(sourceCode)
	if err != nil {
		return err
	}

	if !sourceContainer.IsApproved {
		return fmt.Errorf("%s is not approved", strings.ToLower(sourceContainer.Code))
	}

	if !sourceContainer.IsFull() {
		return fmt.Errorf("%s is not marked full", strings.ToLower(sourceContainer.Code))
	}

	if !destinationContainer.Info().ContainsType(sourceContainer.ContainerType.String()) {
		return fmt.Errorf("%s can not contains %s", strings.ToLower(destinationContainer.ContainerType.String()), strings.ToLower(sourceContainer.ContainerType.String()))
	}

	resourceType := helpers.GetNameOfTheVariable(sourceContainer)
	if destinationContainer.IsConnectedWithResource() && !destinationContainer.IsResourceMatch(sourceContainer.ID, resourceType) {
		return fmt.Errorf("%s already attached with different item", strings.ToLower(destinationContainer.ContainerType.String()))
	}

	count, err := s.ContainerRepo.AttachedCount(sourceContainer.ID, resourceType)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("%s already attached with different container", strings.ToLower(sourceContainer.ContainerType.String()))
	}

	destinationContainer.IncreamentStock(1, sourceContainer.ID, resourceType)
	err = s.InventoryRepo.AttachContainer(destinationContainer, sourceContainer)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *InventoryServiceImpl) DeattachRackContainer(rackCode string, requestID int64, requestName string) error {
	rackContainer, err := s.ContainerRepo.GetByCode(rackCode)
	if err != nil {
		return err
	}
	if !rackContainer.ContainerType.IsRackType() {
		return fmt.Errorf("%s is not rack container", rackCode)
	}
	if rackContainer.Level.IsEmpty() {
		return fmt.Errorf("%s is empty", strings.ToLower(rackContainer.ContainerType.String()))
	}

	requestName, err = helpers.GetRequestType(requestName)
	if err != nil {
		return err
	}

	rackContainer.DecreamentStock(1)
	err = s.InventoryRepo.DeattachRackContainer(rackContainer, requestID, requestName)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *InventoryServiceImpl) ProcessRawMaterialStockout(palletCode string, quantity int64, requestID int64, requestName string) error {
	palletContainer, err := s.ContainerRepo.GetByCode(palletCode)
	if err != nil {
		return err
	}
	if !palletContainer.ContainerType.IsPalletType() {
		return fmt.Errorf("%s is not pallet container", palletCode)
	}
	if palletContainer.Level.IsEmpty() {
		return fmt.Errorf("%s is empty", strings.ToLower(palletContainer.ContainerType.String()))
	}
	if quantity > palletContainer.ItemsCount {
		return fmt.Errorf("%s have maximum %d items", palletContainer.Code, palletContainer.ItemsCount)
	}
	palletContainer.DecreamentStock(quantity)

	requestName, err = helpers.GetRequestType(requestName)
	if err != nil {
		return err
	}

	domainRequisition, err := s.RequisitionRepo.GetById(requestID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	items, err := s.RequisitionRepo.GetItemsForRequisition(domainRequisition.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	for _, item := range items {
		productDomain, err := s.ProductRepo.GetById(item.ProductID)
		if err != nil {
			log.Printf("%+v\n", err)
			return err
		}
		item.Product = productDomain
	}
	domainRequisition.Items = items

	palletItems, err := s.InventoryRepo.GetLockedInventoryStocksWithPalletForRequest(requestID, requestName)
	if err != nil {
		return err
	}
	for i := 0; i < len(palletItems); i++ {
		if palletItems[i].PalletID == palletContainer.ID {
			if quantity > palletItems[i].LockCount {
				return fmt.Errorf("this pallet %s has only %d stocks to collect", palletContainer.Code, palletItems[i].LockCount)
			}
		}
	}

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
	}

	var palletStockDetail *reports.InventoryPalletStatusDetail
	for i := 0; i < len(palletItems); i++ {
		if palletItems[i].PalletID == palletContainer.ID {
			palletStockDetail = palletItems[i]
			break
		}
	}

	if palletStockDetail == nil {
		return fmt.Errorf("this pallet %s is not part of the requisition %s", palletContainer.Code, domainRequisition.OrderNo)
	}
	if quantity > palletStockDetail.RequiredStocks {
		return fmt.Errorf("you can withdraw only %d from this pallet %s", palletStockDetail.RequiredStocks, palletContainer.Code)
	}
	if palletStockDetail.StockDispatchingCount < palletStockDetail.StockOutCount+quantity {
		palletContainer.Clear()
	}

	err = s.InventoryRepo.StockoutRawMaterial(palletContainer, quantity, requestID, requestName)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *InventoryServiceImpl) ProcessFinishedGoodStockout(barcode string, requestID int64, requestName string) error {
	requestName, err := helpers.GetRequestType(requestName)
	if err != nil {
		return err
	}

	domainOutwardRequest, err := s.OutwardRequestRepo.GetById(requestID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	items, err := s.OutwardRequestRepo.GetItemsForOutwardRequest(domainOutwardRequest.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	for _, item := range items {
		productDomain, err := s.ProductRepo.GetById(item.ProductID)
		if err != nil {
			log.Printf("%+v\n", err)
			return err
		}
		item.Product = productDomain
	}
	domainOutwardRequest.Items = items

	domainStock, err := s.StockRepo.GetByBarcode(barcode)
	if err != nil {
		log.Printf("%+v\n", err)
		return fmt.Errorf("given barcode '%s' not found", barcode)
	}
	if !domainStock.BinID.Valid {
		return fmt.Errorf("given barcode '%s' not attached with bin", barcode)
	}

	binContainer, err := s.ContainerRepo.GetById(domainStock.BinID.Int64)
	if err != nil {
		return err
	}
	if !binContainer.ContainerType.IsBinType() {
		return fmt.Errorf("%s is not bin container", binContainer.Code)
	}
	if binContainer.Level.IsEmpty() {
		return fmt.Errorf("%s is empty", strings.ToLower(binContainer.ContainerType.String()))
	}
	if domainStock.Quantity > binContainer.ItemsCount {
		return fmt.Errorf("%s have maximum %d items", binContainer.Code, binContainer.ItemsCount)
	}
	binContainer.DecreamentStock(domainStock.Quantity)

	binItems, err := s.InventoryRepo.GetLockedInventoryStocksWithBinForRequest(requestID, requestName)
	if err != nil {
		return err
	}
	for i := 0; i < len(binItems); i++ {
		if binItems[i].BinID == binContainer.ID {
			if domainStock.Quantity > binItems[i].LockCount {
				return fmt.Errorf("this bin %s has only %d stocks to collect", binContainer.Code, binItems[i].LockCount)
			}
		}
	}

	for _, outwardRequestItem := range domainOutwardRequest.Items {
		requiredQty := outwardRequestItem.Quantity
		for i := 0; i < len(binItems) && requiredQty > 0; i++ {
			if binItems[i].ProductID == outwardRequestItem.ProductID {
				requiredQty -= binItems[i].StockOutCount
				stockOutQty := min(requiredQty, binItems[i].LockCount)
				binItems[i].RequiredStocks = stockOutQty
				requiredQty -= stockOutQty
			}
		}
	}

	var binStockDetail *reports.InventoryBinStatusDetail
	for i := 0; i < len(binItems); i++ {
		if binItems[i].BinID == binContainer.ID {
			binStockDetail = binItems[i]
			break
		}
	}

	if binStockDetail == nil {
		return fmt.Errorf("this bin %s is not part of the outward request %s", binContainer.Code, domainOutwardRequest.OrderNo)
	}
	if domainStock.Quantity > binStockDetail.RequiredStocks {
		return fmt.Errorf("you can withdraw only %d from this bin %s", binStockDetail.RequiredStocks, binContainer.Code)
	}
	if binStockDetail.StockDispatchingCount < binStockDetail.StockOutCount+domainStock.Quantity {
		binContainer.Clear()
	}

	err = s.InventoryRepo.FinishedGoodMaterial(binContainer, barcode)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}
