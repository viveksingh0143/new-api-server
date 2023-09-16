package product

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/base/dto"
	"github.com/vamika-digital/wms-api-server/core/base/validators"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/product"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type ProductRestHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(s service.ProductService) *ProductRestHandler {
	return &ProductRestHandler{ProductService: s}
}

func (h *ProductRestHandler) GetAllProducts(c *gin.Context) {
	var filter = &product.ProductFilterDto{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	var pageProps dto.PaginationProps
	if err := c.ShouldBindQuery(&pageProps); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	pageNumber, rowsPerPage, sort := pageProps.GetValues()
	data, totalCount, err := h.ProductService.GetAllProducts(pageNumber, rowsPerPage, sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	pageResponse := dto.GetPaginatedRestResponse(data, totalCount, pageNumber, rowsPerPage)
	c.JSON(http.StatusOK, pageResponse)
}

func (h *ProductRestHandler) CreateProduct(c *gin.Context) {
	var formDTO = &product.ProductCreateDto{}
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(formDTO); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, formDTO)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	if err := h.ProductService.CreateProduct(formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.Status(http.StatusCreated)
}

func (h *ProductRestHandler) GetProductByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	product, err := h.ProductService.GetProductByID(id)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusNotFound, "Resource not found", nil))
		return
	}
	c.JSON(http.StatusOK, product)
}

// UpdateProduct updates an existing product
func (handler *ProductRestHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	var formDTO = &product.ProductUpdateDto{}
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(formDTO); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, formDTO)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	if err := handler.ProductService.UpdateProduct(id, formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteProduct deletes a product by its ID
func (handler *ProductRestHandler) DeleteProduct(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := handler.ProductService.DeleteProduct(id); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteProduct deletes a product by its IDs
func (handler *ProductRestHandler) DeleteProductByIDs(c *gin.Context) {
	formDTO := &dto.BatchDeleteDTO{}

	// Bind JSON payload to formDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Delete products using ProductService
	if err := handler.ProductService.DeleteProductByIDs(formDTO.IDs); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}
