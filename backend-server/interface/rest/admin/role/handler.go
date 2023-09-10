package role

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/base/dto"
	"github.com/vamika-digital/wms-api-server/core/base/validators"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/role"
	"github.com/vamika-digital/wms-api-server/core/business/admin/service"
)

type RoleRestHandler struct {
	RoleService service.RoleService
}

func NewRoleHandler(s service.RoleService) *RoleRestHandler {
	return &RoleRestHandler{RoleService: s}
}

func (h *RoleRestHandler) GetAllRoles(c *gin.Context) {
	var filter = &role.RoleFilterDto{}
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
	data, totalCount, err := h.RoleService.GetAllRoles(pageNumber, rowsPerPage, sort, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	pageResponse := dto.GetPaginatedRestResponse(data, totalCount, pageNumber, rowsPerPage)
	c.JSON(http.StatusOK, pageResponse)
}

func (h *RoleRestHandler) CreateRole(c *gin.Context) {
	var formDTO = &role.RoleCreateDto{}
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(formDTO); err != nil {
		errors := validators.GetAllErrors(err, formDTO)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	if err := h.RoleService.CreateRole(formDTO); err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}
	c.Status(http.StatusCreated)
}

func (h *RoleRestHandler) GetRoleByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	role, err := h.RoleService.GetRoleByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, role)
}

// UpdateRole updates an existing role
func (handler *RoleRestHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	var formDTO = &role.RoleUpdateDto{}
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := validators.Validate.Struct(formDTO); err != nil {
		errors := validators.GetAllErrors(err, formDTO)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, "Please fill the form correctly", errors))
		return
	}

	if err := handler.RoleService.UpdateRole(id, formDTO); err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteRole deletes a role by its ID
func (handler *RoleRestHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := handler.RoleService.DeleteRole(id); err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteRole deletes a role by its IDs
func (handler *RoleRestHandler) DeleteRoleByIDs(c *gin.Context) {
	formDTO := &dto.BatchDeleteDTO{}

	// Bind JSON payload to formDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Delete roles using RoleService
	if err := handler.RoleService.DeleteRoleByIDs(formDTO.IDs); err != nil {
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}
