package container

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/base/dto"
	"github.com/vamika-digital/wms-api-server/core/base/validators"
	"github.com/vamika-digital/wms-api-server/core/business/master/dto/container"
	"github.com/vamika-digital/wms-api-server/core/business/master/service"
)

type ContainerRestHandler struct {
	ContainerService service.ContainerService
}

func NewContainerHandler(s service.ContainerService) *ContainerRestHandler {
	return &ContainerRestHandler{ContainerService: s}
}

func (h *ContainerRestHandler) GetAllContainers(c *gin.Context) {
	var filter = &container.ContainerFilterDto{}
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
	data, totalCount, err := h.ContainerService.GetAllContainers(pageNumber, rowsPerPage, sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	pageResponse := dto.GetPaginatedRestResponse(data, totalCount, pageNumber, rowsPerPage)
	c.JSON(http.StatusOK, pageResponse)
}

func (h *ContainerRestHandler) CreateContainer(c *gin.Context) {
	var formDTO = &container.ContainerCreateDto{}
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

	if err := h.ContainerService.CreateContainer(formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.Status(http.StatusCreated)
}

func (h *ContainerRestHandler) GetContainerByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	container, err := h.ContainerService.GetContainerByID(id)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusNotFound, "Resource not found", nil))
		return
	}
	c.JSON(http.StatusOK, container)
}

// UpdateContainer updates an existing container
func (handler *ContainerRestHandler) UpdateContainer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	var formDTO = &container.ContainerUpdateDto{}
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

	if err := handler.ContainerService.UpdateContainer(id, formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteContainer deletes a container by its ID
func (handler *ContainerRestHandler) DeleteContainer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := handler.ContainerService.DeleteContainer(id); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteContainer deletes a container by its IDs
func (handler *ContainerRestHandler) DeleteContainerByIDs(c *gin.Context) {
	formDTO := &dto.BatchDeleteDTO{}

	// Bind JSON payload to formDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Delete containers using ContainerService
	if err := handler.ContainerService.DeleteContainerByIDs(formDTO.IDs); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

func (handler *ContainerRestHandler) GetContainerCodeInfoDto(c *gin.Context) {
	containerInfo, err := handler.ContainerService.GetContainerCodeInfoDto()
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusNotFound, "No information found", nil))
		return
	}

	c.JSON(http.StatusOK, containerInfo)
}
