package service

import (
	"fmt"
	"log"
	"time"

	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	masterDomain "github.com/vamika-digital/wms-api-server/core/business/master/domain"
	masterRepository "github.com/vamika-digital/wms-api-server/core/business/master/repository"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/converter"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/batchlabel"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/labelsticker"
	warehouseRepository "github.com/vamika-digital/wms-api-server/core/business/warehouse/repository"
)

type BatchLabelServiceImpl struct {
	BatchLabelRepo        warehouseRepository.BatchLabelRepository
	StickerRepo           warehouseRepository.LabelStickerRepository
	ProductRepository     masterRepository.ProductRepository
	MachineRepository     masterRepository.MachineRepository
	CustomerRepository    masterRepository.CustomerRepository
	BatchLabelConverter   *converter.BatchLabelConverter
	LabelStickerConverter *converter.LabelStickerConverter
}

func NewBatchLabelService(batchLabelRepo warehouseRepository.BatchLabelRepository, stickerRepo warehouseRepository.LabelStickerRepository, productRepository masterRepository.ProductRepository, machineRepository masterRepository.MachineRepository, customerRepository masterRepository.CustomerRepository, batchLabelConverter *converter.BatchLabelConverter, labelStickerConverter *converter.LabelStickerConverter) BatchLabelService {
	return &BatchLabelServiceImpl{BatchLabelRepo: batchLabelRepo, StickerRepo: stickerRepo, ProductRepository: productRepository, MachineRepository: machineRepository, CustomerRepository: customerRepository, BatchLabelConverter: batchLabelConverter, LabelStickerConverter: labelStickerConverter}
}

func (s *BatchLabelServiceImpl) GetAllBatchLabels(page int16, pageSize int16, sort string, filter *batchlabel.BatchLabelFilterDto) ([]*batchlabel.BatchLabelDto, int64, error) {
	totalCount, err := s.BatchLabelRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainBatchLabels, err := s.BatchLabelRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	if len(domainBatchLabels) > 0 {
		uniqueCustomerIDsMap := make(map[int64]bool)
		var uniqueCustomerIDs []int64

		for _, batchlabel := range domainBatchLabels {
			if _, exists := uniqueCustomerIDsMap[batchlabel.CustomerID]; !exists {
				uniqueCustomerIDsMap[batchlabel.CustomerID] = true
				uniqueCustomerIDs = append(uniqueCustomerIDs, batchlabel.CustomerID)
			}
		}

		customers, err := s.CustomerRepository.GetByIds(uniqueCustomerIDs)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, 0, err
		}

		customerMap := make(map[int64]*masterDomain.Customer)
		for _, customer := range customers {
			customerMap[customer.ID] = customer
		}

		for i, batchlabel := range domainBatchLabels {
			if customer, ok := customerMap[batchlabel.CustomerID]; ok {
				domainBatchLabels[i].Customer = customer
			}
		}
	}

	if len(domainBatchLabels) > 0 {
		uniqueProductIDsMap := make(map[int64]bool)
		var uniqueProductIDs []int64

		for _, batchlabel := range domainBatchLabels {
			if _, exists := uniqueProductIDsMap[batchlabel.ProductID]; !exists {
				uniqueProductIDsMap[batchlabel.ProductID] = true
				uniqueProductIDs = append(uniqueProductIDs, batchlabel.ProductID)
			}
		}

		products, err := s.ProductRepository.GetByIds(uniqueProductIDs)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, 0, err
		}

		productMap := make(map[int64]*masterDomain.Product)
		for _, product := range products {
			productMap[product.ID] = product
		}

		for i, batchlabel := range domainBatchLabels {
			if product, ok := productMap[batchlabel.ProductID]; ok {
				domainBatchLabels[i].Product = product
			}
		}
	}

	if len(domainBatchLabels) > 0 {
		uniqueMachineIDsMap := make(map[int64]bool)
		var uniqueMachineIDs []int64

		for _, batchlabel := range domainBatchLabels {
			if _, exists := uniqueMachineIDsMap[batchlabel.MachineID]; !exists {
				uniqueMachineIDsMap[batchlabel.MachineID] = true
				uniqueMachineIDs = append(uniqueMachineIDs, batchlabel.MachineID)
			}
		}

		machines, err := s.MachineRepository.GetByIds(uniqueMachineIDs)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, 0, err
		}

		machineMap := make(map[int64]*masterDomain.Machine)
		for _, machine := range machines {
			machineMap[machine.ID] = machine
		}

		for i, batchlabel := range domainBatchLabels {
			if machine, ok := machineMap[batchlabel.MachineID]; ok {
				domainBatchLabels[i].Machine = machine
			}
		}
	}

	if len(domainBatchLabels) > 0 {
		var batchIDs []int64
		for _, batchLabel := range domainBatchLabels {
			batchIDs = append(batchIDs, batchLabel.ID)
		}

		labelStickers, err := s.StickerRepo.GetStickerCountByIds(batchIDs)
		if err != nil {
			log.Printf("%+v\n", err)
			return nil, 0, err
		}

		for i, batchlabel := range domainBatchLabels {
			if stickerCount, ok := labelStickers[batchlabel.ID]; ok {
				domainBatchLabels[i].TotalPrinted = stickerCount
			}
		}
	}

	var batchlabelDtos []*batchlabel.BatchLabelDto = s.BatchLabelConverter.ToDtoSlice(domainBatchLabels)
	return batchlabelDtos, int64(totalCount), nil
}

