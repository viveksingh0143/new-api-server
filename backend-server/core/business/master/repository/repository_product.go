package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
)

type ProductRepository interface {
	Create(product *domain.Product) error
	Update(product *domain.Product) error
	Delete(productID int64) error
	DeleteByIDs(productIDs []int64) error
	GetById(productID int64) (*domain.Product, error)
	GetByCode(productCode string) (*domain.Product, error)
	GetTotalCount(filter *product.ProductFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *product.ProductFilterDto) ([]*domain.Product, error)
}
