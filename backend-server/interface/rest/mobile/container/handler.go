package container

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/base/customtypes"
	"github.com/vamika-digital/wms-api-server/core/base/dto"
	"github.com/vamika-digital/wms-api-server/core/base/validators"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type ContainerRestHandler struct {
	ContainerService service.ContainerService
}

func NewContainerHandler(s service.ContainerService) *ContainerRestHandler {
	return &ContainerRestHandler{ContainerService: s}
}

type MarkContainerParams struct {
	Code string `form:"code" binding:"required"`
}

type ContainerParams struct {
	Code          string `form:"code" binding:"required"`
	ContainerType string `form:"type" binding:"required,oneof=Pallet Bin Rack pallet bin rack PALLET BIN RACK"`
}

func (h *ContainerRestHandler) GetOneContainerByCodeAndType(c *gin.Context) {
	var containerParams ContainerParams
	if err := c.ShouldBindQuery(&containerParams); err == nil {
		containerParams.ContainerType = strings.ToLower(containerParams.ContainerType)
	} else {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Invalid Inputs", nil))
		return
	}

	// Validate the struct
	if err := validators.Validate.Struct(containerParams); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, containerParams)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	result, err := h.ContainerService.GetOneActiveContainerByCodeAndType(containerParams.Code, customtypes.ContainerType(containerParams.ContainerType))
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ContainerRestHandler) MarkContainerFullByCode(c *gin.Context) {
	var containerParams = &MarkContainerParams{}
	if err := c.ShouldBindJSON(&containerParams); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Invalid Inputs", nil))
		return
	}

	if err := validators.Validate.Struct(containerParams); err != nil {
		log.Printf("%+v\n", err)
		errors := validators.GetAllErrors(err, containerParams)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	err := h.ContainerService.MarkContainerFullByCode(containerParams.Code)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, dto.GetRestResponse(http.StatusOK, "Successfully marked full..."))
}
