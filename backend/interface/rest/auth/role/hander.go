package role

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/common/decoders"
	commonDto "github.com/vamika-digital/wms-api-server/common/dto"
	"github.com/vamika-digital/wms-api-server/common/validators"
	"github.com/vamika-digital/wms-api-server/internal/auth/domain"
	"github.com/vamika-digital/wms-api-server/internal/auth/dto/role"
	"github.com/vamika-digital/wms-api-server/internal/auth/repository"
	"github.com/vamika-digital/wms-api-server/internal/auth/service"
)

type RoleRestHandler struct {
	RoleService service.RoleService
}

func NewRoleHandler(service service.RoleService) *RoleRestHandler {
	return &RoleRestHandler{RoleService: service}
}

func (handler *RoleRestHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	pageProps := new(commonDto.PaginationProps)
	decoder := decoders.CreateRequestDecoder()
	if err := decoder.Decode(&pageProps, r.Form); err != nil {
		errResp := commonDto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil)
		_ = json.NewEncoder(w).Encode(errResp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filter := new(repository.RoleFilterOptions)
	if err := decoder.Decode(filter, r.Form); err != nil {
		errResp := commonDto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil)
		_ = json.NewEncoder(w).Encode(errResp)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pageNumber, rowsPerPage, sort := pageProps.GetValues()
	data, totalCount, err := handler.RoleService.GetAllRoles(pageNumber, rowsPerPage, sort, *filter)

	if err != nil {
		errResp := commonDto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), nil)
		_ = json.NewEncoder(w).Encode(errResp)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	pageResponse := commonDto.GetPaginatedRestResponse(data, totalCount, pageNumber, rowsPerPage)
	if err := json.NewEncoder(w).Encode(pageResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (handler *RoleRestHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var formDTO role.RoleCreateDto

	if err := json.NewDecoder(r.Body).Decode(&formDTO); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	validate := validators.NewValidator()
	err := validate.Struct(formDTO)
	if err != nil {
		errors := validators.GetAllErrors(err)
		errResp := commonDto.GetErrorRestResponse(http.StatusBadRequest, err.Error(), errors)
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(errResp)
		return
	}

	if err := handler.RoleService.CreateRole(formDTO); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// func (handler *RoleRestHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
// 	id, err := strconv.ParseInt(params["id"], 10, 64)
// 	if err != nil {
// 		http.Error(w, "Invalid Role ID", http.StatusBadRequest)
// 		return
// 	}

// 	role, err := handler.UseCase.GetRoleByID(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	if err := json.NewEncoder(w).Encode(role); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

func (handler *RoleRestHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Invalid Role ID", http.StatusBadRequest)
		return
	}

	var role domain.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validateRole(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	role.ID = int64(id)
	if err := handler.UseCase.UpdateRole(&role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (handler *RoleRestHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		http.Error(w, "invalid Role ID", http.StatusBadRequest)
		return
	}

	if err := handler.UseCase.DeleteRole(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// func validateRole(role *domain.Role) error {
// typeErr := role.ValidateType()
// if typeErr != nil {
// 	return errors.New("type should be valid")
// }
// if role.Code == "" {
// 	return errors.New("code is required")
// }
// if role.Name == "" {
// 	return errors.New("name is required")
// }
// if role.Status == "" {
// 	return errors.New("status is required")
// }
// 	return nil
// }
