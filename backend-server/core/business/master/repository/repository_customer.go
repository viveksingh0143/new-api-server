package repository

import (
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
)

type CustomerRepository interface {
	Create(customer *domain.Customer) error
	Update(customer *domain.Customer) error
	Delete(customerID int64) error
	DeleteByIDs(customerIDs []int64) error
	GetById(customerID int64) (*domain.Customer, error)
	GetTotalCount(filter *customer.CustomerFilterDto) (int, error)
	GetAll(page int, pageSize int, sort string, filter *customer.CustomerFilterDto) ([]*domain.Customer, error)
}
