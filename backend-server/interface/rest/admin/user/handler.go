package user

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/base/dto"
	"github.com/vamika-digital/wms-api-server/core/base/validators"
	"github.com/vamika-digital/wms-api-server/core/business/admin/dto/user"
	"github.com/vamika-digital/wms-api-server/core/business/admin/service"
)

type UserRestHandler struct {
	UserService service.UserService
}

func NewUserHandler(s service.UserService) *UserRestHandler {
	return &UserRestHandler{UserService: s}
}

func (h *UserRestHandler) GetAllUsers(c *gin.Context) {
	var filter = &user.UserFilterDto{}
	if err := c.ShouldBindQuery(filter); err != nil {
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
	data, totalCount, err := h.UserService.GetAllUsers(pageNumber, rowsPerPage, sort, filter)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}
	pageResponse := dto.GetPaginatedRestResponse(data, totalCount, pageNumber, rowsPerPage)
	c.JSON(http.StatusOK, pageResponse)
}

func (h *UserRestHandler) CreateUser(c *gin.Context) {
	var formDTO = &user.UserCreateDto{}
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

	if err := h.UserService.CreateUser(formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}
	c.Status(http.StatusCreated)
}

func (h *UserRestHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	user, err := h.UserService.GetUserByID(id)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusNotFound, dto.GetErrorRestResponse(http.StatusNotFound, err.Error(), nil))
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser updates an existing user
func (handler *UserRestHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	var formDTO = &user.UserUpdateDto{}
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

	if err := handler.UserService.UpdateUser(id, formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteUser deletes a user by its ID
func (handler *UserRestHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	if err := handler.UserService.DeleteUser(id); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusInternalServerError, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}

// DeleteUser deletes a role by its IDs
func (handler *UserRestHandler) DeleteUserByIDs(c *gin.Context) {
	formDTO := &dto.BatchDeleteDTO{}

	// Bind JSON payload to formDTO
	if err := c.ShouldBindJSON(&formDTO); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusBadRequest, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	// Delete roles using UserService
	if err := handler.UserService.DeleteUserByIDs(formDTO.IDs); err != nil {
		log.Printf("%+v\n", err)
		c.JSON(http.StatusInternalServerError, dto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil))
		return
	}

	c.Status(http.StatusOK)
}
