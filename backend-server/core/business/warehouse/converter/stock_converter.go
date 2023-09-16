package converter

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	masterDomain "github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/stock"
)

type StockConverter struct {
	productConverter    converter.ProductConverter
	storeConverter      converter.StoreConverter
	containerConverter  converter.ContainerConverter
	batchLabelConverter BatchLabelConverter
}

func NewStockConverter(productConv converter.ProductConverter, storeConv converter.StoreConverter, containerConv converter.ContainerConverter, batchLabelConv BatchLabelConverter) *StockConverter {
	return &StockConverter{
		productConverter:    productConv,
		storeConverter:      storeConv,
		containerConverter:  containerConv,
		batchLabelConverter: batchLabelConv,
	}
}

func (c *StockConverter) ToDto(domainStock *domain.Stock) *stock.StockDto {
	stockDto := &stock.StockDto{
		ID:            domainStock.ID,
		ProductID:     domainStock.ProductID,
		StoreID:       domainStock.StoreID,
		BinID:         domainStock.BinID,
		PalletID:      domainStock.PalletID,
		RackID:        domainStock.RackID,
		BatchLabelID:  domainStock.BatchLabelID,
		Barcode:       domainStock.Barcode,
		BatchNo:       domainStock.BatchNo,
		UnitWeight:    domainStock.UnitWeight,
		Quantity:      domainStock.Quantity,
		MachineCode:   domainStock.MachineCode,
		StockInAt:     domainStock.StockInAt,
		StockOutAt:    domainStock.StockOutAt,
		Status:        domainStock.Status,
		CreatedAt:     domainStock.CreatedAt,
		UpdatedAt:     domainStock.UpdatedAt,
		LastUpdatedBy: domainStock.LastUpdatedBy,
	}
	return stockDto
}

func (c *StockConverter) ToDetailDto(
	domainStock *domain.Stock,
	domainProduct *masterDomain.Product,
	domainStore *masterDomain.Store,
	domainPallet *masterDomain.Container,
	domainBin *masterDomain.Container,
	domainRack *masterDomain.Container,
	domainBatchLabel *domain.BatchLabel,
) *stock.StockDetailDto {
	stockDto := &stock.StockDetailDto{
		ID:            domainStock.ID,
		ProductID:     domainStock.ProductID,
		StoreID:       domainStock.StoreID,
		BinID:         domainStock.BinID,
		PalletID:      domainStock.PalletID,
		RackID:        domainStock.RackID,
		BatchLabelID:  domainStock.BatchLabelID,
		Barcode:       domainStock.Barcode,
		BatchNo:       domainStock.BatchNo,
		UnitWeight:    domainStock.UnitWeight,
		Quantity:      domainStock.Quantity,
		MachineCode:   domainStock.MachineCode,
		StockInAt:     domainStock.StockInAt,
		StockOutAt:    domainStock.StockOutAt,
		Status:        domainStock.Status,
		CreatedAt:     domainStock.CreatedAt,
		UpdatedAt:     domainStock.UpdatedAt,
		LastUpdatedBy: domainStock.LastUpdatedBy,
	}

	if domainProduct != nil {
		stockDto.Product = c.productConverter.ToMinimalDto(domainProduct)
	}
	if domainStore != nil {
		stockDto.Store = c.storeConverter.ToMinimalDto(domainStore)
	}
	if domainPallet != nil {
		stockDto.Pallet = c.containerConverter.ToMinimalDto(domainPallet)
	}
	if domainBin != nil {
		stockDto.Bin = c.containerConverter.ToMinimalDto(domainBin)
	}
	if domainRack != nil {
		stockDto.Rack = c.containerConverter.ToMinimalDto(domainRack)
	}
	if domainBatchLabel != nil {
		stockDto.BatchLabel = c.batchLabelConverter.ToMinimalDto(domainBatchLabel)
	}

	return stockDto
}

func (c *StockConverter) ToDtoSlice(domainStocks []*domain.Stock) []*stock.StockDto {
	var stockDtos = make([]*stock.StockDto, 0)
	for _, domainStock := range domainStocks {
		stockDtos = append(stockDtos, c.ToDto(domainStock))
	}
	return stockDtos
}
