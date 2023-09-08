package auth

import (
	"github.com/gorilla/mux"
	"github.com/vamika-digital/wms-api-server/interface/rest/auth/role"
	"github.com/vamika-digital/wms-api-server/pkg/database"
)

type AuthRestModule struct {
	RoleModule *role.RoleRestModule
	// RoleModule  *role.RoleModule
	// StoreModule *rest.StoreModule
}

func NewAuthRestModule(db database.Connection) *AuthRestModule {
	roleModule := role.NewRoleRestModule(db)
	return &AuthRestModule{RoleModule: roleModule}
}

func (m *AuthRestModule) RegisterRoutes(r *mux.Router) {
	m.RoleModule.RegisterRoutes(r.PathPrefix("/admin").Subrouter())
}
