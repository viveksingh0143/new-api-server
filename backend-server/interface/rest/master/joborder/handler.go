package joborder

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/base/dto"
	"github.com/vamika-digital/wms-api-server/core/base/validators"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/joborder"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type JobOrderRestHandler struct {
	JobOrderService service.JobOrderService
}

func NewJobOrderHandler(s service.JobOrderService) *JobOrderRestHandler {
	return &JobOrderRestHandler{JobOrderService: s}
}

func (h *JobOrderRestHandler) GetAllJobOrders(c *gin.Context) {
	var filter = &joborder.JobOrderFilterDto{}
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
	data, totalCount, err := h.JobOrderService.GetAllJobOrders(pageNumber, rowsPerPage, sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	pageResponse := dto.GetPaginatedRestResponse(data, totalCount, pageNumber, rowsPerPage)
	c.JSON(http.StatusOK, pageResponse)
}

func (h *JobOrderRestHandler) CreateJobOrder(c *gin.Context) {
	var formDTO = &joborder.JobOrderCreateDto{}
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

	if err := h.JobOrderService.CreateJobOrder(formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.Status(http.StatusCreated)
}

func (h *JobOrderRestHandler) GetJobOrderByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	joborder, err := h.JobOrderService.GetJobOrderByID(id)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusNotFound, "Resource not found", nil))
		return
	}
	c.JSON(http.StatusOK, joborder)
}

// UpdateJobOrder updates an existing joborder
func (handler *JobOrderRestHandler) UpdateJobOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	var formDTO = &joborder.JobOrderUpdateDto{}
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

	if err := handler.JobOrderService.UpdateJobOrder(id, formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteJobOrder deletes a joborder by its ID
func (handler *JobOrderRestHandler) DeleteJobOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := handler.JobOrderService.DeleteJobOrder(id); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteJobOrder deletes a joborder by its IDs
func (handler *JobOrderRestHandler) DeleteJobOrderByIDs(c *gin.Context) {
	formDTO := &dto.BatchDeleteDTO{}

	// Bind JSON payload to formDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Delete joborders using JobOrderService
	if err := handler.JobOrderService.DeleteJobOrderByIDs(formDTO.IDs); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}
