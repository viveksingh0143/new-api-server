package service

import (
	"log"

	masterDomain "github.com/vamika-digital/wms-api-server/core/business/master/domain"
	masterRepository "github.com/vamika-digital/wms-api-server/core/business/master/repository"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/converter"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/stock"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/repository"
)

type StockServiceImpl struct {
	StockRepo      repository.StockRepository
	ProductRepo    masterRepository.ProductRepository
	StoreRepo      masterRepository.StoreRepository
	ContainerRepo  masterRepository.ContainerRepository
	BatchLabelRepo repository.BatchLabelRepository
	StockConverter *converter.StockConverter
}

func NewStockService(stockRepo repository.StockRepository, productRepo masterRepository.ProductRepository, storeRepo masterRepository.StoreRepository, containerRepo masterRepository.ContainerRepository, batchLabelRepo repository.BatchLabelRepository, stockConverter *converter.StockConverter) StockService {
	return &StockServiceImpl{StockRepo: stockRepo, ProductRepo: productRepo, StoreRepo: storeRepo, ContainerRepo: containerRepo, BatchLabelRepo: batchLabelRepo, StockConverter: stockConverter}
}

func (s *StockServiceImpl) GetAllStocks(page int16, pageSize int16, sort string, filter *stock.StockFilterDto) ([]*stock.StockDto, int64, error) {
	totalCount, err := s.StockRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainStocks, err := s.StockRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}

	var stockDtos []*stock.StockDto = s.StockConverter.ToDtoSlice(domainStocks)
	return stockDtos, int64(totalCount), nil
}

func (s *StockServiceImpl) GetStockByID(stockID int64) (*stock.StockDetailDto, error) {
	domainStock, err := s.StockRepo.GetById(stockID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	var domainProduct *masterDomain.Product
	var domainStore *masterDomain.Store
	var domainPallet *masterDomain.Container
	var domainBin *masterDomain.Container
	var domainRack *masterDomain.Container
	var domainBatchLabel *domain.BatchLabel

	domainProduct, err = s.ProductRepo.GetById(domainStock.ProductID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainStore, err = s.StoreRepo.GetById(domainStock.StoreID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	if domainStock.PalletID.Valid {
		domainPallet, err = s.ContainerRepo.GetById(domainStock.PalletID.Int64)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
	}
	if domainStock.BinID.Valid {
		domainBin, err = s.ContainerRepo.GetById(domainStock.BinID.Int64)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
	}
	if domainStock.RackID.Valid {
		domainRack, err = s.ContainerRepo.GetById(domainStock.RackID.Int64)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
	}
	if domainStock.BatchLabelID.Valid {
		domainBatchLabel, err = s.BatchLabelRepo.GetById(domainStock.BatchLabelID.Int64)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, err
		}
	}

	return s.StockConverter.ToDetailDto(domainStock, domainProduct, domainStore, domainPallet, domainBin, domainRack, domainBatchLabel), nil
}
