package customer

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type CustomerRestModule struct {
	Handler *CustomerRestHandler
}

func NewCustomerRestModule(customerService service.CustomerService) *CustomerRestModule {
	customerHandler := NewCustomerHandler(customerService)
	return &CustomerRestModule{Handler: customerHandler}
}

func (m *CustomerRestModule) RegisterRoutes(r *gin.RouterGroup) {
	customerGroup := r.Group("/customers")
	{
		customerGroup.POST("", m.Handler.CreateCustomer)
		customerGroup.GET("", m.Handler.GetAllCustomers)
		customerGroup.GET("/:id", m.Handler.GetCustomerByID)
		customerGroup.PUT("/:id", m.Handler.UpdateCustomer)
		customerGroup.DELETE("/:id", m.Handler.DeleteCustomer)
	}
}
