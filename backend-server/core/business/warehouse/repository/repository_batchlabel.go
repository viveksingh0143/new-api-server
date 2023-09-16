package repository

import (
	"time"

	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/batchlabel"
)

type BatchLabelRepository interface {
	Create(batchlabel *domain.BatchLabel) error
	Update(batchlabel *domain.BatchLabel) error
	Delete(batchlabelID int64) error
	DeleteByIDs(batchlabelIDs []int64) error
	GetByBarcode(barcode string) (*domain.BatchLabel, *domain.LabelSticker, error)
	GetById(batchlabelID int64) (*domain.BatchLabel, error)
	GetByIds(batchlabelID []int64) ([]*domain.BatchLabel, error)
	GetByBatchNumber(batchlabelBatchNumber string) (*domain.BatchLabel, error)
	GetTotalCount(filter *batchlabel.BatchLabelFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *batchlabel.BatchLabelFilterDto) ([]*domain.BatchLabel, error)
	GetTotalStickers(batchlabelID int64) (int64, error)
	GetTotalStickersForShift(batchlabelID int64, shift string, createdAt time.Time) (int64, error)
}
