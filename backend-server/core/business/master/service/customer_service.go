package service

import "github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"

type CustomerService interface {
	GetAllCustomers(page int64, pageSize int64, sort string, filter *customer.CustomerFilterDto) ([]*customer.CustomerDto, int64, error)
	CreateCustomer(customerDto *customer.CustomerCreateDto) error
	GetCustomerByID(customerID int64) (*customer.CustomerDto, error)
	UpdateCustomer(customerID int64, customer *customer.CustomerUpdateDto) error
	DeleteCustomer(customerID int64) error
	DeleteCustomerByIDs(customerIDs []int64) error
}
