package service

import (
	"log"

	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
	"github.com/vamika-digital/wms-api-server/core/business/master/repository"
)

type ProductServiceImpl struct {
	ProductRepo      repository.ProductRepository
	ProductConverter *converter.ProductConverter
}

func NewProductService(productRepo repository.ProductRepository, productConverter *converter.ProductConverter) ProductService {
	return &ProductServiceImpl{ProductRepo: productRepo, ProductConverter: productConverter}
}

func (s *ProductServiceImpl) GetAllProducts(page int16, pageSize int16, sort string, filter *product.ProductFilterDto) ([]*product.ProductDto, int64, error) {
	totalCount, err := s.ProductRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainProducts, err := s.ProductRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	// Convert domain products to DTOs. You can do this based on your requirements.
	var productDtos []*product.ProductDto = s.ProductConverter.ToDtoSlice(domainProducts)
	return productDtos, int64(totalCount), nil
}

func (s *ProductServiceImpl) CreateProduct(productDto *product.ProductCreateDto) error {
	var newProduct *domain.Product = s.ProductConverter.ToDomain(productDto)
	err := s.ProductRepo.Create(newProduct)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *ProductServiceImpl) GetProductByID(productID int64) (*product.ProductDto, error) {
	domainProduct, err := s.ProductRepo.GetById(productID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.ProductConverter.ToDto(domainProduct), nil
}

func (s *ProductServiceImpl) GetMinimalProductByID(productID int64) (*product.ProductMinimalDto, error) {
	domainProduct, err := s.ProductRepo.GetById(productID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.ProductConverter.ToMinimalDto(domainProduct), nil
}

func (s *ProductServiceImpl) GetMinimalProductByIds(productIDs []int64) ([]*product.ProductMinimalDto, error) {
	domainProducts, err := s.ProductRepo.GetByIds(productIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.ProductConverter.ToMinimalDtoSlice(domainProducts), nil
}

func (s *ProductServiceImpl) GetProductByCode(productCode string) (*product.ProductDto, error) {
	domainProduct, err := s.ProductRepo.GetByCode(productCode)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.ProductConverter.ToDto(domainProduct), nil
}

func (s *ProductServiceImpl) UpdateProduct(productID int64, productDto *product.ProductUpdateDto) error {
	existingProduct, err := s.ProductRepo.GetById(productID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	s.ProductConverter.ToUpdateDomain(existingProduct, productDto)
	if err := s.ProductRepo.Update(existingProduct); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *ProductServiceImpl) DeleteProduct(productID int64) error {
	if err := s.ProductRepo.Delete(productID); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *ProductServiceImpl) DeleteProductByIDs(productIDs []int64) error {
	if err := s.ProductRepo.DeleteByIDs(productIDs); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}
