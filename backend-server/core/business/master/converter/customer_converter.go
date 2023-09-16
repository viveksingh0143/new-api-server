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
	if domainCustomer == nil {
		return nil
	}
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
			Address1: domainCustomer.BillingAddress1,
			Address2: domainCustomer.BillingAddress2,
			State:    domainCustomer.BillingState,
			Country:  domainCustomer.BillingCountry,
			Pincode:  domainCustomer.BillingPincode,
		},
		ShippingAddress: customer.ShippingAddress{
			Address1: domainCustomer.ShippingAddress1,
			Address2: domainCustomer.ShippingAddress2,
			State:    domainCustomer.ShippingState,
			Country:  domainCustomer.ShippingCountry,
			Pincode:  domainCustomer.ShippingPincode,
		},
		Status:        domainCustomer.Status,
		CreatedAt:     customtypes.NewValidNullTime(domainCustomer.CreatedAt),
		UpdatedAt:     customtypes.GetNullTime(domainCustomer.UpdatedAt),
		LastUpdatedBy: domainCustomer.LastUpdatedBy,
	}
	return customerDto
}

func (c *CustomerConverter) ToDtoSlice(domainCustomers []*domain.Customer) []*customer.CustomerDto {
	var customerDtos = make([]*customer.CustomerDto, 0)
	for _, domainCustomer := range domainCustomers {
		customerDtos = append(customerDtos, c.ToDto(domainCustomer))
	}
	return customerDtos
}

func (c *CustomerConverter) ToMinimalDtoSlice(domainCustomers []*domain.Customer) []*customer.CustomerMinimalDto {
	var customerDtos = make([]*customer.CustomerMinimalDto, 0)
	for _, domainCustomer := range domainCustomers {
		customerDtos = append(customerDtos, c.ToMinimalDto(domainCustomer))
	}
	return customerDtos
}

func (c *CustomerConverter) ToDomain(customerDto *customer.CustomerCreateDto) *domain.Customer {
	domainCustomer := &domain.Customer{
		Code:             customerDto.Code,
		Name:             customerDto.Name,
		ContactPerson:    customerDto.ContactPerson,
		BillingAddress1:  customerDto.BillingAddress.Address1,
		BillingAddress2:  customerDto.BillingAddress.Address2,
		BillingState:     customerDto.BillingAddress.State,
		BillingCountry:   customerDto.BillingAddress.Country,
		BillingPincode:   customerDto.BillingAddress.Pincode,
		ShippingAddress1: customerDto.ShippingAddress.Address1,
		ShippingAddress2: customerDto.ShippingAddress.Address2,
		ShippingState:    customerDto.ShippingAddress.State,
		ShippingCountry:  customerDto.ShippingAddress.Country,
		ShippingPincode:  customerDto.ShippingAddress.Pincode,
		Status:           customerDto.Status,
		LastUpdatedBy:    customerDto.LastUpdatedBy,
	}
	return domainCustomer
}

func (c *CustomerConverter) ToUpdateDomain(domainCustomer *domain.Customer, customerDto *customer.CustomerUpdateDto) {
	domainCustomer.Code = customerDto.Code
	domainCustomer.Name = customerDto.Name
	domainCustomer.ContactPerson = customerDto.ContactPerson
	domainCustomer.BillingAddress1 = customerDto.BillingAddress.Address1
	domainCustomer.BillingAddress2 = customerDto.BillingAddress.Address2
	domainCustomer.BillingState = customerDto.BillingAddress.State
	domainCustomer.BillingCountry = customerDto.BillingAddress.Country
	domainCustomer.BillingPincode = customerDto.BillingAddress.Pincode
	domainCustomer.ShippingAddress1 = customerDto.ShippingAddress.Address1
	domainCustomer.ShippingAddress2 = customerDto.ShippingAddress.Address2
	domainCustomer.ShippingState = customerDto.ShippingAddress.State
	domainCustomer.ShippingCountry = customerDto.ShippingAddress.Country
	domainCustomer.ShippingPincode = customerDto.ShippingAddress.Pincode
	if customerDto.Status.IsValid() {
		domainCustomer.Status = customerDto.Status
	}
	domainCustomer.LastUpdatedBy = customerDto.LastUpdatedBy
}
