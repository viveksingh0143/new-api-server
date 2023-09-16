package service

import "github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"

type CustomerService interface {
	GetAllCustomers(page int16, pageSize int16, sort string, filter *customer.CustomerFilterDto) ([]*customer.CustomerDto, int64, error)
	CreateCustomer(customerDto *customer.CustomerCreateDto) error
	GetCustomerByID(customerID int64) (*customer.CustomerDto, error)
	GetMinimalCustomerByID(customerID int64) (*customer.CustomerMinimalDto, error)
	GetMinimalCustomerByIds(customerIDs []int64) ([]*customer.CustomerMinimalDto, error)
	UpdateCustomer(customerID int64, customer *customer.CustomerUpdateDto) error
	DeleteCustomer(customerID int64) error
	DeleteCustomerByIDs(customerIDs []int64) error
}
