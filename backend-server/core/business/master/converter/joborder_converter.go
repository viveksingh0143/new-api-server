package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/joborder"
)

type JobOrderConverter struct {
	ProductConv  ProductConverter
	CustomerConv CustomerConverter
}

func NewJobOrderConverter(productConv ProductConverter, customerConv CustomerConverter) *JobOrderConverter {
	return &JobOrderConverter{ProductConv: productConv, CustomerConv: customerConv}
}

func (c *JobOrderConverter) ToMinimalDto(domainJobOrder *domain.JobOrder) *joborder.JobOrderMinimalDto {
	joborderDto := &joborder.JobOrderMinimalDto{
		ID:         domainJobOrder.ID,
		IssuedDate: domainJobOrder.IssuedDate,
		OrderNo:    domainJobOrder.OrderNo,
		POCategory: domainJobOrder.POCategory,
	}
	return joborderDto
}

func (c *JobOrderConverter) ToDto(domainJobOrder *domain.JobOrder) *joborder.JobOrderDto {
	joborderDto := &joborder.JobOrderDto{
		ID:            domainJobOrder.ID,
		IssuedDate:    customtypes.NewValidNullTime(domainJobOrder.IssuedDate),
		OrderNo:       domainJobOrder.OrderNo,
		POCategory:    domainJobOrder.POCategory,
		CustomerID:    domainJobOrder.CustomerID,
		Status:        domainJobOrder.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainJobOrder.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainJobOrder.UpdatedAt),
		LastUpdatedBy: domainJobOrder.LastUpdatedBy,
		Customer:      c.CustomerConv.ToMinimalDto(domainJobOrder.Customer),
	}
	var joborderItems = make([]*joborder.JobOrderItemDto, 0)
	if domainJobOrder.Items != nil && len(domainJobOrder.Items) > 0 {
		for _, domainJobOrderItem := range domainJobOrder.Items {
			joborderItems = append(joborderItems, &joborder.JobOrderItemDto{
				ID:         domainJobOrderItem.ID,
				JobOrderID: domainJobOrderItem.JobOrderID,
				ProductID:  domainJobOrderItem.ProductID,
				Quantity:   domainJobOrderItem.Quantity,
				Product:    c.ProductConv.ToMinimalDto(domainJobOrderItem.Product),
			})
		}
	}
	joborderDto.Items = joborderItems
	return joborderDto
}

func (c *JobOrderConverter) ToDtoSlice(domainJobOrders []*domain.JobOrder) []*joborder.JobOrderDto {
	var joborderDtos = make([]*joborder.JobOrderDto, 0)
	for _, domainJobOrder := range domainJobOrders {
		joborderDtos = append(joborderDtos, c.ToDto(domainJobOrder))
	}
	return joborderDtos
}

func (c *JobOrderConverter) ToMinimalDtoSlice(domainJobOrders []*domain.JobOrder) []*joborder.JobOrderMinimalDto {
	var joborderDtos = make([]*joborder.JobOrderMinimalDto, 0)
	for _, domainJobOrder := range domainJobOrders {
		joborderDtos = append(joborderDtos, c.ToMinimalDto(domainJobOrder))
	}
	return joborderDtos
}

func (c *JobOrderConverter) ToDomain(joborderDto *joborder.JobOrderCreateDto) *domain.JobOrder {
	domainJobOrder := &domain.JobOrder{
		IssuedDate:    joborderDto.IssuedDate,
		OrderNo:       joborderDto.OrderNo,
		POCategory:    joborderDto.POCategory,
		CustomerID:    joborderDto.CustomerID,
		Status:        customtypes.Enable,
		LastUpdatedBy: joborderDto.LastUpdatedBy,
	}
	if joborderDto.Customer != nil {
		domainJobOrder.CustomerID = joborderDto.Customer.ID
	}
	var domainJoborderItems = make([]*domain.JobOrderItem, 0)
	if joborderDto.Items != nil && len(joborderDto.Items) > 0 {
		for _, joborderDtoItem := range joborderDto.Items {
			domainJobOrderItem := &domain.JobOrderItem{
				JobOrderID: joborderDtoItem.JobOrderID,
				ProductID:  joborderDtoItem.ProductID,
				Quantity:   joborderDtoItem.Quantity,
			}
			if joborderDtoItem.Product != nil {
				domainJobOrderItem.ProductID = joborderDtoItem.Product.ID
			}
			domainJoborderItems = append(domainJoborderItems, domainJobOrderItem)
		}
	}
	domainJobOrder.Items = domainJoborderItems
	return domainJobOrder
}

func (c *JobOrderConverter) ToUpdateDomain(domainJobOrder *domain.JobOrder, joborderDto *joborder.JobOrderUpdateDto) {
	domainJobOrder.IssuedDate = joborderDto.IssuedDate
	domainJobOrder.OrderNo = joborderDto.OrderNo
	domainJobOrder.POCategory = joborderDto.POCategory
	domainJobOrder.CustomerID = joborderDto.CustomerID
	domainJobOrder.Status = joborderDto.Status
	domainJobOrder.LastUpdatedBy = joborderDto.LastUpdatedBy

	if joborderDto.Customer != nil {
		domainJobOrder.CustomerID = joborderDto.Customer.ID
	}

	var domainJoborderItems = make([]*domain.JobOrderItem, 0)
	if joborderDto.Items != nil && len(joborderDto.Items) > 0 {
		for _, joborderDtoItem := range joborderDto.Items {
			domainJobOrderItem := &domain.JobOrderItem{
				JobOrderID: joborderDtoItem.JobOrderID,
				ProductID:  joborderDtoItem.ProductID,
				Quantity:   joborderDtoItem.Quantity,
			}
			if joborderDtoItem.Product != nil {
				domainJobOrderItem.ProductID = joborderDtoItem.Product.ID
			}
			domainJoborderItems = append(domainJoborderItems, domainJobOrderItem)
		}
	}
	domainJobOrder.Items = domainJoborderItems
}
