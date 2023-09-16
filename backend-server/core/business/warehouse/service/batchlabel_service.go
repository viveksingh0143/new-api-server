package service

import (
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/batchlabel"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/labelsticker"
)

type BatchLabelService interface {
	GetAllBatchLabels(page int16, pageSize int16, sort string, filter *batchlabel.BatchLabelFilterDto) ([]*batchlabel.BatchLabelDto, int64, error)
	CreateBatchLabel(batchlabelDto *batchlabel.BatchLabelCreateDto) error
	GetBatchLabelByID(batchlabelID int64) (*batchlabel.BatchLabelDto, error)
	GetBatchLabelByBarcode(barcode string) (*batchlabel.BatchLabelDto, *labelsticker.LabelStickerMinimalDto, error)
	GenerateBatchLabelStickers(batchlabelID int64, form *batchlabel.BatchLabelStickersCreateDto) error
	GetBatchLabelByBatchNumber(batchlabelBatchNumber string) (*batchlabel.BatchLabelDto, error)
	UpdateBatchLabel(batchlabelID int64, batchlabel *batchlabel.BatchLabelUpdateDto) error
	DeleteBatchLabel(batchlabelID int64) error
	DeleteBatchLabelByIDs(batchlabelIDs []int64) error
}
