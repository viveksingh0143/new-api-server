package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/outwardrequest"
)

type OutwardRequestConverter struct {
	ProductConv  ProductConverter
	CustomerConv CustomerConverter
}

func NewOutwardRequestConverter(productConv ProductConverter, customerConv CustomerConverter) *OutwardRequestConverter {
	return &OutwardRequestConverter{ProductConv: productConv, CustomerConv: customerConv}
}

func (c *OutwardRequestConverter) ToMinimalDto(domainOutwardRequest *domain.OutwardRequest) *outwardrequest.OutwardRequestMinimalDto {
	outwardrequestDto := &outwardrequest.OutwardRequestMinimalDto{
		ID:         domainOutwardRequest.ID,
		IssuedDate: domainOutwardRequest.IssuedDate,
		OrderNo:    domainOutwardRequest.OrderNo,
	}
	return outwardrequestDto
}

func (c *OutwardRequestConverter) ToDto(domainOutwardRequest *domain.OutwardRequest) *outwardrequest.OutwardRequestDto {
	outwardrequestDto := &outwardrequest.OutwardRequestDto{
		ID:            domainOutwardRequest.ID,
		IssuedDate:    customtypes.NewValidNullTime(domainOutwardRequest.IssuedDate),
		OrderNo:       domainOutwardRequest.OrderNo,
		CustomerID:    domainOutwardRequest.CustomerID,
		Status:        domainOutwardRequest.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainOutwardRequest.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainOutwardRequest.UpdatedAt),
		LastUpdatedBy: domainOutwardRequest.LastUpdatedBy,
		Customer:      c.CustomerConv.ToMinimalDto(domainOutwardRequest.Customer),
	}
	var outwardrequestItems = make([]*outwardrequest.OutwardRequestItemDto, 0)
	if domainOutwardRequest.Items != nil && len(domainOutwardRequest.Items) > 0 {
		for _, domainOutwardRequestItem := range domainOutwardRequest.Items {
			outwardrequestItems = append(outwardrequestItems, &outwardrequest.OutwardRequestItemDto{
				ID:               domainOutwardRequestItem.ID,
				OutwardRequestID: domainOutwardRequestItem.OutwardRequestID,
				ProductID:        domainOutwardRequestItem.ProductID,
				Quantity:         domainOutwardRequestItem.Quantity,
				Product:          c.ProductConv.ToMinimalDto(domainOutwardRequestItem.Product),
				LockedQuantity:   domainOutwardRequestItem.LockedQuantity,
			})
		}
	}
	outwardrequestDto.Items = outwardrequestItems
	return outwardrequestDto
}

func (c *OutwardRequestConverter) ToDtoSlice(domainOutwardRequests []*domain.OutwardRequest) []*outwardrequest.OutwardRequestDto {
	var outwardrequestDtos = make([]*outwardrequest.OutwardRequestDto, 0)
	for _, domainOutwardRequest := range domainOutwardRequests {
		outwardrequestDtos = append(outwardrequestDtos, c.ToDto(domainOutwardRequest))
	}
	return outwardrequestDtos
}

func (c *OutwardRequestConverter) ToMinimalDtoSlice(domainOutwardRequests []*domain.OutwardRequest) []*outwardrequest.OutwardRequestMinimalDto {
	var outwardrequestDtos = make([]*outwardrequest.OutwardRequestMinimalDto, 0)
	for _, domainOutwardRequest := range domainOutwardRequests {
		outwardrequestDtos = append(outwardrequestDtos, c.ToMinimalDto(domainOutwardRequest))
	}
	return outwardrequestDtos
}

func (c *OutwardRequestConverter) ToDomain(outwardrequestDto *outwardrequest.OutwardRequestCreateDto) *domain.OutwardRequest {
	domainOutwardRequest := &domain.OutwardRequest{
		IssuedDate:    outwardrequestDto.IssuedDate,
		OrderNo:       outwardrequestDto.OrderNo,
		CustomerID:    outwardrequestDto.CustomerID,
		Status:        customtypes.Enable,
		LastUpdatedBy: outwardrequestDto.LastUpdatedBy,
	}
	if outwardrequestDto.Customer != nil {
		domainOutwardRequest.CustomerID = outwardrequestDto.Customer.ID
	}
	var domainJoborderItems = make([]*domain.OutwardRequestItem, 0)
	if outwardrequestDto.Items != nil && len(outwardrequestDto.Items) > 0 {
		for _, outwardrequestDtoItem := range outwardrequestDto.Items {
			domainOutwardRequestItem := &domain.OutwardRequestItem{
				OutwardRequestID: outwardrequestDtoItem.OutwardRequestID,
				ProductID:        outwardrequestDtoItem.ProductID,
				Quantity:         outwardrequestDtoItem.Quantity,
				LockedQuantity:   outwardrequestDtoItem.LockedQuantity,
			}
			if outwardrequestDtoItem.Product != nil {
				domainOutwardRequestItem.ProductID = outwardrequestDtoItem.Product.ID
			}
			domainJoborderItems = append(domainJoborderItems, domainOutwardRequestItem)
		}
	}
	domainOutwardRequest.Items = domainJoborderItems
	return domainOutwardRequest
}

func (c *OutwardRequestConverter) ToUpdateDomain(domainOutwardRequest *domain.OutwardRequest, outwardrequestDto *outwardrequest.OutwardRequestUpdateDto) {
	domainOutwardRequest.IssuedDate = outwardrequestDto.IssuedDate
	domainOutwardRequest.OrderNo = outwardrequestDto.OrderNo
	domainOutwardRequest.CustomerID = outwardrequestDto.CustomerID
	domainOutwardRequest.Status = outwardrequestDto.Status
	domainOutwardRequest.LastUpdatedBy = outwardrequestDto.LastUpdatedBy

	if outwardrequestDto.Customer != nil {
		domainOutwardRequest.CustomerID = outwardrequestDto.Customer.ID
	}

	var domainJoborderItems = make([]*domain.OutwardRequestItem, 0)
	if outwardrequestDto.Items != nil && len(outwardrequestDto.Items) > 0 {
		for _, outwardrequestDtoItem := range outwardrequestDto.Items {
			domainOutwardRequestItem := &domain.OutwardRequestItem{
				OutwardRequestID: outwardrequestDtoItem.OutwardRequestID,
				ProductID:        outwardrequestDtoItem.ProductID,
				Quantity:         outwardrequestDtoItem.Quantity,
				LockedQuantity:   outwardrequestDtoItem.LockedQuantity,
			}
			if outwardrequestDtoItem.Product != nil {
				domainOutwardRequestItem.ProductID = outwardrequestDtoItem.Product.ID
			}
			domainJoborderItems = append(domainJoborderItems, domainOutwardRequestItem)
		}
	}
	domainOutwardRequest.Items = domainJoborderItems
}
