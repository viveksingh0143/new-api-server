package converter

import (
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/business/master/domain"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/customer"
)

type CustomerConverter struct {
}

func NewCustomerConverter() *CustomerConverter {
	return &CustomerConverter{}
}

func (c *CustomerConverter) ToMinimalDto(domainCustomer *domain.Customer) *customer.CustomerMinimalDto {
	customerDto := &customer.CustomerMinimalDto{
		ID:     domainCustomer.ID,
		Code:   domainCustomer.Code,
		Name:   domainCustomer.Name,
		Status: domainCustomer.Status,
	}
	return customerDto
}

func (c *CustomerConverter) ToDto(domainCustomer *domain.Customer) *customer.CustomerDto {
	customerDto := &customer.CustomerDto{
		ID:            domainCustomer.ID,
		Code:          domainCustomer.Code,
		Name:          domainCustomer.Name,
		ContactPerson: domainCustomer.ContactPerson,
		BillingAddress: customer.BillingAddress{
			Address1: domainCustomer.BillingAddress.Address1,
			Address2: domainCustomer.BillingAddress.Address2,
			State:    domainCustomer.BillingAddress.State,
			Country:  domainCustomer.BillingAddress.Country,
			Pincode:  domainCustomer.BillingAddress.Pincode,
		},
		ShippingAddress: customer.ShippingAddress{
			Address1: domainCustomer.ShippingAddress.Address1,
			Address2: domainCustomer.ShippingAddress.Address2,
			State:    domainCustomer.ShippingAddress.State,
			Country:  domainCustomer.ShippingAddress.Country,
			Pincode:  domainCustomer.ShippingAddress.Pincode,
		},
		Status:        domainCustomer.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainCustomer.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainCustomer.UpdatedAt),
		LastUpdatedBy: domainCustomer.LastUpdatedBy,
	}
	return customerDto
}

func (c *CustomerConverter) ToDtoSlice(domainCustomers []*domain.Customer) []*customer.CustomerDto {
	var customerDtos []*customer.CustomerDto
	for _, domainCustomer := range domainCustomers {
		customerDtos = append(customerDtos, c.ToDto(domainCustomer))
	}
	return customerDtos
}

func (c *CustomerConverter) ToDomain(customerDto *customer.CustomerCreateDto) *domain.Customer {
	domainCustomer := &domain.Customer{
		Code:          customerDto.Code,
		Name:          customerDto.Name,
		ContactPerson: customerDto.ContactPerson,
		BillingAddress: domain.BillingAddress{
			Address1: customerDto.BillingAddress.Address1,
			Address2: customerDto.BillingAddress.Address2,
			State:    customerDto.BillingAddress.State,
			Country:  customerDto.BillingAddress.Country,
			Pincode:  customerDto.BillingAddress.Pincode,
		},
		ShippingAddress: domain.ShippingAddress{
			Address1: customerDto.ShippingAddress.Address1,
			Address2: customerDto.ShippingAddress.Address2,
			State:    customerDto.ShippingAddress.State,
			Country:  customerDto.ShippingAddress.Country,
			Pincode:  customerDto.ShippingAddress.Pincode,
		},
		Status:        customerDto.Status,
		LastUpdatedBy: customerDto.LastUpdatedBy,
	}
	return domainCustomer
}

func (c *CustomerConverter) ToUpdateDomain(domainCustomer *domain.Customer, customerDto *customer.CustomerUpdateDto) {
	domainCustomer.Code = customerDto.Code
	domainCustomer.Name = customerDto.Name
	domainCustomer.ContactPerson = customerDto.ContactPerson
	domainCustomer.BillingAddress = domain.BillingAddress{
		Address1: customerDto.BillingAddress.Address1,
		Address2: customerDto.BillingAddress.Address2,
		State:    customerDto.BillingAddress.State,
		Country:  customerDto.BillingAddress.Country,
		Pincode:  customerDto.BillingAddress.Pincode,
	}
	domainCustomer.ShippingAddress = domain.ShippingAddress{
		Address1: customerDto.ShippingAddress.Address1,
		Address2: customerDto.ShippingAddress.Address2,
		State:    customerDto.ShippingAddress.State,
		Country:  customerDto.ShippingAddress.Country,
		Pincode:  customerDto.ShippingAddress.Pincode,
	}
	if customerDto.Status.IsValid() {
		domainCustomer.Status = customerDto.Status
	}
	domainCustomer.LastUpdatedBy = customerDto.LastUpdatedBy
}
