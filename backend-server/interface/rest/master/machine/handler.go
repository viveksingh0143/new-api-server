package machine

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/base/dto"
	"github.com/vamika-digital/wms-api-server/core/base/validators"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/machine"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type MachineRestHandler struct {
	MachineService service.MachineService
}

func NewMachineHandler(s service.MachineService) *MachineRestHandler {
	return &MachineRestHandler{MachineService: s}
}

func (h *MachineRestHandler) GetAllMachines(c *gin.Context) {
	var filter = &machine.MachineFilterDto{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	var pageProps dto.PaginationProps
	if err := c.ShouldBindQuery(&pageProps); err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	pageNumber, rowsPerPage, sort := pageProps.GetValues()
	data, totalCount, err := h.MachineService.GetAllMachines(pageNumber, rowsPerPage, sort, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	pageResponse := dto.GetPaginatedRestResponse(data, totalCount, pageNumber, rowsPerPage)
	c.JSON(http.StatusOK, pageResponse)
}

func (h *MachineRestHandler) CreateMachine(c *gin.Context) {
	var formDTO = &machine.MachineCreateDto{}
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(formDTO); err != nil {
		errors := validators.GetAllErrors(err, formDTO)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	if err := h.MachineService.CreateMachine(formDTO); err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.Status(http.StatusCreated)
}

func (h *MachineRestHandler) GetMachineByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	machine, err := h.MachineService.GetMachineByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, machine)
}

// UpdateMachine updates an existing machine
func (handler *MachineRestHandler) UpdateMachine(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	var formDTO = &machine.MachineUpdateDto{}
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(formDTO); err != nil {
		errors := validators.GetAllErrors(err, formDTO)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	if err := handler.MachineService.UpdateMachine(id, formDTO); err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteMachine deletes a machine by its ID
func (handler *MachineRestHandler) DeleteMachine(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := handler.MachineService.DeleteMachine(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteMachine deletes a machine by its IDs
func (handler *MachineRestHandler) DeleteMachineByIDs(c *gin.Context) {
	formDTO := &dto.BatchDeleteDTO{}

	// Bind JSON payload to formDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Delete machines using MachineService
	if err := handler.MachineService.DeleteMachineByIDs(formDTO.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}
