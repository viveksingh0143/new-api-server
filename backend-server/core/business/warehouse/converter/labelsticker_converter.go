package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/labelsticker"
)

type LabelStickerConverter struct {
	batchlabelConverter *BatchLabelConverter
}

func NewLabelStickerConverter(batchlabelConverter *BatchLabelConverter) *LabelStickerConverter {
	return &LabelStickerConverter{batchlabelConverter: batchlabelConverter}
}

func (c *LabelStickerConverter) ToMinimalDto(domainLabelSticker *domain.LabelSticker) *labelsticker.LabelStickerMinimalDto {
	labelstickerDto := &labelsticker.LabelStickerMinimalDto{
		ID:          domainLabelSticker.ID,
		UUIDCode:    domainLabelSticker.UUIDCode,
		PrintCount:  domainLabelSticker.PrintCount,
		Shift:       domainLabelSticker.Shift,
		ProductLine: domainLabelSticker.ProductLine,
		IsUsed:      domainLabelSticker.IsUsed,
		BatchNo:     domainLabelSticker.BatchNo,
	}

	return labelstickerDto
}

func (c *LabelStickerConverter) ToDto(domainLabelSticker *domain.LabelSticker) *labelsticker.LabelStickerDto {
	labelstickerDto := &labelsticker.LabelStickerDto{
		ID:            domainLabelSticker.ID,
		UUIDCode:      domainLabelSticker.UUIDCode,
		PacketNo:      domainLabelSticker.PacketNo,
		PrintCount:    domainLabelSticker.PrintCount,
		Shift:         domainLabelSticker.Shift,
		ProductLine:   domainLabelSticker.ProductLine,
		BatchNo:       domainLabelSticker.BatchNo,
		UnitWeight:    domainLabelSticker.UnitWeight,
		Quantity:      domainLabelSticker.Quantity,
		MachineNo:     domainLabelSticker.MachineNo,
		CreatedAt:     customtypes.NewValidNullTime(domainLabelSticker.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainLabelSticker.UpdatedAt),
		LastUpdatedBy: domainLabelSticker.LastUpdatedBy,
		IsUsed:        domainLabelSticker.IsUsed,
		BatchLabelID:  domainLabelSticker.BatchLabelID,
	}
	return labelstickerDto
}

func (c *LabelStickerConverter) ToDtoSlice(domainLabelStickers []*domain.LabelSticker) []*labelsticker.LabelStickerDto {
	var labelstickerDtos = make([]*labelsticker.LabelStickerDto, 0)
	for _, domainLabelSticker := range domainLabelStickers {
		labelstickerDtos = append(labelstickerDtos, c.ToDto(domainLabelSticker))
	}
	return labelstickerDtos
}