func (s *BatchLabelServiceImpl) CreateBatchLabel(batchlabelDto *batchlabel.BatchLabelCreateDto) error {
	var newBatchLabel *domain.BatchLabel = s.BatchLabelConverter.ToDomain(batchlabelDto)
	if err := s.BatchLabelRepo.Create(newBatchLabel); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *BatchLabelServiceImpl) GetBatchLabelByID(batchlabelID int64) (*batchlabel.BatchLabelDto, error) {
	domainBatchLabel, err := s.BatchLabelRepo.GetById(batchlabelID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	customer, err := s.CustomerRepository.GetById(domainBatchLabel.CustomerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainBatchLabel.Customer = customer

	product, err := s.ProductRepository.GetById(domainBatchLabel.ProductID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainBatchLabel.Product = product

	machine, err := s.MachineRepository.GetById(domainBatchLabel.MachineID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainBatchLabel.Machine = machine

	return s.BatchLabelConverter.ToDto(domainBatchLabel), nil
}

func (s *BatchLabelServiceImpl) GenerateBatchLabelStickers(batchlabelID int64, form *batchlabel.BatchLabelStickersCreateDto) error {
	domainBatchLabel, err := s.BatchLabelRepo.GetById(batchlabelID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	product, err := s.ProductRepository.GetById(domainBatchLabel.ProductID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	domainBatchLabel.Product = product

	machine, err := s.MachineRepository.GetById(domainBatchLabel.MachineID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	domainBatchLabel.Machine = machine

	totalCount, err := s.BatchLabelRepo.GetTotalStickers(domainBatchLabel.ID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	stickersAvailable := domainBatchLabel.GetStickerCountToPrint() - totalCount
	if stickersAvailable < form.StickersToGenerate {
		return fmt.Errorf("only %d stickers are left, please reduce the sticker count", stickersAvailable)
	}

	shiftTotalCount, err := s.BatchLabelRepo.GetTotalStickersForShift(domainBatchLabel.ID, form.WorkingShift, time.Now())
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	var stickers []*domain.LabelSticker = make([]*domain.LabelSticker, 0, form.StickersToGenerate)
	for i := int64(0); i < form.StickersToGenerate; i++ {
		nextSerialNumber := fmt.Sprintf("%04d", totalCount+i+1)
		packetNo := shiftTotalCount + i + 1
		packetNoString := fmt.Sprintf("%04d", packetNo)

		labelSticker := &domain.LabelSticker{
			UUIDCode:      "",
			PacketNo:      packetNoString,
			PrintCount:    0,
			Shift:         form.WorkingShift,
			ProductLine:   domainBatchLabel.Product.Name,
			BatchNo:       domainBatchLabel.BatchNo,
			UnitWeight:    fmt.Sprintf("%f", domainBatchLabel.UnitWeight),
			Quantity:      fmt.Sprintf("%d", domainBatchLabel.TargetQuantity),
			MachineNo:     domainBatchLabel.Machine.Name,
			LastUpdatedBy: customtypes.NewNullString(form.LastUpdatedBy),
			BatchLabelID:  domainBatchLabel.ID,
		}
		labelSticker.UUIDCode = fmt.Sprintf("%s%s%s%s", labelSticker.BatchNo, labelSticker.Shift, labelSticker.PacketNo, nextSerialNumber)
		stickers = append(stickers, labelSticker)
	}

	err = s.StickerRepo.Create(stickers)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	return nil
}

func (s *BatchLabelServiceImpl) GetBatchLabelByBatchNumber(batchlabelBatchNumber string) (*batchlabel.BatchLabelDto, error) {
	domainBatchLabel, err := s.BatchLabelRepo.GetByBatchNumber(batchlabelBatchNumber)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}

	customer, err := s.CustomerRepository.GetById(domainBatchLabel.CustomerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainBatchLabel.Customer = customer

	product, err := s.ProductRepository.GetById(domainBatchLabel.ProductID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainBatchLabel.Product = product

	machine, err := s.MachineRepository.GetById(domainBatchLabel.MachineID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	domainBatchLabel.Machine = machine

	return s.BatchLabelConverter.ToDto(domainBatchLabel), nil
}

func (s *BatchLabelServiceImpl) UpdateBatchLabel(batchlabelID int64, batchlabelDto *batchlabel.BatchLabelUpdateDto) error {
	existingBatchLabel, err := s.BatchLabelRepo.GetById(batchlabelID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	s.BatchLabelConverter.ToUpdateDomain(existingBatchLabel, batchlabelDto)
	if err := s.BatchLabelRepo.Update(existingBatchLabel); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *BatchLabelServiceImpl) DeleteBatchLabel(batchlabelID int64) error {
	if err := s.BatchLabelRepo.Delete(batchlabelID); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *BatchLabelServiceImpl) DeleteBatchLabelByIDs(batchlabelIDs []int64) error {
	if err := s.BatchLabelRepo.DeleteByIDs(batchlabelIDs); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *BatchLabelServiceImpl) GetBatchLabelByBarcode(barcode string) (*batchlabel.BatchLabelDto, *labelsticker.LabelStickerMinimalDto, error) {

	domainBatchLabel, domainSticker, err := s.BatchLabelRepo.GetByBarcode(barcode)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, err
	}

	customer, err := s.CustomerRepository.GetById(domainBatchLabel.CustomerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, err
	}
	domainBatchLabel.Customer = customer

	product, err := s.ProductRepository.GetById(domainBatchLabel.ProductID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, err
	}
	domainBatchLabel.Product = product

	machine, err := s.MachineRepository.GetById(domainBatchLabel.MachineID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, nil, err
	}
	domainBatchLabel.Machine = machine

	return s.BatchLabelConverter.ToDto(domainBatchLabel), s.LabelStickerConverter.ToMinimalDto(domainSticker), nil
}
