package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
)

type ProductConverter struct{}

func NewProductConverter() *ProductConverter {
	return &ProductConverter{}
}

func (c *ProductConverter) ToMinimalDto(domainProduct *domain.Product) *product.ProductMinimalDto {
	productDto := &product.ProductMinimalDto{
		ID:             domainProduct.ID,
		ProductType:    domainProduct.ProductType,
		Code:           domainProduct.Code,
		LinkCode:       domainProduct.LinkCode,
		Name:           domainProduct.Name,
		Description:    domainProduct.Description,
		UnitType:       domainProduct.UnitType,
		UnitWeight:     domainProduct.UnitWeight,
		UnitWeightType: domainProduct.UnitWeightType,
		Status:         domainProduct.Status,
	}
	return productDto
}

func (c *ProductConverter) ToDto(domainProduct *domain.Product) *product.ProductDto {
	productDto := &product.ProductDto{
		ID:             domainProduct.ID,
		ProductType:    domainProduct.ProductType,
		Code:           domainProduct.Code,
		LinkCode:       domainProduct.LinkCode,
		Name:           domainProduct.Name,
		Description:    domainProduct.Description,
		UnitType:       domainProduct.UnitType,
		UnitWeight:     domainProduct.UnitWeight,
		UnitWeightType: domainProduct.UnitWeightType,
		Status:         domainProduct.Status,
		CreatedAt:      customtypes.NewValidNullTime(domainProduct.CreatedAt),
		UpdatedAt:      customtypes.GetNullTime(domainProduct.UpdatedAt),
		LastUpdatedBy:  domainProduct.LastUpdatedBy,
	}
	return productDto
}

func (c *ProductConverter) ToDtoSlice(domainProducts []*domain.Product) []*product.ProductDto {
	var productDtos = make([]*product.ProductDto, 0)
	for _, domainProduct := range domainProducts {
		productDtos = append(productDtos, c.ToDto(domainProduct))
	}
	return productDtos
}

func (c *ProductConverter) ToDomain(productDto *product.ProductCreateDto) *domain.Product {
	domainProduct := &domain.Product{
		ProductType:    productDto.ProductType,
		Code:           productDto.Code,
		LinkCode:       productDto.LinkCode,
		Name:           productDto.Name,
		Description:    productDto.Description,
		UnitType:       productDto.UnitType,
		UnitWeight:     productDto.UnitWeight,
		UnitWeightType: productDto.UnitWeightType,
		Status:         productDto.Status,
		LastUpdatedBy:  productDto.LastUpdatedBy,
	}
	if domainProduct.UnitType != customtypes.UnitPiece {
		domainProduct.UnitWeight = 0
		domainProduct.UnitWeightType = customtypes.UnitWeightGram
	}
	return domainProduct
}

func (c *ProductConverter) ToUpdateDomain(domainProduct *domain.Product, productDto *product.ProductUpdateDto) {
	domainProduct.ProductType = productDto.ProductType
	domainProduct.Code = productDto.Code
	domainProduct.LinkCode = productDto.LinkCode
	domainProduct.Name = productDto.Name
	domainProduct.Description = productDto.Description
	domainProduct.UnitType = productDto.UnitType
	domainProduct.UnitWeight = productDto.UnitWeight
	domainProduct.UnitWeightType = productDto.UnitWeightType
	domainProduct.LastUpdatedBy = productDto.LastUpdatedBy
	if productDto.Status.IsValid() {
		domainProduct.Status = productDto.Status
	}
	if domainProduct.UnitType != customtypes.UnitPiece {
		domainProduct.UnitWeight = 0
		domainProduct.UnitWeightType = customtypes.UnitWeightGram
	}
}
