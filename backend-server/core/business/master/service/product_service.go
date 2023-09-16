package service

import "github.com/vamika-digital/wms-api-server/core/business/master/dto/product"

type ProductService interface {
	GetAllProducts(page int16, pageSize int16, sort string, filter *product.ProductFilterDto) ([]*product.ProductDto, int64, error)
	CreateProduct(productDto *product.ProductCreateDto) error
	GetProductByID(productID int64) (*product.ProductDto, error)
	GetMinimalProductByID(productID int64) (*product.ProductMinimalDto, error)
	GetMinimalProductByIds(productIDs []int64) ([]*product.ProductMinimalDto, error)
	GetProductByCode(productCode string) (*product.ProductDto, error)
	UpdateProduct(productID int64, product *product.ProductUpdateDto) error
	DeleteProduct(productID int64) error
	DeleteProductByIDs(productIDs []int64) error
}
