package service

import (
	"log"

	"github.com/vamika-digital/wms-api-server/core/business/warehouse/converter"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/labelsticker"
	warehouseRepository "github.com/vamika-digital/wms-api-server/core/business/warehouse/repository"
)

type LabelStickerServiceImpl struct {
	LabelStickerRepo      warehouseRepository.LabelStickerRepository
	BatchLabelRepo        warehouseRepository.BatchLabelRepository
	LabelStickerConverter *converter.LabelStickerConverter
}

func NewLabelStickerService(labelstickerRepo warehouseRepository.LabelStickerRepository, batchLabelRepo warehouseRepository.BatchLabelRepository, labelstickerConverter *converter.LabelStickerConverter) LabelStickerService {
	return &LabelStickerServiceImpl{LabelStickerRepo: labelstickerRepo, BatchLabelRepo: batchLabelRepo, LabelStickerConverter: labelstickerConverter}
}

func (s *LabelStickerServiceImpl) GetAllLabelStickers(page int16, pageSize int16, sort string, filter *labelsticker.LabelStickerFilterDto) ([]*labelsticker.LabelStickerDto, int64, error) {
	totalCount, err := s.LabelStickerRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainLabelStickers, err := s.LabelStickerRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	if len(domainLabelStickers) > 0 {
		uniqueLabelStickerIDsMap := make(map[int64]bool)
		var uniqueLabelStickerIDs []int64

		for _, labelsticker := range domainLabelStickers {
			if _, exists := uniqueLabelStickerIDsMap[labelsticker.BatchLabelID]; !exists {
				uniqueLabelStickerIDsMap[labelsticker.BatchLabelID] = true
				uniqueLabelStickerIDs = append(uniqueLabelStickerIDs, labelsticker.BatchLabelID)
			}
		}

		batchlabels, err := s.BatchLabelRepo.GetByIds(uniqueLabelStickerIDs)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, 0, err
		}

		batchlabelMap := make(map[int64]*domain.BatchLabel)
		for _, batchlabel := range batchlabels {
			batchlabelMap[batchlabel.ID] = batchlabel
		}

		for i, labelsticker := range domainLabelStickers {
			if batchlabel, ok := batchlabelMap[labelsticker.BatchLabelID]; ok {
				domainLabelStickers[i].BatchLabel = batchlabel
			}
		}
	}

	var labelstickerDtos []*labelsticker.LabelStickerDto = s.LabelStickerConverter.ToDtoSlice(domainLabelStickers)
	return labelstickerDtos, int64(totalCount), nil
}

func (s *LabelStickerServiceImpl) GetLabelStickerByID(labelstickerID int64) (*labelsticker.LabelStickerDto, error) {
	domainLabelSticker, err := s.LabelStickerRepo.GetById(labelstickerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	batchlabel, err := s.BatchLabelRepo.GetById(domainLabelSticker.BatchLabelID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainLabelSticker.BatchLabel = batchlabel

	return s.LabelStickerConverter.ToDto(domainLabelSticker), nil
}
