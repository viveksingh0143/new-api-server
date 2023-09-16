package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/labelsticker"
)

type LabelStickerRepository interface {
	Create(labelstickers []*domain.LabelSticker) error
	GetById(labelstickerID int64) (*domain.LabelSticker, error)
	GetTotalCount(filter *labelsticker.LabelStickerFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *labelsticker.LabelStickerFilterDto) ([]*domain.LabelSticker, error)
	GetStickerCountByIds(batchIDs []int64) (map[int64]int64, error)
}
