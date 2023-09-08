package role

import (
	"github.com/gin-gonic/gin"
	"github.com/vamika-digital/wms-api-server/core/business/admin/service"
)

type RoleRestModule struct {
	Handler *RoleRestHandler
}

func NewRoleRestModule(roleService service.RoleService) *RoleRestModule {
	roleHandler := NewRoleHandler(roleService)
	return &RoleRestModule{Handler: roleHandler}
}

func (m *RoleRestModule) RegisterRoutes(r *gin.RouterGroup) {
	roleGroup := r.Group("/roles")
	{
		roleGroup.POST("", m.Handler.CreateRole)
		roleGroup.GET("", m.Handler.GetAllRoles)
		roleGroup.GET("/:id", m.Handler.GetRoleByID)
		roleGroup.PUT("/:id", m.Handler.UpdateRole)
		roleGroup.DELETE("/:id", m.Handler.DeleteRole)
	}
}
