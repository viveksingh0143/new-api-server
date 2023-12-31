package service

import (
	"log"

	"github.com/vamika-digital/wms-api-server/core/business/master/converter"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
	"github.com/vamika-digital/wms-api-server/core/business/master/repository"
)

type CustomerServiceImpl struct {
	CustomerRepo      repository.CustomerRepository
	CustomerConverter *converter.CustomerConverter
}

func NewCustomerService(customerRepo repository.CustomerRepository, customerConverter *converter.CustomerConverter) CustomerService {
	return &CustomerServiceImpl{CustomerRepo: customerRepo, CustomerConverter: customerConverter}
}

func (s *CustomerServiceImpl) GetAllCustomers(page int16, pageSize int16, sort string, filter *customer.CustomerFilterDto) ([]*customer.CustomerDto, int64, error) {
	totalCount, err := s.CustomerRepo.GetTotalCount(filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}
	domainCustomers, err := s.CustomerRepo.GetAll(int(page), int(pageSize), sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, 0, err
	}

	var customerDtos []*customer.CustomerDto = s.CustomerConverter.ToDtoSlice(domainCustomers)
	return customerDtos, int64(totalCount), nil
}

func (s *CustomerServiceImpl) CreateCustomer(customerDto *customer.CustomerCreateDto) error {
	var newCustomer *domain.Customer = s.CustomerConverter.ToDomain(customerDto)
	if err := s.CustomerRepo.Create(newCustomer); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *CustomerServiceImpl) GetCustomerByID(customerID int64) (*customer.CustomerDto, error) {
	domainCustomer, err := s.CustomerRepo.GetById(customerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.CustomerConverter.ToDto(domainCustomer), nil
}

func (s *CustomerServiceImpl) GetMinimalCustomerByID(customerID int64) (*customer.CustomerMinimalDto, error) {
	domainCustomer, err := s.CustomerRepo.GetById(customerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.CustomerConverter.ToMinimalDto(domainCustomer), nil
}

func (s *CustomerServiceImpl) GetMinimalCustomerByIds(customerIDs []int64) ([]*customer.CustomerMinimalDto, error) {
	domainCustomers, err := s.CustomerRepo.GetByIds(customerIDs)
	if err != nil {
		log.Printf("%+v\n", err)
		return nil, err
	}
	return s.CustomerConverter.ToMinimalDtoSlice(domainCustomers), nil
}

func (s *CustomerServiceImpl) UpdateCustomer(customerID int64, customerDto *customer.CustomerUpdateDto) error {
	existingCustomer, err := s.CustomerRepo.GetById(customerID)
	if err != nil {
		log.Printf("%+v\n", err)
		return err
	}

	s.CustomerConverter.ToUpdateDomain(existingCustomer, customerDto)
	if err := s.CustomerRepo.Update(existingCustomer); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *CustomerServiceImpl) DeleteCustomer(customerID int64) error {
	if err := s.CustomerRepo.Delete(customerID); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}

func (s *CustomerServiceImpl) DeleteCustomerByIDs(customerIDs []int64) error {
	if err := s.CustomerRepo.DeleteByIDs(customerIDs); err != nil {
		log.Printf("%+v\n", err)
		return err
	}
	return nil
}
