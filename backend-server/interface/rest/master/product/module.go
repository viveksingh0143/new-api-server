package product

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type ProductRestModule struct {
	Handler *ProductRestHandler
}

func NewProductRestModule(productService service.ProductService) *ProductRestModule {
	productHandler := NewProductHandler(productService)
	return &ProductRestModule{Handler: productHandler}
}

func (m *ProductRestModule) RegisterRoutes(r *gin.RouterGroup) {
	productGroup := r.Group("/products")
	{
		productGroup.POST("", m.Handler.CreateProduct)
		productGroup.GET("", m.Handler.GetAllProducts)
		productGroup.POST("/bulkdelete", m.Handler.DeleteProductByIDs)
		productGroup.GET("/:id", m.Handler.GetProductByID)
		productGroup.PUT("/:id", m.Handler.UpdateProduct)
		productGroup.DELETE("/:id", m.Handler.DeleteProduct)
	}
}
