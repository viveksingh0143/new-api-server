package role

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/internal/auth/service"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type RoleRestModule struct {
	Handler *RoleRestHandler
}

func NewRoleRestModule(db database.Connection) *RoleRestModule {
	// roleRepo := repository.NewRoleRepository(db)
	roleService := service.NewRoleService()
	roleHandler := NewRoleHandler(roleService)
	return &RoleRestModule{Handler: roleHandler}
}

func (m *RoleRestModule) RegisterRoutes(r *mux.Router) {
	subRouter := r.PathPrefix("/roles").Subrouter()
	subRouter.HandleFunc("", m.Handler.CreateRole).Methods(http.MethodPost)
	subRouter.HandleFunc("", m.Handler.GetAllRoles).Methods(http.MethodGet, http.MethodOptions)
	// subRouter.HandleFunc("/{id}", m.Handler.GetRoleByID).Methods(http.MethodGet, http.MethodOptions)
	// subRouter.HandleFunc("/{id}", m.Handler.UpdateRole).Methods(http.MethodPut)
	// subRouter.HandleFunc("/{id}", m.Handler.DeleteRole).Methods(http.MethodDelete)
}
