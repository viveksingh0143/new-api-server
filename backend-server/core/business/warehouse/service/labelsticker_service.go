package service

import "github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/labelsticker"

type LabelStickerService interface {
	GetAllLabelStickers(page int16, pageSize int16, sort string, filter *labelsticker.LabelStickerFilterDto) ([]*labelsticker.LabelStickerDto, int64, error)
	GetLabelStickerByID(labelstickerID int64) (*labelsticker.LabelStickerDto, error)
}
