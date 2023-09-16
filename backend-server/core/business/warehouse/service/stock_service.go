package service

import (
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/stock"
)

type StockService interface {
	GetAllStocks(page int16, pageSize int16, sort string, filter *stock.StockFilterDto) ([]*stock.StockDto, int64, error)
	GetStockByID(stockID int64) (*stock.StockDetailDto, error)
}
