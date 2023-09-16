package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	adminConverter "github.com/vamika-digital/wms-api-server/core/business/admin/converter"
	masterConverter "github.com/vamika-digital/wms-api-server/core/business/master/converter"
	masterDomain "github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/domain"
	"github.com/vamika-digital/wms-api-server/core/business/warehouse/dto/batchlabel"
)

type BatchLabelConverter struct {
	userConverter     *adminConverter.UserConverter
	productConverter  *masterConverter.ProductConverter
	customerConverter *masterConverter.CustomerConverter
	machineConverter  *masterConverter.MachineConverter
}

func NewBatchLabelConverter(userConverter *adminConverter.UserConverter, productConverter *masterConverter.ProductConverter, customerConverter *masterConverter.CustomerConverter, machineConverter *masterConverter.MachineConverter) *BatchLabelConverter {
	return &BatchLabelConverter{userConverter: userConverter, productConverter: productConverter, customerConverter: customerConverter, machineConverter: machineConverter}
}

func (c *BatchLabelConverter) ToMinimalDto(domainBatchLabel *domain.BatchLabel) *batchlabel.BatchLabelMinimalDto {
	batchlabelDto := &batchlabel.BatchLabelMinimalDto{
		ID:         domainBatchLabel.ID,
		BatchDate:  customtypes.NewValidNullTime(domainBatchLabel.BatchDate),
		BatchNo:    domainBatchLabel.BatchNo,
		SoNumber:   domainBatchLabel.SoNumber,
		PoCategory: domainBatchLabel.PoCategory,
		Status:     domainBatchLabel.Status,
		Customer:   c.customerConverter.ToMinimalDto(domainBatchLabel.Customer),
		Product:    c.productConverter.ToMinimalDto(domainBatchLabel.Product),
	}
	return batchlabelDto
}

func (c *BatchLabelConverter) ToDto(domainBatchLabel *domain.BatchLabel) *batchlabel.BatchLabelDto {
	batchlabelDto := &batchlabel.BatchLabelDto{
		ID:              domainBatchLabel.ID,
		BatchDate:       customtypes.NewValidNullTime(domainBatchLabel.BatchDate),
		BatchNo:         domainBatchLabel.BatchNo,
		SoNumber:        domainBatchLabel.SoNumber,
		PoCategory:      domainBatchLabel.PoCategory,
		TargetQuantity:  domainBatchLabel.TargetQuantity,
		PackageQuantity: domainBatchLabel.PackageQuantity,
		UnitWeight:      domainBatchLabel.UnitWeight,
		UnitWeightType:  domainBatchLabel.UnitWeightType,
		Status:          domainBatchLabel.Status,
		CreatedAt:       customtypes.NewValidNullTime(domainBatchLabel.CreatedAt),
		UpdatedAt:       customtypes.GetNullTime(domainBatchLabel.UpdatedAt),
		LastUpdatedBy:   domainBatchLabel.LastUpdatedBy,
		Customer:        c.customerConverter.ToMinimalDto(domainBatchLabel.Customer),
		Product:         c.productConverter.ToMinimalDto(domainBatchLabel.Product),
		Machine:         c.machineConverter.ToMinimalDto(domainBatchLabel.Machine),
	}
	batchlabelDto.LabelsToPrint = domainBatchLabel.GetStickerCountToPrint()
	batchlabelDto.TotalPrinted = domainBatchLabel.TotalPrinted
	return batchlabelDto
}

func (c *BatchLabelConverter) ToDtoSlice(domainBatchLabels []*domain.BatchLabel) []*batchlabel.BatchLabelDto {
	var batchlabelDtos = make([]*batchlabel.BatchLabelDto, 0)
	for _, domainBatchLabel := range domainBatchLabels {
		batchlabelDtos = append(batchlabelDtos, c.ToDto(domainBatchLabel))
	}
	return batchlabelDtos
}

func (c *BatchLabelConverter) ToDomain(batchlabelDto *batchlabel.BatchLabelCreateDto) *domain.BatchLabel {
	domainBatchLabel := &domain.BatchLabel{
		BatchDate:       batchlabelDto.BatchDate,
		BatchNo:         batchlabelDto.BatchNo,
		SoNumber:        batchlabelDto.SoNumber,
		PoCategory:      batchlabelDto.PoCategory,
		TargetQuantity:  batchlabelDto.TargetQuantity,
		PackageQuantity: batchlabelDto.PackageQuantity,
		UnitWeight:      batchlabelDto.UnitWeight,
		UnitWeightType:  batchlabelDto.UnitWeightType,
		Status:          batchlabelDto.Status,
	}
	if batchlabelDto.Customer != nil && batchlabelDto.Customer.ID != 0 {
		domainBatchLabel.Customer = &masterDomain.Customer{ID: batchlabelDto.Customer.ID}
	}

	if batchlabelDto.Product != nil && batchlabelDto.Product.ID != 0 {
		domainBatchLabel.Product = &masterDomain.Product{ID: batchlabelDto.Product.ID}
	}

	if batchlabelDto.Machine != nil && batchlabelDto.Machine.ID != 0 {
		domainBatchLabel.Machine = &masterDomain.Machine{ID: batchlabelDto.Machine.ID}
	}
	return domainBatchLabel
}

func (c *BatchLabelConverter) ToUpdateDomain(domainBatchLabel *domain.BatchLabel, batchlabelDto *batchlabel.BatchLabelUpdateDto) {
	domainBatchLabel.BatchDate = batchlabelDto.BatchDate
	domainBatchLabel.BatchNo = batchlabelDto.BatchNo
	domainBatchLabel.SoNumber = batchlabelDto.SoNumber
	domainBatchLabel.PoCategory = batchlabelDto.PoCategory
	domainBatchLabel.TargetQuantity = batchlabelDto.TargetQuantity
	domainBatchLabel.PackageQuantity = batchlabelDto.PackageQuantity
	domainBatchLabel.UnitWeight = batchlabelDto.UnitWeight
	domainBatchLabel.UnitWeightType = batchlabelDto.UnitWeightType

	if batchlabelDto.Status.IsValid() {
		domainBatchLabel.Status = batchlabelDto.Status
	}
	if batchlabelDto.Customer != nil && batchlabelDto.Customer.ID != 0 {
		domainBatchLabel.Customer = &masterDomain.Customer{ID: batchlabelDto.Customer.ID}
	}

	if batchlabelDto.Product != nil && batchlabelDto.Product.ID != 0 {
		domainBatchLabel.Product = &masterDomain.Product{ID: batchlabelDto.Product.ID}
	}

	if batchlabelDto.Machine != nil && batchlabelDto.Machine.ID != 0 {
		domainBatchLabel.Machine = &masterDomain.Machine{ID: batchlabelDto.Machine.ID}
	}
	domainBatchLabel.LastUpdatedBy = batchlabelDto.LastUpdatedBy
}
