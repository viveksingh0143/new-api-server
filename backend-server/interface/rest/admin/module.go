package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/admin/service"
	"github.com/vamika-digital/wms-api-server/interface/rest/admin/role"
	"github.com/vamika-digital/wms-api-server/interface/rest/admin/user"
)

type AdminRestModule struct {
	RoleModule *role.RoleRestModule
	UserModule *user.UserRestModule
}

func NewAdminRestModule(roleService service.RoleService, userService service.UserService) *AdminRestModule {
	roleModule := role.NewRoleRestModule(roleService)
	userModule := user.NewUserRestModule(userService)
	return &AdminRestModule{RoleModule: roleModule, UserModule: userModule}
}

func (m *AdminRestModule) RegisterRoutes(r *gin.RouterGroup) {
	adminGroup := r.Group("/admin")
	m.RoleModule.RegisterRoutes(adminGroup)
	m.UserModule.RegisterRoutes(adminGroup)
}
