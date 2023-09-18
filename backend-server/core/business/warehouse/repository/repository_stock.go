package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/stock"
)

type StockRepository interface {
	GetByBarcode(stockBarcode string) (*domain.Stock, error)
	GetById(stockID int64) (*domain.Stock, error)
	GetTotalCount(filter *stock.StockFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *stock.StockFilterDto) ([]*domain.Stock, error)
}
