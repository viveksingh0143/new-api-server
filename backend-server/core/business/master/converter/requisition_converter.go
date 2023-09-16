package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/requisition"
)

type RequisitionConverter struct {
	ProductConv ProductConverter
	StoreConv   StoreConverter
}

func NewRequisitionConverter(productConv ProductConverter, storeConv StoreConverter) *RequisitionConverter {
	return &RequisitionConverter{ProductConv: productConv, StoreConv: storeConv}
}

func (c *RequisitionConverter) ToMinimalDto(domainRequisition *domain.Requisition) *requisition.RequisitionMinimalDto {
	requisitionDto := &requisition.RequisitionMinimalDto{
		ID:         domainRequisition.ID,
		IssuedDate: domainRequisition.IssuedDate,
		OrderNo:    domainRequisition.OrderNo,
		Department: domainRequisition.Department,
	}
	return requisitionDto
}

func (c *RequisitionConverter) ToDto(domainRequisition *domain.Requisition) *requisition.RequisitionDto {
	requisitionDto := &requisition.RequisitionDto{
		ID:            domainRequisition.ID,
		IssuedDate:    customtypes.NewValidNullTime(domainRequisition.IssuedDate),
		OrderNo:       domainRequisition.OrderNo,
		Department:    domainRequisition.Department,
		StoreID:       domainRequisition.StoreID,
		Status:        domainRequisition.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainRequisition.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainRequisition.UpdatedAt),
		LastUpdatedBy: domainRequisition.LastUpdatedBy,
		Store:         c.StoreConv.ToMinimalDto(domainRequisition.Store),
	}
	var requisitionItems = make([]*requisition.RequisitionItemDto, 0)
	if domainRequisition.Items != nil && len(domainRequisition.Items) > 0 {
		for _, domainRequisitionItem := range domainRequisition.Items {
			requisitionItems = append(requisitionItems, &requisition.RequisitionItemDto{
				ID:            domainRequisitionItem.ID,
				RequisitionID: domainRequisitionItem.RequisitionID,
				ProductID:     domainRequisitionItem.ProductID,
				Quantity:      domainRequisitionItem.Quantity,
				Product:       c.ProductConv.ToMinimalDto(domainRequisitionItem.Product),
			})
		}
	}
	requisitionDto.Items = requisitionItems
	return requisitionDto
}

func (c *RequisitionConverter) ToDtoSlice(domainRequisitions []*domain.Requisition) []*requisition.RequisitionDto {
	var requisitionDtos = make([]*requisition.RequisitionDto, 0)
	for _, domainRequisition := range domainRequisitions {
		requisitionDtos = append(requisitionDtos, c.ToDto(domainRequisition))
	}
	return requisitionDtos
}

func (c *RequisitionConverter) ToMinimalDtoSlice(domainRequisitions []*domain.Requisition) []*requisition.RequisitionMinimalDto {
	var requisitionDtos = make([]*requisition.RequisitionMinimalDto, 0)
	for _, domainRequisition := range domainRequisitions {
		requisitionDtos = append(requisitionDtos, c.ToMinimalDto(domainRequisition))
	}
	return requisitionDtos
}

func (c *RequisitionConverter) ToDomain(requisitionDto *requisition.RequisitionCreateDto) *domain.Requisition {
	domainRequisition := &domain.Requisition{
		IssuedDate:    requisitionDto.IssuedDate,
		OrderNo:       requisitionDto.OrderNo,
		Department:    requisitionDto.Department,
		StoreID:       requisitionDto.StoreID,
		Status:        customtypes.Enable,
		LastUpdatedBy: requisitionDto.LastUpdatedBy,
	}
	if requisitionDto.Store != nil {
		domainRequisition.StoreID = requisitionDto.Store.ID
	}
	var domainJoborderItems = make([]*domain.RequisitionItem, 0)
	if requisitionDto.Items != nil && len(requisitionDto.Items) > 0 {
		for _, requisitionDtoItem := range requisitionDto.Items {
			domainRequisitionItem := &domain.RequisitionItem{
				RequisitionID: requisitionDtoItem.RequisitionID,
				ProductID:     requisitionDtoItem.ProductID,
				Quantity:      requisitionDtoItem.Quantity,
			}
			if requisitionDtoItem.Product != nil {
				domainRequisitionItem.ProductID = requisitionDtoItem.Product.ID
			}
			domainJoborderItems = append(domainJoborderItems, domainRequisitionItem)
		}
	}
	domainRequisition.Items = domainJoborderItems
	return domainRequisition
}

func (c *RequisitionConverter) ToUpdateDomain(domainRequisition *domain.Requisition, requisitionDto *requisition.RequisitionUpdateDto) {
	domainRequisition.IssuedDate = requisitionDto.IssuedDate
	domainRequisition.OrderNo = requisitionDto.OrderNo
	domainRequisition.Department = requisitionDto.Department
	domainRequisition.StoreID = requisitionDto.StoreID
	domainRequisition.Status = requisitionDto.Status
	domainRequisition.LastUpdatedBy = requisitionDto.LastUpdatedBy

	if requisitionDto.Store != nil {
		domainRequisition.StoreID = requisitionDto.Store.ID
	}

	var domainJoborderItems = make([]*domain.RequisitionItem, 0)
	if requisitionDto.Items != nil && len(requisitionDto.Items) > 0 {
		for _, requisitionDtoItem := range requisitionDto.Items {
			domainRequisitionItem := &domain.RequisitionItem{
				RequisitionID: requisitionDtoItem.RequisitionID,
				ProductID:     requisitionDtoItem.ProductID,
				Quantity:      requisitionDtoItem.Quantity,
			}
			if requisitionDtoItem.Product != nil {
				domainRequisitionItem.ProductID = requisitionDtoItem.Product.ID
			}
			domainJoborderItems = append(domainJoborderItems, domainRequisitionItem)
		}
	}
	domainRequisition.Items = domainJoborderItems
}
